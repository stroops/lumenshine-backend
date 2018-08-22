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

// Signer is an object representing the database table.
type Signer struct {
	Accountid string `boil:"accountid" json:"accountid" toml:"accountid" yaml:"accountid"`
	Publickey string `boil:"publickey" json:"publickey" toml:"publickey" yaml:"publickey"`
	Weight    int    `boil:"weight" json:"weight" toml:"weight" yaml:"weight"`

	R *signerR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L signerL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SignerColumns = struct {
	Accountid string
	Publickey string
	Weight    string
}{
	Accountid: "accountid",
	Publickey: "publickey",
	Weight:    "weight",
}

// SignerRels is where relationship names are stored.
var SignerRels = struct {
}{}

// signerR is where relationships are stored.
type signerR struct {
}

// NewStruct creates a new relationship struct
func (*signerR) NewStruct() *signerR {
	return &signerR{}
}

// signerL is where Load methods for each relationship are stored.
type signerL struct{}

var (
	signerColumns               = []string{"accountid", "publickey", "weight"}
	signerColumnsWithoutDefault = []string{"accountid", "publickey", "weight"}
	signerColumnsWithDefault    = []string{}
	signerPrimaryKeyColumns     = []string{"accountid", "publickey"}
)

type (
	// SignerSlice is an alias for a slice of pointers to Signer.
	// This should generally be used opposed to []Signer.
	SignerSlice []*Signer
	// SignerHook is the signature for custom Signer hook methods
	SignerHook func(boil.Executor, *Signer) error

	signerQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	signerType                 = reflect.TypeOf(&Signer{})
	signerMapping              = queries.MakeStructMapping(signerType)
	signerPrimaryKeyMapping, _ = queries.BindMapping(signerType, signerMapping, signerPrimaryKeyColumns)
	signerInsertCacheMut       sync.RWMutex
	signerInsertCache          = make(map[string]insertCache)
	signerUpdateCacheMut       sync.RWMutex
	signerUpdateCache          = make(map[string]updateCache)
	signerUpsertCacheMut       sync.RWMutex
	signerUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var signerBeforeInsertHooks []SignerHook
var signerBeforeUpdateHooks []SignerHook
var signerBeforeDeleteHooks []SignerHook
var signerBeforeUpsertHooks []SignerHook

var signerAfterInsertHooks []SignerHook
var signerAfterSelectHooks []SignerHook
var signerAfterUpdateHooks []SignerHook
var signerAfterDeleteHooks []SignerHook
var signerAfterUpsertHooks []SignerHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Signer) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range signerBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Signer) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range signerBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Signer) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range signerBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Signer) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range signerBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Signer) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range signerAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Signer) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range signerAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Signer) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range signerAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Signer) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range signerAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Signer) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range signerAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSignerHook registers your hook function for all future operations.
func AddSignerHook(hookPoint boil.HookPoint, signerHook SignerHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		signerBeforeInsertHooks = append(signerBeforeInsertHooks, signerHook)
	case boil.BeforeUpdateHook:
		signerBeforeUpdateHooks = append(signerBeforeUpdateHooks, signerHook)
	case boil.BeforeDeleteHook:
		signerBeforeDeleteHooks = append(signerBeforeDeleteHooks, signerHook)
	case boil.BeforeUpsertHook:
		signerBeforeUpsertHooks = append(signerBeforeUpsertHooks, signerHook)
	case boil.AfterInsertHook:
		signerAfterInsertHooks = append(signerAfterInsertHooks, signerHook)
	case boil.AfterSelectHook:
		signerAfterSelectHooks = append(signerAfterSelectHooks, signerHook)
	case boil.AfterUpdateHook:
		signerAfterUpdateHooks = append(signerAfterUpdateHooks, signerHook)
	case boil.AfterDeleteHook:
		signerAfterDeleteHooks = append(signerAfterDeleteHooks, signerHook)
	case boil.AfterUpsertHook:
		signerAfterUpsertHooks = append(signerAfterUpsertHooks, signerHook)
	}
}

// OneG returns a single signer record from the query using the global executor.
func (q signerQuery) OneG() (*Signer, error) {
	return q.One(boil.GetDB())
}

