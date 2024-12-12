package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var MySQL *sqlx.DB

func SetupMySQL() error {
	var err error
	MySQL, err = sqlx.Open("mysql", viper.GetString("mysql.dsn"))
	MySQL.SetMaxIdleConns(viper.GetInt("mysql.pool.idle"))
	MySQL.SetMaxOpenConns(viper.GetInt("mysql.pool.open"))

	err = MySQL.Ping()
	return err
}
