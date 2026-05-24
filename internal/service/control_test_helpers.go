package service

import (
	"testing"

	hostservice "github.com/nextdns/nextdns/host/service"
)

type controlTestEnv struct {
	baseService         func() (hostservice.Service, error)
	activateSystemDNS   func() error
	deactivateSystemDNS func()
}

func withControlTestEnv(t *testing.T, env controlTestEnv) {
	t.Helper()

	oldBase := baseServiceFunc
	oldActivate := activateSystemDNSFunc
	oldDeactivate := deactivateSystemDNSFunc

	t.Cleanup(func() {
		baseServiceFunc = oldBase
		activateSystemDNSFunc = oldActivate
		deactivateSystemDNSFunc = oldDeactivate
	})

	if env.baseService != nil {
		baseServiceFunc = env.baseService
	}
	if env.activateSystemDNS != nil {
		activateSystemDNSFunc = env.activateSystemDNS
	}
	if env.deactivateSystemDNS != nil {
		deactivateSystemDNSFunc = env.deactivateSystemDNS
	}
}
