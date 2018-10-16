// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package horizon

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// HistoryTrade is an object representing the database table.
type HistoryTrade struct {
	HistoryOperationID int64      `boil:"history_operation_id" json:"history_operation_id" toml:"history_operation_id" yaml:"history_operation_id"`
	Order              int        `boil:"order" json:"order" toml:"order" yaml:"order"`
	LedgerClosedAt     time.Time  `boil:"ledger_closed_at" json:"ledger_closed_at" toml:"ledger_closed_at" yaml:"ledger_closed_at"`
	OfferID            int64      `boil:"offer_id" json:"offer_id" toml:"offer_id" yaml:"offer_id"`
	BaseAccountID      int64      `boil:"base_account_id" json:"base_account_id" toml:"base_account_id" yaml:"base_account_id"`
	BaseAssetID        int64      `boil:"base_asset_id" json:"base_asset_id" toml:"base_asset_id" yaml:"base_asset_id"`
	BaseAmount         int64      `boil:"base_amount" json:"base_amount" toml:"base_amount" yaml:"base_amount"`
	CounterAccountID   int64      `boil:"counter_account_id" json:"counter_account_id" toml:"counter_account_id" yaml:"counter_account_id"`
	CounterAssetID     int64      `boil:"counter_asset_id" json:"counter_asset_id" toml:"counter_asset_id" yaml:"counter_asset_id"`
	CounterAmount      int64      `boil:"counter_amount" json:"counter_amount" toml:"counter_amount" yaml:"counter_amount"`
	BaseIsSeller       null.Bool  `boil:"base_is_seller" json:"base_is_seller,omitempty" toml:"base_is_seller" yaml:"base_is_seller,omitempty"`
	PriceN             null.Int64 `boil:"price_n" json:"price_n,omitempty" toml:"price_n" yaml:"price_n,omitempty"`
	PriceD             null.Int64 `boil:"price_d" json:"price_d,omitempty" toml:"price_d" yaml:"price_d,omitempty"`

	R *historyTradeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L historyTradeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var HistoryTradeColumns = struct {
	HistoryOperationID string
	Order              string
	LedgerClosedAt     string
	OfferID            string
	BaseAccountID      string
	BaseAssetID        string
	BaseAmount         string
	CounterAccountID   string
	CounterAssetID     string
	CounterAmount      string
	BaseIsSeller       string
	PriceN             string
	PriceD             string
}{
	HistoryOperationID: "history_operation_id",
	Order:              "order",
	LedgerClosedAt:     "ledger_closed_at",
	OfferID:            "offer_id",
	BaseAccountID:      "base_account_id",
	BaseAssetID:        "base_asset_id",
	BaseAmount:         "base_amount",
	CounterAccountID:   "counter_account_id",
	CounterAssetID:     "counter_asset_id",
	CounterAmount:      "counter_amount",
	BaseIsSeller:       "base_is_seller",
	PriceN:             "price_n",
	PriceD:             "price_d",
}

// HistoryTradeRels is where relationship names are stored.
var HistoryTradeRels = struct {
	BaseAccount    string
	CounterAccount string
}{
	BaseAccount:    "BaseAccount",
	CounterAccount: "CounterAccount",
}

// historyTradeR is where relationships are stored.
type historyTradeR struct {
	BaseAccount    *HistoryAccount
	CounterAccount *HistoryAccount
}

// NewStruct creates a new relationship struct
func (*historyTradeR) NewStruct() *historyTradeR {
	return &historyTradeR{}
}

// historyTradeL is where Load methods for each relationship are stored.
type historyTradeL struct{}

var (
	historyTradeColumns               = []string{"history_operation_id", "order", "ledger_closed_at", "offer_id", "base_account_id", "base_asset_id", "base_amount", "counter_account_id", "counter_asset_id", "counter_amount", "base_is_seller", "price_n", "price_d"}
	historyTradeColumnsWithoutDefault = []string{"history_operation_id", "order", "ledger_closed_at", "offer_id", "base_account_id", "base_asset_id", "base_amount", "counter_account_id", "counter_asset_id", "counter_amount", "base_is_seller", "price_n", "price_d"}
	historyTradeColumnsWithDefault    = []string{}
	historyTradePrimaryKeyColumns     = []string{"history_operation_id", "order"}
)

type (
	// HistoryTradeSlice is an alias for a slice of pointers to HistoryTrade.
	// This should generally be used opposed to []HistoryTrade.
	HistoryTradeSlice []*HistoryTrade
	// HistoryTradeHook is the signature for custom HistoryTrade hook methods
	HistoryTradeHook func(boil.Executor, *HistoryTrade) error

	historyTradeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	historyTradeType                 = reflect.TypeOf(&HistoryTrade{})
	historyTradeMapping              = queries.MakeStructMapping(historyTradeType)
	historyTradePrimaryKeyMapping, _ = queries.BindMapping(historyTradeType, historyTradeMapping, historyTradePrimaryKeyColumns)
	historyTradeInsertCacheMut       sync.RWMutex
	historyTradeInsertCache          = make(map[string]insertCache)
	historyTradeUpdateCacheMut       sync.RWMutex
	historyTradeUpdateCache          = make(map[string]updateCache)
	historyTradeUpsertCacheMut       sync.RWMutex
	historyTradeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var historyTradeBeforeInsertHooks []HistoryTradeHook
var historyTradeBeforeUpdateHooks []HistoryTradeHook
var historyTradeBeforeDeleteHooks []HistoryTradeHook
var historyTradeBeforeUpsertHooks []HistoryTradeHook

var historyTradeAfterInsertHooks []HistoryTradeHook
var historyTradeAfterSelectHooks []HistoryTradeHook
var historyTradeAfterUpdateHooks []HistoryTradeHook
var historyTradeAfterDeleteHooks []HistoryTradeHook
var historyTradeAfterUpsertHooks []HistoryTradeHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *HistoryTrade) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *HistoryTrade) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *HistoryTrade) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *HistoryTrade) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *HistoryTrade) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *HistoryTrade) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *HistoryTrade) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *HistoryTrade) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *HistoryTrade) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range historyTradeAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddHistoryTradeHook registers your hook function for all future operations.
func AddHistoryTradeHook(hookPoint boil.HookPoint, historyTradeHook HistoryTradeHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		historyTradeBeforeInsertHooks = append(historyTradeBeforeInsertHooks, historyTradeHook)
	case boil.BeforeUpdateHook:
		historyTradeBeforeUpdateHooks = append(historyTradeBeforeUpdateHooks, historyTradeHook)
	case boil.BeforeDeleteHook:
		historyTradeBeforeDeleteHooks = append(historyTradeBeforeDeleteHooks, historyTradeHook)
	case boil.BeforeUpsertHook:
		historyTradeBeforeUpsertHooks = append(historyTradeBeforeUpsertHooks, historyTradeHook)
	case boil.AfterInsertHook:
		historyTradeAfterInsertHooks = append(historyTradeAfterInsertHooks, historyTradeHook)
	case boil.AfterSelectHook:
		historyTradeAfterSelectHooks = append(historyTradeAfterSelectHooks, historyTradeHook)
	case boil.AfterUpdateHook:
		historyTradeAfterUpdateHooks = append(historyTradeAfterUpdateHooks, historyTradeHook)
	case boil.AfterDeleteHook:
		historyTradeAfterDeleteHooks = append(historyTradeAfterDeleteHooks, historyTradeHook)
	case boil.AfterUpsertHook:
		historyTradeAfterUpsertHooks = append(historyTradeAfterUpsertHooks, historyTradeHook)
	}
}

