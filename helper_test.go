package envar

import (
	"net/url"
	"os"
	"time"
)

type Config struct {
	String     string    `env:"STRING"`
	StringPtr  *string   `env:"STRING"`
	Strings    []string  `env:"STRINGS"`
	StringPtrs []*string `env:"STRINGS"`

	Bool     bool    `env:"BOOL"`
	BoolPtr  *bool   `env:"BOOL"`
	Bools    []bool  `env:"BOOLS"`
	BoolPtrs []*bool `env:"BOOLS"`

	Int     int    `env:"INT"`
	IntPtr  *int   `env:"INT"`
	Ints    []int  `env:"INTS"`
	IntPtrs []*int `env:"INTS"`

	Int8     int8    `env:"INT8"`
	Int8Ptr  *int8   `env:"INT8"`
	Int8s    []int8  `env:"INT8S"`
	Int8Ptrs []*int8 `env:"INT8S"`

	Int16     int16    `env:"INT16"`
	Int16Ptr  *int16   `env:"INT16"`
	Int16s    []int16  `env:"INT16S"`
	Int16Ptrs []*int16 `env:"INT16S"`

	Int32     int32    `env:"INT32"`
	Int32Ptr  *int32   `env:"INT32"`
	Int32s    []int32  `env:"INT32S"`
	Int32Ptrs []*int32 `env:"INT32S"`

	Int64     int64    `env:"INT64"`
	Int64Ptr  *int64   `env:"INT64"`
	Int64s    []int64  `env:"INT64S"`
	Int64Ptrs []*int64 `env:"INT64S"`

	Uint     uint    `env:"UINT"`
	UintPtr  *uint   `env:"UINT"`
	Uints    []uint  `env:"UINTS"`
	UintPtrs []*uint `env:"UINTS"`

	Uint8     uint8    `env:"UINT8"`
	Uint8Ptr  *uint8   `env:"UINT8"`
	Uint8s    []uint8  `env:"UINT8S"`
	Uint8Ptrs []*uint8 `env:"UINT8S"`

	Uint16     uint16    `env:"UINT16"`
	Uint16Ptr  *uint16   `env:"UINT16"`
	Uint16s    []uint16  `env:"UINT16S"`
	Uint16Ptrs []*uint16 `env:"UINT16S"`

	Uint32     uint32    `env:"UINT32"`
	Uint32Ptr  *uint32   `env:"UINT32"`
	Uint32s    []uint32  `env:"UINT32S"`
	Uint32Ptrs []*uint32 `env:"UINT32S"`

	Uint64     uint64    `env:"UINT64"`
	Uint64Ptr  *uint64   `env:"UINT64"`
	Uint64s    []uint64  `env:"UINT64S"`
	Uint64Ptrs []*uint64 `env:"UINT64S"`

	Float32     float32    `env:"FLOAT32"`
	Float32Ptr  *float32   `env:"FLOAT32"`
	Float32s    []float32  `env:"FLOAT32S"`
	Float32Ptrs []*float32 `env:"FLOAT32S"`

	Float64     float64    `env:"FLOAT64"`
	Float64Ptr  *float64   `env:"FLOAT64"`
	Float64s    []float64  `env:"FLOAT64S"`
	Float64Ptrs []*float64 `env:"FLOAT64S"`

	Duration     time.Duration    `env:"DURATION"`
	Durations    []time.Duration  `env:"DURATIONS"`
	DurationPtr  *time.Duration   `env:"DURATION"`
	DurationPtrs []*time.Duration `env:"DURATIONS"`

	Unmarshaler     unmarshaler    `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler   `env:"UNMARSHALER"`
	Unmarshalers    []unmarshaler  `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler `env:"UNMARSHALERS"`

	NestedStruct     nestedStruct    `env:",nested"`
	NestedStructPtr  *nestedStruct   `env:",nested"`
	NestedStructs    []nestedStruct  `env:",nested"`
	NestedStructPtrs []*nestedStruct `env:",nested"`

	EmbeddedStruct               `env:",nested"`
	*EmbeddedStructPtr           `env:",nested"`
	embeddedStructUnexported     `env:",nested"`
	*embeddedStructUnexportedPtr `env:",nested"`

	URL     url.URL    `env:"URL"`
	URLPtr  *url.URL   `env:"URL"`
	URLs    []url.URL  `env:"URLS"`
	URLPtrs []*url.URL `env:"URLS"`

	File     os.File    `env:"FILE"`
	FilePtr  *os.File   `env:"FILE"`
	Files    []os.File  `env:"FILES"`
	FilePtrs []*os.File `env:"FILES"`

	NotTagged struct {
		String string `env:"PARENT_STRUCT_NOT_TAGGED"`
	}

	NotAnEnv         string
	unexportedString string `env:"UNEXPORTED_STRING"`
}

