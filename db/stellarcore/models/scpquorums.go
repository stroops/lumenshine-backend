// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package stellarcore

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

// Scpquorum is an object representing the database table.
type Scpquorum struct {
	Qsethash      string `boil:"qsethash" json:"qsethash" toml:"qsethash" yaml:"qsethash"`
	Lastledgerseq int    `boil:"lastledgerseq" json:"lastledgerseq" toml:"lastledgerseq" yaml:"lastledgerseq"`
	Qset          string `boil:"qset" json:"qset" toml:"qset" yaml:"qset"`

	R *scpquorumR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L scpquorumL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ScpquorumColumns = struct {
	Qsethash      string
	Lastledgerseq string
	Qset          string
}{
	Qsethash:      "qsethash",
	Lastledgerseq: "lastledgerseq",
	Qset:          "qset",
}

// ScpquorumRels is where relationship names are stored.
var ScpquorumRels = struct {
}{}

// scpquorumR is where relationships are stored.
type scpquorumR struct {
}

// NewStruct creates a new relationship struct
func (*scpquorumR) NewStruct() *scpquorumR {
	return &scpquorumR{}
}

// scpquorumL is where Load methods for each relationship are stored.
type scpquorumL struct{}

var (
	scpquorumColumns               = []string{"qsethash", "lastledgerseq", "qset"}
	scpquorumColumnsWithoutDefault = []string{"qsethash", "lastledgerseq", "qset"}
	scpquorumColumnsWithDefault    = []string{}
	scpquorumPrimaryKeyColumns     = []string{"qsethash"}
)

type (
	// ScpquorumSlice is an alias for a slice of pointers to Scpquorum.
	// This should generally be used opposed to []Scpquorum.
	ScpquorumSlice []*Scpquorum
	// ScpquorumHook is the signature for custom Scpquorum hook methods
	ScpquorumHook func(boil.Executor, *Scpquorum) error

	scpquorumQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	scpquorumType                 = reflect.TypeOf(&Scpquorum{})
	scpquorumMapping              = queries.MakeStructMapping(scpquorumType)
	scpquorumPrimaryKeyMapping, _ = queries.BindMapping(scpquorumType, scpquorumMapping, scpquorumPrimaryKeyColumns)
	scpquorumInsertCacheMut       sync.RWMutex
	scpquorumInsertCache          = make(map[string]insertCache)
	scpquorumUpdateCacheMut       sync.RWMutex
	scpquorumUpdateCache          = make(map[string]updateCache)
	scpquorumUpsertCacheMut       sync.RWMutex
	scpquorumUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var scpquorumBeforeInsertHooks []ScpquorumHook
var scpquorumBeforeUpdateHooks []ScpquorumHook
var scpquorumBeforeDeleteHooks []ScpquorumHook
var scpquorumBeforeUpsertHooks []ScpquorumHook

var scpquorumAfterInsertHooks []ScpquorumHook
var scpquorumAfterSelectHooks []ScpquorumHook
var scpquorumAfterUpdateHooks []ScpquorumHook
var scpquorumAfterDeleteHooks []ScpquorumHook
var scpquorumAfterUpsertHooks []ScpquorumHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Scpquorum) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Scpquorum) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Scpquorum) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Scpquorum) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Scpquorum) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Scpquorum) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Scpquorum) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Scpquorum) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Scpquorum) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range scpquorumAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddScpquorumHook registers your hook function for all future operations.
func AddScpquorumHook(hookPoint boil.HookPoint, scpquorumHook ScpquorumHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		scpquorumBeforeInsertHooks = append(scpquorumBeforeInsertHooks, scpquorumHook)
	case boil.BeforeUpdateHook:
		scpquorumBeforeUpdateHooks = append(scpquorumBeforeUpdateHooks, scpquorumHook)
	case boil.BeforeDeleteHook:
		scpquorumBeforeDeleteHooks = append(scpquorumBeforeDeleteHooks, scpquorumHook)
	case boil.BeforeUpsertHook:
		scpquorumBeforeUpsertHooks = append(scpquorumBeforeUpsertHooks, scpquorumHook)
	case boil.AfterInsertHook:
		scpquorumAfterInsertHooks = append(scpquorumAfterInsertHooks, scpquorumHook)
	case boil.AfterSelectHook:
		scpquorumAfterSelectHooks = append(scpquorumAfterSelectHooks, scpquorumHook)
	case boil.AfterUpdateHook:
		scpquorumAfterUpdateHooks = append(scpquorumAfterUpdateHooks, scpquorumHook)
	case boil.AfterDeleteHook:
		scpquorumAfterDeleteHooks = append(scpquorumAfterDeleteHooks, scpquorumHook)
	case boil.AfterUpsertHook:
		scpquorumAfterUpsertHooks = append(scpquorumAfterUpsertHooks, scpquorumHook)
	}
}

// OneG returns a single scpquorum record from the query using the global executor.
func (q scpquorumQuery) OneG() (*Scpquorum, error) {
	return q.One(boil.GetDB())
}

