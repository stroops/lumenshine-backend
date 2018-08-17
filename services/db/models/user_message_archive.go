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

// UserMessageArchive is an object representing the database table.
type UserMessageArchive struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID    int       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Title     string    `boil:"title" json:"title" toml:"title" yaml:"title"`
	Message   string    `boil:"message" json:"message" toml:"message" yaml:"message"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	UpdatedBy string    `boil:"updated_by" json:"updated_by" toml:"updated_by" yaml:"updated_by"`

	R *userMessageArchiveR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userMessageArchiveL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserMessageArchiveColumns = struct {
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

// userMessageArchiveR is where relationships are stored.
type userMessageArchiveR struct {
	User *UserProfile
}

// userMessageArchiveL is where Load methods for each relationship are stored.
type userMessageArchiveL struct{}

var (
	userMessageArchiveColumns               = []string{"id", "user_id", "title", "message", "created_at", "updated_at", "updated_by"}
	userMessageArchiveColumnsWithoutDefault = []string{"user_id", "title", "message", "updated_by"}
	userMessageArchiveColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	userMessageArchivePrimaryKeyColumns     = []string{"id"}
)

type (
	// UserMessageArchiveSlice is an alias for a slice of pointers to UserMessageArchive.
	// This should generally be used opposed to []UserMessageArchive.
	UserMessageArchiveSlice []*UserMessageArchive
	// UserMessageArchiveHook is the signature for custom UserMessageArchive hook methods
	UserMessageArchiveHook func(boil.Executor, *UserMessageArchive) error

	userMessageArchiveQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userMessageArchiveType                 = reflect.TypeOf(&UserMessageArchive{})
	userMessageArchiveMapping              = queries.MakeStructMapping(userMessageArchiveType)
	userMessageArchivePrimaryKeyMapping, _ = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, userMessageArchivePrimaryKeyColumns)
	userMessageArchiveInsertCacheMut       sync.RWMutex
	userMessageArchiveInsertCache          = make(map[string]insertCache)
	userMessageArchiveUpdateCacheMut       sync.RWMutex
	userMessageArchiveUpdateCache          = make(map[string]updateCache)
	userMessageArchiveUpsertCacheMut       sync.RWMutex
	userMessageArchiveUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var userMessageArchiveBeforeInsertHooks []UserMessageArchiveHook
var userMessageArchiveBeforeUpdateHooks []UserMessageArchiveHook
var userMessageArchiveBeforeDeleteHooks []UserMessageArchiveHook
var userMessageArchiveBeforeUpsertHooks []UserMessageArchiveHook

var userMessageArchiveAfterInsertHooks []UserMessageArchiveHook
var userMessageArchiveAfterSelectHooks []UserMessageArchiveHook
var userMessageArchiveAfterUpdateHooks []UserMessageArchiveHook
var userMessageArchiveAfterDeleteHooks []UserMessageArchiveHook
var userMessageArchiveAfterUpsertHooks []UserMessageArchiveHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserMessageArchive) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserMessageArchive) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserMessageArchive) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserMessageArchive) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserMessageArchive) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserMessageArchive) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserMessageArchive) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserMessageArchive) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserMessageArchive) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userMessageArchiveAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserMessageArchiveHook registers your hook function for all future operations.
func AddUserMessageArchiveHook(hookPoint boil.HookPoint, userMessageArchiveHook UserMessageArchiveHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userMessageArchiveBeforeInsertHooks = append(userMessageArchiveBeforeInsertHooks, userMessageArchiveHook)
	case boil.BeforeUpdateHook:
		userMessageArchiveBeforeUpdateHooks = append(userMessageArchiveBeforeUpdateHooks, userMessageArchiveHook)
	case boil.BeforeDeleteHook:
		userMessageArchiveBeforeDeleteHooks = append(userMessageArchiveBeforeDeleteHooks, userMessageArchiveHook)
	case boil.BeforeUpsertHook:
		userMessageArchiveBeforeUpsertHooks = append(userMessageArchiveBeforeUpsertHooks, userMessageArchiveHook)
	case boil.AfterInsertHook:
		userMessageArchiveAfterInsertHooks = append(userMessageArchiveAfterInsertHooks, userMessageArchiveHook)
	case boil.AfterSelectHook:
		userMessageArchiveAfterSelectHooks = append(userMessageArchiveAfterSelectHooks, userMessageArchiveHook)
	case boil.AfterUpdateHook:
		userMessageArchiveAfterUpdateHooks = append(userMessageArchiveAfterUpdateHooks, userMessageArchiveHook)
	case boil.AfterDeleteHook:
		userMessageArchiveAfterDeleteHooks = append(userMessageArchiveAfterDeleteHooks, userMessageArchiveHook)
	case boil.AfterUpsertHook:
		userMessageArchiveAfterUpsertHooks = append(userMessageArchiveAfterUpsertHooks, userMessageArchiveHook)
	}
}

// OneP returns a single userMessageArchive record from the query, and panics on error.
func (q userMessageArchiveQuery) OneP() *UserMessageArchive {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single userMessageArchive record from the query.
func (q userMessageArchiveQuery) One() (*UserMessageArchive, error) {
	o := &UserMessageArchive{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_message_archive")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all UserMessageArchive records from the query, and panics on error.
func (q userMessageArchiveQuery) AllP() UserMessageArchiveSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all UserMessageArchive records from the query.
func (q userMessageArchiveQuery) All() (UserMessageArchiveSlice, error) {
	var o []*UserMessageArchive

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserMessageArchive slice")
	}

	if len(userMessageArchiveAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all UserMessageArchive records in the query, and panics on error.
func (q userMessageArchiveQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all UserMessageArchive records in the query.
func (q userMessageArchiveQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_message_archive rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q userMessageArchiveQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q userMessageArchiveQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_message_archive exists")
	}

	return count > 0, nil
}

// UserG pointed to by the foreign key.
func (o *UserMessageArchive) UserG(mods ...qm.QueryMod) userProfileQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *UserMessageArchive) User(exec boil.Executor, mods ...qm.QueryMod) userProfileQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := UserProfiles(exec, queryMods...)
	queries.SetFrom(query.Query, "\"user_profile\"")

	return query
} // LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userMessageArchiveL) LoadUser(e boil.Executor, singular bool, maybeUserMessageArchive interface{}) error {
	var slice []*UserMessageArchive
	var object *UserMessageArchive

	count := 1
	if singular {
		object = maybeUserMessageArchive.(*UserMessageArchive)
	} else {
		slice = *maybeUserMessageArchive.(*[]*UserMessageArchive)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userMessageArchiveR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userMessageArchiveR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from \"user_profile\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load UserProfile")
	}
	defer results.Close()

	var resultSlice []*UserProfile
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice UserProfile")
	}

	if len(userMessageArchiveAfterSelectHooks) != 0 {
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
		object.R.User = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetUserG of the user_message_archive to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessageArchives.
// Uses the global database handle.
func (o *UserMessageArchive) SetUserG(insert bool, related *UserProfile) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the user_message_archive to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessageArchives.
// Panics on error.
func (o *UserMessageArchive) SetUserP(exec boil.Executor, insert bool, related *UserProfile) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the user_message_archive to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessageArchives.
// Uses the global database handle and panics on error.
func (o *UserMessageArchive) SetUserGP(insert bool, related *UserProfile) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the user_message_archive to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserUserMessageArchives.
func (o *UserMessageArchive) SetUser(exec boil.Executor, insert bool, related *UserProfile) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_message_archive\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userMessageArchivePrimaryKeyColumns),
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
		o.R = &userMessageArchiveR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userProfileR{
			UserUserMessageArchives: UserMessageArchiveSlice{o},
		}
	} else {
		related.R.UserUserMessageArchives = append(related.R.UserUserMessageArchives, o)
	}

	return nil
}

// UserMessageArchivesG retrieves all records.
func UserMessageArchivesG(mods ...qm.QueryMod) userMessageArchiveQuery {
	return UserMessageArchives(boil.GetDB(), mods...)
}

// UserMessageArchives retrieves all the records using an executor.
func UserMessageArchives(exec boil.Executor, mods ...qm.QueryMod) userMessageArchiveQuery {
	mods = append(mods, qm.From("\"user_message_archive\""))
	return userMessageArchiveQuery{NewQuery(exec, mods...)}
}

// FindUserMessageArchiveG retrieves a single record by ID.
func FindUserMessageArchiveG(id int, selectCols ...string) (*UserMessageArchive, error) {
	return FindUserMessageArchive(boil.GetDB(), id, selectCols...)
}

// FindUserMessageArchiveGP retrieves a single record by ID, and panics on error.
func FindUserMessageArchiveGP(id int, selectCols ...string) *UserMessageArchive {
	retobj, err := FindUserMessageArchive(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindUserMessageArchive retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserMessageArchive(exec boil.Executor, id int, selectCols ...string) (*UserMessageArchive, error) {
	userMessageArchiveObj := &UserMessageArchive{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_message_archive\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(userMessageArchiveObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_message_archive")
	}

	return userMessageArchiveObj, nil
}

// FindUserMessageArchiveP retrieves a single record by ID with an executor, and panics on error.
func FindUserMessageArchiveP(exec boil.Executor, id int, selectCols ...string) *UserMessageArchive {
	retobj, err := FindUserMessageArchive(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *UserMessageArchive) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *UserMessageArchive) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *UserMessageArchive) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *UserMessageArchive) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no user_message_archive provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(userMessageArchiveColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	userMessageArchiveInsertCacheMut.RLock()
	cache, cached := userMessageArchiveInsertCache[key]
	userMessageArchiveInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			userMessageArchiveColumns,
			userMessageArchiveColumnsWithDefault,
			userMessageArchiveColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_message_archive\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_message_archive\" DEFAULT VALUES"
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
		return errors.Wrap(err, "models: unable to insert into user_message_archive")
	}

	if !cached {
		userMessageArchiveInsertCacheMut.Lock()
		userMessageArchiveInsertCache[key] = cache
		userMessageArchiveInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single UserMessageArchive record. See Update for
// whitelist behavior description.
func (o *UserMessageArchive) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single UserMessageArchive record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *UserMessageArchive) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the UserMessageArchive, and panics on error.
// See Update for whitelist behavior description.
func (o *UserMessageArchive) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the UserMessageArchive.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *UserMessageArchive) Update(exec boil.Executor, whitelist ...string) error {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	userMessageArchiveUpdateCacheMut.RLock()
	cache, cached := userMessageArchiveUpdateCache[key]
	userMessageArchiveUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			userMessageArchiveColumns,
			userMessageArchivePrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update user_message_archive, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_message_archive\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userMessageArchivePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, append(wl, userMessageArchivePrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update user_message_archive row")
	}

	if !cached {
		userMessageArchiveUpdateCacheMut.Lock()
		userMessageArchiveUpdateCache[key] = cache
		userMessageArchiveUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q userMessageArchiveQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q userMessageArchiveQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for user_message_archive")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserMessageArchiveSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o UserMessageArchiveSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o UserMessageArchiveSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserMessageArchiveSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessageArchivePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_message_archive\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userMessageArchivePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in userMessageArchive slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *UserMessageArchive) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *UserMessageArchive) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *UserMessageArchive) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *UserMessageArchive) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no user_message_archive provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userMessageArchiveColumnsWithDefault, o)

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

	userMessageArchiveUpsertCacheMut.RLock()
	cache, cached := userMessageArchiveUpsertCache[key]
	userMessageArchiveUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			userMessageArchiveColumns,
			userMessageArchiveColumnsWithDefault,
			userMessageArchiveColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			userMessageArchiveColumns,
			userMessageArchivePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert user_message_archive, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userMessageArchivePrimaryKeyColumns))
			copy(conflict, userMessageArchivePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"user_message_archive\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userMessageArchiveType, userMessageArchiveMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert user_message_archive")
	}

	if !cached {
		userMessageArchiveUpsertCacheMut.Lock()
		userMessageArchiveUpsertCache[key] = cache
		userMessageArchiveUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single UserMessageArchive record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *UserMessageArchive) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single UserMessageArchive record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *UserMessageArchive) DeleteG() error {
	if o == nil {
		return errors.New("models: no UserMessageArchive provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single UserMessageArchive record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *UserMessageArchive) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single UserMessageArchive record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserMessageArchive) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no UserMessageArchive provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userMessageArchivePrimaryKeyMapping)
	sql := "DELETE FROM \"user_message_archive\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from user_message_archive")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q userMessageArchiveQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q userMessageArchiveQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no userMessageArchiveQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from user_message_archive")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o UserMessageArchiveSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o UserMessageArchiveSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no UserMessageArchive slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o UserMessageArchiveSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserMessageArchiveSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no UserMessageArchive slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(userMessageArchiveBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessageArchivePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_message_archive\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userMessageArchivePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from userMessageArchive slice")
	}

	if len(userMessageArchiveAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *UserMessageArchive) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *UserMessageArchive) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *UserMessageArchive) ReloadG() error {
	if o == nil {
		return errors.New("models: no UserMessageArchive provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserMessageArchive) Reload(exec boil.Executor) error {
	ret, err := FindUserMessageArchive(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserMessageArchiveSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserMessageArchiveSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserMessageArchiveSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty UserMessageArchiveSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserMessageArchiveSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	userMessageArchives := UserMessageArchiveSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userMessageArchivePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_message_archive\".* FROM \"user_message_archive\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userMessageArchivePrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&userMessageArchives)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserMessageArchiveSlice")
	}

	*o = userMessageArchives

	return nil
}

// UserMessageArchiveExists checks if the UserMessageArchive row exists.
func UserMessageArchiveExists(exec boil.Executor, id int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_message_archive\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_message_archive exists")
	}

	return exists, nil
}

// UserMessageArchiveExistsG checks if the UserMessageArchive row exists.
func UserMessageArchiveExistsG(id int) (bool, error) {
	return UserMessageArchiveExists(boil.GetDB(), id)
}

// UserMessageArchiveExistsGP checks if the UserMessageArchive row exists. Panics on error.
func UserMessageArchiveExistsGP(id int) bool {
	e, err := UserMessageArchiveExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// UserMessageArchiveExistsP checks if the UserMessageArchive row exists. Panics on error.
func UserMessageArchiveExistsP(exec boil.Executor, id int) bool {
	e, err := UserMessageArchiveExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}