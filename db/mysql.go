package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var MySQL *sqlx.DB

func SetupMySQL() error {
	var err error
	if MySQL, err = sqlx.Open("mysql", viper.GetString("mysql.dsn")); err != nil {
		return errors.Wrapf(err, "[DB.MySQL] open mysql failed")
	}

	MySQL.SetMaxIdleConns(viper.GetInt("mysql.pool.idle"))
	MySQL.SetMaxOpenConns(viper.GetInt("mysql.pool.open"))

	return MySQL.Ping()
}
