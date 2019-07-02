package products

import (
	"fmt"
	"github.com/ravik-karn/Dataweave/constants"
	"github.com/ravik-karn/Dataweave/utils"
	"github.com/sirupsen/logrus"

	"github.com/ravik-karn/Dataweave/config"
	"github.com/ravik-karn/Dataweave/repository"
)

type Fetcher interface {
	Fetch() (map[string][]interface{}, error)
}

type DbFetcher struct {
	config *config.Config
	logger logrus.Logger
}

func New(config *config.Config, logger logrus.Logger) Fetcher {
	return &DbFetcher{
		config: config,
		logger: logger,
	}
}

func (fetcher *DbFetcher) Fetch() (map[string][]interface{}, error) {
	logger := fetcher.logger
	repo, err := repository.New(fetcher.config.ConnectionString, fetcher.logger)
	if err != nil {
		logger.Errorf("db conn: %s", err)
		return nil, err
	}
	query := fmt.Sprintf("SELECT * from %s;", constants.TableName)
	rows, err := repo.Query(query)
	if err != nil {
		logger.Errorf("db conn: %s", err)
		return nil, err
	}

	products, err := utils.MakeStructJSON(*rows)
	return products, nil
}
