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

// AdminKnownInflationDestination is an object representing the database table.
type AdminKnownInflationDestination struct {
	ID               int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name             string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	IssuerPublicKey  string    `boil:"issuer_public_key" json:"issuer_public_key" toml:"issuer_public_key" yaml:"issuer_public_key"`
	ShortDescription string    `boil:"short_description" json:"short_description" toml:"short_description" yaml:"short_description"`
	LongDescription  string    `boil:"long_description" json:"long_description" toml:"long_description" yaml:"long_description"`
	OrderIndex       int       `boil:"order_index" json:"order_index" toml:"order_index" yaml:"order_index"`
	CreatedAt        time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt        time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	UpdatedBy        string    `boil:"updated_by" json:"updated_by" toml:"updated_by" yaml:"updated_by"`

	R *adminKnownInflationDestinationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L adminKnownInflationDestinationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AdminKnownInflationDestinationColumns = struct {
	ID               string
	Name             string
	IssuerPublicKey  string
	ShortDescription string
	LongDescription  string
	OrderIndex       string
	CreatedAt        string
	UpdatedAt        string
	UpdatedBy        string
}{
	ID:               "id",
	Name:             "name",
	IssuerPublicKey:  "issuer_public_key",
	ShortDescription: "short_description",
	LongDescription:  "long_description",
	OrderIndex:       "order_index",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	UpdatedBy:        "updated_by",
}

// AdminKnownInflationDestinationRels is where relationship names are stored.
var AdminKnownInflationDestinationRels = struct {
}{}

// adminKnownInflationDestinationR is where relationships are stored.
type adminKnownInflationDestinationR struct {
}

// NewStruct creates a new relationship struct
func (*adminKnownInflationDestinationR) NewStruct() *adminKnownInflationDestinationR {
	return &adminKnownInflationDestinationR{}
}

// adminKnownInflationDestinationL is where Load methods for each relationship are stored.
type adminKnownInflationDestinationL struct{}

var (
	adminKnownInflationDestinationColumns               = []string{"id", "name", "issuer_public_key", "short_description", "long_description", "order_index", "created_at", "updated_at", "updated_by"}
	adminKnownInflationDestinationColumnsWithoutDefault = []string{"name", "issuer_public_key", "short_description", "long_description", "order_index", "updated_by"}
	adminKnownInflationDestinationColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	adminKnownInflationDestinationPrimaryKeyColumns     = []string{"id"}
)

type (
	// AdminKnownInflationDestinationSlice is an alias for a slice of pointers to AdminKnownInflationDestination.
	// This should generally be used opposed to []AdminKnownInflationDestination.
	AdminKnownInflationDestinationSlice []*AdminKnownInflationDestination
	// AdminKnownInflationDestinationHook is the signature for custom AdminKnownInflationDestination hook methods
	AdminKnownInflationDestinationHook func(boil.Executor, *AdminKnownInflationDestination) error

	adminKnownInflationDestinationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	adminKnownInflationDestinationType                 = reflect.TypeOf(&AdminKnownInflationDestination{})
	adminKnownInflationDestinationMapping              = queries.MakeStructMapping(adminKnownInflationDestinationType)
	adminKnownInflationDestinationPrimaryKeyMapping, _ = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, adminKnownInflationDestinationPrimaryKeyColumns)
	adminKnownInflationDestinationInsertCacheMut       sync.RWMutex
	adminKnownInflationDestinationInsertCache          = make(map[string]insertCache)
	adminKnownInflationDestinationUpdateCacheMut       sync.RWMutex
	adminKnownInflationDestinationUpdateCache          = make(map[string]updateCache)
	adminKnownInflationDestinationUpsertCacheMut       sync.RWMutex
	adminKnownInflationDestinationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
)

var adminKnownInflationDestinationBeforeInsertHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationBeforeUpdateHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationBeforeDeleteHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationBeforeUpsertHooks []AdminKnownInflationDestinationHook

var adminKnownInflationDestinationAfterInsertHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationAfterSelectHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationAfterUpdateHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationAfterDeleteHooks []AdminKnownInflationDestinationHook
var adminKnownInflationDestinationAfterUpsertHooks []AdminKnownInflationDestinationHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *AdminKnownInflationDestination) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *AdminKnownInflationDestination) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *AdminKnownInflationDestination) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *AdminKnownInflationDestination) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *AdminKnownInflationDestination) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *AdminKnownInflationDestination) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *AdminKnownInflationDestination) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *AdminKnownInflationDestination) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *AdminKnownInflationDestination) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range adminKnownInflationDestinationAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAdminKnownInflationDestinationHook registers your hook function for all future operations.
func AddAdminKnownInflationDestinationHook(hookPoint boil.HookPoint, adminKnownInflationDestinationHook AdminKnownInflationDestinationHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		adminKnownInflationDestinationBeforeInsertHooks = append(adminKnownInflationDestinationBeforeInsertHooks, adminKnownInflationDestinationHook)
	case boil.BeforeUpdateHook:
		adminKnownInflationDestinationBeforeUpdateHooks = append(adminKnownInflationDestinationBeforeUpdateHooks, adminKnownInflationDestinationHook)
	case boil.BeforeDeleteHook:
		adminKnownInflationDestinationBeforeDeleteHooks = append(adminKnownInflationDestinationBeforeDeleteHooks, adminKnownInflationDestinationHook)
	case boil.BeforeUpsertHook:
		adminKnownInflationDestinationBeforeUpsertHooks = append(adminKnownInflationDestinationBeforeUpsertHooks, adminKnownInflationDestinationHook)
	case boil.AfterInsertHook:
		adminKnownInflationDestinationAfterInsertHooks = append(adminKnownInflationDestinationAfterInsertHooks, adminKnownInflationDestinationHook)
	case boil.AfterSelectHook:
		adminKnownInflationDestinationAfterSelectHooks = append(adminKnownInflationDestinationAfterSelectHooks, adminKnownInflationDestinationHook)
	case boil.AfterUpdateHook:
		adminKnownInflationDestinationAfterUpdateHooks = append(adminKnownInflationDestinationAfterUpdateHooks, adminKnownInflationDestinationHook)
	case boil.AfterDeleteHook:
		adminKnownInflationDestinationAfterDeleteHooks = append(adminKnownInflationDestinationAfterDeleteHooks, adminKnownInflationDestinationHook)
	case boil.AfterUpsertHook:
		adminKnownInflationDestinationAfterUpsertHooks = append(adminKnownInflationDestinationAfterUpsertHooks, adminKnownInflationDestinationHook)
	}
}

// OneG returns a single adminKnownInflationDestination record from the query using the global executor.
func (q adminKnownInflationDestinationQuery) OneG() (*AdminKnownInflationDestination, error) {
	return q.One(boil.GetDB())
}

// One returns a single adminKnownInflationDestination record from the query.
func (q adminKnownInflationDestinationQuery) One(exec boil.Executor) (*AdminKnownInflationDestination, error) {
	o := &AdminKnownInflationDestination{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for admin_known_inflation_destinations")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all AdminKnownInflationDestination records from the query using the global executor.
func (q adminKnownInflationDestinationQuery) AllG() (AdminKnownInflationDestinationSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all AdminKnownInflationDestination records from the query.
func (q adminKnownInflationDestinationQuery) All(exec boil.Executor) (AdminKnownInflationDestinationSlice, error) {
	var o []*AdminKnownInflationDestination

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to AdminKnownInflationDestination slice")
	}

	if len(adminKnownInflationDestinationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all AdminKnownInflationDestination records in the query, and panics on error.
func (q adminKnownInflationDestinationQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all AdminKnownInflationDestination records in the query.
func (q adminKnownInflationDestinationQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count admin_known_inflation_destinations rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q adminKnownInflationDestinationQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q adminKnownInflationDestinationQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if admin_known_inflation_destinations exists")
	}

	return count > 0, nil
}

// AdminKnownInflationDestinations retrieves all the records using an executor.
func AdminKnownInflationDestinations(mods ...qm.QueryMod) adminKnownInflationDestinationQuery {
	mods = append(mods, qm.From("\"admin_known_inflation_destinations\""))
	return adminKnownInflationDestinationQuery{NewQuery(mods...)}
}

// FindAdminKnownInflationDestinationG retrieves a single record by ID.
func FindAdminKnownInflationDestinationG(iD int, selectCols ...string) (*AdminKnownInflationDestination, error) {
	return FindAdminKnownInflationDestination(boil.GetDB(), iD, selectCols...)
}

// FindAdminKnownInflationDestination retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAdminKnownInflationDestination(exec boil.Executor, iD int, selectCols ...string) (*AdminKnownInflationDestination, error) {
	adminKnownInflationDestinationObj := &AdminKnownInflationDestination{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"admin_known_inflation_destinations\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, adminKnownInflationDestinationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from admin_known_inflation_destinations")
	}

	return adminKnownInflationDestinationObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *AdminKnownInflationDestination) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AdminKnownInflationDestination) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no admin_known_inflation_destinations provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(adminKnownInflationDestinationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	adminKnownInflationDestinationInsertCacheMut.RLock()
	cache, cached := adminKnownInflationDestinationInsertCache[key]
	adminKnownInflationDestinationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			adminKnownInflationDestinationColumns,
			adminKnownInflationDestinationColumnsWithDefault,
			adminKnownInflationDestinationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"admin_known_inflation_destinations\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"admin_known_inflation_destinations\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into admin_known_inflation_destinations")
	}

	if !cached {
		adminKnownInflationDestinationInsertCacheMut.Lock()
		adminKnownInflationDestinationInsertCache[key] = cache
		adminKnownInflationDestinationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single AdminKnownInflationDestination record using the global executor.
// See Update for more documentation.
func (o *AdminKnownInflationDestination) UpdateG(columns boil.Columns) (int64, error) {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the AdminKnownInflationDestination.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AdminKnownInflationDestination) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	currTime := time.Now().In(boil.GetLocation())

	o.UpdatedAt = currTime

	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	adminKnownInflationDestinationUpdateCacheMut.RLock()
	cache, cached := adminKnownInflationDestinationUpdateCache[key]
	adminKnownInflationDestinationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			adminKnownInflationDestinationColumns,
			adminKnownInflationDestinationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update admin_known_inflation_destinations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"admin_known_inflation_destinations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, adminKnownInflationDestinationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, append(wl, adminKnownInflationDestinationPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update admin_known_inflation_destinations row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for admin_known_inflation_destinations")
	}

	if !cached {
		adminKnownInflationDestinationUpdateCacheMut.Lock()
		adminKnownInflationDestinationUpdateCache[key] = cache
		adminKnownInflationDestinationUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(exec)
}