// OneG returns a single historyTrade record from the query using the global executor.
func (q historyTradeQuery) OneG() (*HistoryTrade, error) {
	return q.One(boil.GetDB())
}

// One returns a single historyTrade record from the query.
func (q historyTradeQuery) One(exec boil.Executor) (*HistoryTrade, error) {
	o := &HistoryTrade{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: failed to execute a one query for history_trades")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all HistoryTrade records from the query using the global executor.
func (q historyTradeQuery) AllG() (HistoryTradeSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all HistoryTrade records from the query.
func (q historyTradeQuery) All(exec boil.Executor) (HistoryTradeSlice, error) {
	var o []*HistoryTrade

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "horizon: failed to assign all query results to HistoryTrade slice")
	}

	if len(historyTradeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all HistoryTrade records in the query, and panics on error.
func (q historyTradeQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all HistoryTrade records in the query.
func (q historyTradeQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to count history_trades rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q historyTradeQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q historyTradeQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "horizon: failed to check if history_trades exists")
	}

	return count > 0, nil
}

// BaseAccount pointed to by the foreign key.
func (o *HistoryTrade) BaseAccount(mods ...qm.QueryMod) historyAccountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.BaseAccountID),
	}

	queryMods = append(queryMods, mods...)

	query := HistoryAccounts(queryMods...)
	queries.SetFrom(query.Query, "\"history_accounts\"")

	return query
}

