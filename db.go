package go_mysql_tools

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(config *DBConfig) (*sql.DB, error) {
	return sql.Open(config.Driver, config.GetConnString())
}

func UseDb(config *DBConfig, operate func(db *sql.DB) error) error {
	db, err := Connect(config)
	if err != nil {
		return err
	}
	defer db.Close()
	return operate(db)
}

func InTx(config *DBConfig, ops []func(tx *sql.Tx) error) error {
	return UseDb(config, func(db *sql.DB) error {
		return InTxWithDB(db, ops)
	})
}

func InTxWithDB(db *sql.DB, ops []func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return WithInTx(tx, ops)
}

func WithInTx(tx *sql.Tx, ops []func(tx *sql.Tx) error) error {
	for _, op := range ops {
		if err := op(tx); err != nil {
			return err
		}
	}
	return nil
}

func Query(config *DBConfig, query string, args []interface{}, handler func(rowNo int, rows *sql.Rows) error) error {
	return UseDb(config, func(db *sql.DB) error {
		return QueryDb(db, query, args, handler)
	})
}

func QueryDb(db *sql.DB, query string, args []interface{}, handler func(rowNo int, rows *sql.Rows) error) error {
	rows, err := db.Query(query, args...)
	return iterateRows(rows, err, handler)

}

func QueryTx(tx *sql.DB, query string, args[] interface{}, handler func(rowNo int, rows *sql.Rows) error) error {
	rows, err := tx.Query(query, args...)
	return iterateRows(rows, err, handler)
}

func iterateRows(rows *sql.Rows, err error, handler func(rowNo int, rows *sql.Rows) error) error{
	if err != nil {
		return err
	}
	defer rows.Close()
	var current = 0
	for rows.Next() {
		err = handler(current, rows)
		if err != nil {
			return err
		}
		current += 1
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}