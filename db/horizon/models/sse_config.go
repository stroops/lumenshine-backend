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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// SseConfig is an object representing the database table.
type SseConfig struct {
	ID             int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	SourceReceiver string    `boil:"source_receiver" json:"source_receiver" toml:"source_receiver" yaml:"source_receiver"`
	StellarAccount string    `boil:"stellar_account" json:"stellar_account" toml:"stellar_account" yaml:"stellar_account"`
	OperationTypes int64     `boil:"operation_types" json:"operation_types" toml:"operation_types" yaml:"operation_types"`
	WithResume     bool      `boil:"with_resume" json:"with_resume" toml:"with_resume" yaml:"with_resume"`
	CreatedAt      time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ReturnData     bool      `boil:"return_data" json:"return_data" toml:"return_data" yaml:"return_data"`

	R *sseConfigR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sseConfigL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SseConfigColumns = struct {
	ID             string
	SourceReceiver string
	StellarAccount string
	OperationTypes string
	WithResume     string
	CreatedAt      string
	UpdatedAt      string
	ReturnData     string
}{
	ID:             "id",
	SourceReceiver: "source_receiver",
	StellarAccount: "stellar_account",
	OperationTypes: "operation_types",
	WithResume:     "with_resume",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	ReturnData:     "return_data",
}

// SseConfigRels is where relationship names are stored.
var SseConfigRels = struct {
}{}

// sseConfigR is where relationships are stored.
type sseConfigR struct {
}

// NewStruct creates a new relationship struct
func (*sseConfigR) NewStruct() *sseConfigR {
	return &sseConfigR{}
}

// sseConfigL is where Load methods for each relationship are stored.
type sseConfigL struct{}

var (
	sseConfigColumns               = []string{"id", "source_receiver", "stellar_account", "operation_types", "with_resume", "created_at", "updated_at", "return_data"}
	sseConfigColumnsWithoutDefault = []string{"source_receiver", "stellar_account", "operation_types"}
	sseConfigColumnsWithDefault    = []string{"id", "with_resume", "created_at", "updated_at", "return_data"}
	sseConfigPrimaryKeyColumns     = []string{"id"}
)

type (
	// SseConfigSlice is an alias for a slice of pointers to SseConfig.
	// This should generally be used opposed to []SseConfig.
	SseConfigSlice []*SseConfig
	// SseConfigHook is the signature for custom SseConfig hook methods
	SseConfigHook func(boil.Executor, *SseConfig) error

	sseConfigQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sseConfigType                 = reflect.TypeOf(&SseConfig{})
	sseConfigMapping              = queries.MakeStructMapping(sseConfigType)
	sseConfigPrimaryKeyMapping, _ = queries.BindMapping(sseConfigType, sseConfigMapping, sseConfigPrimaryKeyColumns)
	sseConfigInsertCacheMut       sync.RWMutex
	sseConfigInsertCache          = make(map[string]insertCache)
	sseConfigUpdateCacheMut       sync.RWMutex
	sseConfigUpdateCache          = make(map[string]updateCache)
	sseConfigUpsertCacheMut       sync.RWMutex
	sseConfigUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var sseConfigBeforeInsertHooks []SseConfigHook
var sseConfigBeforeUpdateHooks []SseConfigHook
var sseConfigBeforeDeleteHooks []SseConfigHook
var sseConfigBeforeUpsertHooks []SseConfigHook

var sseConfigAfterInsertHooks []SseConfigHook
var sseConfigAfterSelectHooks []SseConfigHook
var sseConfigAfterUpdateHooks []SseConfigHook
var sseConfigAfterDeleteHooks []SseConfigHook
var sseConfigAfterUpsertHooks []SseConfigHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SseConfig) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SseConfig) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SseConfig) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SseConfig) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SseConfig) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SseConfig) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SseConfig) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SseConfig) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SseConfig) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range sseConfigAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSseConfigHook registers your hook function for all future operations.
func AddSseConfigHook(hookPoint boil.HookPoint, sseConfigHook SseConfigHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		sseConfigBeforeInsertHooks = append(sseConfigBeforeInsertHooks, sseConfigHook)
	case boil.BeforeUpdateHook:
		sseConfigBeforeUpdateHooks = append(sseConfigBeforeUpdateHooks, sseConfigHook)
	case boil.BeforeDeleteHook:
		sseConfigBeforeDeleteHooks = append(sseConfigBeforeDeleteHooks, sseConfigHook)
	case boil.BeforeUpsertHook:
		sseConfigBeforeUpsertHooks = append(sseConfigBeforeUpsertHooks, sseConfigHook)
	case boil.AfterInsertHook:
		sseConfigAfterInsertHooks = append(sseConfigAfterInsertHooks, sseConfigHook)
	case boil.AfterSelectHook:
		sseConfigAfterSelectHooks = append(sseConfigAfterSelectHooks, sseConfigHook)
	case boil.AfterUpdateHook:
		sseConfigAfterUpdateHooks = append(sseConfigAfterUpdateHooks, sseConfigHook)
	case boil.AfterDeleteHook:
		sseConfigAfterDeleteHooks = append(sseConfigAfterDeleteHooks, sseConfigHook)
	case boil.AfterUpsertHook:
		sseConfigAfterUpsertHooks = append(sseConfigAfterUpsertHooks, sseConfigHook)
	}
}

