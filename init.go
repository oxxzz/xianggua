package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("cfg.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("load config file: " + viper.ConfigFileUsed())

	// if err := db.SetupMySQL(); err != nil {
	// 	panic(err)
	// }

	logrus.Debugf("[DB.MySQL] setup mysql success")
}
