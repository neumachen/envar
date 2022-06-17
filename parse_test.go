package envar

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/neumachen/errorx"
	"github.com/neumachen/gohelpers"
	"github.com/stretchr/testify/require"
)

type nestedStruct struct {
	Duration time.Duration `env:"NESTED_DURATION"`
}

type EmbeddedStruct struct {
	Value string `env:"EMBEDDED_VALUE"`
}

type EmbeddedStructPtr struct {
	Value string `env:"EMBEDDED_VALUE"`
}

type embeddedStructUnexported struct {
	Value string `env:"EMBEDDED_VALUE_UNEXPORTED"`
}

type embeddedStructUnexportedPtr struct {
	Value string `env:"EMBEDDED_VALUE_UNEXPORTED"`
}

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler.
func (u *unmarshaler) UnmarshalText(data []byte) error {
	if len(data) < 1 {
		return nil
	}
	parsedDuration, err := time.ParseDuration(string(data))
	if err != nil {
		return errorx.New(err)
	}
	*u = unmarshaler{parsedDuration}

	return nil
}

type testDatum struct {
	envKey    string
	envValue  string
	assertion func(t *testing.T, getter ConfigGetter)
}

func TestParse(t *testing.T) {
	tempDir := t.TempDir()
	testData := []testDatum{
		{
			envKey:   "STRING",
			envValue: "test",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "test", getter.GetString())
				require.Equal(t, "test", *getter.GetStringPtr())
			},
		},
		{
			envKey:   "STRINGS",
			envValue: "test1,test2,test3",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedStrs := []string{"test1", "test2", "test3"}
				require.Equal(t, expectedStrs, getter.GetStrings())
				for i := range getter.GetStringPtrs() {
					require.Equal(t, expectedStrs[i], *getter.GetStringPtrs()[i])
				}
			},
		},
		{
			envKey:   "BOOL",
			envValue: "true",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, true, getter.GetBool())
				require.Equal(t, true, *getter.GetBoolPtr())
			},
		},
		{
			envKey:   "BOOLS",
			envValue: "true,false,true,true,false",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedBools := []bool{true, false, true, true, false}
				require.Equal(t, expectedBools, getter.GetBools())
				for i := range getter.GetBoolPtrs() {
					require.Equal(t, expectedBools[i], *getter.GetBoolPtrs()[i])
				}
			},
		},
		{
			envKey:   "INT",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, 1, getter.GetInt())
				require.Equal(t, 1, *getter.GetIntPtr())
			},
		},
		{
			envKey:   "INTS",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedInts := []int{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedInts, getter.GetInts())
				for i := range getter.GetIntPtrs() {
					require.Equal(t, expectedInts[i], *getter.GetIntPtrs()[i])
				}
			},
		},
		{
			envKey:   "INT8",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, int8(1), getter.GetInt8())
				require.Equal(t, int8(1), *getter.GetInt8Ptr())
			},
		},
		{
			envKey:   "INT8S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedInt8s := []int8{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedInt8s, getter.GetInt8s())
				for i := range getter.GetInt8Ptrs() {
					require.Equal(t, expectedInt8s[i], *getter.GetInt8Ptrs()[i])
				}
			},
		},
		{
			envKey:   "INT16",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, int16(1), getter.GetInt16())
				require.Equal(t, int16(1), *getter.GetInt16Ptr())
			},
		},
		{
			envKey:   "INT16S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedInt16s := []int16{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedInt16s, getter.GetInt16s())
				for i := range getter.GetInt16Ptrs() {
					require.Equal(t, expectedInt16s[i], *getter.GetInt16Ptrs()[i])
				}
			},
		},
		{
			envKey:   "INT32",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, int32(1), getter.GetInt32())
				require.Equal(t, int32(1), *getter.GetInt32Ptr())
			},
		},
		{
			envKey:   "INT32S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedInt32s := []int32{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedInt32s, getter.GetInt32s())
				for i := range getter.GetInt32Ptrs() {
					require.Equal(t, expectedInt32s[i], *getter.GetInt32Ptrs()[i])
				}
			},
		},
		{
			envKey:   "INT64",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, int64(1), getter.GetInt64())
				require.Equal(t, int64(1), *getter.GetInt64Ptr())
			},
		},
		{
			envKey:   "INT64S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedInt64s := []int64{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedInt64s, getter.GetInt64s())
				for i := range getter.GetInt64Ptrs() {
					require.Equal(t, expectedInt64s[i], *getter.GetInt64Ptrs()[i])
				}
			},
		},
		{
			envKey:   "UINT",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, uint(1), getter.GetUint())
				require.Equal(t, uint(1), *getter.GetUintPtr())
			},
		},
		{
			envKey:   "UINTS",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedUints := []uint{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedUints, getter.GetUints())
				for i := range getter.GetUintPtrs() {
					require.Equal(t, expectedUints[i], *getter.GetUintPtrs()[i])
				}
			},
		},
		{
			envKey:   "UINT8",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, uint8(1), getter.GetUint8())
				require.Equal(t, uint8(1), *getter.GetUint8Ptr())
			},
		},
		{
			envKey:   "UINT8S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedUint8s := []uint8{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedUint8s, getter.GetUint8s())
				for i := range getter.GetUint8Ptrs() {
					require.Equal(t, expectedUint8s[i], *getter.GetUint8Ptrs()[i])
				}
			},
		},
		{
			envKey:   "UINT16",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, uint16(1), getter.GetUint16())
				require.Equal(t, uint16(1), *getter.GetUint16Ptr())
			},
		},
		{
			envKey:   "UINT16S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedUint16s := []uint16{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedUint16s, getter.GetUint16s())
				for i := range getter.GetUint16Ptrs() {
					require.Equal(t, expectedUint16s[i], *getter.GetUint16Ptrs()[i])
				}
			},
		},
		{
			envKey:   "UINT32",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, uint32(1), getter.GetUint32())
				require.Equal(t, uint32(1), *getter.GetUint32Ptr())
			},
		},
		{
			envKey:   "UINT32S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedUint32s := []uint32{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedUint32s, getter.GetUint32s())
				for i := range getter.GetUint32Ptrs() {
					require.Equal(t, expectedUint32s[i], *getter.GetUint32Ptrs()[i])
				}
			},
		},
		{
			envKey:   "UINT64",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, uint64(1), getter.GetUint64())
				require.Equal(t, uint64(1), *getter.GetUint64Ptr())
			},
		},
		{
			envKey:   "UINT64S",
			envValue: "1,2,3,4,5,6",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedUint64s := []uint64{1, 2, 3, 4, 5, 6}
				require.Equal(t, expectedUint64s, getter.GetUint64s())
				for i := range getter.GetUint64Ptrs() {
					require.Equal(t, expectedUint64s[i], *getter.GetUint64Ptrs()[i])
				}
			},
		},
		{
			envKey:   "FLOAT32",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, float32(1), getter.GetFloat32())
				require.Equal(t, float32(1), *getter.GetFloat32Ptr())
			},
		},
		{
			envKey:   "FLOAT32S",
			envValue: "1.0,2.0,3.0,4.0,5.0,6.0",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedFloat32s := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
				require.Equal(t, expectedFloat32s, getter.GetFloat32s())
				for i := range getter.GetFloat32Ptrs() {
					require.Equal(t, expectedFloat32s[i], *getter.GetFloat32Ptrs()[i])
				}
			},
		},
		{
			envKey:   "FLOAT64",
			envValue: "1",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, float64(1), getter.GetFloat64())
				require.Equal(t, float64(1), *getter.GetFloat64Ptr())
			},
		},
		{
			envKey:   "FLOAT64S",
			envValue: "1.0,2.0,3.0,4.0,5.0,6.0",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedFloat64s := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
				require.Equal(t, expectedFloat64s, getter.GetFloat64s())
				for i := range getter.GetFloat64Ptrs() {
					require.Equal(t, expectedFloat64s[i], *getter.GetFloat64Ptrs()[i])
				}
			},
		},
		{
			envKey:   "DURATION",
			envValue: "10s",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, time.Second*10, getter.GetDuration())
				require.Equal(t, time.Second*10, *getter.GetDurationPtr())
			},
		},
		{
			envKey:   "DURATIONS",
			envValue: "5ms, 10s,15m,20h",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedDurations := []time.Duration{
					time.Millisecond * 5,
					time.Second * 10,
					time.Minute * 15,
					time.Hour * 20,
				}
				require.Equal(t, expectedDurations, getter.GetDurations())
				for i := range getter.GetDurationPtrs() {
					require.Equal(t, expectedDurations[i], *getter.GetDurationPtrs()[i])
				}
			},
		},
		{
			envKey:   "UNMARSHALER",
			envValue: "10s",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, time.Second*10, getter.GetUnmarshaler().Duration)
				c := *getter.GetUnmarshalerPtr()
				require.Equal(t, time.Second*10, c.Duration)
			},
		},
		{
			envKey:   "UNMARSHALERS",
			envValue: "5ms,10s,15m,20h",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedDurations := []time.Duration{
					time.Millisecond * 5,
					time.Second * 10,
					time.Minute * 15,
					time.Hour * 20,
				}
				for i := range getter.GetUnmarshalers() {
					require.Equal(t, expectedDurations[i], getter.GetUnmarshalers()[i].Duration)
				}
				for i := range getter.GetUnmarshalerPtrs() {
					require.Equal(t, expectedDurations[i], (*getter.GetUnmarshalerPtrs()[i]).Duration)
				}
			},
		},
		{
			envKey:   "NESTED_DURATION",
			envValue: "10s",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedDuration := time.Second * 10
				require.Equal(t, expectedDuration, getter.GetNestedStruct().Duration)
				require.Equal(t, expectedDuration, getter.GetNestedStructPtr().Duration)
				for i := range getter.GetNestedStructs() {
					require.Equal(t, expectedDuration, getter.GetNestedStructs()[i].Duration)
				}
				for i := range getter.GetNestedStructPtrs() {
					require.Equal(t, expectedDuration, (*getter.GetNestedStructPtrs()[i]).Duration)
				}
			},
		},
		{
			envKey:   "EMBEDDED_VALUE",
			envValue: "test",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedValue := "test"
				require.Equal(t, expectedValue, getter.GetEmbeddedStruct().Value)
				require.Equal(t, expectedValue, getter.GetEmbeddedStructPtr().Value)
			},
		},
		{
			envKey:   "EMBEDDED_VALUE_UNEXPORTED",
			envValue: "test",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedValue := ""
				require.Equal(t, expectedValue, getter.GetEmbeddedStructUnexported().Value)
				require.Nil(t, getter.GetEmbeddedStructUnexportedPtr())
			},
		},
		{
			envKey:   "URL",
			envValue: "https://envartest.com",
			assertion: func(t *testing.T, getter ConfigGetter) {
				url := getter.GetURL()
				require.Equal(t, "https://envartest.com", url.String())
				require.Equal(t, "https://envartest.com", getter.GetURLPtr().String())
			},
		},
		{
			envKey:   "URLS",
			envValue: "https://envartest.com,https://testing.app",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedURLs := []string{
					"https://envartest.com",
					"https://testing.app",
				}
				for i := range getter.GetURLs() {
					require.Equal(t, expectedURLs[i], getter.GetURLs()[i].String())
				}
				for i := range getter.GetURLPtrs() {
					require.Equal(t, expectedURLs[i], getter.GetURLPtrs()[i].String())
				}
			},
		},
		{
			envKey: "FILE",
			envValue: func() string {
				file := filepath.Join(tempDir, "temp_file")
				require.NoError(t, ioutil.WriteFile(file, []byte("secret"), 0o660))
				return file
			}(),
			assertion: func(t *testing.T, getter ConfigGetter) {
				file := getter.GetFile()
				require.Equal(t, filepath.Join(tempDir, "temp_file"), file.Name())
				require.Equal(t, filepath.Join(tempDir, "temp_file"), getter.GetFilePtr().Name())
			},
		},
		{
			envKey: "FILES",
			envValue: func() string {
				files := []string{
					"temp_file_1",
					"temp_file_2",
					"temp_file_3",
					"temp_file_4",
				}
				fileStrs := make([]string, 0)
				for i := range files {
					file := filepath.Join(tempDir, files[i])
					require.NoError(t, ioutil.WriteFile(file, []byte("secret"), 0o660))
					fileStrs = append(fileStrs, file)
				}
				return strings.Join(fileStrs, ",")
			}(),
			assertion: func(t *testing.T, getter ConfigGetter) {
				fileNames := []string{
					filepath.Join(tempDir, "temp_file_1"),
					filepath.Join(tempDir, "temp_file_2"),
					filepath.Join(tempDir, "temp_file_3"),
					filepath.Join(tempDir, "temp_file_4"),
				}

				for i := range getter.GetFiles() {
					require.Equal(t, fileNames[i], getter.GetFiles()[i].Name())
				}
				for i := range getter.GetFilePtrs() {
					require.Equal(t, fileNames[i], getter.GetFilePtrs()[i].Name())
				}
			},
		},
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetNotTagged().String)
			},
		},
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetNotTagged().String)
			},
		},
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetNotAnEnv())
			},
		},
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetUnexportedString())
			},
		},
	}

	t.Run("embedded", func(t *testing.T) {
		for i := range testData {
			if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
				continue
			}
			os.Setenv(testData[i].envKey, testData[i].envValue)
		}

		defer func() {
			for i := range testData {
				os.Unsetenv(testData[i].envKey)
			}
		}()

		type embedded struct {
			Config `env:",nested"`
		}

		e := &embedded{}

		parserCtx, err := Parse(e)
		require.NoError(t, err)
		require.NotNil(t, parserCtx)

		for i := range testData {
			testData[i].assertion(t, e)
		}
	})

	t.Run("without prefix", func(t *testing.T) {
		for i := range testData {
			if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
				continue
			}
			os.Setenv(testData[i].envKey, testData[i].envValue)
		}

		defer func() {
			for i := range testData {
				os.Unsetenv(testData[i].envKey)
			}
		}()

		cfg := Config{}

		parserCtx, err := Parse(&cfg)
		require.NoError(t, err)
		require.NotNil(t, parserCtx)

		for i := range testData {
			testData[i].assertion(t, &cfg)
		}
	})

	t.Run("with prefix", func(t *testing.T) {
		prefix := "MOO"
		for i := range testData {
			if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
				continue
			}
			os.Setenv(strings.Join([]string{prefix, testData[i].envKey}, "_"), testData[i].envValue)
		}

		defer func() {
			for i := range testData {
				os.Unsetenv(strings.Join([]string{prefix, testData[i].envKey}, "_"))
			}
		}()

		cfg := Config{}

		parserCtx, err := Parse(&cfg, SetEnvPrefix(prefix))
		require.NoError(t, err)
		require.NotNil(t, parserCtx)

		for i := range testData {
			testData[i].assertion(t, &cfg)
		}
	})
}

