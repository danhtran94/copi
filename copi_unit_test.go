package copi_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/danhtran94/copi"
)

// Type definition for component schema "Question"
type QuestionRes struct {
	ID          *int64              `json:"ID,omitempty"`
	CreatedAt   *time.Time          `json:"createdAt,omitempty"`
	Description *string             `json:"description,omitempty"`
	FootNote    *string             `json:"footNote,omitempty"`
	FormID      int64               `json:"formID"`
	Max         *int32              `json:"max,omitempty"`
	MaxLabel    *string             `json:"maxLabel,omitempty"`
	MaxLength   *int32              `json:"maxLength,omitempty"`
	Min         *int32              `json:"min,omitempty"`
	MinLabel    *string             `json:"minLabel,omitempty"`
	Options     []QuestionOptionRes `json:"options,omitempty"`
	PlaceHolder *string             `json:"placeHolder,omitempty"`
	Required    *bool               `json:"required,omitempty"`
	Sequence    *int32              `json:"sequence,omitempty"`
	Step        *int32              `json:"step,omitempty"`
	Title       string              `json:"title"`
	Type        string              `json:"type"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
}

type QuestionType string

// Question for form, each question belong to only one form
type Question struct {
	ID          int64             `json:"ID,omitempty" db:"id,omitempty"`
	FormID      int64             `json:"formID" db:"form_id"`
	Title       string            `json:"title" db:"title,omitempty"`
	Description string            `json:"description" db:"description"`
	Type        QuestionType      `json:"questionType" db:"type"`
	Min         int               `json:"min" db:"min"`
	MinLabel    string            `json:"minLabel" db:"min_label"`
	Max         int               `json:"max" db:"max"`
	MaxLabel    string            `json:"maxLabel" db:"max_label"`
	Step        int               `json:"step" db:"step"`
	Required    bool              `json:"required" db:"required"`
	Sequence    int               `json:"sequence" db:"sequence"`
	MaxLength   int               `json:"maxLength" db:"max_length"`
	PlaceHolder string            `json:"placeHolder" db:"place_holder"`
	FootNote    string            `json:"footNote" db:"foot_note"`
	CreatedAt   time.Time         `json:"createdAt" db:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updatedAt" db:"updated_at,omitempty"`
	Options     []*QuestionOption `json:"options" db:"-"`
}

type QuestionOption struct {
	ID         int64  `json:"ID" db:"id,omitempty"`
	Label      string `json:"label" db:"label"`
	Value      string `json:"value" db:"value"`
	QuestionID int64  `json:"questionID" db:"question_id"`
	Sequence   int    `json:"sequence" db:"sequence"`
}

// Type definition for component schema "QuestionOption"
type QuestionOptionRes struct {
	ID         *int64  `json:"ID,omitempty"`
	Label      *string `json:"label,omitempty"`
	QuestionID *int64  `json:"questionID,omitempty"`
	Sequence   *int64  `json:"sequence,omitempty"`
	Value      string  `json:"value"`
}

func TestBasicType(t *testing.T) {
	t.Run("Simple basic types", func(t *testing.T) {
		// t.Skip()
		ts := []struct {
			Source []*Question
			Dest   []QuestionRes
			Expect []QuestionRes
		}{
			// {
			// 	Source: 10,
			// 	Dest:   0,
			// 	Expect: 10,
			// },
			// {
			// 	Source: "copi sample",
			// 	Dest:   "",
			// 	Expect: "copi sample",
			// },
			// {
			// 	Source: []int{},
			// 	Dest:   []int{},
			// 	Expect: []int{},
			// },
			{
				Source: []*Question{&Question{
					// Options: []*QuestionOption{},
				}},
				Dest: []QuestionRes{},
				Expect: []QuestionRes{{
					Options: nil,
				}},
			},
			// {
			// 	Source: true,
			// 	Dest:   false,
			// 	Expect: true,
			// },
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			// fmt.Printf("s1 (addr: %p): %+8v\n", tt.Dest[0].Options, tt.Expect[0].Options)
			fmt.Printf("s1 (addr: %p): %+8v\n", &tt.Dest[0].Options, *(*reflect.SliceHeader)(unsafe.Pointer(&tt.Dest[0].Options)))
			fmt.Printf("s2 (addr: %p): %+8v\n", &tt.Expect[0].Options, *(*reflect.SliceHeader)(unsafe.Pointer(&tt.Expect[0].Options)))
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect[0].Options, tt.Dest[0].Options)
		}
	})

	// t.Run("Simple basic types, dest is poiter", func(t *testing.T) {
	// 	// t.Skip()
	// 	ts := []struct {
	// 		Source interface{}
	// 		Dest   interface{}
	// 		Expect interface{}
	// 	}{
	// 		{
	// 			Source: 10,
	// 			Dest:   new(int),
	// 			Expect: 10,
	// 		},
	// 		{
	// 			Source: "copi sample",
	// 			Dest:   new(string),
	// 			Expect: "copi sample",
	// 		},
	// 		{
	// 			Source: true,
	// 			Dest:   new(bool),
	// 			Expect: true,
	// 		},
	// 	}
	// 	for _, tt := range ts {
	// 		err := copi.Dup(tt.Source, tt.Dest)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, reflect.ValueOf(tt.Expect).Interface(), reflect.Indirect(reflect.ValueOf(tt.Dest)).Interface())
	// 	}
	// })

	// t.Run("Simple basic types, dest is nil", func(t *testing.T) {
	// 	// t.Skip()
	// 	ts := []struct {
	// 		Source interface{}
	// 		Dest   interface{}
	// 		Expect interface{}
	// 	}{
	// 		{
	// 			Source: 10,
	// 			Dest:   nil,
	// 			Expect: nil,
	// 		},
	// 		{
	// 			Source: "copi sample",
	// 			Dest:   nil,
	// 			Expect: nil,
	// 		},
	// 		{
	// 			Source: true,
	// 			Dest:   nil,
	// 			Expect: nil,
	// 		},
	// 	}
	// 	for _, tt := range ts {
	// 		err := copi.Dup(tt.Source, tt.Dest)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, tt.Expect, tt.Dest)
	// 	}
	// })
}

