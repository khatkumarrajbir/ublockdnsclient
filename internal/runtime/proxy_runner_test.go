package runtime

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/nextdns/nextdns/proxy"
)

func TestProxyRunnerStopWithoutStart(t *testing.T) {
	t.Parallel()

	runner := &proxyRunner{}
	if err := runner.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
}

func TestProxyRunnerStartReturnsImmediateError(t *testing.T) {
	t.Parallel()

	runner := &proxyRunner{
		proxy: proxy.Proxy{
			Addrs: []string{"not-a-valid-listen-address"},
		},
	}

	err := runner.Start()
	if err == nil {
		t.Fatal("expected Start() to fail for invalid listen address")
	}
	if errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProxyRunnerStartFailureSkipsOnReady(t *testing.T) {
	t.Parallel()

	var readyCalled atomic.Bool
	runner := &proxyRunner{
		proxy: proxy.Proxy{
			Addrs: []string{"not-a-valid-listen-address"},
		},
		onReady: []func(context.Context){
			func(context.Context) { readyCalled.Store(true) },
		},
	}

	if err := runner.Start(); err == nil {
		t.Fatal("expected Start() to fail for invalid listen address")
	}
	if readyCalled.Load() {
		t.Fatal("onReady hooks must not run when startup fails immediately")
	}
}
