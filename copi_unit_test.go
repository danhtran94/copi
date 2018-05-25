package copi_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danhtran94/copi"
)

func TestBasicType(t *testing.T) {
	t.Run("Simple basic types", func(t *testing.T) {
		// t.Skip()
		ts := []struct {
			Source interface{}
			Dest   interface{}
			Expect interface{}
		}{
			{
				Source: 10,
				Dest:   0,
				Expect: 10,
			},
			{
				Source: "copi sample",
				Dest:   "",
				Expect: "copi sample",
			},
			{
				Source: true,
				Dest:   false,
				Expect: true,
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, &tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})

	t.Run("Simple basic types, dest is poiter", func(t *testing.T) {
		// t.Skip()
		ts := []struct {
			Source interface{}
			Dest   interface{}
			Expect interface{}
		}{
			{
				Source: 10,
				Dest:   new(int),
				Expect: 10,
			},
			{
				Source: "copi sample",
				Dest:   new(string),
				Expect: "copi sample",
			},
			{
				Source: true,
				Dest:   new(bool),
				Expect: true,
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, reflect.ValueOf(tt.Expect).Interface(), reflect.Indirect(reflect.ValueOf(tt.Dest)).Interface())
		}
	})

	t.Run("Simple basic types, dest is nil", func(t *testing.T) {
		// t.Skip()
		ts := []struct {
			Source interface{}
			Dest   interface{}
			Expect interface{}
		}{
			{
				Source: 10,
				Dest:   nil,
				Expect: nil,
			},
			{
				Source: "copi sample",
				Dest:   nil,
				Expect: nil,
			},
			{
				Source: true,
				Dest:   nil,
				Expect: nil,
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expect, tt.Dest)
		}
	})
}

func TestStruct(t *testing.T) {
	copi.Debugging()
	t.Run("Simple struct with basic type", func(t *testing.T) {
		// t.Skip()
		type Source struct {
			Num   int
			Text  string
			Quest bool
			Point float32
		}
		type Dest struct {
			Num   int
			Text  string
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
					Num:   10,
					Text:  "sample",
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
			Num   *int
			Text  string
			Quest bool
			Point float32
		}
		type Dest struct {
			Num   int
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
					Num:   nil,
					Text:  "sample",
					Quest: *ptrQuest,
					Point: *ptrPoint,
				},
				Dest: Dest{
					Dummy: "foobar",
				},
				Expect: Dest{
					Num:   0,
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
}