// CounterAccount pointed to by the foreign key.
func (o *HistoryTrade) CounterAccount(mods ...qm.QueryMod) historyAccountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.CounterAccountID),
	}

	queryMods = append(queryMods, mods...)

	query := HistoryAccounts(queryMods...)
	queries.SetFrom(query.Query, "\"history_accounts\"")

	return query
}

// LoadBaseAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (historyTradeL) LoadBaseAccount(e boil.Executor, singular bool, maybeHistoryTrade interface{}, mods queries.Applicator) error {
	var slice []*HistoryTrade
	var object *HistoryTrade

	if singular {
		object = maybeHistoryTrade.(*HistoryTrade)
	} else {
		slice = *maybeHistoryTrade.(*[]*HistoryTrade)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &historyTradeR{}
		}
		args = append(args, object.BaseAccountID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &historyTradeR{}
			}

			for _, a := range args {
				if a == obj.BaseAccountID {
					continue Outer
				}
			}

			args = append(args, obj.BaseAccountID)
		}
	}

	query := NewQuery(qm.From(`history_accounts`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load HistoryAccount")
	}

	var resultSlice []*HistoryAccount
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice HistoryAccount")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for history_accounts")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for history_accounts")
	}

	if len(historyTradeAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.BaseAccount = foreign
		if foreign.R == nil {
			foreign.R = &historyAccountR{}
		}
		foreign.R.BaseAccountHistoryTrades = append(foreign.R.BaseAccountHistoryTrades, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BaseAccountID == foreign.ID {
				local.R.BaseAccount = foreign
				if foreign.R == nil {
					foreign.R = &historyAccountR{}
				}
				foreign.R.BaseAccountHistoryTrades = append(foreign.R.BaseAccountHistoryTrades, local)
				break
			}
		}
	}

	return nil
}

// LoadCounterAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (historyTradeL) LoadCounterAccount(e boil.Executor, singular bool, maybeHistoryTrade interface{}, mods queries.Applicator) error {
	var slice []*HistoryTrade
	var object *HistoryTrade

	if singular {
		object = maybeHistoryTrade.(*HistoryTrade)
	} else {
		slice = *maybeHistoryTrade.(*[]*HistoryTrade)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &historyTradeR{}
		}
		args = append(args, object.CounterAccountID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &historyTradeR{}
			}

			for _, a := range args {
				if a == obj.CounterAccountID {
					continue Outer
				}
			}

			args = append(args, obj.CounterAccountID)
		}
	}

	query := NewQuery(qm.From(`history_accounts`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load HistoryAccount")
	}

	var resultSlice []*HistoryAccount
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice HistoryAccount")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for history_accounts")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for history_accounts")
	}

	if len(historyTradeAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.CounterAccount = foreign
		if foreign.R == nil {
			foreign.R = &historyAccountR{}
		}
		foreign.R.CounterAccountHistoryTrades = append(foreign.R.CounterAccountHistoryTrades, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CounterAccountID == foreign.ID {
				local.R.CounterAccount = foreign
				if foreign.R == nil {
					foreign.R = &historyAccountR{}
				}
				foreign.R.CounterAccountHistoryTrades = append(foreign.R.CounterAccountHistoryTrades, local)
				break
			}
		}
	}

	return nil
}

