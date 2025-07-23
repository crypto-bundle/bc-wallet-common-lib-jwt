// MIT NON-AI License
//
// Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated
// documentation files (the "Software"),to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// In addition, the following restrictions apply:
//
// 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine
// learning algorithms, including but not limited to artificial intelligence, natural language processing,
// or data mining.This condition applies to any derivatives, modifications, or updates based on the
// Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
//
// 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
// including but not limited to artificial intelligence, natural language processing, or data mining.
//
// 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and
// may be held liable for any damages resulting from such use.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package jwt

import (
	"math"
	"time"
)

//go:generate easyjson types_field.go

// A FieldType indicates which member of the Field union struct should be used
// and how it should be serialized.
type FieldType uint8

const (
	// UnknownType is the default field type. Attempting to add it to an encoder will panic.
	UnknownType FieldType = iota
	// BinaryType indicates that the field carries an opaque binary blob.
	BinaryType
	// BoolType indicates that the field carries a bool.
	BoolType
	// ByteStringType indicates that the field carries UTF-8 encoded bytes.
	ByteStringType
	// Complex128Type indicates that the field carries a complex128.
	Complex128Type
	// Complex64Type indicates that the field carries a complex64.
	Complex64Type
	// DurationType indicates that the field carries a time.Duration.
	DurationType
	// Float64Type indicates that the field carries a float64.
	Float64Type
	// Float32Type indicates that the field carries a float32.
	Float32Type
	// Int64Type indicates that the field carries an int64.
	Int64Type
	// Int32Type indicates that the field carries an int32.
	Int32Type
	// Int16Type indicates that the field carries an int16.
	Int16Type
	// Int8Type indicates that the field carries an int8.
	Int8Type
	// StringType indicates that the field carries a string.
	StringType
	// StringsType indicates that the field carries a list of strings.
	StringsType
	// TimeType indicates that the field carries a time.Time that is
	// representable by a UnixNano() stored as an int64.
	TimeType
	// TimeFullType indicates that the field carries a time.Time stored as-is.
	TimeFullType
	// Uint64Type indicates that the field carries an uint64.
	Uint64Type
	// Uint32Type indicates that the field carries an uint32.
	Uint32Type
	// Uint16Type indicates that the field carries an uint16.
	Uint16Type
	// Uint8Type indicates that the field carries an uint8.
	Uint8Type
	ErrorType
	// SkipType indicates that the field is a no-op.
	SkipType
)

// A Field is a marshaling operation used to add a key-value pair to a logger's
// context. Most fields are lazily marshaled, so it's inexpensive to add fields
// to disabled debug-level log statements.
// easyjson:json
type Field struct {
	Interface interface{} `json:"interface,omitempty"`
	Key       string      `json:"key"`
	String    string      `json:"string"`
	Integer   int64       `json:"integer"`
	Type      FieldType   `json:"type"`
}

func (f *Field) Scan(target any) error {
	switch f.Type {
	case StringType:
		return ScanString(f, target)
	case StringsType:
		return ScanStrings(f, target)
	case TimeType:
		return ScanTime(f, target)
	default:
		return ErrUnableToScanDataMismatchType
	}
}

// easyjson:json
type stringArray []string

// ScanString scan field to string pointer.
func ScanString(f *Field, target any) error {
	switch t := target.(type) {
	case *string:
		*t = f.String

		return nil
	case string:
		return ErrUnableToScanDataRequiredPointer
	default:
		return ErrUnableToScanDataMismatchType
	}
}

// ScanStrings scan field to string pointer.
func ScanStrings(f *Field, target any) error {
	values, isCasted := f.Interface.([]interface{})
	if !isCasted {
		return ErrUnableToScanDataMismatchType
	}

	result := make([]string, len(values))

	for i := range values {
		value := values[i]

		strValue, isOk := value.(string)
		if !isOk {
			return ErrUnableToScanDataMismatchType
		}

		result[i] = strValue
	}

	switch t := target.(type) {
	case *[]string:
		*t = result

		return nil
	case []string:
		return ErrUnableToScanDataRequiredPointer
	default:
		return ErrUnableToScanDataMismatchType
	}
}

// ScanTime scan field to string or string pointer.
func ScanTime(f *Field, target any) error {
	var value = time.Unix(0, f.Integer)

	switch t := target.(type) {
	case *time.Time:
		*t = value

		return nil
	case time.Time:
		return ErrUnableToScanDataRequiredPointer
	default:
		return ErrUnableToScanDataMismatchType
	}
}

//nolint:gochecknoglobals
var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

// Skip constructs a no-op field, which is often useful when handling invalid
// inputs in other Field constructors.
func Skip() Field {
	return Field{Type: SkipType}
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: BinaryType, Interface: val}
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}

	return Field{Key: key, Type: BoolType, Integer: ival}
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: ByteStringType, Interface: val}
}

// Complex128 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex128 to
// interface{}).
func Complex128(key string, val complex128) Field {
	return Field{Key: key, Type: Complex128Type, Interface: val}
}

// Complex64 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex64 to
// interface{}).
func Complex64(key string, val complex64) Field {
	return Field{Key: key, Type: Complex64Type, Interface: val}
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
//
//nolint:gosec //todo: migrate to big.Int
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, Integer: int64(math.Float64bits(val))}
}

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) Field {
	return Field{Key: key, Type: Float32Type, Integer: int64(math.Float32bits(val))}
}

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: val}
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) Field {
	return Field{Key: key, Type: Int32Type, Integer: int64(val)}
}

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) Field {
	return Field{Key: key, Type: Int16Type, Integer: int64(val)}
}

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) Field {
	return Field{Key: key, Type: Int8Type, Integer: int64(val)}
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Type: StringType, String: val}
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) Field {
	return Field{Key: key, Type: StringsType, Interface: stringArray(ss)}
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

// Uint64 constructs a field with the given key and value.
//
//nolint:gosec //todo: migrate to big.Int
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: Uint64Type, Integer: int64(val)}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) Field {
	return Field{Key: key, Type: Uint32Type, Integer: int64(val)}
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) Field {
	return Field{Key: key, Type: Uint16Type, Integer: int64(val)}
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) Field {
	return Field{Key: key, Type: Uint8Type, Integer: int64(val)}
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	if val.Before(_minTimeInt64) || val.After(_maxTimeInt64) {
		return Field{Key: key, Type: TimeFullType, Interface: val}
	}

	return Field{Key: key, Type: TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, Integer: int64(val)}
}