// UpdateAll updates all rows with the specified column values.
func (q adminKnownInflationDestinationQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for admin_known_inflation_destinations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for admin_known_inflation_destinations")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o AdminKnownInflationDestinationSlice) UpdateAllG(cols M) (int64, error) {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AdminKnownInflationDestinationSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), adminKnownInflationDestinationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"admin_known_inflation_destinations\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, adminKnownInflationDestinationPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in adminKnownInflationDestination slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all adminKnownInflationDestination")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *AdminKnownInflationDestination) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AdminKnownInflationDestination) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no admin_known_inflation_destinations provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
	o.UpdatedAt = currTime

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(adminKnownInflationDestinationColumnsWithDefault, o)

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

	adminKnownInflationDestinationUpsertCacheMut.RLock()
	cache, cached := adminKnownInflationDestinationUpsertCache[key]
	adminKnownInflationDestinationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			adminKnownInflationDestinationColumns,
			adminKnownInflationDestinationColumnsWithDefault,
			adminKnownInflationDestinationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			adminKnownInflationDestinationColumns,
			adminKnownInflationDestinationPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("models: unable to upsert admin_known_inflation_destinations, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(adminKnownInflationDestinationPrimaryKeyColumns))
			copy(conflict, adminKnownInflationDestinationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"admin_known_inflation_destinations\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(adminKnownInflationDestinationType, adminKnownInflationDestinationMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert admin_known_inflation_destinations")
	}

	if !cached {
		adminKnownInflationDestinationUpsertCacheMut.Lock()
		adminKnownInflationDestinationUpsertCache[key] = cache
		adminKnownInflationDestinationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single AdminKnownInflationDestination record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *AdminKnownInflationDestination) DeleteG() (int64, error) {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single AdminKnownInflationDestination record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AdminKnownInflationDestination) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no AdminKnownInflationDestination provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), adminKnownInflationDestinationPrimaryKeyMapping)
	sql := "DELETE FROM \"admin_known_inflation_destinations\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from admin_known_inflation_destinations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for admin_known_inflation_destinations")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q adminKnownInflationDestinationQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no adminKnownInflationDestinationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from admin_known_inflation_destinations")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for admin_known_inflation_destinations")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o AdminKnownInflationDestinationSlice) DeleteAllG() (int64, error) {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AdminKnownInflationDestinationSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no AdminKnownInflationDestination slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(adminKnownInflationDestinationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), adminKnownInflationDestinationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"admin_known_inflation_destinations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, adminKnownInflationDestinationPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from adminKnownInflationDestination slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for admin_known_inflation_destinations")
	}

	if len(adminKnownInflationDestinationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *AdminKnownInflationDestination) ReloadG() error {
	if o == nil {
		return errors.New("models: no AdminKnownInflationDestination provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AdminKnownInflationDestination) Reload(exec boil.Executor) error {
	ret, err := FindAdminKnownInflationDestination(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AdminKnownInflationDestinationSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty AdminKnownInflationDestinationSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AdminKnownInflationDestinationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AdminKnownInflationDestinationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), adminKnownInflationDestinationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"admin_known_inflation_destinations\".* FROM \"admin_known_inflation_destinations\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, adminKnownInflationDestinationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AdminKnownInflationDestinationSlice")
	}

	*o = slice

	return nil
}

// AdminKnownInflationDestinationExistsG checks if the AdminKnownInflationDestination row exists.
func AdminKnownInflationDestinationExistsG(iD int) (bool, error) {
	return AdminKnownInflationDestinationExists(boil.GetDB(), iD)
}

// AdminKnownInflationDestinationExists checks if the AdminKnownInflationDestination row exists.
func AdminKnownInflationDestinationExists(exec boil.Executor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"admin_known_inflation_destinations\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if admin_known_inflation_destinations exists")
	}

	return exists, nil
}