// SetBaseAccountG of the historyTrade to the related item.
// Sets o.R.BaseAccount to related.
// Adds o to related.R.BaseAccountHistoryTrades.
// Uses the global database handle.
func (o *HistoryTrade) SetBaseAccountG(insert bool, related *HistoryAccount) error {
	return o.SetBaseAccount(boil.GetDB(), insert, related)
}

// SetBaseAccount of the historyTrade to the related item.
// Sets o.R.BaseAccount to related.
// Adds o to related.R.BaseAccountHistoryTrades.
func (o *HistoryTrade) SetBaseAccount(exec boil.Executor, insert bool, related *HistoryAccount) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"history_trades\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"base_account_id"}),
		strmangle.WhereClause("\"", "\"", 2, historyTradePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.HistoryOperationID, o.Order}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.BaseAccountID = related.ID
	if o.R == nil {
		o.R = &historyTradeR{
			BaseAccount: related,
		}
	} else {
		o.R.BaseAccount = related
	}

	if related.R == nil {
		related.R = &historyAccountR{
			BaseAccountHistoryTrades: HistoryTradeSlice{o},
		}
	} else {
		related.R.BaseAccountHistoryTrades = append(related.R.BaseAccountHistoryTrades, o)
	}

	return nil
}

// SetCounterAccountG of the historyTrade to the related item.
// Sets o.R.CounterAccount to related.
// Adds o to related.R.CounterAccountHistoryTrades.
// Uses the global database handle.
func (o *HistoryTrade) SetCounterAccountG(insert bool, related *HistoryAccount) error {
	return o.SetCounterAccount(boil.GetDB(), insert, related)
}

// SetCounterAccount of the historyTrade to the related item.
// Sets o.R.CounterAccount to related.
// Adds o to related.R.CounterAccountHistoryTrades.
func (o *HistoryTrade) SetCounterAccount(exec boil.Executor, insert bool, related *HistoryAccount) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"history_trades\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"counter_account_id"}),
		strmangle.WhereClause("\"", "\"", 2, historyTradePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.HistoryOperationID, o.Order}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CounterAccountID = related.ID
	if o.R == nil {
		o.R = &historyTradeR{
			CounterAccount: related,
		}
	} else {
		o.R.CounterAccount = related
	}

	if related.R == nil {
		related.R = &historyAccountR{
			CounterAccountHistoryTrades: HistoryTradeSlice{o},
		}
	} else {
		related.R.CounterAccountHistoryTrades = append(related.R.CounterAccountHistoryTrades, o)
	}

	return nil
}

// HistoryTrades retrieves all the records using an executor.
func HistoryTrades(mods ...qm.QueryMod) historyTradeQuery {
	mods = append(mods, qm.From("\"history_trades\""))
	return historyTradeQuery{NewQuery(mods...)}
}

// FindHistoryTradeG retrieves a single record by ID.
func FindHistoryTradeG(historyOperationID int64, order int, selectCols ...string) (*HistoryTrade, error) {
	return FindHistoryTrade(boil.GetDB(), historyOperationID, order, selectCols...)
}