type ConfigDefault struct {
	Config

	String     string    `env:"STRING,default=test"`
	StringPtr  *string   `env:"STRING,default=test"`
	Strings    []string  `env:"STRINGS,default=test1|test2|test3"`
	StringPtrs []*string `env:"STRINGS,default=test1|test2|test3"`
}

func (c *ConfigDefault) GetString() string {
	return c.String
}

func (c *ConfigDefault) GetStringPtr() *string {
	return c.StringPtr
}

func (c *ConfigDefault) GetStrings() []string {
	return c.Strings
}

func (c *ConfigDefault) GetStringPtrs() []*string {
	return c.StringPtrs
}

func TestParse_Default(t *testing.T) {
	testData := []testDatum{
		{
			envKey:   "STRING",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "test", getter.GetString())
				require.Equal(t, "test", *getter.GetStringPtr())
			},
		},
		{
			envKey:   "STRINGS",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedStrs := []string{"test1", "test2", "test3"}
				require.Equal(t, expectedStrs, getter.GetStrings())
				for i := range getter.GetStringPtrs() {
					require.Equal(t, expectedStrs[i], *getter.GetStringPtrs()[i])
				}
			},
		},
	}
	for i := range testData {
		if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
			continue
		}
		os.Setenv(testData[i].envKey, testData[i].envValue)
	}

	defer func() {
		for i := range testData {
			os.Unsetenv(testData[i].envKey)
		}
	}()

	cfg := ConfigDefault{}

	parserCtx, err := Parse(&cfg)
	require.NoError(t, err)
	require.NotNil(t, parserCtx)

	for i := range testData {
		testData[i].assertion(t, &cfg)
	}
}