// One returns a single signer record from the query.
func (q signerQuery) One(exec boil.Executor) (*Signer, error) {
	o := &Signer{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "stellarcore: failed to execute a one query for signers")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Signer records from the query using the global executor.
func (q signerQuery) AllG() (SignerSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all Signer records from the query.
func (q signerQuery) All(exec boil.Executor) (SignerSlice, error) {
	var o []*Signer

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "stellarcore: failed to assign all query results to Signer slice")
	}

	if len(signerAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Signer records in the query, and panics on error.
func (q signerQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all Signer records in the query.
func (q signerQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to count signers rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q signerQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q signerQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "stellarcore: failed to check if signers exists")
	}

	return count > 0, nil
}

// Signers retrieves all the records using an executor.
func Signers(mods ...qm.QueryMod) signerQuery {
	mods = append(mods, qm.From("\"signers\""))
	return signerQuery{NewQuery(mods...)}
}

// FindSignerG retrieves a single record by ID.
func FindSignerG(accountid string, publickey string, selectCols ...string) (*Signer, error) {
	return FindSigner(boil.GetDB(), accountid, publickey, selectCols...)
}

// FindSigner retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSigner(exec boil.Executor, accountid string, publickey string, selectCols ...string) (*Signer, error) {
	signerObj := &Signer{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"signers\" where \"accountid\"=$1 AND \"publickey\"=$2", sel,
	)

	q := queries.Raw(query, accountid, publickey)

	err := q.Bind(nil, exec, signerObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "stellarcore: unable to select from signers")
	}

	return signerObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Signer) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Signer) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("stellarcore: no signers provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(signerColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	signerInsertCacheMut.RLock()
	cache, cached := signerInsertCache[key]
	signerInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			signerColumns,
			signerColumnsWithDefault,
			signerColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(signerType, signerMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(signerType, signerMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"signers\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"signers\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "stellarcore: unable to insert into signers")
	}

	if !cached {
		signerInsertCacheMut.Lock()
		signerInsertCache[key] = cache
		signerInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Signer record using the global executor.
// See Update for more documentation.
func (o *Signer) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the Signer.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Signer) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	signerUpdateCacheMut.RLock()
	cache, cached := signerUpdateCache[key]
	signerUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			signerColumns,
			signerPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("stellarcore: unable to update signers, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"signers\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, signerPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(signerType, signerMapping, append(wl, signerPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "stellarcore: unable to update signers row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by update for signers")
	}

	if !cached {
		signerUpdateCacheMut.Lock()
		signerUpdateCache[key] = cache
		signerUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q signerQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to update all for signers")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to retrieve rows affected for signers")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SignerSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SignerSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), signerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"signers\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, signerPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to update all in signer slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to retrieve rows affected all in update all signer")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Signer) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Signer) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("stellarcore: no signers provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(signerColumnsWithDefault, o)

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

	signerUpsertCacheMut.RLock()
	cache, cached := signerUpsertCache[key]
	signerUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			signerColumns,
			signerColumnsWithDefault,
			signerColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			signerColumns,
			signerPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("stellarcore: unable to upsert signers, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(signerPrimaryKeyColumns))
			copy(conflict, signerPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"signers\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(signerType, signerMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(signerType, signerMapping, ret)
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
		return errors.Wrap(err, "stellarcore: unable to upsert signers")
	}

	if !cached {
		signerUpsertCacheMut.Lock()
		signerUpsertCache[key] = cache
		signerUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single Signer record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Signer) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single Signer record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Signer) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("stellarcore: no Signer provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), signerPrimaryKeyMapping)
	sql := "DELETE FROM \"signers\" WHERE \"accountid\"=$1 AND \"publickey\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete from signers")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by delete for signers")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q signerQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("stellarcore: no signerQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete all from signers")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by deleteall for signers")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o SignerSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SignerSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("stellarcore: no Signer slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(signerBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), signerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"signers\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, signerPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: unable to delete all from signer slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "stellarcore: failed to get rows affected by deleteall for signers")
	}

	if len(signerAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Signer) ReloadG() error {
	if o == nil {
		return errors.New("stellarcore: no Signer provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Signer) Reload(exec boil.Executor) error {
	ret, err := FindSigner(exec, o.Accountid, o.Publickey)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SignerSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("stellarcore: empty SignerSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SignerSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SignerSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), signerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"signers\".* FROM \"signers\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, signerPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "stellarcore: unable to reload all in SignerSlice")
	}

	*o = slice

	return nil
}

// SignerExistsG checks if the Signer row exists.
func SignerExistsG(accountid string, publickey string) (bool, error) {
	return SignerExists(boil.GetDB(), accountid, publickey)
}

// SignerExists checks if the Signer row exists.
func SignerExists(exec boil.Executor, accountid string, publickey string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"signers\" where \"accountid\"=$1 AND \"publickey\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, accountid, publickey)
	}

	row := exec.QueryRow(sql, accountid, publickey)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "stellarcore: unable to check if signers exists")
	}

	return exists, nil
}