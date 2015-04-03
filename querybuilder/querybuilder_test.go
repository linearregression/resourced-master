package querybuilder

import (
	"testing"
)

func TestParseTags(t *testing.T) {
	toBeTested := []string{
		`Tags = ["aaa","bbb","ccc"]`,
		`Tags=["aaa","bbb","ccc"]`,
		`tags = ["aaa","bbb","ccc"]`,
	}

	for _, testString := range toBeTested {
		output := Parse(testString)
		if output != `tags ?& array["aaa","bbb","ccc"]` {
			t.Errorf("Failed to generate tags query. Output: %v", output)
		}
	}
}

func TestParseNameExact(t *testing.T) {
	toBeTested := []string{
		`Name = "Awesome Sauce"`,
		`Name="Awesome Sauce"`,
		`name = "Awesome Sauce"`,
	}

	for _, testString := range toBeTested {
		output := Parse(testString)
		if output != `name = "Awesome Sauce"` {
			t.Errorf("Failed to generate name query. Output: %v", output)
		}
	}
}

func TestParseNameStartsWith(t *testing.T) {
	toBeTested := []string{
		`Name ~^ "brotato"`,
		`Name~^"brotato"`,
		`name ~^ "brotato"`,
	}

	for _, testString := range toBeTested {
		output := Parse(testString)
		if output != `name LIKE "brotato%"` {
			t.Errorf("Failed to generate name query. Output: %v", output)
		}
	}
}

func TestParseAnd(t *testing.T) {
	toBeTested := `tags = ["aaa","bbb","ccc"] AND Name~^"brotato"`

	output := Parse(toBeTested)
	if output != `tags ?& array["aaa","bbb","ccc"] AND name LIKE "brotato%"` {
		t.Errorf("Failed to generate name query. Output: %v", output)
	}
}