type ConfigValidateRequired struct {
	Config

	String     string    `env:"STRING,validate=required"`
	StringPtr  *string   `env:"STRING,validate=required"`
	Strings    []string  `env:"STRINGS,validate=required"`
	StringPtrs []*string `env:"STRINGS,validate=required"`

	File     os.File    `env:"FILE,validate=required"`
	FilePtr  *os.File   `env:"FILE,validate=required"`
	Files    []os.File  `env:"FILES,validate=required"`
	FilePtrs []*os.File `env:"FILES,validate=required"`
}

func (c *ConfigValidateRequired) GetString() string {
	return c.String
}

func (c *ConfigValidateRequired) GetStringPtr() *string {
	return c.StringPtr
}

func (c *ConfigValidateRequired) GetStrings() []string {
	return c.Strings
}

func (c *ConfigValidateRequired) GetStringPtrs() []*string {
	return c.StringPtrs
}

func (c *ConfigValidateRequired) GetFile() os.File {
	return c.File
}

func (c *ConfigValidateRequired) GetFilePtr() *os.File {
	return c.FilePtr
}

func (c *ConfigValidateRequired) GetFiles() []os.File {
	return c.Files
}

func (c *ConfigValidateRequired) GetFilePtrs() []*os.File {
	return c.FilePtrs
}

