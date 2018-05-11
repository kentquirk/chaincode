package main

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// bcheck evaluates a stream of bytes against a string of hex
func bcheck(t *testing.T, b []byte, s string) {
	shouldbe := []byte{}
	// convert s to array of bytes
	// first, strip all whitespace
	p := regexp.MustCompile("[ \t\r\n]")
	s = p.ReplaceAllString(s, "")

	assert.Equal(t, len(b), len(s)/2, "expected string is wrong length")
	// then split it up into bytes
	for i := 0; i < len(s); i += 2 {
		sb := s[i : i+2]
		sv, err := strconv.ParseInt(sb, 16, 16)
		assert.Nil(t, err)
		shouldbe = append(shouldbe, byte(sv))
	}
	assert.Equal(t, b, shouldbe)
}

// checkParse makes sure that the result of a parse is a given stream of bytes
func checkParse(t *testing.T, name string, code string, result string) {
	sn, err := Parse(name, []byte(code))
	if err != nil {
		fmt.Println(describeErrors(err, code))
	}
	assert.Nil(t, err)
	b := sn.(Script).bytes()
	bcheck(t, b, result)
}