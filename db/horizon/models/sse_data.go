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

// SseDatum is an object representing the database table.
type SseDatum struct {
	ID             int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	SourceReceiver string    `boil:"source_receiver" json:"source_receiver" toml:"source_receiver" yaml:"source_receiver"`
	Status         string    `boil:"status" json:"status" toml:"status" yaml:"status"`
	StellarAccount string    `boil:"stellar_account" json:"stellar_account" toml:"stellar_account" yaml:"stellar_account"`
	OperationType  int       `boil:"operation_type" json:"operation_type" toml:"operation_type" yaml:"operation_type"`
	OperationData  null.JSON `boil:"operation_data" json:"operation_data,omitempty" toml:"operation_data" yaml:"operation_data,omitempty"`
	TransactionID  int64     `boil:"transaction_id" json:"transaction_id" toml:"transaction_id" yaml:"transaction_id"`
	OperationID    int64     `boil:"operation_id" json:"operation_id" toml:"operation_id" yaml:"operation_id"`
	LedgerID       int64     `boil:"ledger_id" json:"ledger_id" toml:"ledger_id" yaml:"ledger_id"`
	CreatedAt      time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *sseDatumR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sseDatumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SseDatumColumns = struct {
	ID             string
	SourceReceiver string
	Status         string
	StellarAccount string
	OperationType  string
	OperationData  string
	TransactionID  string
	OperationID    string
	LedgerID       string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	SourceReceiver: "source_receiver",
	Status:         "status",
	StellarAccount: "stellar_account",
	OperationType:  "operation_type",
	OperationData:  "operation_data",
	TransactionID:  "transaction_id",
	OperationID:    "operation_id",
	LedgerID:       "ledger_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// SseDatumRels is where relationship names are stored.
var SseDatumRels = struct {
}{}

// sseDatumR is where relationships are stored.
type sseDatumR struct {
}

// NewStruct creates a new relationship struct
func (*sseDatumR) NewStruct() *sseDatumR {
	return &sseDatumR{}
}

// sseDatumL is where Load methods for each relationship are stored.
type sseDatumL struct{}

var (
	sseDatumColumns               = []string{"id", "source_receiver", "status", "stellar_account", "operation_type", "operation_data", "transaction_id", "operation_id", "ledger_id", "created_at", "updated_at"}
	sseDatumColumnsWithoutDefault = []string{"source_receiver", "status", "stellar_account", "operation_type", "operation_data", "transaction_id", "operation_id", "ledger_id"}
	sseDatumColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	sseDatumPrimaryKeyColumns     = []string{"id"}
)

type (
	// SseDatumSlice is an alias for a slice of pointers to SseDatum.
	// This should generally be used opposed to []SseDatum.
	SseDatumSlice []*SseDatum
	// SseDatumHook is the signature for custom SseDatum hook methods
	SseDatumHook func(boil.Executor, *SseDatum) error

	sseDatumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sseDatumType                 = reflect.TypeOf(&SseDatum{})
	sseDatumMapping              = queries.MakeStructMapping(sseDatumType)
	sseDatumPrimaryKeyMapping, _ = queries.BindMapping(sseDatumType, sseDatumMapping, sseDatumPrimaryKeyColumns)
	sseDatumInsertCacheMut       sync.RWMutex
	sseDatumInsertCache          = make(map[string]insertCache)
	sseDatumUpdateCacheMut       sync.RWMutex
	sseDatumUpdateCache          = make(map[string]updateCache)
	sseDatumUpsertCacheMut       sync.RWMutex
	sseDatumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var sseDatumBeforeInsertHooks []SseDatumHook
var sseDatumBeforeUpdateHooks []SseDatumHook
var sseDatumBeforeDeleteHooks []SseDatumHook
var sseDatumBeforeUpsertHooks []SseDatumHook

var sseDatumAfterInsertHooks []SseDatumHook
var sseDatumAfterSelectHooks []SseDatumHook
var sseDatumAfterUpdateHooks []SseDatumHook
var sseDatumAfterDeleteHooks []SseDatumHook
var sseDatumAfterUpsertHooks []SseDatumHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SseDatum) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SseDatum) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SseDatum) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SseDatum) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SseDatum) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SseDatum) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SseDatum) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SseDatum) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SseDatum) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseDatumAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSseDatumHook registers your hook function for all future operations.
func AddSseDatumHook(hookPoint boil.HookPoint, sseDatumHook SseDatumHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		sseDatumBeforeInsertHooks = append(sseDatumBeforeInsertHooks, sseDatumHook)
	case boil.BeforeUpdateHook:
		sseDatumBeforeUpdateHooks = append(sseDatumBeforeUpdateHooks, sseDatumHook)
	case boil.BeforeDeleteHook:
		sseDatumBeforeDeleteHooks = append(sseDatumBeforeDeleteHooks, sseDatumHook)
	case boil.BeforeUpsertHook:
		sseDatumBeforeUpsertHooks = append(sseDatumBeforeUpsertHooks, sseDatumHook)
	case boil.AfterInsertHook:
		sseDatumAfterInsertHooks = append(sseDatumAfterInsertHooks, sseDatumHook)
	case boil.AfterSelectHook:
		sseDatumAfterSelectHooks = append(sseDatumAfterSelectHooks, sseDatumHook)
	case boil.AfterUpdateHook:
		sseDatumAfterUpdateHooks = append(sseDatumAfterUpdateHooks, sseDatumHook)
	case boil.AfterDeleteHook:
		sseDatumAfterDeleteHooks = append(sseDatumAfterDeleteHooks, sseDatumHook)
	case boil.AfterUpsertHook:
		sseDatumAfterUpsertHooks = append(sseDatumAfterUpsertHooks, sseDatumHook)
	}
}

