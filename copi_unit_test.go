package copi_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danhtran94/copi"
)

func Test(t *testing.T) {
	t.Run("Simple basic types", func(t *testing.T) {
		t.Skip()
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
		t.Skip()
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
		ts := []struct {
			Source interface{}
			Dest   interface{}
			Expect interface{}
		}{
			{
				Source: 10,
				Dest:   nil,
				Expect: 10,
			},
			{
				Source: "copi sample",
				Dest:   nil,
				Expect: "copi sample",
			},
			{
				Source: true,
				Dest:   nil,
				Expect: true,
			},
		}
		for _, tt := range ts {
			err := copi.Dup(tt.Source, tt.Dest)
			assert.NoError(t, err)
			assert.Equal(t, reflect.ValueOf(tt.Expect).Interface(), reflect.Indirect(reflect.ValueOf(tt.Dest)).Interface())
		}
	})
}
