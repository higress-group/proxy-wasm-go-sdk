// Copyright 2020-2024 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
)

const tickMilliseconds uint32 = 1000

func main() {}
func init() {
	proxywasm.SetVMContext(&vmContext{})
}

// vmContext implements types.VMContext.
type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext implements types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{
		contextID: contextID,
		postponed: make([]uint32, 0, 1024),
	}
}

// pluginContext implements types.PluginContext.
type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	contextID uint32
	postponed []uint32
}

// OnPluginStart implements types.PluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}

	return types.OnPluginStartStatusOK
}

// OnTick implements types.PluginContext.
func (ctx *pluginContext) OnTick() {
	for len(ctx.postponed) > 0 {
		httpCtxId, tail := ctx.postponed[0], ctx.postponed[1:]
		proxywasm.LogInfof("resume request with contextID=%v", httpCtxId)
		_ = proxywasm.SetEffectiveContext(httpCtxId)
		_ = proxywasm.ResumeHttpRequest()
		ctx.postponed = tail
	}
}

// NewHttpContext implements types.PluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID: contextID,
		pluginCtx: ctx,
	}
}

// httpContext implements types.HttpContext.
type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
	pluginCtx *pluginContext
}

// OnHttpRequestHeaders implements types.HttpContext.
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfof("postpone request with contextID=%d", ctx.contextID)
	ctx.pluginCtx.postponed = append(ctx.pluginCtx.postponed, ctx.contextID)
	return types.ActionPause
}
