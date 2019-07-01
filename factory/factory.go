package factory

import "../config"

type Factory interface {

}

type DataWeaveFactory struct {

}

func New(conf *config.Config) Factory{
	return &DataWeaveFactory {}
}
