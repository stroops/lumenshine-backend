package main

import (
	"time"

	"github.com/Soneso/lumenshine-backend/db/querying"

	m "github.com/Soneso/lumenshine-backend/db/horizon/models"
	"github.com/Soneso/lumenshine-backend/services/sse/db"
	"github.com/Soneso/lumenshine-backend/services/sse/environment"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"database/sql"
)

var (
	//last processed transaction id
	startupOrderID int64
)

//StellarProcessor processes the stellar data
type StellarProcessor struct {
	env *environment.Environment
	db  *db.DB
	log *logrus.Entry
}

//NewStellarProcessor returns a new StellarProcessor
func NewStellarProcessor(log *logrus.Entry) *StellarProcessor {
	// start from actual opID
	c := new(StellarProcessor)
	c.db = environment.Env.DBH
	c.env = environment.Env
	c.log = log

	//we need to read the current latest OrderID
	//we need this information, in order to know, which operations must be signled and which not.
	//this is configured via the with_resume flag in the DB

	err := queries.Raw("select max("+m.HistoryOperationColumns.ID+") from "+m.TableNames.HistoryOperations).Bind(nil, c.db, &startupOrderID)
	if err != nil {
		c.log.WithError(err).Error("Error reading intial index")
		startupOrderID = 0
	}

	return c
}

type nextID struct {
	NextID int64 `boil:"value" json:"next_id"`
}

//StartProcessing starts processing the stellar data
//will wait for last_op_participant_id to be set to != 0
//should be started as a goroutine
//Set the value in the db for the last index to one lower the one you want to process. the value will be incremented on first run
func (s *StellarProcessor) StartProcessing() {
	//read latest processed op id
	//we wait until the value is set != 0

	sseIndex, err := m.SseIndices(qm.Where(m.SseIndexColumns.Name+"=?", m.SseIndexNamesNextOperationID)).One(s.db)
	if err != nil {
		s.log.WithError(err).Error("Error reading index")
	}
	if sseIndex.Value == 0 {
		for {
			s.log.Info("No index set. Waiting...")
			time.Sleep(5 * time.Second)

			err := sseIndex.Reload(s.db)
			if err != nil {
				s.log.WithError(err).Error("Error reading index in loop")
			}
			if sseIndex.Value != 0 {
				break
			}
		}
	}

	sqlStr := querying.GetSQLKeyString(`update @sse_index set value=value+1 where @name=$1 returning value`,
		map[string]string{
			"@sse_index": m.TableNames.SseIndex,
			"@name":      m.SseIndexColumns.Name,
		})

	var ni nextID
	ni.NextID = sseIndex.Value

	for {
		err = s.processOperation(ni.NextID)
		if err != nil {
			if err != sql.ErrNoRows {
				s.log.WithError(err).WithField("operation_id", ni.NextID).Error("Error processing operation")
			}
			time.Sleep(5 * time.Second) // ledger closes roughly every 5 seconds
			continue                    // try againe
		}

		//read next operation_id to process
		err := queries.Raw(sqlStr, m.SseIndexNamesNextOperationID).Bind(nil, s.db, &ni)
		if err != nil {
			if err != sql.ErrNoRows {
				s.log.WithError(err).Error("Error selecting next id")
			}
			time.Sleep(2 * time.Second)
		}
	}
}

type sseConfigData struct {
	m.SseConfig
	m.HistoryOperation
}

//{ CREATE_ACCOUNT = 0, PAYMENT = 1, PATH_PAYMENT = 2, MANAGE_OFFER = 3, CREATE_PASSIVE_OFFER = 4,
//	SET_OPTIONS = 5, CHANGE_TRUST = 6, ALLOW_TRUST = 7, ACCOUNT_MERGE = 8, INFLATION = 9, MANAGE_DATA = 10, BUMP_SEQUENCE = 11 }
//processOperationParticipant
func (s *StellarProcessor) processOperation(id int64) error {
	s.log.Printf("prosessing order %d", id)
	cO := m.HistoryOperationColumns
	sC := m.SseConfigColumns
	mT := m.TableNames

	q := []qm.QueryMod{
		qm.From(mT.SseConfig),
		qm.Select(mT.SseConfig + ".*, " + mT.HistoryOperations + ".*"),

		qm.InnerJoin(m.TableNames.HistoryOperations + " on cast(" + cO.Details + "->>'to' as character varying)=" + sC.StellarAccount),
		qm.Where(cO.ID+"=?", id),
		qm.And("2<<" + cO.Type + "&" + sC.OperationTypes + "=" + sC.OperationTypes),

		//for payments and paymentPath, we only check the receivers
		qm.And("case when " + cO.Type + "=1 or " + cO.Type + "=2 then cast(" + cO.Details + "->>'to' as character varying)=" + sC.StellarAccount + " else true end"),

		//for creates we only check the generated address
		qm.And("case when " + cO.Type + "=0 then cast(" + cO.Details + "->>'account' as character varying)=" + sC.StellarAccount + " else true end"),

		//if with_resume not set on config, we will check that the id bigger than the startupOrderID
		qm.And("case when "+sC.WithResume+"=false then "+mT.HistoryOperations+"."+cO.ID+">=? else true end", startupOrderID),
	}

	var sCs []sseConfigData
	err := m.NewQuery(q...).Bind(nil, s.db, sCs)
	if err != nil {
		return err
	}

	for _, sC := range sCs {
		var sD m.SseDatum
		sD.SseConfigID = sC.SseConfig.ID
		sD.SourceReceiver = sC.SseConfig.SourceReceiver
		sD.Status = m.SseDataStatusNew
		sD.StellarAccount = sC.SseConfig.StellarAccount
		sD.OperationType = sC.HistoryOperation.Type
		sD.OperationData = sC.HistoryOperation.Details
	}

	return nil
}
