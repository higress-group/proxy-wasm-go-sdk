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

//go:build !wasm

package internal

import (
	"fmt"
	"strings"
	"testing"
)

// TestRegisterMockWasmHostFailsFastWhenTheLockIsHeld pins the behaviour that
// keeps a test binary from getting wedged.
//
// The lock that guards the current mock host is process-wide and is held for
// the whole lifetime of the host, so a host that was created and never
// released blocks every later RegisterMockWasmHost call - forever, and
// silently: the process sits there with no output, and not even the "go test"
// timeout panic can be printed, because the runtime is stuck waiting for the
// lock.
//
// So registering a host must fail **immediately** with a message that says
// what to do about it, and releasing a host must be safe to do twice.
func TestRegisterMockWasmHostFailsFastWhenTheLockIsHeld(t *testing.T) {
	release := RegisterMockWasmHost(DefaultProxyWAMSHost{})

	// A second registration must fail right away instead of waiting.
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("registering a second host must fail immediately, " +
					"otherwise the whole test binary wedges")
			}
			if msg := fmt.Sprint(r); !strings.Contains(msg, "never released") {
				t.Fatalf("the message must explain the cause, got: %s", msg)
			}
		}()
		RegisterMockWasmHost(DefaultProxyWAMSHost{})
	}()

	// Releasing gives the lock back.
	release()
	release2 := RegisterMockWasmHost(DefaultProxyWAMSHost{})

	// Releasing twice must be harmless: unlocking a mutex that is not held is
	// a fatal error, and "reset the host twice" is an easy mistake to make.
	release2()
	release2()
}
