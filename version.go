// The rbxver package handles parsing and formatting of Roblox version strings.
package rbxver

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Format determines how a version is parsed and formatted.
type Format int

const (
	// Parse by guessing separator. Format as `0.0.0.0`.
	Any Format = iota
	// Parse with dot as separator. Format as `0.0.0.0`.
	Dot
	// Parse with comma as separator. Format as `0, 0, 0, 0`.
	Comma
)

// Version represents the version of a Roblox build. Versions can be compared
// for equality.
type Version struct {
	Generation int // The first component.
	Version    int // The second component.
	Patch      int // The third component.
	Commit     int // The fourth component.
}

// Formats i, writing to b. Writes 0 if i is less than 0.
func formatInt(b *strings.Builder, i int) {
	if i <= 0 {
		b.WriteByte('0')
		return
	}
	b.Write(strconv.AppendInt(nil, int64(i), 10))
}

// Format formats v according to f.
//
// Panics if f is not valid format.
func (v Version) Format(f Format) string {
	var sep string
	switch f {
	case Any, Dot:
		sep = "."
	case Comma:
		sep = ", "
	default:
		panic("invalid format")
	}
	var b strings.Builder
	formatInt(&b, v.Generation)
	b.WriteString(sep)
	formatInt(&b, v.Version)
	b.WriteString(sep)
	formatInt(&b, v.Patch)
	b.WriteString(sep)
	formatInt(&b, v.Commit)
	return b.String()
}

// String returns v as a string in the default format (dot).
func (v Version) String() string {
	return v.Format(Any)
}

// Less returns true if v is semantically lower than u, and false otherwise.
func (v Version) Less(u Version) bool {
	if v.Generation < u.Generation {
		return true
	}
	if v.Version < u.Version {
		return true
	}
	if v.Patch < u.Patch {
		return true
	}
	if v.Commit < u.Commit {
		return true
	}
	return false
}

// Parses an integer from b to comp. Returns false if an error occurred when
// parsing the integer, or the value is less than 0. b is set to the index after
// the parsed value.
func parseInt(comp *int, b *[]byte) bool {
	i := 0
	for ; len(*b) > i && '0' <= (*b)[i] && (*b)[i] <= '9'; i++ {
	}
	n, err := strconv.ParseInt(string((*b)[:i]), 10, strconv.IntSize)
	if err != nil || n < 0 {
		return false
	}
	*comp = int(n)
	*b = (*b)[i:]
	return true
}

// Expects sep at the start of b. If *sep is nil, then the separator will be
// guessed, and sep is set to the guessed separator. b is set to the index after
// the parsed separator and any whitespace.
func parseSep(sep *[]byte, b *[]byte) error {
	if len(*b) == 0 {
		return io.ErrUnexpectedEOF
	}
	if *sep == nil {
		// Guess separator. This will be used for subsequent separators.
		switch (*b)[0] {
		case '.':
			*sep = (*b)[:1]
		case ',':
			if len(*b) < 2 || (*b)[1] != ' ' {
				return ErrSyntax
			}
			*sep = (*b)[:2]
		default:
			return ErrSyntax
		}
	} else {
		if len(*b) < len(*sep) {
			return io.ErrUnexpectedEOF
		}
		if !bytes.Equal((*b)[:len(*sep)], *sep) {
			return ErrSyntax
		}
	}
	*b = (*b)[len(*sep):]
	return nil
}

// ErrSyntax indicates a syntax error while parsing a version string.
var ErrSyntax = errors.New("invalid syntax")

// Parse parses a version from b according to f. Leading and trailing whitespace
// is ignored, as well as whitespace between components. err will be ErrSyntax
// if the syntax is invalid, or io.ErrUnexpectedEOF if b does not have enough
// bytes to correctly parse the version.
//
// Panics if f is not valid format.
func Parse(b []byte, f Format) (v Version, n int, err error) {
	var sep []byte
	switch f {
	case Any:
	case Dot:
		sep = []byte{'.'}
	case Comma:
		sep = []byte{',', ' '}
	default:
		panic("invalid format")
	}

	l := len(b)
	if len(b) == 0 {
		return v, l - len(b), io.ErrUnexpectedEOF
	}
	if !parseInt(&v.Generation, &b) {
		return v, l - len(b), ErrSyntax
	}
	if err := parseSep(&sep, &b); err != nil {
		return v, l - len(b), err
	}
	if !parseInt(&v.Version, &b) {
		return v, l - len(b), ErrSyntax
	}
	if err := parseSep(&sep, &b); err != nil {
		return v, l - len(b), err
	}
	if !parseInt(&v.Patch, &b) {
		return v, l - len(b), ErrSyntax
	}
	if err := parseSep(&sep, &b); err != nil {
		return v, l - len(b), err
	}
	if !parseInt(&v.Commit, &b) {
		return v, l - len(b), ErrSyntax
	}
	return v, l - len(b), nil
}

// Parses a version from s according to f. Returns the zero version if a version
// could not be parsed.
func ParseString(s string, f Format) Version {
	if v, n, err := Parse([]byte(s), f); err == nil && n == len(s) {
		return v
	}
	return Version{}
}
