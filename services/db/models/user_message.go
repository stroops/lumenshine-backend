// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

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

// UserMessage is an object representing the database table.
type UserMessage struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID    int       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Title     string    `boil:"title" json:"title" toml:"title" yaml:"title"`
	Message   string    `boil:"message" json:"message" toml:"message" yaml:"message"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	UpdatedBy string    `boil:"updated_by" json:"updated_by" toml:"updated_by" yaml:"updated_by"`

	R *userMessageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userMessageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserMessageColumns = struct {
	ID        string
	UserID    string
	Title     string
	Message   string
	CreatedAt string
	UpdatedAt string
	UpdatedBy string
}{
	ID:        "id",
	UserID:    "user_id",
	Title:     "title",
	Message:   "message",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	UpdatedBy: "updated_by",
}

// UserMessageRels is where relationship names are stored.
var UserMessageRels = struct {
	User string
}{
	User: "User",
}

// userMessageR is where relationships are stored.
type userMessageR struct {
	User *UserProfile
}

// NewStruct creates a new relationship struct
func (*userMessageR) NewStruct() *userMessageR {
	return &userMessageR{}
}

// userMessageL is where Load methods for each relationship are stored.
type userMessageL struct{}

var (
	userMessageColumns               = []string{"id", "user_id", "title", "message", "created_at", "updated_at", "updated_by"}
	userMessageColumnsWithoutDefault = []string{"user_id", "title", "message", "updated_by"}
	userMessageColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	userMessagePrimaryKeyColumns     = []string{"id"}
)

type (
	// UserMessageSlice is an alias for a slice of pointers to UserMessage.
	// This should generally be used opposed to []UserMessage.
	UserMessageSlice []*UserMessage
	// UserMessageHook is the signature for custom UserMessage hook methods
	UserMessageHook func(boil.Executor, *UserMessage) error

	userMessageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userMessageType                 = reflect.TypeOf(&UserMessage{})
	userMessageMapping              = queries.MakeStructMapping(userMessageType)
	userMessagePrimaryKeyMapping, _ = queries.BindMapping(userMessageType, userMessageMapping, userMessagePrimaryKeyColumns)
	userMessageInsertCacheMut       sync.RWMutex
	userMessageInsertCache          = make(map[string]insertCache)
	userMessageUpdateCacheMut       sync.RWMutex
	userMessageUpdateCache          = make(map[string]updateCache)
	userMessageUpsertCacheMut       sync.RWMutex
	userMessageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var userMessageBeforeInsertHooks []UserMessageHook
var userMessageBeforeUpdateHooks []UserMessageHook
var userMessageBeforeDeleteHooks []UserMessageHook
var userMessageBeforeUpsertHooks []UserMessageHook

var userMessageAfterInsertHooks []UserMessageHook
var userMessageAfterSelectHooks []UserMessageHook
var userMessageAfterUpdateHooks []UserMessageHook
var userMessageAfterDeleteHooks []UserMessageHook
var userMessageAfterUpsertHooks []UserMessageHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserMessage) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserMessage) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserMessage) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserMessage) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserMessage) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserMessage) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserMessage) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserMessage) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserMessage) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserMessageHook registers your hook function for all future operations.
func AddUserMessageHook(hookPoint boil.HookPoint, userMessageHook UserMessageHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userMessageBeforeInsertHooks = append(userMessageBeforeInsertHooks, userMessageHook)
	case boil.BeforeUpdateHook:
		userMessageBeforeUpdateHooks = append(userMessageBeforeUpdateHooks, userMessageHook)
	case boil.BeforeDeleteHook:
		userMessageBeforeDeleteHooks = append(userMessageBeforeDeleteHooks, userMessageHook)
	case boil.BeforeUpsertHook:
		userMessageBeforeUpsertHooks = append(userMessageBeforeUpsertHooks, userMessageHook)
	case boil.AfterInsertHook:
		userMessageAfterInsertHooks = append(userMessageAfterInsertHooks, userMessageHook)
	case boil.AfterSelectHook:
		userMessageAfterSelectHooks = append(userMessageAfterSelectHooks, userMessageHook)
	case boil.AfterUpdateHook:
		userMessageAfterUpdateHooks = append(userMessageAfterUpdateHooks, userMessageHook)
	case boil.AfterDeleteHook:
		userMessageAfterDeleteHooks = append(userMessageAfterDeleteHooks, userMessageHook)
	case boil.AfterUpsertHook:
		userMessageAfterUpsertHooks = append(userMessageAfterUpsertHooks, userMessageHook)
	}
}

// OneG returns a single userMessage record from the query using the global executor.
func (q userMessageQuery) OneG() (*UserMessage, error) {
	return q.One(boil.GetDB())
}

// One returns a single userMessage record from the query.
func (q userMessageQuery) One(exec boil.Executor) (*UserMessage, error) {
	o := &UserMessage{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_message")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all UserMessage records from the query using the global executor.
func (q userMessageQuery) AllG() (UserMessageSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all UserMessage records from the query.
func (q userMessageQuery) All(exec boil.Executor) (UserMessageSlice, error) {
	var o []*UserMessage

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserMessage slice")
	}

	if len(userMessageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all UserMessage records in the query, and panics on error.
func (q userMessageQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all UserMessage records in the query.
func (q userMessageQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_message rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q userMessageQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q userMessageQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_message exists")
	}

	return count > 0, nil
}

// User pointed to by the foreign key.
func (o *UserMessage) User(mods ...qm.QueryMod) userProfileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := UserProfiles(queryMods...)
	queries.SetFrom(query.Query, "\"user_profile\"")

	return query
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userMessageL) LoadUser(e boil.Executor, singular bool, maybeUserMessage interface{}, mods queries.Applicator) error {
	var slice []*UserMessage
	var object *UserMessage

	if singular {
		object = maybeUserMessage.(*UserMessage)
	} else {
		slice = *maybeUserMessage.(*[]*UserMessage)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userMessageR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userMessageR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	query := NewQuery(qm.From(`user_profile`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UserProfile")
	}

	var resultSlice []*UserProfile
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UserProfile")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for user_profile")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_profile")
	}

	if len(userMessageAfterSelectHooks) != 0 {
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
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userProfileR{}
		}
		foreign.R.UserUserMessages = append(foreign.R.UserUserMessages, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userProfileR{}
				}
				foreign.R.UserUserMessages = append(foreign.R.UserUserMessages, local)
				break
			}
		}
	}

	return nil
}

// SetUserG of the userMessage to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessages.
// Uses the global database handle.
func (o *UserMessage) SetUserG(insert bool, related *UserProfile) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUser of the userMessage to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessages.
func (o *UserMessage) SetUser(exec boil.Executor, insert bool, related *UserProfile) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_message\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userMessagePrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &userMessageR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userProfileR{
			UserUserMessages: UserMessageSlice{o},
		}
	} else {
		related.R.UserUserMessages = append(related.R.UserUserMessages, o)
	}

	return nil
}

// UserMessages retrieves all the records using an executor.
func UserMessages(mods ...qm.QueryMod) userMessageQuery {
	mods = append(mods, qm.From("\"user_message\""))
	return userMessageQuery{NewQuery(mods...)}
}

// FindUserMessageG retrieves a single record by ID.
func FindUserMessageG(iD int, selectCols ...string) (*UserMessage, error) {
	return FindUserMessage(boil.GetDB(), iD, selectCols...)
}

// FindUserMessage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserMessage(exec boil.Executor, iD int, selectCols ...string) (*UserMessage, error) {
	userMessageObj := &UserMessage{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_message\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, userMessageObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_message")
	}

	return userMessageObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *UserMessage) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserMessage) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_message provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(userMessageColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userMessageInsertCacheMut.RLock()
	cache, cached := userMessageInsertCache[key]
	userMessageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userMessageColumns,
			userMessageColumnsWithDefault,
			userMessageColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userMessageType, userMessageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userMessageType, userMessageMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_message\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_message\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into user_message")
	}

	if !cached {
		userMessageInsertCacheMut.Lock()
		userMessageInsertCache[key] = cache
		userMessageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single UserMessage record using the global executor.
// See Update for more documentation.
func (o *UserMessage) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the UserMessage.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserMessage) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userMessageUpdateCacheMut.RLock()
	cache, cached := userMessageUpdateCache[key]
	userMessageUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userMessageColumns,
			userMessagePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update user_message, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_message\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userMessagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userMessageType, userMessageMapping, append(wl, userMessagePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update user_message row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for user_message")
	}

	if !cached {
		userMessageUpdateCacheMut.Lock()
		userMessageUpdateCache[key] = cache
		userMessageUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q userMessageQuery) UpdateAllG(cols M) (int64, error) {
	return q.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q userMessageQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for user_message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for user_message")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserMessageSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserMessageSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_message\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userMessagePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in userMessage slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all userMessage")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *UserMessage) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserMessage) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_message provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userMessageColumnsWithDefault, o)

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

	userMessageUpsertCacheMut.RLock()
	cache, cached := userMessageUpsertCache[key]
	userMessageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userMessageColumns,
			userMessageColumnsWithDefault,
			userMessageColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userMessageColumns,
			userMessagePrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("models: unable to upsert user_message, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userMessagePrimaryKeyColumns))
			copy(conflict, userMessagePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_message\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userMessageType, userMessageMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userMessageType, userMessageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert user_message")
	}

	if !cached {
		userMessageUpsertCacheMut.Lock()
		userMessageUpsertCache[key] = cache
		userMessageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single UserMessage record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *UserMessage) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single UserMessage record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserMessage) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserMessage provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userMessagePrimaryKeyMapping)
	sql := "DELETE FROM \"user_message\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from user_message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for user_message")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userMessageQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no userMessageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from user_message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_message")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o UserMessageSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserMessageSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserMessage slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(userMessageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_message\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userMessagePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from userMessage slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_message")
	}

	if len(userMessageAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *UserMessage) ReloadG() error {
	if o == nil {
		return errors.New("models: no UserMessage provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserMessage) Reload(exec boil.Executor) error {
	ret, err := FindUserMessage(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserMessageSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty UserMessageSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserMessageSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserMessageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_message\".* FROM \"user_message\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userMessagePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserMessageSlice")
	}

	*o = slice

	return nil
}

// UserMessageExistsG checks if the UserMessage row exists.
func UserMessageExistsG(iD int) (bool, error) {
	return UserMessageExists(boil.GetDB(), iD)
}

// UserMessageExists checks if the UserMessage row exists.
func UserMessageExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_message\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_message exists")
	}

	return exists, nil
}
