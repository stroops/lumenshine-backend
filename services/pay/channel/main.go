package channel

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Soneso/lumenshine-backend/db/querying"
	"github.com/Soneso/lumenshine-backend/helpers"
	m "github.com/Soneso/lumenshine-backend/services/db/models"
	"github.com/Soneso/lumenshine-backend/services/pay/config"
	"github.com/Soneso/lumenshine-backend/services/pay/db"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/support/log"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var (
	distMutex = &sync.Mutex{}
)

// Configurator is holding the configuration for the channel manager
type Configurator struct {
	DB      *db.DB
	log     *logrus.Entry
	cnf     *config.Config
	Horizon *horizon.Client
}

//Manager is used for handling the channels
type Manager struct {
	*Configurator
}

//Channel represents one channel
type Channel struct {
	ID   int
	PK   string
	Seed string
}

//NewChanneler initialises the channel manager
func NewChanneler(DB *db.DB, cnf *config.Config) *Manager {
	var err error

	l := new(Configurator)
	l.DB = DB
	l.cnf = cnf
	l.log = helpers.GetDefaultLog("Channel", "")

	l.Horizon = &horizon.Client{
		URL: cnf.Stellar.Horizon,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	root, err := l.Horizon.Root()
	if err != nil {
		log.WithField("err", err).Error("Error loading Horizon root")
		os.Exit(-1)
	}

	if root.NetworkPassphrase != cnf.Stellar.NetworkPassphrase {
		log.Errorf("Invalid network passphrase (have=%s, want=%s)", root.NetworkPassphrase, cnf.Stellar.NetworkPassphrase)
		os.Exit(-1)
	}

	l.log.Info("Channeler created")
	mgr := &Manager{Configurator: l}

	return mgr
}

//GetChannel returns a free channel from the channel-pool
func (mg *Manager) GetChannel(DistPK string, PreSignerSeed string, PostSignerSeed string) (*Channel, error) {
	var v m.Channel
	sqlStr := querying.GetSQLKeyString(`update @channel set @status=$1 where id =
			(select id from @channel where @status=$2 and seed <> '' limit 1 for update) returning *`,
		map[string]string{
			"@channel": m.TableNames.Channels,
			"@status":  m.ChannelColumns.Status,
		})

	err := queries.Raw(sqlStr, m.ChannelStatusInUse, m.ChannelStatusFree).Bind(nil, mg.DB, &v)
	if err != nil {
		if err == sql.ErrNoRows {
			//create the channel and use this one
			return mg.createChannel(DistPK, PreSignerSeed, PostSignerSeed)
		}
		return nil, err
	}

	return &Channel{ID: v.ID, PK: v.PK, Seed: v.Seed}, nil
}

//createChannel creates the channel in the network and returns it
//DistPK is the distribution pk for the ico-phase. we need this info in order to lock the dist-pk as long, as we create a new channel
//we will try to lock the dist pk for a defined time, and if we don't succeed, we will return an error
func (mg *Manager) createChannel(DistPK string, PreSignerSeed string, PostSignerSeed string) (*Channel, error) {
	defer mg.ReleaseChannel(DistPK)
	err := mg.lockDistChannel(DistPK)
	if err != nil {
		return nil, fmt.Errorf("could not lock dist-pk %s", DistPK)
	}

	kp, err := keypair.Random()
	if err != nil {
		return nil, err
	}

	muts := []build.TransactionMutator{
		build.SourceAccount{AddressOrSeed: DistPK},
		build.Network{Passphrase: mg.cnf.Stellar.NetworkPassphrase},
		build.AutoSequence{SequenceProvider: mg.Horizon},
		build.CreateAccount(
			build.Destination{AddressOrSeed: kp.Address()},
			build.NativeAmount{Amount: mg.cnf.StellarChannelStartBalanceXLM},
		),
	}

	tx, err := build.Transaction(muts...)
	if err != nil {
		return nil, err
	}

	txe, err := tx.Sign(PreSignerSeed, PostSignerSeed)
	if err != nil {
		return nil, err
	}
	txs, err := txe.Base64()
	if err != nil {
		return nil, err
	}

	_, err = mg.Horizon.SubmitTransaction(txs)
	if err != nil {
		mg.log.WithError(err).Error("Could not submit transaction")
		return nil, err
	}

	//insert channel in DB
	var c m.Channel
	c.PK = kp.Address()
	c.Seed = kp.Seed()
	c.Status = m.ChannelStatusInUse

	err = c.Insert(mg.DB, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &Channel{ID: c.ID, PK: c.PK, Seed: c.Seed}, err
}

//lockDistChannel creates or locks the dist channel in the database
//the dist-channel is created if it does not exist yet in the db
func (mg *Manager) lockDistChannel(DistPK string) error {
	defer distMutex.Unlock()
	distMutex.Lock()

	cnt, err := m.Channels(qm.Where(m.ChannelColumns.PK+"=?", DistPK)).Count(mg.DB)
	if err != nil {
		return err
	}
	if cnt == 0 {
		//insert channel into db
		var cDB m.Channel
		cDB.PK = DistPK
		cDB.Status = m.ChannelStatusInUse
		err := cDB.Insert(mg.DB, boil.Infer())
		if err != nil {
			return err
		}
		return nil
	}

	c1 := make(chan error)
	go func(pk string) {
		var v m.Channel
		sqlStr := querying.GetSQLKeyString(`update @channel set @status=$1 where id =
			(select id from @channel where @status=$2 and @pk=$3 limit 1 for update) returning *`,
			map[string]string{
				"@channel": m.TableNames.Channels,
				"@status":  m.ChannelColumns.Status,
				"@pk":      m.ChannelColumns.PK,
			})

		for {
			err := queries.Raw(sqlStr, m.ChannelStatusInUse, m.ChannelStatusFree, pk).Bind(nil, mg.DB, &v)
			if err != nil && err != sql.ErrNoRows {
				mg.log.WithError(err).WithField("dist_pk", pk).Error("Error selecting dist channel")
				c1 <- err
				return
			}
			if err == nil {
				c1 <- nil
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}(DistPK)

	select {
	case res := <-c1:
		return res
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout. could not lock dist pk '%s'", DistPK)
	}
}

//ReleaseChannel fill release the channel in the database
func (mg *Manager) ReleaseChannel(ChannelPK string) error {
	_, err := m.Channels(qm.Where(m.ChannelColumns.PK+"=?", ChannelPK)).UpdateAll(mg.DB, m.M{m.ChannelColumns.Status: m.ChannelStatusFree})
	return err
}