func TestParse_Validate_Required(t *testing.T) {
	testData := []testDatum{
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetString())
				require.Equal(t, (*string)(nil), getter.GetStringPtr())
			},
		},
		{
			envKey:   "",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedStrs := []string(nil)
				require.Equal(t, expectedStrs, getter.GetStrings())
				for i := range getter.GetStringPtrs() {
					require.Equal(t, expectedStrs[i], *getter.GetStringPtrs()[i])
				}
			},
		},
	}

	for i := range testData {
		if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
			continue
		}
		os.Setenv(testData[i].envKey, testData[i].envValue)
	}

	defer func() {
		for i := range testData {
			os.Unsetenv(testData[i].envKey)
		}
	}()

	cfg := ConfigValidateRequired{}

	parserCtx, err := Parse(&cfg)
	require.NoError(t, err)
	require.NotNil(t, parserCtx)

	for i := range testData {
		testData[i].assertion(t, &cfg)
	}

	valErrors := parserCtx.GetValidationErrors()
	require.NotEmpty(t, valErrors)
	fields := []string{
		"String",
		"StringPtr",
		"Strings",
		"StringPtrs",
		"File",
		"FilePtr",
		"Files",
		"FilePtrs",
	}
	require.Equal(t, len(fields), valErrors.GetLength())

	for i := range fields {
		require.True(t, valErrors.HasErrors(fields[i]))
	}
}