// OneG returns a single sseConfig record from the query using the global executor.
func (q sseConfigQuery) OneG() (*SseConfig, error) {
	return q.One(boil.GetDB())
}

// One returns a single sseConfig record from the query.
func (q sseConfigQuery) One(exec boil.Executor) (*SseConfig, error) {
	o := &SseConfig{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: failed to execute a one query for sse_config")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all SseConfig records from the query using the global executor.
func (q sseConfigQuery) AllG() (SseConfigSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all SseConfig records from the query.
func (q sseConfigQuery) All(exec boil.Executor) (SseConfigSlice, error) {
	var o []*SseConfig

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "horizon: failed to assign all query results to SseConfig slice")
	}

	if len(sseConfigAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all SseConfig records in the query, and panics on error.
func (q sseConfigQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all SseConfig records in the query.
func (q sseConfigQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to count sse_config rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q sseConfigQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q sseConfigQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "horizon: failed to check if sse_config exists")
	}

	return count > 0, nil
}

// SseConfigs retrieves all the records using an executor.
func SseConfigs(mods ...qm.QueryMod) sseConfigQuery {
	mods = append(mods, qm.From("\"sse_config\""))
	return sseConfigQuery{NewQuery(mods...)}
}

// FindSseConfigG retrieves a single record by ID.
func FindSseConfigG(iD int, selectCols ...string) (*SseConfig, error) {
	return FindSseConfig(boil.GetDB(), iD, selectCols...)
}

// FindSseConfig retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSseConfig(exec boil.Executor, iD int, selectCols ...string) (*SseConfig, error) {
	sseConfigObj := &SseConfig{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sse_config\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, sseConfigObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "horizon: unable to select from sse_config")
	}

	return sseConfigObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SseConfig) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SseConfig) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no sse_config provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(sseConfigColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sseConfigInsertCacheMut.RLock()
	cache, cached := sseConfigInsertCache[key]
	sseConfigInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sseConfigColumns,
			sseConfigColumnsWithDefault,
			sseConfigColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sseConfigType, sseConfigMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sseConfigType, sseConfigMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sse_config\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sse_config\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "horizon: unable to insert into sse_config")
	}

	if !cached {
		sseConfigInsertCacheMut.Lock()
		sseConfigInsertCache[key] = cache
		sseConfigInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single SseConfig record using the global executor.
// See Update for more documentation.
func (o *SseConfig) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the SseConfig.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SseConfig) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	sseConfigUpdateCacheMut.RLock()
	cache, cached := sseConfigUpdateCache[key]
	sseConfigUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sseConfigColumns,
			sseConfigPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("horizon: unable to update sse_config, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sse_config\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sseConfigPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sseConfigType, sseConfigMapping, append(wl, sseConfigPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "horizon: unable to update sse_config row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by update for sse_config")
	}

	if !cached {
		sseConfigUpdateCacheMut.Lock()
		sseConfigUpdateCache[key] = cache
		sseConfigUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q sseConfigQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all for sse_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected for sse_config")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SseConfigSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SseConfigSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sse_config\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sseConfigPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to update all in sseConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to retrieve rows affected all in update all sseConfig")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SseConfig) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SseConfig) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("horizon: no sse_config provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sseConfigColumnsWithDefault, o)

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

	sseConfigUpsertCacheMut.RLock()
	cache, cached := sseConfigUpsertCache[key]
	sseConfigUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			sseConfigColumns,
			sseConfigColumnsWithDefault,
			sseConfigColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			sseConfigColumns,
			sseConfigPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("horizon: unable to upsert sse_config, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sseConfigPrimaryKeyColumns))
			copy(conflict, sseConfigPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sse_config\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(sseConfigType, sseConfigMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sseConfigType, sseConfigMapping, ret)
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
		return errors.Wrap(err, "horizon: unable to upsert sse_config")
	}

	if !cached {
		sseConfigUpsertCacheMut.Lock()
		sseConfigUpsertCache[key] = cache
		sseConfigUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single SseConfig record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SseConfig) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single SseConfig record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SseConfig) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no SseConfig provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sseConfigPrimaryKeyMapping)
	sql := "DELETE FROM \"sse_config\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete from sse_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by delete for sse_config")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sseConfigQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("horizon: no sseConfigQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from sse_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for sse_config")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o SseConfigSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SseConfigSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("horizon: no SseConfig slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(sseConfigBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sse_config\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sseConfigPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "horizon: unable to delete all from sseConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "horizon: failed to get rows affected by deleteall for sse_config")
	}

	if len(sseConfigAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SseConfig) ReloadG() error {
	if o == nil {
		return errors.New("horizon: no SseConfig provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SseConfig) Reload(exec boil.Executor) error {
	ret, err := FindSseConfig(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SseConfigSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("horizon: empty SseConfigSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SseConfigSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SseConfigSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sseConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sse_config\".* FROM \"sse_config\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sseConfigPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "horizon: unable to reload all in SseConfigSlice")
	}

	*o = slice

	return nil
}

// SseConfigExistsG checks if the SseConfig row exists.
func SseConfigExistsG(iD int) (bool, error) {
	return SseConfigExists(boil.GetDB(), iD)
}

// SseConfigExists checks if the SseConfig row exists.
func SseConfigExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sse_config\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "horizon: unable to check if sse_config exists")
	}

	return exists, nil
}
