package tokenizer

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {
	input := `
"stuff\"" 'and\'' «things"'«
$xpath-with-dash @_underscored_xpath $xpath2
1234 23_45 34-56 12.34 2_3.4_5 3-4.5-6
ident-with-dash _underscored_ident ident2
"&<>'" '&<>"'
"\n"
你好世界
ident·with·interpunkt
ʻokina
`
	input = strings.ReplaceAll(input, "«", "`")

	expectTokens := []string{
		`DQ("stuff\"")`,
		`SQ("and'")`,
		"BQ(\"things\\\"'\")",
		`XP("$xpath-with-dash")`,
		`XP("@_underscored_xpath")`,
		`XP("$xpath2")`,
		`NUM("1234")`,
		`NUM("2345")`,
		`NUM("3456")`,
		`NUM("12.34")`,
		`NUM("23.45")`,
		`NUM("34.56")`,
		`IDENT("ident-with-dash")`,
		`IDENT("_underscored_ident")`,
		`IDENT("ident2")`,
		`DQ("&<>'")`,
		`SQ("&<>\"")`,
		`DQ("\n")`,
		`IDENT("你好世界")`,
		`IDENT("ident·with·interpunkt")`,
		`IDENT("ʻokina")`,
	}

	r := &Reader{
		S: bufio.NewScanner(strings.NewReader(input)),
	}

	for i, expect := range expectTokens {
		got, err := r.ReadToken()
		if err != nil {
			t.Fatalf("token %d %s: unexpected error: %v", i, got, err)
		}

		if got.String() != expect {
			t.Errorf("token %d was %s, but expected %s", i, got, expect)
		}
	}

	got, err := r.ReadToken()
	if err != nil && err != io.EOF {
		t.Fatalf("final ReadToken gave error %q, but expected io.EOF", err)
	}
	if expect := `EOF("")`; got.String() != expect {
		t.Errorf("final ReadToken was %s, but expected %s", got, expect)
	}
}