type ConfigValidateNotEmpty struct {
	Config

	String     string    `env:"STRING,validate=not_empty"`
	StringPtr  *string   `env:"STRING,validate=not_empty"`
	Strings    []string  `env:"STRINGS,validate=not_empty"`
	StringPtrs []*string `env:"STRINGS,validate=not_empty"`
}

func (c *ConfigValidateNotEmpty) GetString() string {
	return c.String
}

func (c *ConfigValidateNotEmpty) GetStringPtr() *string {
	return c.StringPtr
}

func (c *ConfigValidateNotEmpty) GetStrings() []string {
	return c.Strings
}

func (c *ConfigValidateNotEmpty) GetStringPtrs() []*string {
	return c.StringPtrs
}

func TestParse_Validate_NotEmpty(t *testing.T) {
	testData := []testDatum{
		{
			envKey:   "STRING",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "", getter.GetString())
				require.Equal(t, (*string)(nil), getter.GetStringPtr())
			},
		},
		{
			envKey:   "STRINGS",
			envValue: "",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedStrs := []string(nil)
				require.Equal(t, expectedStrs, getter.GetStrings())
				for i := range getter.GetStringPtrs() {
					require.Equal(t, expectedStrs[i], *getter.GetStringPtrs()[i])
				}
			},
		},
	}

	for i := range testData {
		if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
			continue
		}
		os.Setenv(testData[i].envKey, testData[i].envValue)
	}

	defer func() {
		for i := range testData {
			os.Unsetenv(testData[i].envKey)
		}
	}()

	cfg := ConfigValidateNotEmpty{}

	parserCtx, err := Parse(&cfg)
	require.NoError(t, err)
	require.NotNil(t, parserCtx)

	for i := range testData {
		testData[i].assertion(t, &cfg)
	}

	valErrors := parserCtx.GetValidationErrors()
	require.NotEmpty(t, valErrors)

	fields := []string{
		"String",
		"StringPtr",
		"Strings",
		"StringPtrs",
	}
	require.Equal(t, len(fields), valErrors.GetLength())

	for i := range fields {
		require.True(t, valErrors.HasErrors(fields[i]))
	}
}

