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
)

// Mail is an object representing the database table.
type Mail struct {
	ID               int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	MailFrom         string    `boil:"mail_from" json:"mail_from" toml:"mail_from" yaml:"mail_from"`
	MailTo           string    `boil:"mail_to" json:"mail_to" toml:"mail_to" yaml:"mail_to"`
	MailSubject      string    `boil:"mail_subject" json:"mail_subject" toml:"mail_subject" yaml:"mail_subject"`
	MailBody         string    `boil:"mail_body" json:"mail_body" toml:"mail_body" yaml:"mail_body"`
	ExternalStatus   string    `boil:"external_status" json:"external_status" toml:"external_status" yaml:"external_status"`
	ExternalStatusID string    `boil:"external_status_id" json:"external_status_id" toml:"external_status_id" yaml:"external_status_id"`
	InternalStatus   int64     `boil:"internal_status" json:"internal_status" toml:"internal_status" yaml:"internal_status"`
	CreatedAt        time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedBy        string    `boil:"updated_by" json:"updated_by" toml:"updated_by" yaml:"updated_by"`

	R *mailR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L mailL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MailColumns = struct {
	ID               string
	MailFrom         string
	MailTo           string
	MailSubject      string
	MailBody         string
	ExternalStatus   string
	ExternalStatusID string
	InternalStatus   string
	CreatedAt        string
	UpdatedBy        string
}{
	ID:               "id",
	MailFrom:         "mail_from",
	MailTo:           "mail_to",
	MailSubject:      "mail_subject",
	MailBody:         "mail_body",
	ExternalStatus:   "external_status",
	ExternalStatusID: "external_status_id",
	InternalStatus:   "internal_status",
	CreatedAt:        "created_at",
	UpdatedBy:        "updated_by",
}

// mailR is where relationships are stored.
type mailR struct {
}

// mailL is where Load methods for each relationship are stored.
type mailL struct{}

var (
	mailColumns               = []string{"id", "mail_from", "mail_to", "mail_subject", "mail_body", "external_status", "external_status_id", "internal_status", "created_at", "updated_by"}
	mailColumnsWithoutDefault = []string{"mail_from", "mail_to", "mail_subject", "mail_body", "external_status", "external_status_id", "internal_status", "updated_by"}
	mailColumnsWithDefault    = []string{"id", "created_at"}
	mailPrimaryKeyColumns     = []string{"id"}
)

type (
	// MailSlice is an alias for a slice of pointers to Mail.
	// This should generally be used opposed to []Mail.
	MailSlice []*Mail
	// MailHook is the signature for custom Mail hook methods
	MailHook func(boil.Executor, *Mail) error

	mailQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	mailType                 = reflect.TypeOf(&Mail{})
	mailMapping              = queries.MakeStructMapping(mailType)
	mailPrimaryKeyMapping, _ = queries.BindMapping(mailType, mailMapping, mailPrimaryKeyColumns)
	mailInsertCacheMut       sync.RWMutex
	mailInsertCache          = make(map[string]insertCache)
	mailUpdateCacheMut       sync.RWMutex
	mailUpdateCache          = make(map[string]updateCache)
	mailUpsertCacheMut       sync.RWMutex
	mailUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var mailBeforeInsertHooks []MailHook
var mailBeforeUpdateHooks []MailHook
var mailBeforeDeleteHooks []MailHook
var mailBeforeUpsertHooks []MailHook

var mailAfterInsertHooks []MailHook
var mailAfterSelectHooks []MailHook
var mailAfterUpdateHooks []MailHook
var mailAfterDeleteHooks []MailHook
var mailAfterUpsertHooks []MailHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Mail) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range mailBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Mail) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range mailBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Mail) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range mailBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Mail) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range mailBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Mail) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range mailAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Mail) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range mailAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Mail) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range mailAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Mail) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range mailAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Mail) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range mailAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMailHook registers your hook function for all future operations.
func AddMailHook(hookPoint boil.HookPoint, mailHook MailHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		mailBeforeInsertHooks = append(mailBeforeInsertHooks, mailHook)
	case boil.BeforeUpdateHook:
		mailBeforeUpdateHooks = append(mailBeforeUpdateHooks, mailHook)
	case boil.BeforeDeleteHook:
		mailBeforeDeleteHooks = append(mailBeforeDeleteHooks, mailHook)
	case boil.BeforeUpsertHook:
		mailBeforeUpsertHooks = append(mailBeforeUpsertHooks, mailHook)
	case boil.AfterInsertHook:
		mailAfterInsertHooks = append(mailAfterInsertHooks, mailHook)
	case boil.AfterSelectHook:
		mailAfterSelectHooks = append(mailAfterSelectHooks, mailHook)
	case boil.AfterUpdateHook:
		mailAfterUpdateHooks = append(mailAfterUpdateHooks, mailHook)
	case boil.AfterDeleteHook:
		mailAfterDeleteHooks = append(mailAfterDeleteHooks, mailHook)
	case boil.AfterUpsertHook:
		mailAfterUpsertHooks = append(mailAfterUpsertHooks, mailHook)
	}
}

