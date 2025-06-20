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

package internal

import (
	"time"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
)

//go:wasmexport proxy_on_new_connection
func proxyOnNewConnection(contextID uint32) types.Action {
	if recordTiming {
		defer logTiming("proxyOnNewConnection", time.Now())
	}
	ctx, ok := currentState.tcpContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnNewConnection()
}

//go:wasmexport proxy_on_downstream_data
func proxyOnDownstreamData(contextID uint32, dataSize int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnDownstreamData", time.Now())
	}
	ctx, ok := currentState.tcpContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnDownstreamData(int(dataSize), endOfStream)
}

//go:wasmexport proxy_on_downstream_connection_close
func proxyOnDownstreamConnectionClose(contextID uint32, pType types.PeerType) {
	if recordTiming {
		defer logTiming("proxyOnDownstreamConnectionClose", time.Now())
	}
	ctx, ok := currentState.tcpContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	ctx.OnDownstreamClose(pType)
}

//go:wasmexport proxy_on_upstream_data
func proxyOnUpstreamData(contextID uint32, dataSize int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnUpstreamData", time.Now())
	}
	ctx, ok := currentState.tcpContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnUpstreamData(int(dataSize), endOfStream)
}

//go:wasmexport proxy_on_upstream_connection_close
func proxyOnUpstreamConnectionClose(contextID uint32, pType types.PeerType) {
	if recordTiming {
		defer logTiming("proxyOnUpstreamConnectionClose", time.Now())
	}
	ctx, ok := currentState.tcpContexts[contextID]
	if !ok {
		panic("invalid context")
	}
	currentState.setActiveContextID(contextID)
	ctx.OnUpstreamClose(pType)
}
