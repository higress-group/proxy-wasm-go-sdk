// The framework emulates the expected behavior of Envoyproxy, and you can test your extensions without running Envoy and with
// the standard Go CLI. To run tests, simply run
// go test ./...

package main

import (
	"os"
	"testing"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/proxytest"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpAuthRandom_OnHttpRequestHeaders(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// Initialize context.
		contextID := host.InitializeHttpContext()

		// Call OnHttpRequestHeaders.
		action := host.CallOnRequestHeaders(contextID,
			[][2]string{{"key", "value"}}, false)
		require.Equal(t, types.ActionPause, action)

		// Verify DispatchHttpCall is called.
		attrs := host.GetCalloutAttributesFromContext(contextID)
		require.Equal(t, len(attrs), 1)
		require.Equal(t, "httpbin", attrs[0].Upstream)
		// Check if the current action is pause.
		require.Equal(t, types.ActionPause,
			host.GetCurrentHttpStreamAction(contextID))

		// Check Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, "http call dispatched to "+clusterName)
		require.Contains(t, logs, "request header: key: value")
	})
}

func TestHttpAuthRandom_OnHttpCallResponse(t *testing.T) {
	vmTest(t, func(t *testing.T, vm types.VMContext) {
		opt := proxytest.NewEmulatorOption().WithVMContext(vm)
		host, reset := proxytest.NewHostEmulator(opt)
		defer reset()

		// http://httpbin.org/uuid
		headers := [][2]string{
			{"HTTP/1.1", "200 OK"}, {"Date:", "Thu, 17 Sep 2020 02:47:07 GMT"},
			{"Content-Type", "application/json"}, {"Content-Length", "53"},
			{"Connection", "keep-alive"}, {"Server", "gunicorn/19.9.0"},
			{"Access-Control-Allow-Origin", "*"}, {"Access-Control-Allow-Credentials", "true"},
		}

		// Access granted case -> Local response must not be sent.
		contextID := host.InitializeHttpContext()
		// Call OnHttpRequestHeaders.
		action := host.CallOnRequestHeaders(contextID, nil,
			false)
		require.Equal(t, types.ActionPause, action)
		// Verify DispatchHttpCall is called.
		attrs := host.GetCalloutAttributesFromContext(contextID)
		require.Equal(t, len(attrs), 1)
		// Call OnHttpCallResponse.
		body := []byte(`{"uuid": "7b10a67a-1c67-4199-835b-cbefcd4a63d4"}`)
		host.CallOnHttpCallResponse(attrs[0].CalloutID, headers, nil, body)
		// Check local response.
		assert.Nil(t, host.GetSentLocalResponse(contextID))
		// CHeck Envoy logs.
		logs := host.GetInfoLogs()
		require.Contains(t, logs, "access granted")

		// Access denied case -> Local response must be sent.
		contextID = host.InitializeHttpContext()
		// Call OnHttpRequestHeaders.
		action = host.CallOnRequestHeaders(contextID, nil, false)
		require.Equal(t, types.ActionPause, action)
		// Verify DispatchHttpCall is called.
		attrs = host.GetCalloutAttributesFromContext(contextID)
		require.Equal(t, len(attrs), 1)
		// Call OnHttpCallResponse.
		body = []byte(`{"uuid": "aaaaaaaa-1c67-4199-835b-cbefcd4a63d4"}`)
		host.CallOnHttpCallResponse(attrs[0].CalloutID, headers, nil, body)
		// Check local response.
		localResponse := host.GetSentLocalResponse(contextID)
		assert.NotNil(t, localResponse)
		require.Equal(t, uint32(403), localResponse.StatusCode)
		require.Equal(t, []byte("access forbidden"), localResponse.Data)
		require.Len(t, localResponse.Headers, 1)
		require.Equal(t, "powered-by", localResponse.Headers[0][0])
		require.Equal(t, "proxy-wasm-go-sdk!!", localResponse.Headers[0][1])
		// Check Envoy logs.
		logs = host.GetInfoLogs()
		require.Contains(t, logs, "access forbidden")
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
