package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"math/rand"
	"time"
)

func main() {
	proxywasm.SetNewRootContext(newRootContext)
}

type rootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
}

func newRootContext(uint32) proxywasm.RootContext { return &rootContext{} }

func (ctx *rootContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
	rand.Seed(time.Now().UnixNano())

	proxywasm.LogInfo("proxy_on_vm_start from Go!")

	return types.OnVMStartStatusOK
}

// Override DefaultRootContext.
func (*rootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultHttpContext
	contextID uint32
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	err := proxywasm.SetHttpRequestHeader("test", "best")
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("request header --> %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("response header <-- %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
