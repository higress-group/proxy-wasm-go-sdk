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

package proxywasm

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/internal"
	"github.com/stretchr/testify/require"
)

type logHost struct {
	internal.DefaultProxyWAMSHost
	t           *testing.T
	expMessage  string
	expLogLevel internal.LogLevel
}

type logLevelHost struct {
	internal.DefaultProxyWAMSHost
	levelData [4]byte
	logCalls  int
}

func newLogLevelHost(level internal.LogLevel) *logLevelHost {
	host := &logLevelHost{}
	binary.LittleEndian.PutUint32(host.levelData[:], uint32(level))
	return host
}

func (l *logLevelHost) ProxyCallForeignFunction(funcNamePtr *byte, funcNameSize int32, _ *byte, _ int32, returnData unsafe.Pointer, returnSize *int32) internal.Status {
	if unsafe.String(funcNamePtr, funcNameSize) != "get_log_level" {
		return internal.StatusBadArgument
	}
	*(*unsafe.Pointer)(returnData) = unsafe.Pointer(&l.levelData[0])
	*returnSize = int32(len(l.levelData))
	return internal.StatusOK
}

func (l *logLevelHost) ProxyLog(_ internal.LogLevel, _ *byte, _ int32) internal.Status {
	l.logCalls++
	return internal.StatusOK
}

type countingStringer struct {
	calls *int
}

func (s countingStringer) String() string {
	*s.calls++
	return "large response body"
}

func (l logHost) ProxyLog(logLevel internal.LogLevel, messageData *byte, messageSize int32) internal.Status {
	actual := unsafe.String(messageData, messageSize)
	require.Equal(l.t, l.expMessage, actual)
	require.Equal(l.t, l.expLogLevel, logLevel)
	return internal.StatusOK
}

func TestHostCall_ForeignFunction(t *testing.T) {
	defer internal.RegisterMockWasmHost(internal.DefaultProxyWAMSHost{})()

	ret, err := CallForeignFunction("testFunc", []byte(""))
	require.NoError(t, err)
	require.Equal(t, []byte(nil), ret)
}

func TestHostCall_Logging(t *testing.T) {
	t.Run("trace", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace",
			expLogLevel:          internal.LogLevelTrace,
		})
		defer release()
		LogTrace("trace")
	})

	t.Run("tracef", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "trace: log",
			expLogLevel:          internal.LogLevelTrace,
		})
		defer release()
		LogTracef("trace: %s", "log")
	})

	t.Run("debug", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "abc",
			expLogLevel:          internal.LogLevelDebug,
		})
		defer release()
		LogDebug("abc")
	})

	t.Run("debugf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "debug: log",
			expLogLevel:          internal.LogLevelDebug,
		})
		defer release()
		LogDebugf("debug: %s", "log")
	})

	t.Run("info", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info",
			expLogLevel:          internal.LogLevelInfo,
		})
		defer release()
		LogInfo("info")
	})

	t.Run("infof", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "info: log: 10",
			expLogLevel:          internal.LogLevelInfo,
		})
		defer release()
		LogInfof("info: %s: %d", "log", 10)
	})

	t.Run("warn", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarn("warn")
	})

	t.Run("warnf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn: log: 10",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarnf("warn: %s: %d", "log", 10)
	})

	t.Run("error", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "error",
			expLogLevel:          internal.LogLevelError,
		})
		defer release()
		LogError("error")
	})

	t.Run("warnf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "warn: log: 10",
			expLogLevel:          internal.LogLevelWarn,
		})
		defer release()
		LogWarnf("warn: %s: %d", "log", 10)
	})

	t.Run("critical", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical error",
			expLogLevel:          internal.LogLevelCritical,
		})
		defer release()
		LogCritical("critical error")
	})

	t.Run("criticalf", func(t *testing.T) {
		release := internal.RegisterMockWasmHost(logHost{
			DefaultProxyWAMSHost: internal.DefaultProxyWAMSHost{},
			t:                    t,
			expMessage:           "critical: log: 10",
			expLogLevel:          internal.LogLevelCritical,
		})
		defer release()
		LogCriticalf("critical: %s: %d", "log", 10)
	})
}

func TestLogDebugfSkipsFormattingWhenDebugDisabled(t *testing.T) {
	host := newLogLevelHost(internal.LogLevelInfo)
	release := internal.RegisterMockWasmHost(host)
	defer release()

	formatCalls := 0
	LogDebugf("response body: %s", countingStringer{calls: &formatCalls})

	require.Equal(t, 0, formatCalls)
	require.Equal(t, 0, host.logCalls)
}

func TestLogDebugfFormatsWhenDebugEnabled(t *testing.T) {
	host := newLogLevelHost(internal.LogLevelDebug)
	release := internal.RegisterMockWasmHost(host)
	defer release()

	formatCalls := 0
	LogDebugf("response body: %s", countingStringer{calls: &formatCalls})

	require.Equal(t, 1, formatCalls)
	require.Equal(t, 1, host.logCalls)
}

type metricProxyWasmHost struct {
	internal.DefaultProxyWAMSHost
	idToValue map[uint32]uint64
	idToType  map[uint32]internal.MetricType
	nameToID  map[string]uint32
}

func (m metricProxyWasmHost) ProxyDefineMetric(metricType internal.MetricType,
	metricNameData *byte, metricNameSize int32, returnMetricIDPtr *uint32) internal.Status {
	name := unsafe.String(metricNameData, metricNameSize)
	id, ok := m.nameToID[name]
	if !ok {
		id = uint32(len(m.nameToID))
		m.nameToID[name] = id
		m.idToValue[id] = 0
		m.idToType[id] = metricType
	}
	*returnMetricIDPtr = id
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyIncrementMetric(metricID uint32, offset int64) internal.Status {
	val, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}

	m.idToValue[metricID] = val + uint64(offset)
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyRecordMetric(metricID uint32, value uint64) internal.Status {
	_, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	m.idToValue[metricID] = value
	return internal.StatusOK
}

func (m metricProxyWasmHost) ProxyGetMetric(metricID uint32, returnMetricValue *uint64) internal.Status {
	value, ok := m.idToValue[metricID]
	if !ok {
		return internal.StatusBadArgument
	}
	*returnMetricValue = value
	return internal.StatusOK
}

func TestHostCall_Metric(t *testing.T) {
	host := metricProxyWasmHost{
		internal.DefaultProxyWAMSHost{},
		map[uint32]uint64{},
		map[uint32]internal.MetricType{},
		map[string]uint32{},
	}
	release := internal.RegisterMockWasmHost(host)
	defer release()

	t.Run("counter", func(t *testing.T) {
		for _, c := range []struct {
			name   string
			offset uint64
		}{
			{name: "requests", offset: 100},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineCounterMetric(c.name)

				// increment
				m.Increment(c.offset)

				// get
				require.Equal(t, c.offset, m.Value())
			})
		}
	})

	t.Run("gauge", func(t *testing.T) {
		for _, c := range []struct {
			name   string
			offset int64
		}{
			{name: "rate", offset: -50},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineGaugeMetric(c.name)

				// increment
				m.Add(c.offset)

				// get
				require.Equal(t, c.offset, m.Value())
			})
		}
	})

	t.Run("histogram", func(t *testing.T) {
		for _, c := range []struct {
			name  string
			value uint64
		}{
			{name: "request count", value: 10000},
		} {
			t.Run(c.name, func(t *testing.T) {
				// define metric
				m := DefineHistogramMetric(c.name)

				// record
				m.Record(c.value)

				// get
				require.Equal(t, c.value, m.Value())
			})
		}
	})
}
