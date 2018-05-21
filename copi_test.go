package copi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/danhtran94/copi"
	strfmt "github.com/go-openapi/strfmt"
)

type C struct {
	C1 string
	C2 string
}

type D struct {
	D1 string  `copi:"C1"`
	D2 *string `copi:"C2"`
}

type A struct {
	F1 string
	S1 []string
	C  []C
	T  time.Time
}

type B struct {
	F1 **string
	F2 string    `copi:"F1"`
	S2 *[]string `copi:"S1"`
	D  []D       `copi:"C"`
	E  []*D      `copi:"C"`
	T  strfmt.DateTime
}

type EmA struct {
	A
}

type EmB struct {
	B
}

func TestDup(t *testing.T) {
	copi.Debugging()

	s := ""
	ss := &s
	sss := &ss
	copi.Dup("Danh", &sss)
	fmt.Println(s)

	a := A{
		F1: "Danh",
		S1: []string{
			"Danh",
		},
		C: []C{
			{
				C1: "OK",
				C2: "NO",
			},
		},
		T: time.Now(),
	}
	b := B{}
	copi.Dup(a, &b)
	fmt.Println(**b.F1, b.F2, b.S2, b.D, *b.D[0].D2, *b.E[0], b.T)

	ea := EmA{
		a,
	}
	eb := EmB{}
	copi.Dup(ea, &eb)
	fmt.Println(eb.F1, eb.F2, eb.S2, eb.D, eb.E, eb.T)

	ma := map[int64]int{
		3: 1,
		5: 2,
		6: 3,
	}

	mb := map[int]int{}
	copi.Dup(ma, &mb)
	fmt.Println(mb)

	maa := map[string]A{
		"A": a,
	}

	mbb := map[string]B{}

	copi.Dup(maa, &mbb)
	fmt.Println(mbb)
}
