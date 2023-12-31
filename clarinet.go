package goclarinet

import (
	"errors"
	"fmt"

	"github.com/arobie1992/go-clarinet/config"
	"github.com/arobie1992/go-clarinet/control"
	"github.com/arobie1992/go-clarinet/cryptography"
	"github.com/arobie1992/go-clarinet/log"
	"github.com/arobie1992/go-clarinet/p2p"
	"github.com/arobie1992/go-clarinet/repository"
	"github.com/arobie1992/go-clarinet/reputation"
)

func Start(configPath string) error {
	log.InitLogger()

	config, err := config.LoadConfig(configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to load configuration: %s", err))
	}

	if config.Libp2p.CertPath != "" {
		if err := cryptography.LoadPrivKey(config.Libp2p.CertPath); err != nil {
			return errors.New(fmt.Sprintf("Failed to initialize private keys: %s", err))
		}	
	}

	if err := p2p.InitLibp2pNode(config); err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize libp2p node: %s", err))
	}
	log.Log().Infof("I am %s", p2p.GetFullAddr())

	if err := repository.InitDB(config, &p2p.Connection{}, &p2p.DataMessage{}, &reputation.ReputationInfo{}); err != nil {
		return errors.New(fmt.Sprintf("Failed to initialize database: %s", err))
	}

	// start a http handler so we have some endpoints to trigger behavior through for testing
	if err := control.StartAdminServer(config); err != nil {
		return errors.New(fmt.Sprintf("Failed to start server: %s", err))
	}
	return nil
}