// One returns a single scpquorum record from the query.
func (q scpquorumQuery) One(exec boil.Executor) (*Scpquorum, error) {
	o := &Scpquorum{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "stellarcore: failed to execute a one query for scpquorums")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Scpquorum records from the query using the global executor.
func (q scpquorumQuery) AllG() (ScpquorumSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all Scpquorum records from the query.
func (q scpquorumQuery) All(exec boil.Executor) (ScpquorumSlice, error) {
	var o []*Scpquorum

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "stellarcore: failed to assign all query results to Scpquorum slice")
	}

	if len(scpquorumAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Scpquorum records in the query, and panics on error.
func (q scpquorumQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all Scpquorum records in the query.
func (q scpquorumQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to count scpquorums rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q scpquorumQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q scpquorumQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "stellarcore: failed to check if scpquorums exists")
	}

	return count > 0, nil
}

// Scpquorums retrieves all the records using an executor.
func Scpquorums(mods ...qm.QueryMod) scpquorumQuery {
	mods = append(mods, qm.From("\"scpquorums\""))
	return scpquorumQuery{NewQuery(mods...)}
}

// FindScpquorumG retrieves a single record by ID.
func FindScpquorumG(qsethash string, selectCols ...string) (*Scpquorum, error) {
	return FindScpquorum(boil.GetDB(), qsethash, selectCols...)
}

// FindScpquorum retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindScpquorum(exec boil.Executor, qsethash string, selectCols ...string) (*Scpquorum, error) {
	scpquorumObj := &Scpquorum{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"scpquorums\" where \"qsethash\"=$1", sel,
	)

	q := queries.Raw(query, qsethash)

	err := q.Bind(nil, exec, scpquorumObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "stellarcore: unable to select from scpquorums")
	}

	return scpquorumObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Scpquorum) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Scpquorum) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("stellarcore: no scpquorums provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(scpquorumColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	scpquorumInsertCacheMut.RLock()
	cache, cached := scpquorumInsertCache[key]
	scpquorumInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			scpquorumColumns,
			scpquorumColumnsWithDefault,
			scpquorumColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(scpquorumType, scpquorumMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(scpquorumType, scpquorumMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"scpquorums\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"scpquorums\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "stellarcore: unable to insert into scpquorums")
	}

	if !cached {
		scpquorumInsertCacheMut.Lock()
		scpquorumInsertCache[key] = cache
		scpquorumInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Scpquorum record using the global executor.
// See Update for more documentation.
func (o *Scpquorum) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the Scpquorum.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Scpquorum) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	scpquorumUpdateCacheMut.RLock()
	cache, cached := scpquorumUpdateCache[key]
	scpquorumUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			scpquorumColumns,
			scpquorumPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("stellarcore: unable to update scpquorums, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"scpquorums\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, scpquorumPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(scpquorumType, scpquorumMapping, append(wl, scpquorumPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "stellarcore: unable to update scpquorums row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by update for scpquorums")
	}

	if !cached {
		scpquorumUpdateCacheMut.Lock()
		scpquorumUpdateCache[key] = cache
		scpquorumUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q scpquorumQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to update all for scpquorums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to retrieve rows affected for scpquorums")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ScpquorumSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ScpquorumSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("stellarcore: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), scpquorumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"scpquorums\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, scpquorumPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to update all in scpquorum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to retrieve rows affected all in update all scpquorum")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Scpquorum) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Scpquorum) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("stellarcore: no scpquorums provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(scpquorumColumnsWithDefault, o)

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

	scpquorumUpsertCacheMut.RLock()
	cache, cached := scpquorumUpsertCache[key]
	scpquorumUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			scpquorumColumns,
			scpquorumColumnsWithDefault,
			scpquorumColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			scpquorumColumns,
			scpquorumPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("stellarcore: unable to upsert scpquorums, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(scpquorumPrimaryKeyColumns))
			copy(conflict, scpquorumPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"scpquorums\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(scpquorumType, scpquorumMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(scpquorumType, scpquorumMapping, ret)
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
		return errors.Wrap(err, "stellarcore: unable to upsert scpquorums")
	}

	if !cached {
		scpquorumUpsertCacheMut.Lock()
		scpquorumUpsertCache[key] = cache
		scpquorumUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single Scpquorum record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Scpquorum) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single Scpquorum record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Scpquorum) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("stellarcore: no Scpquorum provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), scpquorumPrimaryKeyMapping)
	sql := "DELETE FROM \"scpquorums\" WHERE \"qsethash\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete from scpquorums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by delete for scpquorums")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q scpquorumQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("stellarcore: no scpquorumQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete all from scpquorums")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by deleteall for scpquorums")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o ScpquorumSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ScpquorumSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("stellarcore: no Scpquorum slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(scpquorumBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), scpquorumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"scpquorums\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, scpquorumPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete all from scpquorum slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by deleteall for scpquorums")
	}

	if len(scpquorumAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Scpquorum) ReloadG() error {
	if o == nil {
		return errors.New("stellarcore: no Scpquorum provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Scpquorum) Reload(exec boil.Executor) error {
	ret, err := FindScpquorum(exec, o.Qsethash)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ScpquorumSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("stellarcore: empty ScpquorumSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ScpquorumSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ScpquorumSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), scpquorumPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"scpquorums\".* FROM \"scpquorums\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, scpquorumPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "stellarcore: unable to reload all in ScpquorumSlice")
	}

	*o = slice

	return nil
}

// ScpquorumExistsG checks if the Scpquorum row exists.
func ScpquorumExistsG(qsethash string) (bool, error) {
	return ScpquorumExists(boil.GetDB(), qsethash)
}

// ScpquorumExists checks if the Scpquorum row exists.
func ScpquorumExists(exec boil.Executor, qsethash string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"scpquorums\" where \"qsethash\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, qsethash)
	}

	row := exec.QueryRow(sql, qsethash)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "stellarcore: unable to check if scpquorums exists")
	}

	return exists, nil
}
