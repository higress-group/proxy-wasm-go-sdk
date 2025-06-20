// Copyright 2024 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxytest

import (
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
)

// EmulatorOption is an option that can be passed to NewHostEmulator.
type EmulatorOption struct {
	pluginConfiguration []byte
	vmConfiguration     []byte
	context             interface{}
	properties          map[string][]byte
}

// NewEmulatorOption creates a new EmulatorOption.
func NewEmulatorOption() *EmulatorOption {
	return &EmulatorOption{context: &types.DefaultVMContext{}}
}

// WithVMContext sets the VMContext.
func (o *EmulatorOption) WithVMContext(context types.VMContext) *EmulatorOption {
	o.context = context
	return o
}

// WithPluginContext sets up the emulator to use the passed func to construct
// PluginContexts with an anonymous VMContext.
func (o *EmulatorOption) WithPluginContext(context types.PluginContextFactory) *EmulatorOption {
	o.context = context
	return o
}

// WithHttpContext sets up the emulator to use the passed func to construct
// HttpContexts an anonymous VMContext and PluginContext.
func (o *EmulatorOption) WithHttpContext(context types.HttpContextFactory) *EmulatorOption {
	o.context = context
	return o
}

// WithTcpContext sets up the emulator to use the passed func to construct
// TcpContexts an anonymous VMContext and PluginContext.
func (o *EmulatorOption) WithTcpContext(context types.TcpContextFactory) *EmulatorOption {
	o.context = context
	return o
}

// WithPluginConfiguration sets the plugin configuration.
func (o *EmulatorOption) WithPluginConfiguration(data []byte) *EmulatorOption {
	o.pluginConfiguration = data
	return o
}

// WithVMConfiguration sets the VM configuration.
func (o *EmulatorOption) WithVMConfiguration(data []byte) *EmulatorOption {
	o.vmConfiguration = data
	return o
}

// WithProperty sets a property. If the property already exists, it will be overwritten.
func (o *EmulatorOption) WithProperty(path []string, value []byte) *EmulatorOption {
	if o.properties == nil {
		o.properties = map[string][]byte{}
	}
	o.properties[string(internal.SerializePropertyPath(path))] = value
	return o
}