// FindHistoryTrade retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindHistoryTrade(exec boil.Executor, historyOperationID int64, order int, selectCols ...string) (*HistoryTrade, error) {
	historyTradeObj := &HistoryTrade{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"history_trades\" where \"history_operation_id\"=$1 AND \"order\"=$2", sel,
	)

	q := queries.Raw(query, historyOperationID, order)

	err := q.Bind(nil, exec, historyTradeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: unable to select from history_trades")
	}

	return historyTradeObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *HistoryTrade) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *HistoryTrade) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no history_trades provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(historyTradeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	historyTradeInsertCacheMut.RLock()
	cache, cached := historyTradeInsertCache[key]
	historyTradeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			historyTradeColumns,
			historyTradeColumnsWithDefault,
			historyTradeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(historyTradeType, historyTradeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(historyTradeType, historyTradeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"history_trades\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"history_trades\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "horizon: unable to insert into history_trades")
	}

	if !cached {
		historyTradeInsertCacheMut.Lock()
		historyTradeInsertCache[key] = cache
		historyTradeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single HistoryTrade record using the global executor.
// See Update for more documentation.
func (o *HistoryTrade) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the HistoryTrade.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *HistoryTrade) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	historyTradeUpdateCacheMut.RLock()
	cache, cached := historyTradeUpdateCache[key]
	historyTradeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			historyTradeColumns,
			historyTradePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("horizon: unable to update history_trades, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"history_trades\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, historyTradePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(historyTradeType, historyTradeMapping, append(wl, historyTradePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update history_trades row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by update for history_trades")
	}

	if !cached {
		historyTradeUpdateCacheMut.Lock()
		historyTradeUpdateCache[key] = cache
		historyTradeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q historyTradeQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all for history_trades")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected for history_trades")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o HistoryTradeSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o HistoryTradeSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("horizon: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), historyTradePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"history_trades\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, historyTradePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all in historyTrade slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected all in update all historyTrade")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *HistoryTrade) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *HistoryTrade) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no history_trades provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(historyTradeColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	historyTradeUpsertCacheMut.RLock()
	cache, cached := historyTradeUpsertCache[key]
	historyTradeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			historyTradeColumns,
			historyTradeColumnsWithDefault,
			historyTradeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			historyTradeColumns,
			historyTradePrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("horizon: unable to upsert history_trades, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(historyTradePrimaryKeyColumns))
			copy(conflict, historyTradePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"history_trades\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(historyTradeType, historyTradeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(historyTradeType, historyTradeMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "horizon: unable to upsert history_trades")
	}

	if !cached {
		historyTradeUpsertCacheMut.Lock()
		historyTradeUpsertCache[key] = cache
		historyTradeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single HistoryTrade record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *HistoryTrade) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single HistoryTrade record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *HistoryTrade) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no HistoryTrade provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), historyTradePrimaryKeyMapping)
	sql := "DELETE FROM \"history_trades\" WHERE \"history_operation_id\"=$1 AND \"order\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete from history_trades")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by delete for history_trades")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q historyTradeQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("horizon: no historyTradeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from history_trades")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for history_trades")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o HistoryTradeSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o HistoryTradeSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no HistoryTrade slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(historyTradeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), historyTradePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"history_trades\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, historyTradePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from historyTrade slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for history_trades")
	}

	if len(historyTradeAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *HistoryTrade) ReloadG() error {
	if o == nil {
		return errors.New("horizon: no HistoryTrade provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *HistoryTrade) Reload(exec boil.Executor) error {
	ret, err := FindHistoryTrade(exec, o.HistoryOperationID, o.Order)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *HistoryTradeSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("horizon: empty HistoryTradeSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *HistoryTradeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := HistoryTradeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), historyTradePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"history_trades\".* FROM \"history_trades\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, historyTradePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "horizon: unable to reload all in HistoryTradeSlice")
	}

	*o = slice

	return nil
}

// HistoryTradeExistsG checks if the HistoryTrade row exists.
func HistoryTradeExistsG(historyOperationID int64, order int) (bool, error) {
	return HistoryTradeExists(boil.GetDB(), historyOperationID, order)
}

// HistoryTradeExists checks if the HistoryTrade row exists.
func HistoryTradeExists(exec boil.Executor, historyOperationID int64, order int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"history_trades\" where \"history_operation_id\"=$1 AND \"order\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, historyOperationID, order)
	}

	row := exec.QueryRow(sql, historyOperationID, order)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "horizon: unable to check if history_trades exists")
	}

	return exists, nil
}
