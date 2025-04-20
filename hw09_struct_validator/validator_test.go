package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidateFunctionalErrors(t *testing.T) {
	tests := []struct {
		name        string
		in          any
		expectedErr error
	}{
		{
			name: "invalid validator value",
			in: struct {
				ID int `validate:"min:one"`
			}{
				ID: 1,
			},
			expectedErr: ErrInvalidValueType,
		},
		{
			name: "invalid regexp",
			in: struct {
				ID string `validate:"regexp:^("`
			}{},
			expectedErr: ErrInvalidRegexp,
		},
		{
			name: "invalid validator format",
			in: struct {
				ID string `validate:"min-4"`
			}{},
			expectedErr: ErrInvalidFormat,
		},
		{
			name: "unknown validator",
			in: struct {
				ID string `validate:"length:4"`
			}{},
			expectedErr: ErrUnknownValidator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestValidateValidStructs(t *testing.T) {
	tests := []struct {
		in any
	}{
		{
			in: User{
				ID:     "qqqqqqqqqqwwwwwwwwwweeeeeeeeee123456",
				Name:   "Ivan",
				Age:    18,
				Email:  "ivan@mail.com",
				Role:   "stuff",
				Phones: []string{"01234567890", "12345678901"},
			},
		},
		{
			in: App{
				Version: "abcde",
			},
		},
		{
			in: Token{
				Header:    []byte("qwerty"),
				Payload:   []byte("test"),
				Signature: nil,
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "qwerty",
			},
		},
		{
			in: "abcde",
		},
		{
			in: 4,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

type (
	MyString string
	MyInt    int
	Nested   struct {
		ID int `validate:"min:1|max:10|in:1,3,5,7,9"`
	}
)

type Container struct {
	Int      int      `validate:"min:1|max:10|in:1,3,5,7,9"`
	Int8     int8     `validate:"min:1|max:10|in:1,3,5,7,9"`
	Int16    int16    `validate:"min:1|max:10|in:1,3,5,7,9"`
	Int32    int32    `validate:"min:1|max:10|in:1,3,5,7,9"`
	Int64    int64    `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uint     uint     `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uint8    uint8    `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uint16   uint16   `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uint32   uint32   `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uint64   uint64   `validate:"min:1|max:10|in:1,3,5,7,9"`
	Uintptr  uintptr  `validate:"min:1|max:10|in:1,3,5,7,9"`
	Float32  float32  `validate:"min:1|max:10|in:1,3,5,7,9"`
	Float64  float64  `validate:"min:1|max:10|in:1,3,5,7,9"`
	String   string   `validate:"len:6|regexp:\\d+|in:qwerty,abcdef,toster"`
	AnySlice []string `validate:"len:4|regexp:\\w+"`
	MyString MyString `validate:"len:4|regexp:\\d+|in:1234,5678"`
	MyInt    MyInt    `validate:"min:1|max:10|in:3,5,7"`
	MyStruct Nested   `validate:"nested"`
}

func TestValidateNotValidValues(t *testing.T) {
	t.Run("struct with not valid values", func(t *testing.T) {
		t.Parallel()
		container := Container{
			Int:      11,            // 2 ошибки
			Int8:     0,             // 2 ошибки
			Int16:    8,             // 1 ошибка
			Int32:    -1,            // 2 ошибки
			Int64:    2,             // 1 ошибка
			Uint:     11,            // 2 ошибки
			Uint8:    0,             // 2 ошибки
			Uint16:   8,             // 1 ошибка
			Uint32:   11,            // 2 ошибки
			Uint64:   11,            // 2 ошибки
			Uintptr:  11,            // 2 ошибки
			Float32:  11,            // 2 ошибки
			Float64:  -0.12345,      // 2 ошибки
			String:   "n",           // 3 ошибки
			AnySlice: []string{":"}, // 2 ошибки
			MyString: "a",           // 3 ошибки
			MyInt:    0,             // 2 ошибки
			MyStruct: Nested{
				ID: 11, // 2 ошибки
			},
		} // итого 35 ошибок

		err := Validate(container)

		var validationErrors ValidationErrors
		require.ErrorAs(t, err, &validationErrors)

		require.Len(t, validationErrors, 35)
	})
}
