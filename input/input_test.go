package input

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input     string
		expectErr bool
		expectObj DataItem
	}{
		{
			input:     `<di><peruna>potato</peruna><porkkana>carrot</porkkana></di>`,
			expectObj: DataItem{Format: FmtXML, Peruna: "potato", Porkkana: "carrot"},
			expectErr: false,
		},
		{
			input:     `{"peruna": "potato", "porkkana": "carrot"}`,
			expectErr: false,
			expectObj: DataItem{Format: FmtJSON, Peruna: "potato", Porkkana: "carrot"},
		},
		{
			input:     ``,
			expectErr: true,
		},
		{
			input:     `{}`,
			expectErr: true,
		},
	}

	for i, c := range cases {
		t.Logf("Input test case '%s'", c.input)
		tmp, err := os.CreateTemp("", fmt.Sprintf("sample_input_test_case_%d", i))
		if err != nil {
			t.Fatalf("Could not create temp file due to %s", err)
		}
		tmp.WriteString(c.input)
		tmp.Close()
		di, err := Parse(c.expectObj.Format, tmp.Name())
		if err != nil {
			if !c.expectErr {
				t.Fatalf("Expected no errors in case, but got %s", err)
			}
		}
		if err == nil && c.expectErr {
			t.Fatalf("Expected case to return error, but got %s instead", di)
		}
		// structs are comparable if all their fields are comparable: https://go.dev/ref/spec#Comparison_operators
		if !c.expectErr && di.String() != c.expectObj.String() {
			t.Fatalf("Case failed (%s) != %s", di, c.expectObj)
		}
		os.Remove(tmp.Name())
	}
}
