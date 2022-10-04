package main

import (
	"Taiki/logger"
)

var log = logger.Log

// server provides a taiki server for handling communications to and from taiki peers.
type server struct {
	Listeners []string
}

// TODO: agentBlacklist, agentWhitelist []string,db database.DB, chainParams *chaincfg.Params
func newServer(listenAddrs []string, interrupt <-chan struct{}) (*server, error) {
	return &server{
		// TODO
		Listeners: listenAddrs,
	}, nil
}

func (s *server) Start() error {
	log.Debug("server started")
	return nil
}

func (s *server) Stop() error {
	log.Debug("server stopped")
	return nil
}
