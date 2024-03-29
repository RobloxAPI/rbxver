package rbxver

import (
	"io"
	"testing"
)

// Tests for ParseBytes and Parse.
var tests = []struct {
	s   string   // Input string.
	f   Format   // Input format.
	v   Version  // Expected version.
	n   int      // Expected read bytes.
	e   error    // Expected error.
	str *Version // If Parse, only compare this.
}{
	{s: "", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: io.ErrUnexpectedEOF},
	{s: "", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: io.ErrUnexpectedEOF},
	{s: "", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: io.ErrUnexpectedEOF},
	{s: " ", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " ", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " ", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "a", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "a", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "a", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "ab", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "ab", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "ab", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ".", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ".", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ".", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ",", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ",", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: ",", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "1", f: Any, v: Version{1, 0, 0, 0, Any}, n: 1, e: io.ErrUnexpectedEOF},
	{s: "1", f: Dot, v: Version{1, 0, 0, 0, Any}, n: 1, e: io.ErrUnexpectedEOF},
	{s: "1", f: Comma, v: Version{1, 0, 0, 0, Any}, n: 1, e: io.ErrUnexpectedEOF},
	{s: " 1", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " 1", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " 1", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "12", f: Any, v: Version{12, 0, 0, 0, Any}, n: 2, e: io.ErrUnexpectedEOF},
	{s: "12", f: Dot, v: Version{12, 0, 0, 0, Any}, n: 2, e: io.ErrUnexpectedEOF},
	{s: "12", f: Comma, v: Version{12, 0, 0, 0, Any}, n: 2, e: io.ErrUnexpectedEOF},
	{s: "12.34.56.78", f: Any, v: Version{12, 34, 56, 78, Dot}, n: 11, e: nil},
	{s: "12.34.56.78", f: Dot, v: Version{12, 34, 56, 78, Dot}, n: 11, e: nil},
	{s: "12.34.56.78", f: Comma, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12.34,56.78", f: Any, v: Version{12, 34, 0, 0, Any}, n: 5, e: ErrSyntax},
	{s: "12.34,56.78", f: Dot, v: Version{12, 34, 0, 0, Any}, n: 5, e: ErrSyntax},
	{s: "12.34,56.78", f: Comma, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12,34,56,78", f: Any, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12,34,56,78", f: Dot, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12,34,56,78", f: Comma, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12, 34, 56, 78", f: Any, v: Version{12, 34, 56, 78, Comma}, n: 14, e: nil},
	{s: "12, 34, 56, 78", f: Dot, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "12, 34, 56, 78", f: Comma, v: Version{12, 34, 56, 78, Comma}, n: 14, e: nil},
	{s: " 12 . 34 . 56 . 78 ", f: Any, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " 12 . 34 . 56 . 78 ", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: " 12 . 34 . 56 . 78 ", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 0, e: ErrSyntax},
	{s: "12.34. 56 . 78 ", f: Any, v: Version{12, 34, 0, 0, Any}, n: 6, e: ErrSyntax},
	{s: "12.34. 56 . 78 ", f: Dot, v: Version{12, 34, 0, 0, Any}, n: 6, e: ErrSyntax},
	{s: "12.34. 56 . 78 ", f: Comma, v: Version{12, 0, 0, 0, Any}, n: 2, e: ErrSyntax},
	{s: "0.123.1.1234567", f: Any, v: Version{0, 123, 1, 1234567, Dot}, n: 15, e: nil},
	{s: "0.123.1.1234567", f: Dot, v: Version{0, 123, 1, 1234567, Dot}, n: 15, e: nil},
	{s: "0.123.1.1234567", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 1, e: ErrSyntax},
	{s: "0 . 123 .1. 1234567", f: Any, v: Version{0, 0, 0, 0, Any}, n: 1, e: ErrSyntax},
	{s: "0 . 123 .1. 1234567", f: Dot, v: Version{0, 0, 0, 0, Any}, n: 1, e: ErrSyntax},
	{s: "0 . 123 .1. 1234567", f: Comma, v: Version{0, 0, 0, 0, Any}, n: 1, e: ErrSyntax},
	{s: "0.123.-1.1234567", f: Any, v: Version{0, 123, 0, 0, Any}, n: 6, e: ErrSyntax},
	{s: "0.123.1.1234567trailingdata", f: Dot, v: Version{0, 123, 1, 1234567, Dot}, n: 15, e: nil, str: &Version{}},
}

var fmtstr = [...]string{
	"Any",
	"Dot",
	"Comma",
}

func TestParseBytes(t *testing.T) {
	for _, test := range tests {
		v, n, err := ParseBytes([]byte(test.s), test.f)
		if v != test.v {
			t.Errorf("ParseBytes(%q, %s): expected version %v, got %v", test.s, fmtstr[test.f], test.v, v)
		}
		if n != test.n {
			t.Errorf("ParseBytes(%q, %s): expected bytes %d, got %d", test.s, fmtstr[test.f], test.n, n)
		}
		if err != test.e {
			t.Errorf("ParseBytes(%q, %s): expected error %v, got %v", test.s, fmtstr[test.f], test.e, err)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		v := Parse(test.s, test.f)
		if test.str != nil {
			if v != *test.str {
				t.Errorf("Parse(%q, %s): expected version %v, got %v", test.s, fmtstr[test.f], *test.str, v)
			}
		} else {
			if test.e == nil {
				if v != test.v {
					t.Errorf("Parse(%q, %s): expected version %v, got %v", test.s, fmtstr[test.f], test.v, v)
				}
			} else {
				if v != (Version{}) {
					t.Errorf("Parse(%q, %s): expected zero version, got %v", test.s, fmtstr[test.f], v)
				}
			}
		}
	}
}