// OneP returns a single mail record from the query, and panics on error.
func (q mailQuery) OneP() *Mail {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single mail record from the query.
func (q mailQuery) One() (*Mail, error) {
	o := &Mail{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for mail")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Mail records from the query, and panics on error.
func (q mailQuery) AllP() MailSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Mail records from the query.
func (q mailQuery) All() (MailSlice, error) {
	var o []*Mail

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Mail slice")
	}

	if len(mailAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Mail records in the query, and panics on error.
func (q mailQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Mail records in the query.
func (q mailQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count mail rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q mailQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q mailQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if mail exists")
	}

	return count > 0, nil
}

// MailsG retrieves all records.
func MailsG(mods ...qm.QueryMod) mailQuery {
	return Mails(boil.GetDB(), mods...)
}

// Mails retrieves all the records using an executor.
func Mails(exec boil.Executor, mods ...qm.QueryMod) mailQuery {
	mods = append(mods, qm.From("\"mail\""))
	return mailQuery{NewQuery(exec, mods...)}
}

// FindMailG retrieves a single record by ID.
func FindMailG(id int, selectCols ...string) (*Mail, error) {
	return FindMail(boil.GetDB(), id, selectCols...)
}

// FindMailGP retrieves a single record by ID, and panics on error.
func FindMailGP(id int, selectCols ...string) *Mail {
	retobj, err := FindMail(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindMail retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMail(exec boil.Executor, id int, selectCols ...string) (*Mail, error) {
	mailObj := &Mail{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"mail\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(mailObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from mail")
	}

	return mailObj, nil
}

// FindMailP retrieves a single record by ID with an executor, and panics on error.
func FindMailP(exec boil.Executor, id int, selectCols ...string) *Mail {
	retobj, err := FindMail(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Mail) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Mail) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Mail) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Mail) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no mail provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mailColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	mailInsertCacheMut.RLock()
	cache, cached := mailInsertCache[key]
	mailInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			mailColumns,
			mailColumnsWithDefault,
			mailColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(mailType, mailMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(mailType, mailMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"mail\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"mail\" DEFAULT VALUES"
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
		return errors.Wrap(err, "models: unable to insert into mail")
	}

	if !cached {
		mailInsertCacheMut.Lock()
		mailInsertCache[key] = cache
		mailInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Mail record. See Update for
// whitelist behavior description.
func (o *Mail) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Mail record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Mail) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Mail, and panics on error.
// See Update for whitelist behavior description.
func (o *Mail) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Mail.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Mail) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	mailUpdateCacheMut.RLock()
	cache, cached := mailUpdateCache[key]
	mailUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			mailColumns,
			mailPrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update mail, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"mail\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, mailPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(mailType, mailMapping, append(wl, mailPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update mail row")
	}

	if !cached {
		mailUpdateCacheMut.Lock()
		mailUpdateCache[key] = cache
		mailUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q mailQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q mailQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for mail")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o MailSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o MailSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o MailSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MailSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"mail\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, mailPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in mail slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Mail) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Mail) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Mail) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Mail) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no mail provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(mailColumnsWithDefault, o)

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

	mailUpsertCacheMut.RLock()
	cache, cached := mailUpsertCache[key]
	mailUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			mailColumns,
			mailColumnsWithDefault,
			mailColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			mailColumns,
			mailPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert mail, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(mailPrimaryKeyColumns))
			copy(conflict, mailPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"mail\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(mailType, mailMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(mailType, mailMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert mail")
	}

	if !cached {
		mailUpsertCacheMut.Lock()
		mailUpsertCache[key] = cache
		mailUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Mail record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Mail) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Mail record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Mail) DeleteG() error {
	if o == nil {
		return errors.New("models: no Mail provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Mail record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Mail) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Mail record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Mail) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Mail provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), mailPrimaryKeyMapping)
	sql := "DELETE FROM \"mail\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from mail")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q mailQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q mailQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no mailQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from mail")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o MailSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o MailSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Mail slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o MailSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MailSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Mail slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(mailBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"mail\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mailPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from mail slice")
	}

	if len(mailAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Mail) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Mail) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Mail) ReloadG() error {
	if o == nil {
		return errors.New("models: no Mail provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Mail) Reload(exec boil.Executor) error {
	ret, err := FindMail(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *MailSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *MailSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MailSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty MailSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MailSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	mails := MailSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mailPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"mail\".* FROM \"mail\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, mailPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&mails)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MailSlice")
	}

	*o = mails

	return nil
}

// MailExists checks if the Mail row exists.
func MailExists(exec boil.Executor, id int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"mail\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if mail exists")
	}

	return exists, nil
}

// MailExistsG checks if the Mail row exists.
func MailExistsG(id int) (bool, error) {
	return MailExists(boil.GetDB(), id)
}

// MailExistsGP checks if the Mail row exists. Panics on error.
func MailExistsGP(id int) bool {
	e, err := MailExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// MailExistsP checks if the Mail row exists. Panics on error.
func MailExistsP(exec boil.Executor, id int) bool {
	e, err := MailExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}