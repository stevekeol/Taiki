package main

import (
	"Taiki/db"
	"Taiki/logger"
)

var log = logger.Log

// server provides a taiki server for handling communications to and from taiki peers.
type server struct {
	Listeners []string
	DB        db.KeyValueStore
	QuitChan  chan struct{}
}

// TODO: agentBlacklist, agentWhitelist []string,db database.DB, chainParams *chaincfg.Params
func newServer(listenAddrs []string, db db.KeyValueStore, interrupt <-chan struct{}) (*server, error) {
	return &server{
		// TODO
		Listeners: listenAddrs,
		DB:        db,
		QuitChan:  interrupt,
	}, nil
}

func (s *server) Start() error {

	log.Info("[3/3]server started...")
	return nil
}

func (s *server) Stop() error {
	log.Info("[1/2]server stopped.")
	return nil
}
