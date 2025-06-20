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

func main() {}
func init() {
	proxywasm.SetHttpContext(func(contextID uint32) types.HttpContext {
		return &properties{contextID: contextID}
	})
}

// properties implements types.HttpContext.
type properties struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

var propertyPrefix = []string{
	"route_metadata",
	"filter_metadata",
	"envoy.filters.http.wasm",
}

// OnHttpRequestHeaders implements types.HttpContext.
func (ctx *properties) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	auth, err := proxywasm.GetProperty(append(propertyPrefix, "auth"))
	if err != nil {
		if err == types.ErrorStatusNotFound {
			proxywasm.LogInfo("no auth header for route")
			return types.ActionContinue
		}
		proxywasm.LogCriticalf("failed to read properties: %v", err)
	}
	proxywasm.LogInfof("auth header is \"%s\"", auth)

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	// Verify authentication header exists
	authHeader := false
	for _, h := range hs {
		if h[0] == string(auth) {
			authHeader = true
			break
		}
	}

	// Reject requests without authentication header
	if !authHeader {
		_ = proxywasm.SendHttpResponse(401, nil, nil, 16)
		return types.ActionPause
	}

	return types.ActionContinue
}

// OnHttpStreamDone implements types.HttpContext.
func (ctx *properties) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
