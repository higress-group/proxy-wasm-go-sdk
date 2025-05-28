// Copyright 2021 Tetrate
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

//go:build go1.20

package internal

import (
	"unsafe"
)

// RawBytePtrToString uses unsafe to convert a raw byte pointer to a string.
func RawBytePtrToString(raw *byte, size int) string {
	return unsafe.String(raw, size)
}

// RawBytePtrToByteSlice uses unsafe to convert a raw byte pointer to a byte slice.
func RawBytePtrToByteSlice(raw *byte, size int) []byte {
	return unsafe.Slice(raw, size)
}

// StringBytePtr uses unsafe to convert a string to a raw byte pointer.
func StringBytePtr(msg string) *byte {
	return unsafe.StringData(msg)
}