type ConfigUnset struct {
	Config

	String     string    `env:"STRING,unset"`
	StringPtr  *string   `env:"STRING,unset"`
	Strings    []string  `env:"STRINGS,unset"`
	StringPtrs []*string `env:"STRINGS,unset"`
}

func (c *ConfigUnset) GetString() string {
	return c.String
}

func (c *ConfigUnset) GetStringPtr() *string {
	return c.StringPtr
}

func (c *ConfigUnset) GetStrings() []string {
	return c.Strings
}

func (c *ConfigUnset) GetStringPtrs() []*string {
	return c.StringPtrs
}

func TestParse_Unset(t *testing.T) {
	testData := []testDatum{
		{
			envKey:   "STRING",
			envValue: "test",
			assertion: func(t *testing.T, getter ConfigGetter) {
				require.Equal(t, "test", getter.GetString())
				require.Equal(t, "test", *getter.GetStringPtr())
				require.Empty(t, os.Getenv("STRING"))
			},
		},
		{
			envKey:   "STRINGS",
			envValue: "test1,test2,test3",
			assertion: func(t *testing.T, getter ConfigGetter) {
				expectedStrs := []string{"test1", "test2", "test3"}
				require.Equal(t, expectedStrs, getter.GetStrings())
				for i := range getter.GetStringPtrs() {
					require.Equal(t, expectedStrs[i], *getter.GetStringPtrs()[i])
				}
				require.Empty(t, os.Getenv("STRINGS"))
			},
		},
	}

	for i := range testData {
		if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
			continue
		}
		os.Setenv(testData[i].envKey, testData[i].envValue)
	}

	defer func() {
		for i := range testData {
			os.Unsetenv(testData[i].envKey)
		}
	}()

	cfg := ConfigUnset{}

	parserCtx, err := Parse(&cfg)
	require.NoError(t, err)
	require.NotNil(t, parserCtx)

	for i := range testData {
		testData[i].assertion(t, &cfg)
	}
}

type ConfigFile struct {
	Config

	File     os.File    `env:"FILE,validate=required|not_empty"`
	FilePtr  *os.File   `env:"FILE,validate=required|not_empty"`
	Files    []os.File  `env:"FILES,validate=required|not_empty"`
	FilePtrs []*os.File `env:"FILES,validate=required|not_empty"`
}

func (c *ConfigFile) GetFile() os.File {
	return c.File
}

func (c *ConfigFile) GetFilePtr() *os.File {
	return c.FilePtr
}

func (c *ConfigFile) GetFiles() []os.File {
	return c.Files
}

func (c *ConfigFile) GetFilePtrs() []*os.File {
	return c.FilePtrs
}

func TestParse_File(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		testData := []testDatum{
			{
				envKey:   "FILE",
				envValue: "i-do-not-exist.txt",
				assertion: func(t *testing.T, getter ConfigGetter) {
					require.Equal(t, os.File{}, getter.GetFile())
					require.Equal(t, (*os.File)(nil), getter.GetFilePtr())
				},
			},
			{
				envKey:   "FILES",
				envValue: "i-do-not-exist-1.txt,i-do-not-exist-2.txt,i-do-not-exist-3.txt",
				assertion: func(t *testing.T, getter ConfigGetter) {
					require.Nil(t, getter.GetFiles())
					require.Nil(t, getter.GetFilePtrs())
				},
			},
		}

		for i := range testData {
			if v := testData[i].envKey; gohelpers.StringIsEmpty(v) {
				continue
			}
			os.Setenv(testData[i].envKey, testData[i].envValue)
		}

		defer func() {
			for i := range testData {
				os.Unsetenv(testData[i].envKey)
			}
		}()

		cfg := ConfigFile{}

		parserCtx, err := Parse(&cfg)
		require.Error(t, err)
		require.Nil(t, parserCtx)

		for i := range testData {
			testData[i].assertion(t, &cfg)
		}
	})
}