func (c *Config) GetString() string {
	return c.String
}

func (c *Config) GetStringPtr() *string {
	return c.StringPtr
}

func (c *Config) GetStrings() []string {
	return c.Strings
}

func (c *Config) GetStringPtrs() []*string {
	return c.StringPtrs
}

func (c *Config) GetBool() bool {
	return c.Bool
}

func (c *Config) GetBoolPtr() *bool {
	return c.BoolPtr
}

func (c *Config) GetBools() []bool {
	return c.Bools
}

func (c *Config) GetBoolPtrs() []*bool {
	return c.BoolPtrs
}

func (c *Config) GetInt() int {
	return c.Int
}

func (c *Config) GetIntPtr() *int {
	return c.IntPtr
}

func (c *Config) GetInts() []int {
	return c.Ints
}

func (c *Config) GetIntPtrs() []*int {
	return c.IntPtrs
}

func (c *Config) GetInt8() int8 {
	return c.Int8
}

func (c *Config) GetInt8Ptr() *int8 {
	return c.Int8Ptr
}

func (c *Config) GetInt8s() []int8 {
	return c.Int8s
}

func (c *Config) GetInt8Ptrs() []*int8 {
	return c.Int8Ptrs
}

func (c *Config) GetInt16() int16 {
	return c.Int16
}

func (c *Config) GetInt16Ptr() *int16 {
	return c.Int16Ptr
}

func (c *Config) GetInt16s() []int16 {
	return c.Int16s
}

func (c *Config) GetInt16Ptrs() []*int16 {
	return c.Int16Ptrs
}

func (c *Config) GetInt32() int32 {
	return c.Int32
}

func (c *Config) GetInt32Ptr() *int32 {
	return c.Int32Ptr
}

func (c *Config) GetInt32s() []int32 {
	return c.Int32s
}

func (c *Config) GetInt32Ptrs() []*int32 {
	return c.Int32Ptrs
}

func (c *Config) GetInt64() int64 {
	return c.Int64
}

func (c *Config) GetInt64Ptr() *int64 {
	return c.Int64Ptr
}

func (c *Config) GetInt64s() []int64 {
	return c.Int64s
}

func (c *Config) GetInt64Ptrs() []*int64 {
	return c.Int64Ptrs
}

func (c *Config) GetUint() uint {
	return c.Uint
}

func (c *Config) GetUintPtr() *uint {
	return c.UintPtr
}

func (c *Config) GetUints() []uint {
	return c.Uints
}

func (c *Config) GetUintPtrs() []*uint {
	return c.UintPtrs
}

func (c *Config) GetUint8() uint8 {
	return c.Uint8
}

func (c *Config) GetUint8Ptr() *uint8 {
	return c.Uint8Ptr
}

func (c *Config) GetUint8s() []uint8 {
	return c.Uint8s
}

func (c *Config) GetUint8Ptrs() []*uint8 {
	return c.Uint8Ptrs
}

