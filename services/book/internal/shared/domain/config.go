package domain

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	TableName string
}

func NewConfiguration() (*Configuration, error) {
	viper.SetEnvPrefix("library")
	viper.AutomaticEnv()

	tableName := viper.GetString("table")
	if tableName == "" {
		// return nil, errors.New("table name not found")
		return &Configuration{"damascus-ebooks"}, nil
	}

	return &Configuration{tableName}, nil
}
