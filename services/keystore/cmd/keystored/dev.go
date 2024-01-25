//go:build !aws
// +build !aws

package main

import (
	"go/services/keystore"
	"go/support/env"
)

func getConfig() *keystore.Config {
	return &keystore.Config{
		DBURL:          env.String("KEYSTORE_DATABASE_URL", "postgres:///keystore?sslmode=disable"),
		MaxIdleDBConns: env.Int("DB_MAX_IDLE_CONNS", 5),
		MaxOpenDBConns: env.Int("DB_MAX_OPEN_CONNS", 5),
		AUTHURL:        env.String("KEYSTORE_AUTHFORWARDING_URL", ""),
		ListenerPort:   env.Int("KEYSTORE_LISTENER_PORT", 8000),
	}
}
