package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/sirupsen/logrus"
)

type conf struct {
	AppPort    string `json:appPort`
	DbServer   string `json:dbServer`
	DbPort     string `json:dbPort`
	DbName     string `json:dbName`
	DbUsername string `json:dbUsername`
	DbPassword string `json:dbPassword`
	DbTimeout  string `json:dbTimeout`
}

type Config struct {
	ConnectionString string
	AppPort          string
}

func New(filePath string) (*Config, error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file '%s': %s", filePath, err)
	}

	conf := &conf{}
	err = json.Unmarshal(contents, &conf)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON configuration: %s", err)
	}
	connectionString, err := buildConnectionString(conf)
	if err != nil {
		logrus.Errorf("build connection string: %s", err)
	}
	return &Config{
		ConnectionString: connectionString,
		AppPort:          conf.AppPort,
	}, nil
}

func buildConnectionString(conf *conf) (string, error) {
	connectionString := fmt.Sprintf("mysql://%s:%s/%s?sslmode=disable&connect_timeout=%s",
		conf.DbServer,
		conf.DbPort,
		conf.DbName,
		conf.DbTimeout,
	)

	u, err := url.Parse(connectionString)
	if err != nil {
		return "", err
	}

	u.User = url.UserPassword(conf.DbUsername, conf.DbPassword)
	return u.String(), nil
}