// OneG returns a single sseDatum record from the query using the global executor.
func (q sseDatumQuery) OneG() (*SseDatum, error) {
	return q.One(boil.GetDB())
}

// One returns a single sseDatum record from the query.
func (q sseDatumQuery) One(exec boil.Executor) (*SseDatum, error) {
	o := &SseDatum{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: failed to execute a one query for sse_data")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all SseDatum records from the query using the global executor.
func (q sseDatumQuery) AllG() (SseDatumSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all SseDatum records from the query.
func (q sseDatumQuery) All(exec boil.Executor) (SseDatumSlice, error) {
	var o []*SseDatum

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "horizon: failed to assign all query results to SseDatum slice")
	}

	if len(sseDatumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all SseDatum records in the query, and panics on error.
func (q sseDatumQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all SseDatum records in the query.
func (q sseDatumQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to count sse_data rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q sseDatumQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q sseDatumQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "horizon: failed to check if sse_data exists")
	}

	return count > 0, nil
}

// SseData retrieves all the records using an executor.
func SseData(mods ...qm.QueryMod) sseDatumQuery {
	mods = append(mods, qm.From("\"sse_data\""))
	return sseDatumQuery{NewQuery(mods...)}
}

// FindSseDatumG retrieves a single record by ID.
func FindSseDatumG(iD int, selectCols ...string) (*SseDatum, error) {
	return FindSseDatum(boil.GetDB(), iD, selectCols...)
}

// FindSseDatum retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSseDatum(exec boil.Executor, iD int, selectCols ...string) (*SseDatum, error) {
	sseDatumObj := &SseDatum{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sse_data\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, sseDatumObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: unable to select from sse_data")
	}

	return sseDatumObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SseDatum) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SseDatum) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no sse_data provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sseDatumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sseDatumInsertCacheMut.RLock()
	cache, cached := sseDatumInsertCache[key]
	sseDatumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sseDatumColumns,
			sseDatumColumnsWithDefault,
			sseDatumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sseDatumType, sseDatumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sseDatumType, sseDatumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sse_data\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sse_data\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "horizon: unable to insert into sse_data")
	}

	if !cached {
		sseDatumInsertCacheMut.Lock()
		sseDatumInsertCache[key] = cache
		sseDatumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single SseDatum record using the global executor.
// See Update for more documentation.
func (o *SseDatum) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the SseDatum.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SseDatum) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	sseDatumUpdateCacheMut.RLock()
	cache, cached := sseDatumUpdateCache[key]
	sseDatumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sseDatumColumns,
			sseDatumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("horizon: unable to update sse_data, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sse_data\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sseDatumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sseDatumType, sseDatumMapping, append(wl, sseDatumPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "horizon: unable to update sse_data row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by update for sse_data")
	}

	if !cached {
		sseDatumUpdateCacheMut.Lock()
		sseDatumUpdateCache[key] = cache
		sseDatumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q sseDatumQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all for sse_data")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected for sse_data")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SseDatumSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SseDatumSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseDatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sse_data\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sseDatumPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all in sseDatum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected all in update all sseDatum")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SseDatum) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SseDatum) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no sse_data provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sseDatumColumnsWithDefault, o)

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

	sseDatumUpsertCacheMut.RLock()
	cache, cached := sseDatumUpsertCache[key]
	sseDatumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			sseDatumColumns,
			sseDatumColumnsWithDefault,
			sseDatumColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			sseDatumColumns,
			sseDatumPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("horizon: unable to upsert sse_data, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sseDatumPrimaryKeyColumns))
			copy(conflict, sseDatumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sse_data\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(sseDatumType, sseDatumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sseDatumType, sseDatumMapping, ret)
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
		return errors.Wrap(err, "horizon: unable to upsert sse_data")
	}

	if !cached {
		sseDatumUpsertCacheMut.Lock()
		sseDatumUpsertCache[key] = cache
		sseDatumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single SseDatum record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SseDatum) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single SseDatum record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SseDatum) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no SseDatum provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sseDatumPrimaryKeyMapping)
	sql := "DELETE FROM \"sse_data\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete from sse_data")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by delete for sse_data")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sseDatumQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("horizon: no sseDatumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from sse_data")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for sse_data")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o SseDatumSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SseDatumSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no SseDatum slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(sseDatumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseDatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sse_data\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sseDatumPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from sseDatum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for sse_data")
	}

	if len(sseDatumAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SseDatum) ReloadG() error {
	if o == nil {
		return errors.New("horizon: no SseDatum provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SseDatum) Reload(exec boil.Executor) error {
	ret, err := FindSseDatum(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SseDatumSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("horizon: empty SseDatumSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SseDatumSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SseDatumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseDatumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sse_data\".* FROM \"sse_data\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sseDatumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "horizon: unable to reload all in SseDatumSlice")
	}

	*o = slice

	return nil
}

// SseDatumExistsG checks if the SseDatum row exists.
func SseDatumExistsG(iD int) (bool, error) {
	return SseDatumExists(boil.GetDB(), iD)
}

// SseDatumExists checks if the SseDatum row exists.
func SseDatumExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sse_data\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "horizon: unable to check if sse_data exists")
	}

	return exists, nil
}
