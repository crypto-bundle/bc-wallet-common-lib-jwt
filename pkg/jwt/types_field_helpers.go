/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package jwt

import "time"

type anyFieldC[T any] func(string, T) Field

func (f anyFieldC[T]) Any(key string, val any) Field {
	v, _ := val.(T)
	// val is guaranteed to be a T, except when it's nil.
	return f(key, v)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	var c interface{ Any(string, any) Field }

	switch value.(type) {
	case bool:
		c = anyFieldC[bool](Bool)
	case complex128:
		c = anyFieldC[complex128](Complex128)
	case complex64:
		c = anyFieldC[complex64](Complex64)
	case float64:
		c = anyFieldC[float64](Float64)
	case float32:
		c = anyFieldC[float32](Float32)
	case int:
		c = anyFieldC[int](Int)
	case int64:
		c = anyFieldC[int64](Int64)
	case int32:
		c = anyFieldC[int32](Int32)
	case int16:
		c = anyFieldC[int16](Int16)
	case int8:
		c = anyFieldC[int8](Int8)
	case string:
		c = anyFieldC[string](String)
	case uint:
		c = anyFieldC[uint](Uint)
	case uint64:
		c = anyFieldC[uint64](Uint64)
	case uint32:
		c = anyFieldC[uint32](Uint32)
	case uint16:
		c = anyFieldC[uint16](Uint16)
	case uint8:
		c = anyFieldC[uint8](Uint8)
	case time.Time:
		c = anyFieldC[time.Time](Time)
	case time.Duration:
		c = anyFieldC[time.Duration](Duration)
	case error:
		c = anyFieldC[error](NamedError)
	}

	return c.Any(key, value)
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) Field {
	return NamedError("error", err)
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: key, Type: ErrorType, Interface: err}
}
