package service

import (
	"fmt"
	"log"

	"github.com/nextdns/nextdns/host/service"
	"github.com/ugzv/ublockdnsclient/internal/core"
	"github.com/ugzv/ublockdnsclient/internal/state"
)

var (
	activateSystemDNSFunc         = core.ActivateSystemDNS
	deactivateSystemDNSFunc       = core.DeactivateSystemDNSBestEffort
	deactivateSystemDNSStrictFunc = core.DeactivateSystemDNS
)

// Uninstall removes the system service and restores DNS.
func Uninstall() error {
	svc, err := baseService()
	if err != nil {
		return err
	}

	_ = stopServiceAndRestoreDNS(svc)

	if err := svc.Uninstall(); err != nil {
		return fmt.Errorf("uninstall service: %w", err)
	}
	_ = state.ClearPersistedTokens()
	_ = state.ClearInstallState()

	return nil
}

// ServiceStart starts the service when needed and strictly re-applies system DNS.
// Start is best-effort because launchd/systemd may report "already running" while
// DNS still needs repair after an OS update or network reset.
func ServiceStart() error {
	svc, err := baseService()
	if err != nil {
		return err
	}
	_ = svc.Start()

	st, err := svc.Status()
	if err != nil {
		return fmt.Errorf("service status: %w", err)
	}
	if st != service.StatusRunning {
		return fmt.Errorf("service is not running")
	}
	if err := activateSystemDNSFunc(); err != nil {
		return fmt.Errorf("activate system DNS: %w", err)
	}
	return nil
}

func ServiceStop() error {
	svc, err := baseService()
	if err != nil {
		return err
	}
	return stopServiceAndRestoreDNS(svc)
}

// stopServiceAndRestoreDNS stops the service and always resets system DNS.
// Legacy daemons do not deactivate on stop; the CLI must restore DNS even when
// manageSystemDNS is absent or stop fails partway through.
func stopServiceAndRestoreDNS(svc service.Service) error {
	err := stopService(svc)
	deactivateSystemDNSFunc()
	return err
}

func stopService(svc service.Service) error {
	if err := svc.Stop(); err != nil {
		log.Printf("Warning: failed to stop service: %v", err)
		return err
	}
	return nil
}
