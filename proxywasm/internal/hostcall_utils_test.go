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

package internal

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestStringBytePtr(t *testing.T) {
	input := "hello"
	expected := []byte("hello")

	result := StringBytePtr(input)
	for i := range expected {
		require.Equal(t, expected[i], *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(result)) + uintptr(i))))
	}
}

func TestRawBytePtrToString(t *testing.T) {
	input := []byte("hello")
	inputPtr := &input[0]
	size := 5
	expected := "hello"

	result := RawBytePtrToString(inputPtr, size)
	require.Equal(t, expected, result)
}

func TestRawBytePtrToByteSlice(t *testing.T) {
	input := []byte("hello")
	inputPtr := &input[0]
	size := 5
	expected := []byte("hello")

	result := RawBytePtrToByteSlice(inputPtr, size)
	require.Equal(t, expected, result)
}
