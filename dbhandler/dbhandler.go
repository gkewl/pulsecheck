package dbhandler

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	. "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/logger"
)

func CreateConnection() (*sqlx.DB, error) {

	dbConfig := &Config{
		User:      "pulse",              //"root",
		Passwd:    "pulsecheck",         //"divith",
		Addr:      "23.236.60.151:3306", //"10.128.0.5:3306", //"localhost:3306",
		Net:       "tcp",
		DBName:    "pulsecheck",
		Collation: "utf8_unicode_ci",
		Loc:       time.UTC,
		ParseTime: true,
	}
	db, err := sqlx.Connect("mysql", dbConfig.FormatDSN())

	if err != nil {
		fmt.Println(err)
		fmt.Println(dbConfig.FormatDSN())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the %s database using: %s\n", config.GetEnv(config.PULSE_DB_TYPE), err)
		return nil, err
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(128)
	d := time.Duration(1800) * time.Second
	db.SetConnMaxLifetime(d)
	return db, nil
}

// CreateTx starts a db connection
func CreateTx(ctx *common.AppContext) (*sqlx.Tx, error) {
	//create transaction object
	ctx.Db.Stats()

	clientctx := context.Background()

	timeout, _ := time.ParseDuration("3600s")

	clientctx, _ = context.WithTimeout(clientctx, timeout)

	tx, err := ctx.Db.BeginTxx(clientctx, nil)
	if err != nil {
		newErr := eh.NewError(eh.ErrDBCreatingTransactions, ""+err.Error())
		return nil, newErr
	}
	return tx, nil
}

// UpdateCheckRowVersion executes an update checking for rowversion. It treats
// zero rows updated as a "deadlock" situation to force retry of whole transaction
func UpdateCheckRowVersion(reqCtx common.RequestContext, updateStmt string, params map[string]interface{}, tablename string) (err error) {
	pkName := "id"

	pkColumnName, hasPkColumnName := params["pkcolumnname"]
	if hasPkColumnName {
		ok := false
		if pkName, ok = pkColumnName.(string); !ok {
			return eh.NewError(eh.ErrInternalAppError, "Primary column name must be a string")
		}
	}

	pk, hasPk := params[pkName]
	rowVersion, hasVersion := params["rowversion"]
	if !hasPk || !hasVersion {
		return eh.NewError(eh.ErrInternalAppError, "Row version update had no primary key %v or version %v: %s",
			hasPk, hasVersion, updateStmt)
	}

	stmt := updateStmt
	if !strings.Contains(stmt, ":rowversion") {
		stmt = stmt + " and rowversion = :rowversion"
	}

	result, err := reqCtx.Tx().NamedExec(stmt, params)
	if err == nil {
		affected, _ := result.RowsAffected()
		if affected == 0 {
			qry := fmt.Sprintf("select rowversion from `%s` where `%s` = ?", tablename, pkName)
			var curVersion int
			err = reqCtx.Tx().Get(&curVersion, qry, pk)
			if err != nil {
				return err
			}
			if curVersion != rowVersion {
				logger.LogInfo(fmt.Sprintf("Detected rowversion out of date table %s id %d in memory %d in DB %d Stack %s",
					tablename, pk, rowVersion, curVersion, debug.Stack()), reqCtx.Xid())
				err = errors.New(eh.ErrDeadlockText)
			}
		}
	}
	return err
}
