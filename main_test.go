package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itchyny/gojq"
	"github.com/stretchr/testify/assert"
)

func readTestData(filename string) string {
	fi, err := os.Open("testdata/" + filename)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func Test_parseJQ(t *testing.T) {
	type args struct {
		inputString string
		query       *gojq.Query
	}
	tests := []struct {
		name           string
		args           args
		wantPrettified string
		assertion      assert.ErrorAssertionFunc
	}{
		{
			name: "simple",
			args: args{
				inputString: readTestData("simple.json"),
				query: func() *gojq.Query {
					query, _ := gojq.Parse(`.glossary.GlossDiv.title`)
					return query
				}(),
			},
			wantPrettified: "\x1b[92m\"S\"\x1b[0m\n",
			assertion:      assert.NoError,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotPrettified, err := parseJQ(tt.args.inputString, tt.args.query)
		tt.assertion(t, err, fmt.Sprintf("%q. parseJQ()", tt.name))
		assert.Equalf(t, tt.wantPrettified, gotPrettified, "%q. parseJQ()", tt.name)
	}
}
