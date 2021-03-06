package vm

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ndau/ndaumath/pkg/address"
)

// MiniAsm is a miniature assembler that has very simple syntax. It's primarily intended for writing
// simple test code.
//
// It takes a single string as input. It then converts everything to lower case and splits it into
// 'words' by whitespace -- all whitespace is equivalent.
//
// If a word matches an opcode, it generates the associated opcode.
// If a word matches a simple pattern for a timestamp, it attempts to parse it as a timestamp.
// All other words are expected to be one-byte hex values.
//
// Any failure in parsing causes a panic; there is no error recovery.
//
// No attempt is made to ensure that opcode parameters or types are correct, and each opcode is individually
// specified (Push1, Push2, etc).
//

// MiniAsm is a mini-assembler for Chaincode. It is primarily for testing but we want it available.
// nolint: deadcode
func MiniAsm(s string) Chaincode {
	// whitespace
	wsp := regexp.MustCompile("[ \t\r\n]|;.+\n")
	// timestamp
	tsp := regexp.MustCompile("[0-9-]+T[0-9:]+Z")
	// address is 48 chars starting with nd and not containing io10
	addrp := regexp.MustCompile("nd[2-9a-km-np-zA-KM-NP-Z]{46}")
	// quoted string without spaces (this is a mini assembler!)
	qsp := regexp.MustCompile(`"[^" ]+"`)
	words := wsp.Split(strings.TrimSpace(s), -1)
	opcodes := []Opcode{}
	for _, w := range words {
		// skip empty words
		if w == "" {
			continue
		}
		// see if it's an opcode
		if op, ok := opcodeMap[strings.ToLower(w)]; ok {
			opcodes = append(opcodes, op)
			continue
		}
		// see if it's a timestamp
		if tsp.MatchString(strings.ToUpper(w)) {
			t, err := ParseTimestamp(strings.ToUpper(w))
			if err != nil {
				panic(err)
			}
			bytes := ToBytes(int64(t.t))
			for _, byt := range bytes {
				opcodes = append(opcodes, Opcode(byt))
			}
			continue
		}
		// see if it's an address
		if addrp.MatchString(w) {
			_, err := address.Validate(w)
			if err != nil {
				panic(err)
			}
			bytes := []byte(w)
			opcodes = append(opcodes, Opcode(len(bytes)))
			for _, byt := range bytes {
				opcodes = append(opcodes, Opcode(byt))
			}
			continue
		}
		// see if it's a quoted string
		if qsp.MatchString(w) {
			bytes := w[1 : len(w)-1]
			opcodes = append(opcodes, Opcode(len(bytes)))
			for _, byt := range bytes {
				opcodes = append(opcodes, Opcode(byt))
			}
			continue
		}
		// otherwise it should be a hex value
		b, err := strconv.ParseUint(w, 16, 8)
		if err != nil {
			panic(err)
		}
		opcodes = append(opcodes, Opcode(b))
	}
	return Chaincode(opcodes)
}

// MiniAsmSafe is a mini-assembler for chaincode which does not panic on errors.
func MiniAsmSafe(s string) (c Chaincode, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch p := r.(type) {
			case error:
				err = p
			case fmt.Stringer:
				err = errors.New(p.String())
			default:
				err = fmt.Errorf("%v", p)
			}
		}
	}()
	c = MiniAsm(s)
	return
}
