package core

import (
	"log"

	"github.com/nextdns/nextdns/host"
)

const LocalDNSAddress = "127.0.0.1"

var (
	setSystemDNSFunc   = host.SetDNS
	resetSystemDNSFunc = host.ResetDNS
)

// ActivateSystemDNS points the host resolver at the local uBlockDNS proxy.
func ActivateSystemDNS() error {
	return setSystemDNSFunc(LocalDNSAddress)
}

// ActivateSystemDNSBestEffort logs and continues when activation fails.
func ActivateSystemDNSBestEffort() {
	if err := ActivateSystemDNS(); err != nil {
		log.Printf("Warning: failed to activate system DNS: %v", err)
	}
}

// DeactivateSystemDNS restores the host resolver to its default configuration.
func DeactivateSystemDNS() error {
	return resetSystemDNSFunc()
}

// DeactivateSystemDNSBestEffort logs and continues when deactivation fails.
func DeactivateSystemDNSBestEffort() {
	if err := DeactivateSystemDNS(); err != nil {
		log.Printf("Warning: failed to deactivate system DNS: %v", err)
	}
}

// SwapSystemDNSFuncs overrides DNS set/reset hooks and returns a restore func.
func SwapSystemDNSFuncs(set func(string) error, reset func() error) (restore func()) {
	oldSet, oldReset := setSystemDNSFunc, resetSystemDNSFunc
	if set != nil {
		setSystemDNSFunc = set
	}
	if reset != nil {
		resetSystemDNSFunc = reset
	}
	return func() {
		setSystemDNSFunc = oldSet
		resetSystemDNSFunc = oldReset
	}
}