func (c *Config) GetUint16() uint16 {
	return c.Uint16
}

func (c *Config) GetUint16Ptr() *uint16 {
	return c.Uint16Ptr
}

func (c *Config) GetUint16s() []uint16 {
	return c.Uint16s
}

func (c *Config) GetUint16Ptrs() []*uint16 {
	return c.Uint16Ptrs
}

func (c *Config) GetUint32() uint32 {
	return c.Uint32
}

func (c *Config) GetUint32Ptr() *uint32 {
	return c.Uint32Ptr
}

func (c *Config) GetUint32s() []uint32 {
	return c.Uint32s
}

func (c *Config) GetUint32Ptrs() []*uint32 {
	return c.Uint32Ptrs
}

func (c *Config) GetUint64() uint64 {
	return c.Uint64
}

func (c *Config) GetUint64Ptr() *uint64 {
	return c.Uint64Ptr
}

func (c *Config) GetUint64s() []uint64 {
	return c.Uint64s
}

func (c *Config) GetUint64Ptrs() []*uint64 {
	return c.Uint64Ptrs
}

func (c *Config) GetFloat32() float32 {
	return c.Float32
}

func (c *Config) GetFloat32Ptr() *float32 {
	return c.Float32Ptr
}

func (c *Config) GetFloat32s() []float32 {
	return c.Float32s
}

func (c *Config) GetFloat32Ptrs() []*float32 {
	return c.Float32Ptrs
}

func (c *Config) GetFloat64() float64 {
	return c.Float64
}

func (c *Config) GetFloat64Ptr() *float64 {
	return c.Float64Ptr
}

func (c *Config) GetFloat64s() []float64 {
	return c.Float64s
}

func (c *Config) GetFloat64Ptrs() []*float64 {
	return c.Float64Ptrs
}

func (c *Config) GetDuration() time.Duration {
	return c.Duration
}

func (c *Config) GetDurationPtr() *time.Duration {
	return c.DurationPtr
}

func (c *Config) GetDurations() []time.Duration {
	return c.Durations
}

func (c *Config) GetDurationPtrs() []*time.Duration {
	return c.DurationPtrs
}

func (c *Config) GetUnmarshaler() unmarshaler {
	return c.Unmarshaler
}

func (c *Config) GetUnmarshalerPtr() *unmarshaler {
	return c.UnmarshalerPtr
}

func (c *Config) GetUnmarshalers() []unmarshaler {
	return c.Unmarshalers
}

func (c *Config) GetUnmarshalerPtrs() []*unmarshaler {
	return c.UnmarshalerPtrs
}

func (c *Config) GetNestedStruct() nestedStruct {
	return c.NestedStruct
}

func (c *Config) GetNestedStructPtr() *nestedStruct {
	return c.NestedStructPtr
}

func (c *Config) GetEmbeddedStruct() EmbeddedStruct {
	return c.EmbeddedStruct
}

func (c *Config) GetEmbeddedStructPtr() *EmbeddedStructPtr {
	return c.EmbeddedStructPtr
}

func (c *Config) GetEmbeddedStructUnexported() embeddedStructUnexported {
	return c.embeddedStructUnexported
}

func (c *Config) GetEmbeddedStructUnexportedPtr() *embeddedStructUnexportedPtr {
	return c.embeddedStructUnexportedPtr
}

func (c *Config) GetNestedStructs() []nestedStruct {
	return c.NestedStructs
}

func (c *Config) GetNestedStructPtrs() []*nestedStruct {
	return c.NestedStructPtrs
}

func (c *Config) GetURL() url.URL {
	return c.URL
}

func (c *Config) GetURLPtr() *url.URL {
	return c.URLPtr
}

func (c *Config) GetURLs() []url.URL {
	return c.URLs
}

func (c *Config) GetURLPtrs() []*url.URL {
	return c.URLPtrs
}

