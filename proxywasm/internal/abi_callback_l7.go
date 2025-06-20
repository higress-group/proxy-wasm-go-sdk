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

//go:wasmexport proxy_on_request_headers
func proxyOnRequestHeaders(contextID uint32, numHeaders int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnRequestHeaders", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_headers")
	}

	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestHeaders(int(numHeaders), endOfStream)
}

//go:wasmexport proxy_on_request_body
func proxyOnRequestBody(contextID uint32, bodySize int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnRequestBody", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_body")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestBody(int(bodySize), endOfStream)
}

//go:wasmexport proxy_on_request_trailers
func proxyOnRequestTrailers(contextID uint32, numTrailers int32) types.Action {
	if recordTiming {
		defer logTiming("proxyOnRequestTrailers", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context on proxy_on_request_trailers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpRequestTrailers(int(numTrailers))
}

//go:wasmexport proxy_on_response_headers
func proxyOnResponseHeaders(contextID uint32, numHeaders int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnResponseHeaders", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseHeaders(int(numHeaders), endOfStream)
}

//go:wasmexport proxy_on_response_body
func proxyOnResponseBody(contextID uint32, bodySize int32, endOfStream bool) types.Action {
	if recordTiming {
		defer logTiming("proxyOnResponseBody", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseBody(int(bodySize), endOfStream)
}

//go:wasmexport proxy_on_response_trailers
func proxyOnResponseTrailers(contextID uint32, numTrailers int32) types.Action {
	if recordTiming {
		defer logTiming("proxyOnResponseTrailers", time.Now())
	}
	ctx, ok := currentState.httpContexts[contextID]
	if !ok {
		panic("invalid context id on proxy_on_response_headers")
	}
	currentState.setActiveContextID(contextID)
	return ctx.OnHttpResponseTrailers(int(numTrailers))
}

//go:wasmexport proxy_on_http_call_response
func proxyOnHttpCallResponse(pluginContextID, calloutID uint32, numHeaders, bodySize, numTrailers int32) {
	if recordTiming {
		defer logTiming("proxyOnHttpCallResponse", time.Now())
	}
	root, ok := currentState.pluginContexts[pluginContextID]
	if !ok {
		panic("http_call_response on invalid plugin context")
	}

	cb := root.httpCallbacks[calloutID]
	if cb == nil {
		panic("invalid callout id")
	}

	ctxID := cb.callerContextID
	currentState.setActiveContextID(ctxID)
	delete(root.httpCallbacks, calloutID)

	// Check if the context is already deleted.
	// For example, if the connection expired before the call response arrival,
	// proxy_on_http_call_response is called AFTER ProxyOnDelete is called for the context id.
	// In that case, if the callback continues response or make local reply, then the subsequent
	// callbacks (for example OnHttpResponseHeaders) would follow and result in calling callback
	// for already-deleted context id. See https://github.com/higress-group/proxy-wasm-go-sdk/issues/261 for detail.
	if _, ok := currentState.contextIDToRootID[ctxID]; ok {
		ProxySetEffectiveContext(ctxID)
		cb.callback(int(numHeaders), int(bodySize), int(numTrailers))
	}
}

//go:wasmexport proxy_on_redis_call_response
func proxyOnRedisCallResponse(pluginContextID, calloutID uint32, status, responseSize int32) {
	root, ok := currentState.pluginContexts[pluginContextID]
	if !ok {
		panic("redis_call_response on invalid plugin context")
	}

	cb := root.redisCallbacks[calloutID]
	if cb == nil {
		panic("invalid callout id")
	}

	ctxID := cb.callerContextID
	currentState.setActiveContextID(ctxID)
	delete(root.redisCallbacks, calloutID)

	// Check if the context is already deleted.
	// For example, if the connection expired before the call response arrival,
	// proxy_on_http_call_response is called AFTER ProxyOnDelete is called for the context id.
	// In that case, if the callback continues response or make local reply, then the subsequent
	// callbacks (for example OnHttpResponseHeaders) would follow and result in calling callback
	// for already-deleted context id. See https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/261 for detail.
	if _, ok := currentState.contextIDToRootID[ctxID]; ok {
		ProxySetEffectiveContext(ctxID)
		cb.callback(int(status), int(responseSize))
	}
}
