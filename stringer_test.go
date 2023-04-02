package antidig_test

import (
	"math/rand"
	"testing"

	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/internal/digtest"
	"github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	type A struct{}
	type B struct{}
	type C struct{}
	type D struct{}

	type in struct {
		dig.In

		A A `name:"foo"`
		B B `optional:"true"`
		C C `name:"bar" optional:"true"`

		Strings []string `group:"baz"`
	}

	type out struct {
		dig.Out

		A A `name:"foo"`
		C C `name:"bar"`
	}

	type stringOut struct {
		dig.Out

		S string `group:"baz"`
	}

	c := digtest.New(t, dig.SetRand(rand.New(rand.NewSource(0))))

	c.RequireProvide(func(i in) D {
		assert.Equal(t, []string{"bar", "baz", "foo"}, i.Strings)
		return D{}
	})

	c.RequireProvide(func() out {
		return out{
			A: A{},
			C: C{},
		}
	})

	c.RequireProvide(func() A { return A{} })
	c.RequireProvide(func() B { return B{} })
	c.RequireProvide(func() C { return C{} })

	c.RequireProvide(func(A) stringOut { return stringOut{S: "foo"} })
	c.RequireProvide(func(B) stringOut { return stringOut{S: "bar"} })
	c.RequireProvide(func(C) stringOut { return stringOut{S: "baz"} })

	c.RequireInvoke(func(D) {
	})

	s := c.String()

	// All nodes
	assert.Contains(t, s, `dig_test.A[name="foo"] -> deps: []`)
	assert.Contains(t, s, "dig_test.A -> deps: []")
	assert.Contains(t, s, "dig_test.B -> deps: []")
	assert.Contains(t, s, "dig_test.C -> deps: []")
	assert.Contains(t, s, `dig_test.C[name="bar"] -> deps: []`)
	assert.Contains(t, s, `dig_test.D -> deps: [dig_test.A[name="foo"] dig_test.B[optional] dig_test.C[optional, name="bar"] string[group="baz"]]`)
	assert.Contains(t, s, `string[group="baz"] -> deps: [dig_test.A]`)
	assert.Contains(t, s, `string[group="baz"] -> deps: [dig_test.B]`)
	assert.Contains(t, s, `string[group="baz"] -> deps: [dig_test.C]`)

	// Values
	assert.Contains(t, s, "dig_test.A => {}")
	assert.Contains(t, s, "dig_test.B => {}")
	assert.Contains(t, s, "dig_test.C => {}")
	assert.Contains(t, s, "dig_test.D => {}")
	assert.Contains(t, s, `dig_test.A[name="foo"] => {}`)
	assert.Contains(t, s, `dig_test.C[name="bar"] => {}`)
	assert.Contains(t, s, `string[group="baz"] => foo`)
	assert.Contains(t, s, `string[group="baz"] => bar`)
	assert.Contains(t, s, `string[group="baz"] => baz`)
}
