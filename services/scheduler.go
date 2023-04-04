package services

import (
	"github.com/btc-scan/config"
	"github.com/btc-scan/types"
	"os"
)

type ServiceScheduler struct {
	conf *config.Config

	db types.IDB

	services []types.IAsyncService

	closeCh <-chan os.Signal
}

func NewServiceScheduler(conf *config.Config, db types.IDB, closeCh <-chan os.Signal) (t *ServiceScheduler, err error) {
	t = &ServiceScheduler{
		conf:     conf,
		closeCh:  closeCh,
		db:       db,
		services: make([]types.IAsyncService, 0),
	}

	return
}

func (t *ServiceScheduler) Start() {
	//create collect service
	collectService := NewCollectService(t.db, t.conf)

	collectService.Run()
}
