package sqldb

import (
	"database/sql"

	// importing mysql driver
	_ "github.com/go-sql-driver/mysql"
	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
)

const (
	// maxIdleConns is maximum idle db connections
	maxIdleConns = 50
	// maxOpenConns is maximum open db connections
	maxOpenConns = 100
	// dbDriver is db driver
	dbDriver = "mysql"
	// dbEngine is db engine
	dbEngine = "InnoDB"
	// dbEncoding is db encoding
	dbEncoding = "UTF8"
)

// DB is db management structure
type DB struct {
	dbMap *gorp.DbMap
}

// DBInterface is db management interface
type DBInterface interface {
	AddTableWithName(i interface{}, name string) *gorp.TableMap
	Insert(trans *gorp.Transaction, list ...interface{}) error
	Update(trans *gorp.Transaction, list ...interface{}) (int64, error)
	Delete(trans *gorp.Transaction, list ...interface{}) (int64, error)
	Get(trans *gorp.Transaction, i interface{}, keys ...interface{}) (interface{}, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectInt(query string, args ...interface{}) (int64, error)
	SelectStr(query string, args ...interface{}) (string, error)
	SelectOne(trans *gorp.Transaction, holder interface{}, query string, args ...interface{}) error
	Exec(trans *gorp.Transaction, query string, args ...interface{}) (sql.Result, error)
	Begin() (*gorp.Transaction, error)
	Rollback(trans *gorp.Transaction) error
	Commit(trans *gorp.Transaction) error
	GetDBMap() *gorp.DbMap
}

// NewDB creates new db management structure
func NewDB(connString string, db *sql.DB, isProduction bool,
	gorpLogger gorp.GorpLogger) (DBInterface, error) {
	var err error

	if db == nil {
		db, err = sql.Open(dbDriver, connString)
		if err != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":        err,
				"db":           connString,
				"isProduction": isProduction,
			}).Error("Error connecting database")
			return nil, err
		}

		if isProduction {
			db.SetMaxIdleConns(maxIdleConns)
			db.SetMaxOpenConns(maxOpenConns)
		}
	}

	err = db.Ping()
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":        err,
			"db":           connString,
			"isProduction": isProduction,
		}).Error("Error pinging database")
		return nil, err
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"connString":   connString,
		"isProduction": isProduction,
	}).Debug("DB connection successfully established")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: dbEngine, Encoding: dbEncoding}}
	if !isProduction {
		dbMap.TraceOn("[gorp]", gorpLogger)
	}

	return &DB{dbMap: dbMap}, nil
}

// AddTableWithName adds table with name to db map
func (db *DB) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	return db.dbMap.AddTableWithName(i, name)
}

// Insert inserts data to the db table
func (db *DB) Insert(trans *gorp.Transaction, list ...interface{}) error {
	var err error

	if trans != nil {
		err = trans.Insert(list...)
	} else {
		err = db.dbMap.Insert(list...)
	}

	return err
}

// Update updates data in the db table
func (db *DB) Update(trans *gorp.Transaction, list ...interface{}) (int64, error) {
	var err error
	var count int64

	if trans != nil {
		count, err = trans.Update(list...)
	} else {
		count, err = db.dbMap.Update(list...)
	}

	return count, err
}

// Delete deletes data from the db table
func (db *DB) Delete(trans *gorp.Transaction, list ...interface{}) (int64, error) {
	var err error
	var count int64

	if trans != nil {
		count, err = trans.Delete(list...)
	} else {
		count, err = db.dbMap.Delete(list...)
	}

	return count, err
}

// Get gets data from the db table
func (db *DB) Get(trans *gorp.Transaction, i interface{}, keys ...interface{}) (interface{}, error) {
	var err error
	var object interface{}

	if trans != nil {
		object, err = trans.Get(i, keys...)
	} else {
		object, err = db.dbMap.Get(i, keys...)
	}

	return object, err
}

// Select selects data from the db table
func (db *DB) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return db.dbMap.Select(i, query, args...)
}

// SelectInt selects int from the db table
func (db *DB) SelectInt(query string, args ...interface{}) (int64, error) {
	return db.dbMap.SelectInt(query, args...)
}

// SelectStr selects string from the db table
func (db *DB) SelectStr(query string, args ...interface{}) (string, error) {
	return db.dbMap.SelectStr(query, args...)
}

// SelectOne selects one row from the db table
func (db *DB) SelectOne(trans *gorp.Transaction, holder interface{}, query string, args ...interface{}) error {
	var err error

	if trans != nil {
		err = trans.SelectOne(holder, query, args...)
	} else {
		err = db.dbMap.SelectOne(holder, query, args...)
	}

	return err
}

// Exec executes statement
func (db *DB) Exec(trans *gorp.Transaction, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	if trans != nil {
		result, err = trans.Exec(query, args...)
	} else {
		result, err = db.dbMap.Exec(query, args...)
	}

	return result, err
}

// Begin begins transaction
func (db *DB) Begin() (*gorp.Transaction, error) {
	return db.dbMap.Begin()
}

// Rollback rollbacks transaction
func (db *DB) Rollback(trans *gorp.Transaction) error {
	return trans.Rollback()
}

// Commit commits transaction
func (db *DB) Commit(trans *gorp.Transaction) error {
	return trans.Commit()
}

// GetDBMap returns dbmap
func (db *DB) GetDBMap() *gorp.DbMap {
	return db.dbMap
}