func TestStruct(t *testing.T) {
	// copi.Debugging()
	t.Run("Simple struct with basic type", func(t *testing.T) {
		// t.Skip()
		type Nest struct {
			Message string
		}

		type Source struct {
			Num   int
			Text  string
			Nest  Nest
			Quest bool
			Point float32
		}

		type Dest struct {
			Num   int
			Text  string
			Nest  Nest
			Quest bool
			Point float32
			Dummy string
		}

		ts := []struct {
			Source Source
			Dest   Dest
			Expect Dest
		}{
			{
				Source: Source{
					Num:  10,
					Text: "sample",
					Nest: Nest{
						Message: "hello",
					},
					Quest: true,
					Point: 3.14,
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:   10,
					Text:  "sample",
					Quest: true,
					Nest: Nest{
						Message: "hello",
					},
					Point: 3.14,
					Dummy: "foobar",
				},
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

	t.Run("Simple struct with basic type dest has pointer", func(t *testing.T) {
		// t.Skip()
		ptrNum := new(int)
		*ptrNum = 10
		ptrText := new(string)
		*ptrText = "sample"
		ptrQuest := new(bool)
		*ptrQuest = true
		ptrPoint := new(float32)
		*ptrPoint = 3.14
		type Source struct {
			Num   int
			Text  string
			Quest bool
			Point float32
		}
		type Dest struct {
			Num   *int
			Text  *string
			Quest *bool
			Point *float32
			Dummy string
		}
		ts := []struct {
			Source Source
			Dest   Dest
			Expect Dest
		}{
			{
				Source: Source{
					Num:   *ptrNum,
					Text:  *ptrText,
					Quest: *ptrQuest,
					Point: *ptrPoint,
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:   ptrNum,
					Text:  ptrText,
					Quest: ptrQuest,
					Point: ptrPoint,
					Dummy: "foobar",
				},
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

	t.Run("Simple struct with basic type dest has interface{}", func(t *testing.T) {
		// t.Skip()
		ptrNum := new(int)
		*ptrNum = 10
		ptrQuest := new(bool)
		*ptrQuest = true
		ptrPoint := new(float32)
		*ptrPoint = 3.14
		type Source struct {
			Num   int
			Text  string
			Quest bool
			Point float32
		}
		type Dest struct {
			Num   *int
			Text  interface{}
			Quest *bool
			Point *float32
			Dummy string
		}
		ts := []struct {
			Source Source
			Dest   Dest
			Expect Dest
		}{
			{
				Source: Source{
					Num:   *ptrNum,
					Text:  "sample",
					Quest: *ptrQuest,
					Point: *ptrPoint,
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:   ptrNum,
					Text:  "sample",
					Quest: ptrQuest,
					Point: ptrPoint,
					Dummy: "foobar",
				},
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

	t.Run("Simple struct with basic type src has nil value", func(t *testing.T) {
		ptrNum := new(int)
		*ptrNum = 10
		ptrQuest := new(bool)
		*ptrQuest = true
		ptrPoint := new(float32)
		*ptrPoint = 3.14

		type Source struct {
			Num     *int
			Text    string
			Quest   bool
			Point   float32
			AnyTag1 string `copi-to:"AnyTag2"`
		}
		type Dest struct {
			Num     int
			Text    interface{}
			Quest   *bool
			Point   *float32
			Dummy   string
			AnyTag2 string
		}
		ts := []struct {
			Source Source
			Dest   Dest
			Expect Dest
		}{
			{
				Source: Source{
					Num:     nil,
					Text:    "sample",
					Quest:   *ptrQuest,
					Point:   *ptrPoint,
					AnyTag1: "abcabc",
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:     0,
					Text:    "sample",
					Quest:   ptrQuest,
					Point:   ptrPoint,
					Dummy:   "foobar",
					AnyTag2: "abcabc",
				},
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

	t.Run("Simple struct with basic type src has nil value to nil", func(t *testing.T) {
		ptrNum := new(time.Time)
		*ptrNum = time.Now()
		ptrQuest := new(bool)
		*ptrQuest = true
		ptrPoint := new(float32)
		*ptrPoint = 3.14

		type Source struct {
			Num     *time.Time
			IntNum  *interface{}
			Text    string
			Quest   bool
			Point   float32
			AnyTag1 string `copi-to:"AnyTag2"`
		}
		type Dest struct {
			Num     *time.Time
			IntNum  *int
			Text    interface{}
			Quest   *bool
			Point   *float32
			Dummy   string
			AnyTag2 string
		}
		ts := []struct {
			Source Source
			Dest   Dest
			Expect Dest
		}{
			{
				Source: Source{
					Num:     nil,
					Text:    "sample",
					Quest:   *ptrQuest,
					Point:   *ptrPoint,
					AnyTag1: "abcabc",
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:     nil,
					Text:    "sample",
					Quest:   ptrQuest,
					Point:   ptrPoint,
					Dummy:   "foobar",
					AnyTag2: "abcabc",
				},
			},
		}
		for _, tt := range ts {
			// copi.Debugging()
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

}
