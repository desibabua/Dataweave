package factory

import (
	"github.com/sirupsen/logrus"

	"github.com/ravik-karn/Dataweave/config"
	"github.com/ravik-karn/Dataweave/products"
)

type Factory struct {
	config *config.Config
	logger logrus.Logger
}

func New(config *config.Config, logger logrus.Logger) Factory {
	return Factory{
		config: config,
		logger: logger,
	}
}

func (f Factory) ProductFetcher() products.Fetcher {
	return products.New(f.config, f.logger)
}
