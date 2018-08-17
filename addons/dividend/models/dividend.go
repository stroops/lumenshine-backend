// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/strmangle"
	"gopkg.in/volatiletech/null.v6"
)

// Dividend is an object representing the database table.
type Dividend struct {
	ID             int        `boil:"id" json:"id" toml:"id" yaml:"id"`
	SnapshotID     int        `boil:"snapshot_id" json:"snapshot_id" toml:"snapshot_id" yaml:"snapshot_id"`
	AccountID      string     `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	BalanceLimit   int64      `boil:"balance_limit" json:"balance_limit" toml:"balance_limit" yaml:"balance_limit"`
	Balance        int64      `boil:"balance" json:"balance" toml:"balance" yaml:"balance"`
	DividendAmount null.Int64 `boil:"dividend_amount" json:"dividend_amount,omitempty" toml:"dividend_amount" yaml:"dividend_amount,omitempty"`

	R *dividendR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L dividendL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var DividendColumns = struct {
	ID             string
	SnapshotID     string
	AccountID      string
	BalanceLimit   string
	Balance        string
	DividendAmount string
}{
	ID:             "id",
	SnapshotID:     "snapshot_id",
	AccountID:      "account_id",
	BalanceLimit:   "balance_limit",
	Balance:        "balance",
	DividendAmount: "dividend_amount",
}

// dividendR is where relationships are stored.
type dividendR struct {
	Snapshot *Snapshot
}

// dividendL is where Load methods for each relationship are stored.
type dividendL struct{}

var (
	dividendColumns               = []string{"id", "snapshot_id", "account_id", "balance_limit", "balance", "dividend_amount"}
	dividendColumnsWithoutDefault = []string{"snapshot_id", "account_id", "balance_limit", "balance", "dividend_amount"}
	dividendColumnsWithDefault    = []string{"id"}
	dividendPrimaryKeyColumns     = []string{"id"}
)

type (
	// DividendSlice is an alias for a slice of pointers to Dividend.
	// This should generally be used opposed to []Dividend.
	DividendSlice []*Dividend
	// DividendHook is the signature for custom Dividend hook methods
	DividendHook func(boil.Executor, *Dividend) error

	dividendQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	dividendType                 = reflect.TypeOf(&Dividend{})
	dividendMapping              = queries.MakeStructMapping(dividendType)
	dividendPrimaryKeyMapping, _ = queries.BindMapping(dividendType, dividendMapping, dividendPrimaryKeyColumns)
	dividendInsertCacheMut       sync.RWMutex
	dividendInsertCache          = make(map[string]insertCache)
	dividendUpdateCacheMut       sync.RWMutex
	dividendUpdateCache          = make(map[string]updateCache)
	dividendUpsertCacheMut       sync.RWMutex
	dividendUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var dividendBeforeInsertHooks []DividendHook
var dividendBeforeUpdateHooks []DividendHook
var dividendBeforeDeleteHooks []DividendHook
var dividendBeforeUpsertHooks []DividendHook

var dividendAfterInsertHooks []DividendHook
var dividendAfterSelectHooks []DividendHook
var dividendAfterUpdateHooks []DividendHook
var dividendAfterDeleteHooks []DividendHook
var dividendAfterUpsertHooks []DividendHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Dividend) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Dividend) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Dividend) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Dividend) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Dividend) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Dividend) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Dividend) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Dividend) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Dividend) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range dividendAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddDividendHook registers your hook function for all future operations.
func AddDividendHook(hookPoint boil.HookPoint, dividendHook DividendHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		dividendBeforeInsertHooks = append(dividendBeforeInsertHooks, dividendHook)
	case boil.BeforeUpdateHook:
		dividendBeforeUpdateHooks = append(dividendBeforeUpdateHooks, dividendHook)
	case boil.BeforeDeleteHook:
		dividendBeforeDeleteHooks = append(dividendBeforeDeleteHooks, dividendHook)
	case boil.BeforeUpsertHook:
		dividendBeforeUpsertHooks = append(dividendBeforeUpsertHooks, dividendHook)
	case boil.AfterInsertHook:
		dividendAfterInsertHooks = append(dividendAfterInsertHooks, dividendHook)
	case boil.AfterSelectHook:
		dividendAfterSelectHooks = append(dividendAfterSelectHooks, dividendHook)
	case boil.AfterUpdateHook:
		dividendAfterUpdateHooks = append(dividendAfterUpdateHooks, dividendHook)
	case boil.AfterDeleteHook:
		dividendAfterDeleteHooks = append(dividendAfterDeleteHooks, dividendHook)
	case boil.AfterUpsertHook:
		dividendAfterUpsertHooks = append(dividendAfterUpsertHooks, dividendHook)
	}
}

// OneP returns a single dividend record from the query, and panics on error.
func (q dividendQuery) OneP() *Dividend {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single dividend record from the query.
func (q dividendQuery) One() (*Dividend, error) {
	o := &Dividend{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for dividend")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Dividend records from the query, and panics on error.
func (q dividendQuery) AllP() DividendSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Dividend records from the query.
func (q dividendQuery) All() (DividendSlice, error) {
	var o []*Dividend

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Dividend slice")
	}

	if len(dividendAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Dividend records in the query, and panics on error.
func (q dividendQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Dividend records in the query.
func (q dividendQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count dividend rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q dividendQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q dividendQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if dividend exists")
	}

	return count > 0, nil
}

// SnapshotG pointed to by the foreign key.
func (o *Dividend) SnapshotG(mods ...qm.QueryMod) snapshotQuery {
	return o.Snapshot(boil.GetDB(), mods...)
}

// Snapshot pointed to by the foreign key.
func (o *Dividend) Snapshot(exec boil.Executor, mods ...qm.QueryMod) snapshotQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.SnapshotID),
	}

	queryMods = append(queryMods, mods...)

	query := Snapshots(exec, queryMods...)
	queries.SetFrom(query.Query, "\"snapshot\"")

	return query
} // LoadSnapshot allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (dividendL) LoadSnapshot(e boil.Executor, singular bool, maybeDividend interface{}) error {
	var slice []*Dividend
	var object *Dividend

	count := 1
	if singular {
		object = maybeDividend.(*Dividend)
	} else {
		slice = *maybeDividend.(*[]*Dividend)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &dividendR{}
		}
		args[0] = object.SnapshotID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &dividendR{}
			}
			args[i] = obj.SnapshotID
		}
	}

	query := fmt.Sprintf(
		"select * from \"snapshot\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Snapshot")
	}
	defer results.Close()

	var resultSlice []*Snapshot
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Snapshot")
	}

	if len(dividendAfterSelectHooks) != 0 {
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
		object.R.Snapshot = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SnapshotID == foreign.ID {
				local.R.Snapshot = foreign
				break
			}
		}
	}

	return nil
}

// SetSnapshotG of the dividend to the related item.
// Sets o.R.Snapshot to related.
// Adds o to related.R.Dividends.
// Uses the global database handle.
func (o *Dividend) SetSnapshotG(insert bool, related *Snapshot) error {
	return o.SetSnapshot(boil.GetDB(), insert, related)
}

// SetSnapshotP of the dividend to the related item.
// Sets o.R.Snapshot to related.
// Adds o to related.R.Dividends.
// Panics on error.
func (o *Dividend) SetSnapshotP(exec boil.Executor, insert bool, related *Snapshot) {
	if err := o.SetSnapshot(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSnapshotGP of the dividend to the related item.
// Sets o.R.Snapshot to related.
// Adds o to related.R.Dividends.
// Uses the global database handle and panics on error.
func (o *Dividend) SetSnapshotGP(insert bool, related *Snapshot) {
	if err := o.SetSnapshot(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSnapshot of the dividend to the related item.
// Sets o.R.Snapshot to related.
// Adds o to related.R.Dividends.
func (o *Dividend) SetSnapshot(exec boil.Executor, insert bool, related *Snapshot) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"dividend\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"snapshot_id"}),
		strmangle.WhereClause("\"", "\"", 2, dividendPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SnapshotID = related.ID

	if o.R == nil {
		o.R = &dividendR{
			Snapshot: related,
		}
	} else {
		o.R.Snapshot = related
	}

	if related.R == nil {
		related.R = &snapshotR{
			Dividends: DividendSlice{o},
		}
	} else {
		related.R.Dividends = append(related.R.Dividends, o)
	}

	return nil
}

// DividendsG retrieves all records.
func DividendsG(mods ...qm.QueryMod) dividendQuery {
	return Dividends(boil.GetDB(), mods...)
}

// Dividends retrieves all the records using an executor.
func Dividends(exec boil.Executor, mods ...qm.QueryMod) dividendQuery {
	mods = append(mods, qm.From("\"dividend\""))
	return dividendQuery{NewQuery(exec, mods...)}
}

// FindDividendG retrieves a single record by ID.
func FindDividendG(id int, selectCols ...string) (*Dividend, error) {
	return FindDividend(boil.GetDB(), id, selectCols...)
}

// FindDividendGP retrieves a single record by ID, and panics on error.
func FindDividendGP(id int, selectCols ...string) *Dividend {
	retobj, err := FindDividend(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindDividend retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDividend(exec boil.Executor, id int, selectCols ...string) (*Dividend, error) {
	dividendObj := &Dividend{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"dividend\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(dividendObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from dividend")
	}

	return dividendObj, nil
}

// FindDividendP retrieves a single record by ID with an executor, and panics on error.
func FindDividendP(exec boil.Executor, id int, selectCols ...string) *Dividend {
	retobj, err := FindDividend(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Dividend) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Dividend) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Dividend) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Dividend) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no dividend provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(dividendColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	dividendInsertCacheMut.RLock()
	cache, cached := dividendInsertCache[key]
	dividendInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			dividendColumns,
			dividendColumnsWithDefault,
			dividendColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(dividendType, dividendMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(dividendType, dividendMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"dividend\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"dividend\" DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
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
		return errors.Wrap(err, "models: unable to insert into dividend")
	}

	if !cached {
		dividendInsertCacheMut.Lock()
		dividendInsertCache[key] = cache
		dividendInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Dividend record. See Update for
// whitelist behavior description.
func (o *Dividend) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Dividend record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Dividend) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Dividend, and panics on error.
// See Update for whitelist behavior description.
func (o *Dividend) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Dividend.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Dividend) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	dividendUpdateCacheMut.RLock()
	cache, cached := dividendUpdateCache[key]
	dividendUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			dividendColumns,
			dividendPrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update dividend, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"dividend\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, dividendPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(dividendType, dividendMapping, append(wl, dividendPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update dividend row")
	}

	if !cached {
		dividendUpdateCacheMut.Lock()
		dividendUpdateCache[key] = cache
		dividendUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q dividendQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q dividendQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for dividend")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o DividendSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o DividendSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o DividendSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DividendSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), dividendPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"dividend\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, dividendPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in dividend slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Dividend) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Dividend) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Dividend) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Dividend) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no dividend provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(dividendColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
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
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	dividendUpsertCacheMut.RLock()
	cache, cached := dividendUpsertCache[key]
	dividendUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			dividendColumns,
			dividendColumnsWithDefault,
			dividendColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			dividendColumns,
			dividendPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert dividend, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(dividendPrimaryKeyColumns))
			copy(conflict, dividendPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"dividend\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(dividendType, dividendMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(dividendType, dividendMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert dividend")
	}

	if !cached {
		dividendUpsertCacheMut.Lock()
		dividendUpsertCache[key] = cache
		dividendUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Dividend record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Dividend) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Dividend record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Dividend) DeleteG() error {
	if o == nil {
		return errors.New("models: no Dividend provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Dividend record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Dividend) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Dividend record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Dividend) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Dividend provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), dividendPrimaryKeyMapping)
	sql := "DELETE FROM \"dividend\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from dividend")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q dividendQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q dividendQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no dividendQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from dividend")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o DividendSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o DividendSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Dividend slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o DividendSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DividendSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Dividend slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(dividendBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), dividendPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"dividend\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, dividendPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from dividend slice")
	}

	if len(dividendAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Dividend) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Dividend) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Dividend) ReloadG() error {
	if o == nil {
		return errors.New("models: no Dividend provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Dividend) Reload(exec boil.Executor) error {
	ret, err := FindDividend(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DividendSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DividendSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DividendSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty DividendSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DividendSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	dividends := DividendSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), dividendPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"dividend\".* FROM \"dividend\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, dividendPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&dividends)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in DividendSlice")
	}

	*o = dividends

	return nil
}

// DividendExists checks if the Dividend row exists.
func DividendExists(exec boil.Executor, id int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"dividend\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if dividend exists")
	}

	return exists, nil
}

// DividendExistsG checks if the Dividend row exists.
func DividendExistsG(id int) (bool, error) {
	return DividendExists(boil.GetDB(), id)
}

// DividendExistsGP checks if the Dividend row exists. Panics on error.
func DividendExistsGP(id int) bool {
	e, err := DividendExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// DividendExistsP checks if the Dividend row exists. Panics on error.
func DividendExistsP(exec boil.Executor, id int) bool {
	e, err := DividendExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}