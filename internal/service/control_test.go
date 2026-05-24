package service

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/nextdns/nextdns/host/service"
)

type fakeService struct {
	startErr  error
	stopErr   error
	status    service.Status
	statusErr error
}

func (f fakeService) Install() error   { return nil }
func (f fakeService) Uninstall() error { return nil }
func (f fakeService) Restart() error   { return nil }
func (f fakeService) SaveConfig(map[string]service.ConfigEntry) error {
	return nil
}
func (f fakeService) LoadConfig(map[string]service.ConfigEntry) error { return nil }
func (f fakeService) Start() error                                    { return f.startErr }
func (f fakeService) Stop() error                                     { return f.stopErr }
func (f fakeService) Status() (service.Status, error)                 { return f.status, f.statusErr }

func TestServiceStartActivatesDNSWhenRunning(t *testing.T) {
	var activated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		activateSystemDNS: func() error {
			activated.Store(true)
			return nil
		},
		baseService: func() (service.Service, error) {
			return fakeService{status: service.StatusRunning}, nil
		},
	})

	if err := ServiceStart(); err != nil {
		t.Fatalf("ServiceStart() error = %v", err)
	}
	if !activated.Load() {
		t.Fatal("expected system DNS activation")
	}
}

func TestServiceStartReturnsErrorWhenNotRunning(t *testing.T) {
	withControlTestEnv(t, controlTestEnv{
		activateSystemDNS: func() error {
			t.Fatal("DNS activation should not run when service is not running")
			return nil
		},
		baseService: func() (service.Service, error) {
			return fakeService{
				startErr: errors.New("start failed"),
				status:   service.StatusStopped,
			}, nil
		},
	})

	err := ServiceStart()
	if err == nil || !strings.Contains(err.Error(), "service is not running") {
		t.Fatalf("ServiceStart() error = %v, want service is not running", err)
	}
}

func TestServiceStartReturnsActivationError(t *testing.T) {
	want := errors.New("activate failed")
	withControlTestEnv(t, controlTestEnv{
		activateSystemDNS: func() error { return want },
		baseService: func() (service.Service, error) {
			return fakeService{status: service.StatusRunning}, nil
		},
	})

	err := ServiceStart()
	if !errors.Is(err, want) {
		t.Fatalf("ServiceStart() error = %v, want %v", err, want)
	}
}

func TestServiceStartRepairsDNSWhenAlreadyRunning(t *testing.T) {
	var activated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		activateSystemDNS: func() error {
			activated.Store(true)
			return nil
		},
		baseService: func() (service.Service, error) {
			return fakeService{
				startErr: errors.New("already loaded"),
				status:   service.StatusRunning,
			}, nil
		},
	})

	if err := ServiceStart(); err != nil {
		t.Fatalf("ServiceStart() error = %v", err)
	}
	if !activated.Load() {
		t.Fatal("expected DNS repair when service is already running")
	}
}

func TestUninstallDeactivatesDNSWhenStopFails(t *testing.T) {
	var deactivated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		deactivateSystemDNS: func() {
			deactivated.Store(true)
		},
		baseService: func() (service.Service, error) {
			return fakeService{stopErr: errors.New("stop failed")}, nil
		},
	})

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if !deactivated.Load() {
		t.Fatal("expected DNS deactivation fallback when stop fails")
	}
}

func TestServiceStopDeactivatesDNSWhenStopFails(t *testing.T) {
	var deactivated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		deactivateSystemDNS: func() {
			deactivated.Store(true)
		},
		baseService: func() (service.Service, error) {
			return fakeService{stopErr: errors.New("stop failed")}, nil
		},
	})

	err := ServiceStop()
	if err == nil {
		t.Fatal("ServiceStop() expected stop error")
	}
	if !deactivated.Load() {
		t.Fatal("expected DNS deactivation when stop fails")
	}
}

func TestServiceStopDeactivatesDNSWhenStopSucceeds(t *testing.T) {
	var deactivated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		deactivateSystemDNS: func() {
			deactivated.Store(true)
		},
		baseService: func() (service.Service, error) {
			return fakeService{}, nil
		},
	})

	if err := ServiceStop(); err != nil {
		t.Fatalf("ServiceStop() error = %v", err)
	}
	if !deactivated.Load() {
		t.Fatal("expected DNS deactivation after stop even when stop succeeds")
	}
}

func TestUninstallDeactivatesDNSWhenStopSucceeds(t *testing.T) {
	var deactivated atomic.Bool
	withControlTestEnv(t, controlTestEnv{
		deactivateSystemDNS: func() {
			deactivated.Store(true)
		},
		baseService: func() (service.Service, error) {
			return fakeService{}, nil
		},
	})

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() error = %v", err)
	}
	if !deactivated.Load() {
		t.Fatal("expected DNS deactivation after uninstall even when stop succeeds")
	}
}
