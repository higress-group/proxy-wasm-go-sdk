// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/stretchr/testify/require"
)

func TestPluginContext_OnTick(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart.
		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		require.Equal(t, tickMilliseconds, host.GetTickPeriod())

		for i := 1; i < 10; i++ {
			host.Tick() // call OnTick
			attrs := host.GetCalloutAttributesFromContext(proxytest.PluginContextID)
			// Verify DispatchHttpCall is called
			require.Equal(t, len(attrs), i)
			// Receive callout response.
			host.CallOnHttpCallResponse(attrs[0].CalloutID, nil, nil, nil)
			// Check Envoy logs.
			logs := host.GetInfoLogs()
			require.Contains(t, logs, fmt.Sprintf("called %d for contextID=%d", i, proxytest.PluginContextID))
		}
	})
}

func TestPluginContext_OnVMStart(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Call OnVMStart.
		require.Equal(t, types.OnPluginStartStatusOK, host.StartPlugin())
		require.Equal(t, tickMilliseconds, host.GetTickPeriod())
	})
}

// vmTest executes f twice, once with a types.VMContext that executes plugin code directly
// in the host, and again by executing the plugin code within the compiled main.wasm binary.
// Execution with main.wasm will be skipped if the file cannot be found.
func vmTest(t *testing.T, f func(*testing.T, types.VMContext)) {
	t.Helper()

	t.Run("go", func(t *testing.T) {
		f(t, &vmContext{})
	})

	t.Run("wasm", func(t *testing.T) {
		wasm, err := os.ReadFile("main.wasm")
		if err != nil {
			t.Skip("wasm not found")
		}
		v, err := proxytest.NewWasmVMContext(wasm)
		require.NoError(t, err)
		defer v.Close()
		f(t, v)
	})
}