func (c *Config) GetFile() os.File {
	return c.File
}

func (c *Config) GetFilePtr() *os.File {
	return c.FilePtr
}

func (c *Config) GetFiles() []os.File {
	return c.Files
}

func (c *Config) GetFilePtrs() []*os.File {
	return c.FilePtrs
}

func (c *Config) GetNotTagged() struct {
	String string `env:"PARENT_STRUCT_NOT_TAGGED"`
} {
	return c.NotTagged
}

func (c *Config) GetNotAnEnv() string {
	return c.NotAnEnv
}

func (c *Config) GetUnexportedString() string {
	return c.unexportedString
}

var _ ConfigGetter = (*Config)(nil)

type ConfigGetter interface {
	GetString() string
	GetStringPtr() *string
	GetStrings() []string
	GetStringPtrs() []*string
	GetBool() bool
	GetBoolPtr() *bool
	GetBools() []bool
	GetBoolPtrs() []*bool
	GetInt() int
	GetIntPtr() *int
	GetInts() []int
	GetIntPtrs() []*int
	GetInt8() int8
	GetInt8Ptr() *int8
	GetInt8s() []int8
	GetInt8Ptrs() []*int8
	GetInt16() int16
	GetInt16Ptr() *int16
	GetInt16s() []int16
	GetInt16Ptrs() []*int16
	GetInt32() int32
	GetInt32Ptr() *int32
	GetInt32s() []int32
	GetInt32Ptrs() []*int32
	GetInt64() int64
	GetInt64Ptr() *int64
	GetInt64s() []int64
	GetInt64Ptrs() []*int64
	GetUint() uint
	GetUintPtr() *uint
	GetUints() []uint
	GetUintPtrs() []*uint
	GetUint8() uint8
	GetUint8Ptr() *uint8
	GetUint8s() []uint8
	GetUint8Ptrs() []*uint8
	GetUint16() uint16
	GetUint16Ptr() *uint16
	GetUint16s() []uint16
	GetUint16Ptrs() []*uint16
	GetUint32() uint32
	GetUint32Ptr() *uint32
	GetUint32s() []uint32
	GetUint32Ptrs() []*uint32
	GetUint64() uint64
	GetUint64Ptr() *uint64
	GetUint64s() []uint64
	GetUint64Ptrs() []*uint64
	GetFloat32() float32
	GetFloat32Ptr() *float32
	GetFloat32s() []float32
	GetFloat32Ptrs() []*float32
	GetFloat64() float64
	GetFloat64Ptr() *float64
	GetFloat64s() []float64
	GetFloat64Ptrs() []*float64
	GetDuration() time.Duration
	GetDurationPtr() *time.Duration
	GetDurations() []time.Duration
	GetDurationPtrs() []*time.Duration
	GetUnmarshaler() unmarshaler
	GetUnmarshalerPtr() *unmarshaler
	GetUnmarshalers() []unmarshaler
	GetUnmarshalerPtrs() []*unmarshaler
	GetNestedStruct() nestedStruct
	GetNestedStructPtr() *nestedStruct
	GetNestedStructs() []nestedStruct
	GetNestedStructPtrs() []*nestedStruct
	GetEmbeddedStruct() EmbeddedStruct
	GetEmbeddedStructPtr() *EmbeddedStructPtr
	GetEmbeddedStructUnexported() embeddedStructUnexported
	GetEmbeddedStructUnexportedPtr() *embeddedStructUnexportedPtr
	GetURL() url.URL
	GetURLPtr() *url.URL
	GetURLs() []url.URL
	GetURLPtrs() []*url.URL
	GetFile() os.File
	GetFilePtr() *os.File
	GetFiles() []os.File
	GetFilePtrs() []*os.File
	GetNotAnEnv() string
	GetUnexportedString() string
	GetNotTagged() struct {
		String string `env:"PARENT_STRUCT_NOT_TAGGED"`
	}
}
