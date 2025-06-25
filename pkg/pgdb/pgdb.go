package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"
)

// Config for DB.
type Config struct {
	URL                   string        `json:"url" yaml:"ulr" toml:"url" mapstructure:"url"`                                                                                 //nolint
	DialTimeout           time.Duration `json:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout" mapstructure:"dial_timeout"`                                             //nolint
	IdleTimeout           time.Duration `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" mapstructure:"idle_timeout"`                                             //nolint
	ReadTimeout           time.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" mapstructure:"read_timeout"`                                             //nolint
	WriteTimeout          time.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" mapstructure:"write_timeout"`                                         //nolint
	RetryStatementTimeout bool          `json:"retry_statement_timeout" yaml:"retry_statement_timeout" toml:"retry_statement_timeout" mapstructure:"retry_statement_timeout"` //nolint
	MaxRetries            int           `json:"max_retries" yaml:"max_retries" toml:"max_retries" mapstructure:"max_retries"`                                                 //nolint
	MaxRetryBackoff       time.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff" toml:"max_retry_backoff" mapstructure:"max_retry_backoff"`                         //nolint
	PoolSize              int           `json:"pool_size" yaml:"pool_size" toml:"pool_size" mapstructure:"pool_size"`                                                         //nolint
	PoolTimeout           time.Duration `json:"pool_timeout" yaml:"pool_timeout" toml:"pool_timeout" mapstructure:"pool_timeout"`                                             //nolint
}

// NewConfig returns new default configurations.
func NewConfig() (conf *Config) {
	conf = new(Config)
	// no default DB URL
	// all other are zero by default: i.e. use defaults
	return
}

// Options for underlying DB ORM by this configurations.
func (c *Config) options() (opts *pg.Options, err error) {
	if opts, err = pg.ParseURL(c.URL); err != nil {
		return
	}
	opts.DialTimeout = c.DialTimeout
	opts.IdleTimeout = c.IdleTimeout
	opts.ReadTimeout = c.ReadTimeout
	opts.WriteTimeout = c.WriteTimeout
	opts.RetryStatementTimeout = c.RetryStatementTimeout
	opts.MaxRetries = c.MaxRetries
	opts.MaxRetryBackoff = c.MaxRetryBackoff
	opts.PoolSize = c.PoolSize
	opts.PoolTimeout = c.PoolTimeout

	return
}

// The DB represents DB.
type DB struct {
	*pg.DB // underlying go-pg DB instance, embed
}

// New DB with given configurations and log.
func New(conf *Config) (db *DB, err error) {

	var opts *pg.Options
	if opts, err = conf.options(); err != nil {
		return
	}

	db = new(DB)
	db.DB = pg.Connect(opts)

	db.DB.AddQueryHook(db)

	// check out DB server and fail if not responding
	if err = db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return
}

// Ping DB server.
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.Ping(ctx)
}

// BeforeQuery hook.
func (db *DB) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (
	context.Context, error) {

	return ctx, nil // no op.
}

// AfterQuery hook.
func (db *DB) AfterQuery(ctx context.Context, qe *pg.QueryEvent) (err error) {
	// var fq []byte
	q, err := qe.FormattedQuery()
	if err != nil {
		return
	}
	// fmt.Println(fq)
	// db.log.Debug(ctx, "query", zap.Duration("after", time.Since(qe.StartTime)),
	// 	zap.String("sql", string(fq)))
	fmt.Println("SQL Query: ", string(q))

	return
}

// Close DB.
func (db *DB) Close() (err error) {
	return db.DB.Close()
}

// NotFoundError returns given notFoundError if the given err is
// sql.ErrNoRows. Otherwise, it return the err. Use it in queries
// with corresponding 'not found error' from your model.
//
//	if user, err = mydb.User(ctx, query); err != nil {
//	    return nil, pgdb.NotFoundError(err, model.ErrUserNotFound)
//	}
//	return
//
// Or any other case.
//
// Don't propagate the database/slq to a controller.
func NotFoundError(err, notFoundError error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pg.ErrNoRows) {
		return notFoundError
	}
	return err
}

// IsNotFoundError returns true if given error represents no SQL result.
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, pg.ErrNoRows)
}

// PostgreSQL related constants
const (
	UniqueViolation     = "23505" // unique violation
	ForeignKeyViolation = "23503" // foreign key violation
)

// IsUniqueViolation returns true if given error represents unique
// constraint violation with given constrain name.
func IsUniqueViolation(err error, name string) bool {
	var pgErr, ok = err.(pg.Error)
	return ok &&
		pgErr.Field('C') == UniqueViolation &&
		pgErr.Field('n') == name
}

// IsForeignKeyViolation returns true if given error represents foreign
// key constraint violation with given constrain name.
func IsForeignKeyViolation(err error, name string) bool {
	var pgErr, ok = err.(pg.Error)
	return ok &&
		pgErr.Field('C') == ForeignKeyViolation &&
		pgErr.Field('n') == name
}

// DuplicateError returns the duplicateError if given error is unique violation
// for given constraint.
func DuplicateError(err error, constraint string, duplicateError error) error {
	if err == nil {
		return nil
	}

	if IsUniqueViolation(err, constraint) {
		return duplicateError
	}

	return err
}

// ForeignKeyError returns the foreignKeyError if given error is foreign key
// violation for given constraint.
func ForeignKeyError(err error, constraint string,
	foreignKeyError error) error {

	if err == nil {
		return nil
	}

	if IsForeignKeyViolation(err, constraint) {
		return foreignKeyError
	}

	return err
}

// DuplicateCase is  a duplicate case
type DuplicateCase struct {
	Constraint string
	Error      error
}

// DuplicateErrors used to switch multiple duplicate errors at once
func DuplicateErrors(err error, cases ...DuplicateCase) error {
	if err == nil {
		return nil
	}

	var pgErr, ok = err.(pg.Error)
	if !ok {
		return err
	}
	if pgErr.Field('C') != UniqueViolation {
		return err
	}

	var constraint = pgErr.Field('n')
	for _, cs := range cases {
		if cs.Constraint == constraint {
			return cs.Error // match
		}
	}

	return err // not match
}
