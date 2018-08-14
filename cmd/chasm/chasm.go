package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Script",
			pos:  position{line: 6, col: 1, offset: 23},
			expr: &actionExpr{
				pos: position{line: 6, col: 11, offset: 33},
				run: (*parser).callonScript1,
				expr: &seqExpr{
					pos: position{line: 6, col: 11, offset: 33},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 6, col: 11, offset: 33},
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 11, offset: 33},
								name: "EOL",
							},
						},
						&labeledExpr{
							pos:   position{line: 6, col: 16, offset: 38},
							label: "code",
							expr: &ruleRefExpr{
								pos:  position{line: 6, col: 21, offset: 43},
								name: "Code",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 6, col: 26, offset: 48},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 8, col: 1, offset: 125},
			expr: &actionExpr{
				pos: position{line: 8, col: 10, offset: 134},
				run: (*parser).callonCode1,
				expr: &seqExpr{
					pos: position{line: 8, col: 10, offset: 134},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 8, col: 10, offset: 134},
							expr: &ruleRefExpr{
								pos:  position{line: 8, col: 10, offset: 134},
								name: "EOL",
							},
						},
						&labeledExpr{
							pos:   position{line: 8, col: 15, offset: 139},
							label: "rs",
							expr: &oneOrMoreExpr{
								pos: position{line: 8, col: 18, offset: 142},
								expr: &ruleRefExpr{
									pos:  position{line: 8, col: 18, offset: 142},
									name: "RoutineDef",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "RoutineDef",
			pos:  position{line: 10, col: 1, offset: 174},
			expr: &choiceExpr{
				pos: position{line: 11, col: 4, offset: 191},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 11, col: 4, offset: 191},
						name: "HandlerDef",
					},
					&ruleRefExpr{
						pos:  position{line: 12, col: 4, offset: 205},
						name: "FunctionDef",
					},
				},
			},
		},
		{
			name: "HandlerDef",
			pos:  position{line: 15, col: 1, offset: 221},
			expr: &actionExpr{
				pos: position{line: 15, col: 15, offset: 235},
				run: (*parser).callonHandlerDef1,
				expr: &seqExpr{
					pos: position{line: 15, col: 15, offset: 235},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 15, col: 15, offset: 235},
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 15, offset: 235},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 15, col: 18, offset: 238},
							val:        "handler",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 28, offset: 248},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 15, col: 30, offset: 250},
							label: "ids",
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 34, offset: 254},
								name: "HandlerIDList",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 15, col: 48, offset: 268},
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 48, offset: 268},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 15, col: 51, offset: 271},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 15, col: 55, offset: 275},
							label: "s",
							expr: &oneOrMoreExpr{
								pos: position{line: 15, col: 57, offset: 277},
								expr: &ruleRefExpr{
									pos:  position{line: 15, col: 57, offset: 277},
									name: "Line",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 15, col: 63, offset: 283},
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 63, offset: 283},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 15, col: 66, offset: 286},
							val:        "}",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 15, col: 70, offset: 290},
							expr: &ruleRefExpr{
								pos:  position{line: 15, col: 70, offset: 290},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionDef",
			pos:  position{line: 19, col: 1, offset: 395},
			expr: &actionExpr{
				pos: position{line: 19, col: 16, offset: 410},
				run: (*parser).callonFunctionDef1,
				expr: &seqExpr{
					pos: position{line: 19, col: 16, offset: 410},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 19, col: 16, offset: 410},
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 16, offset: 410},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 19, col: 19, offset: 413},
							val:        "func",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 19, col: 26, offset: 420},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 19, col: 28, offset: 422},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 30, offset: 424},
								name: "FunctionName",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 19, col: 43, offset: 437},
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 43, offset: 437},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 19, col: 46, offset: 440},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 19, col: 50, offset: 444},
							label: "s",
							expr: &oneOrMoreExpr{
								pos: position{line: 19, col: 52, offset: 446},
								expr: &ruleRefExpr{
									pos:  position{line: 19, col: 52, offset: 446},
									name: "Line",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 19, col: 58, offset: 452},
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 58, offset: 452},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 19, col: 61, offset: 455},
							val:        "}",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 19, col: 65, offset: 459},
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 65, offset: 459},
								name: "EOL",
							},
						},
					},
				},
			},
		},
		{
			name: "HandlerIDList",
			pos:  position{line: 29, col: 1, offset: 736},
			expr: &choiceExpr{
				pos: position{line: 30, col: 4, offset: 756},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 30, col: 4, offset: 756},
						run: (*parser).callonHandlerIDList2,
						expr: &seqExpr{
							pos: position{line: 30, col: 4, offset: 756},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 30, col: 4, offset: 756},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 30, col: 6, offset: 758},
										name: "Value",
									},
								},
								&litMatcher{
									pos:        position{line: 30, col: 12, offset: 764},
									val:        ",",
									ignoreCase: false,
								},
								&zeroOrOneExpr{
									pos: position{line: 30, col: 16, offset: 768},
									expr: &ruleRefExpr{
										pos:  position{line: 30, col: 16, offset: 768},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 30, col: 19, offset: 771},
									label: "h",
									expr: &ruleRefExpr{
										pos:  position{line: 30, col: 21, offset: 773},
										name: "HandlerIDList",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 31, col: 4, offset: 839},
						run: (*parser).callonHandlerIDList11,
						expr: &labeledExpr{
							pos:   position{line: 31, col: 4, offset: 839},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 31, col: 6, offset: 841},
								name: "Value",
							},
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 34, col: 1, offset: 898},
			expr: &choiceExpr{
				pos: position{line: 35, col: 7, offset: 912},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 35, col: 7, offset: 912},
						run: (*parser).callonLine2,
						expr: &seqExpr{
							pos: position{line: 35, col: 7, offset: 912},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 35, col: 7, offset: 912},
									expr: &ruleRefExpr{
										pos:  position{line: 35, col: 7, offset: 912},
										name: "_",
									},
								},
								&labeledExpr{
									pos:   position{line: 35, col: 10, offset: 915},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 35, col: 13, offset: 918},
										name: "Operation",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 35, col: 23, offset: 928},
									name: "EOL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 36, col: 7, offset: 957},
						run: (*parser).callonLine9,
						expr: &ruleRefExpr{
							pos:  position{line: 36, col: 7, offset: 957},
							name: "EOL",
						},
					},
				},
			},
		},
		{
			name: "Operation",
			pos:  position{line: 39, col: 1, offset: 988},
			expr: &choiceExpr{
				pos: position{line: 40, col: 7, offset: 1007},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 40, col: 7, offset: 1007},
						name: "ConstDef",
					},
					&ruleRefExpr{
						pos:  position{line: 41, col: 7, offset: 1022},
						name: "Opcode",
					},
				},
			},
		},
		{
			name: "ConstDef",
			pos:  position{line: 44, col: 1, offset: 1036},
			expr: &actionExpr{
				pos: position{line: 45, col: 7, offset: 1054},
				run: (*parser).callonConstDef1,
				expr: &seqExpr{
					pos: position{line: 45, col: 7, offset: 1054},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 45, col: 7, offset: 1054},
							label: "k",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 9, offset: 1056},
								name: "Constant",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 45, col: 18, offset: 1065},
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 18, offset: 1065},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 45, col: 21, offset: 1068},
							val:        "=",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 45, col: 25, offset: 1072},
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 25, offset: 1072},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 45, col: 28, offset: 1075},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 30, offset: 1077},
								name: "Value",
							},
						},
					},
				},
			},
		},
		{
			name: "Opcode",
			pos:  position{line: 59, col: 1, offset: 1582},
			expr: &choiceExpr{
				pos: position{line: 60, col: 7, offset: 1597},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 60, col: 7, offset: 1597},
						run: (*parser).callonOpcode2,
						expr: &litMatcher{
							pos:        position{line: 60, col: 7, offset: 1597},
							val:        "nop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 4, offset: 1714},
						run: (*parser).callonOpcode4,
						expr: &litMatcher{
							pos:        position{line: 63, col: 4, offset: 1714},
							val:        "zero",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 64, col: 4, offset: 1763},
						run: (*parser).callonOpcode6,
						expr: &seqExpr{
							pos: position{line: 64, col: 4, offset: 1763},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 64, col: 4, offset: 1763},
									val:        "wchoice",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 64, col: 14, offset: 1773},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 64, col: 16, offset: 1775},
									label: "ix",
									expr: &ruleRefExpr{
										pos:  position{line: 64, col: 19, offset: 1778},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 4, offset: 1841},
						run: (*parser).callonOpcode12,
						expr: &seqExpr{
							pos: position{line: 65, col: 4, offset: 1841},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 65, col: 4, offset: 1841},
									val:        "tuck",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 65, col: 11, offset: 1848},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 65, col: 13, offset: 1850},
									label: "offset",
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 20, offset: 1857},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 66, col: 4, offset: 1921},
						run: (*parser).callonOpcode18,
						expr: &litMatcher{
							pos:        position{line: 66, col: 4, offset: 1921},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 67, col: 4, offset: 1970},
						run: (*parser).callonOpcode20,
						expr: &litMatcher{
							pos:        position{line: 67, col: 4, offset: 1970},
							val:        "swap",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 68, col: 4, offset: 2019},
						run: (*parser).callonOpcode22,
						expr: &litMatcher{
							pos:        position{line: 68, col: 4, offset: 2019},
							val:        "sum",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 69, col: 4, offset: 2066},
						run: (*parser).callonOpcode24,
						expr: &litMatcher{
							pos:        position{line: 69, col: 4, offset: 2066},
							val:        "sub",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 70, col: 4, offset: 2113},
						run: (*parser).callonOpcode26,
						expr: &seqExpr{
							pos: position{line: 70, col: 4, offset: 2113},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 70, col: 4, offset: 2113},
									val:        "sort",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 70, col: 11, offset: 2120},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 70, col: 13, offset: 2122},
									label: "ix",
									expr: &ruleRefExpr{
										pos:  position{line: 70, col: 16, offset: 2125},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 71, col: 4, offset: 2185},
						run: (*parser).callonOpcode32,
						expr: &litMatcher{
							pos:        position{line: 71, col: 4, offset: 2185},
							val:        "slice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 72, col: 4, offset: 2236},
						run: (*parser).callonOpcode34,
						expr: &seqExpr{
							pos: position{line: 72, col: 4, offset: 2236},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 72, col: 4, offset: 2236},
									val:        "roll",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 72, col: 11, offset: 2243},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 72, col: 13, offset: 2245},
									label: "offset",
									expr: &ruleRefExpr{
										pos:  position{line: 72, col: 20, offset: 2252},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 73, col: 4, offset: 2316},
						run: (*parser).callonOpcode40,
						expr: &litMatcher{
							pos:        position{line: 73, col: 4, offset: 2316},
							val:        "ret",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 74, col: 4, offset: 2363},
						run: (*parser).callonOpcode42,
						expr: &litMatcher{
							pos:        position{line: 74, col: 4, offset: 2363},
							val:        "rand",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 75, col: 4, offset: 2412},
						run: (*parser).callonOpcode44,
						expr: &seqExpr{
							pos: position{line: 75, col: 4, offset: 2412},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 75, col: 4, offset: 2412},
									val:        "pusht",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 75, col: 12, offset: 2420},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 75, col: 14, offset: 2422},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 75, col: 16, offset: 2424},
										name: "Timestamp",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 76, col: 4, offset: 2477},
						run: (*parser).callonOpcode50,
						expr: &litMatcher{
							pos:        position{line: 76, col: 4, offset: 2477},
							val:        "pushl",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 77, col: 4, offset: 2528},
						run: (*parser).callonOpcode52,
						expr: &seqExpr{
							pos: position{line: 77, col: 4, offset: 2528},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 77, col: 4, offset: 2528},
									val:        "pushb",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 77, col: 12, offset: 2536},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 77, col: 14, offset: 2538},
									label: "ba",
									expr: &ruleRefExpr{
										pos:  position{line: 77, col: 17, offset: 2541},
										name: "Bytes",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 78, col: 4, offset: 2574},
						run: (*parser).callonOpcode58,
						expr: &seqExpr{
							pos: position{line: 78, col: 4, offset: 2574},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 78, col: 4, offset: 2574},
									val:        "pusha",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 78, col: 12, offset: 2582},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 78, col: 14, offset: 2584},
									label: "a",
									expr: &ruleRefExpr{
										pos:  position{line: 78, col: 16, offset: 2586},
										name: "Address",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 79, col: 4, offset: 2632},
						run: (*parser).callonOpcode64,
						expr: &seqExpr{
							pos: position{line: 79, col: 4, offset: 2632},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 79, col: 4, offset: 2632},
									val:        "pick",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 79, col: 11, offset: 2639},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 79, col: 13, offset: 2641},
									label: "offset",
									expr: &ruleRefExpr{
										pos:  position{line: 79, col: 20, offset: 2648},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 80, col: 4, offset: 2712},
						run: (*parser).callonOpcode70,
						expr: &litMatcher{
							pos:        position{line: 80, col: 4, offset: 2712},
							val:        "over",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 81, col: 4, offset: 2761},
						run: (*parser).callonOpcode72,
						expr: &litMatcher{
							pos:        position{line: 81, col: 4, offset: 2761},
							val:        "one",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 82, col: 4, offset: 2808},
						run: (*parser).callonOpcode74,
						expr: &litMatcher{
							pos:        position{line: 82, col: 4, offset: 2808},
							val:        "now",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 83, col: 4, offset: 2855},
						run: (*parser).callonOpcode76,
						expr: &litMatcher{
							pos:        position{line: 83, col: 4, offset: 2855},
							val:        "not",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 4, offset: 2902},
						run: (*parser).callonOpcode78,
						expr: &litMatcher{
							pos:        position{line: 84, col: 4, offset: 2902},
							val:        "neg1",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 85, col: 4, offset: 2951},
						run: (*parser).callonOpcode80,
						expr: &litMatcher{
							pos:        position{line: 85, col: 4, offset: 2951},
							val:        "neg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 86, col: 4, offset: 2998},
						run: (*parser).callonOpcode82,
						expr: &litMatcher{
							pos:        position{line: 86, col: 4, offset: 2998},
							val:        "muldiv",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 87, col: 4, offset: 3051},
						run: (*parser).callonOpcode84,
						expr: &litMatcher{
							pos:        position{line: 87, col: 4, offset: 3051},
							val:        "mul",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 88, col: 4, offset: 3098},
						run: (*parser).callonOpcode86,
						expr: &litMatcher{
							pos:        position{line: 88, col: 4, offset: 3098},
							val:        "mod",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 89, col: 4, offset: 3145},
						run: (*parser).callonOpcode88,
						expr: &litMatcher{
							pos:        position{line: 89, col: 4, offset: 3145},
							val:        "min",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 90, col: 4, offset: 3192},
						run: (*parser).callonOpcode90,
						expr: &litMatcher{
							pos:        position{line: 90, col: 4, offset: 3192},
							val:        "max",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 91, col: 4, offset: 3239},
						run: (*parser).callonOpcode92,
						expr: &litMatcher{
							pos:        position{line: 91, col: 4, offset: 3239},
							val:        "lt",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 92, col: 7, offset: 3287},
						run: (*parser).callonOpcode94,
						expr: &seqExpr{
							pos: position{line: 92, col: 7, offset: 3287},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 92, col: 7, offset: 3287},
									val:        "lookup",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 92, col: 16, offset: 3296},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 92, col: 18, offset: 3298},
									label: "id",
									expr: &ruleRefExpr{
										pos:  position{line: 92, col: 21, offset: 3301},
										name: "FunctionName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 92, col: 34, offset: 3314},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 92, col: 36, offset: 3316},
									label: "count",
									expr: &ruleRefExpr{
										pos:  position{line: 92, col: 42, offset: 3322},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 93, col: 4, offset: 3398},
						run: (*parser).callonOpcode103,
						expr: &litMatcher{
							pos:        position{line: 93, col: 4, offset: 3398},
							val:        "len",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 94, col: 4, offset: 3445},
						run: (*parser).callonOpcode105,
						expr: &litMatcher{
							pos:        position{line: 94, col: 4, offset: 3445},
							val:        "index",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 95, col: 4, offset: 3496},
						run: (*parser).callonOpcode107,
						expr: &litMatcher{
							pos:        position{line: 95, col: 4, offset: 3496},
							val:        "inc",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 96, col: 4, offset: 3543},
						run: (*parser).callonOpcode109,
						expr: &litMatcher{
							pos:        position{line: 96, col: 4, offset: 3543},
							val:        "ifz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 97, col: 4, offset: 3590},
						run: (*parser).callonOpcode111,
						expr: &litMatcher{
							pos:        position{line: 97, col: 4, offset: 3590},
							val:        "ifnz",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 98, col: 4, offset: 3639},
						run: (*parser).callonOpcode113,
						expr: &litMatcher{
							pos:        position{line: 98, col: 4, offset: 3639},
							val:        "gt",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 99, col: 4, offset: 3684},
						run: (*parser).callonOpcode115,
						expr: &seqExpr{
							pos: position{line: 99, col: 4, offset: 3684},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 99, col: 4, offset: 3684},
									val:        "fieldl",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 99, col: 13, offset: 3693},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 99, col: 15, offset: 3695},
									label: "ix",
									expr: &ruleRefExpr{
										pos:  position{line: 99, col: 18, offset: 3698},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 100, col: 4, offset: 3760},
						run: (*parser).callonOpcode121,
						expr: &seqExpr{
							pos: position{line: 100, col: 4, offset: 3760},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 100, col: 4, offset: 3760},
									val:        "field",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 100, col: 12, offset: 3768},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 100, col: 14, offset: 3770},
									label: "ix",
									expr: &ruleRefExpr{
										pos:  position{line: 100, col: 17, offset: 3773},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 101, col: 4, offset: 3834},
						run: (*parser).callonOpcode127,
						expr: &litMatcher{
							pos:        position{line: 101, col: 4, offset: 3834},
							val:        "false",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 102, col: 4, offset: 3885},
						run: (*parser).callonOpcode129,
						expr: &litMatcher{
							pos:        position{line: 102, col: 4, offset: 3885},
							val:        "fail",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 103, col: 4, offset: 3934},
						run: (*parser).callonOpcode131,
						expr: &litMatcher{
							pos:        position{line: 103, col: 4, offset: 3934},
							val:        "extend",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 104, col: 4, offset: 3987},
						run: (*parser).callonOpcode133,
						expr: &litMatcher{
							pos:        position{line: 104, col: 4, offset: 3987},
							val:        "eq",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 105, col: 4, offset: 4032},
						run: (*parser).callonOpcode135,
						expr: &litMatcher{
							pos:        position{line: 105, col: 4, offset: 4032},
							val:        "endif",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 106, col: 4, offset: 4083},
						run: (*parser).callonOpcode137,
						expr: &litMatcher{
							pos:        position{line: 106, col: 4, offset: 4083},
							val:        "else",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 107, col: 4, offset: 4132},
						run: (*parser).callonOpcode139,
						expr: &litMatcher{
							pos:        position{line: 107, col: 4, offset: 4132},
							val:        "dup2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 108, col: 4, offset: 4181},
						run: (*parser).callonOpcode141,
						expr: &litMatcher{
							pos:        position{line: 108, col: 4, offset: 4181},
							val:        "dup",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 109, col: 4, offset: 4228},
						run: (*parser).callonOpcode143,
						expr: &litMatcher{
							pos:        position{line: 109, col: 4, offset: 4228},
							val:        "drop2",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 110, col: 4, offset: 4279},
						run: (*parser).callonOpcode145,
						expr: &litMatcher{
							pos:        position{line: 110, col: 4, offset: 4279},
							val:        "drop",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 111, col: 4, offset: 4328},
						run: (*parser).callonOpcode147,
						expr: &litMatcher{
							pos:        position{line: 111, col: 4, offset: 4328},
							val:        "divmod",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 112, col: 4, offset: 4381},
						run: (*parser).callonOpcode149,
						expr: &litMatcher{
							pos:        position{line: 112, col: 4, offset: 4381},
							val:        "div",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 113, col: 7, offset: 4431},
						run: (*parser).callonOpcode151,
						expr: &seqExpr{
							pos: position{line: 113, col: 7, offset: 4431},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 113, col: 7, offset: 4431},
									val:        "deco",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 113, col: 14, offset: 4438},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 113, col: 16, offset: 4440},
									label: "id",
									expr: &ruleRefExpr{
										pos:  position{line: 113, col: 19, offset: 4443},
										name: "FunctionName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 113, col: 32, offset: 4456},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 113, col: 34, offset: 4458},
									label: "count",
									expr: &ruleRefExpr{
										pos:  position{line: 113, col: 40, offset: 4464},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 114, col: 4, offset: 4538},
						run: (*parser).callonOpcode160,
						expr: &litMatcher{
							pos:        position{line: 114, col: 4, offset: 4538},
							val:        "dec",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 115, col: 4, offset: 4585},
						run: (*parser).callonOpcode162,
						expr: &litMatcher{
							pos:        position{line: 115, col: 4, offset: 4585},
							val:        "choice",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 116, col: 7, offset: 4641},
						run: (*parser).callonOpcode164,
						expr: &seqExpr{
							pos: position{line: 116, col: 7, offset: 4641},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 116, col: 7, offset: 4641},
									val:        "call",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 116, col: 14, offset: 4648},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 116, col: 16, offset: 4650},
									label: "id",
									expr: &ruleRefExpr{
										pos:  position{line: 116, col: 19, offset: 4653},
										name: "FunctionName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 116, col: 32, offset: 4666},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 116, col: 34, offset: 4668},
									label: "count",
									expr: &ruleRefExpr{
										pos:  position{line: 116, col: 40, offset: 4674},
										name: "Value",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 117, col: 4, offset: 4748},
						run: (*parser).callonOpcode173,
						expr: &litMatcher{
							pos:        position{line: 117, col: 4, offset: 4748},
							val:        "avg",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 118, col: 4, offset: 4795},
						run: (*parser).callonOpcode175,
						expr: &litMatcher{
							pos:        position{line: 118, col: 4, offset: 4795},
							val:        "append",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 119, col: 4, offset: 4848},
						run: (*parser).callonOpcode177,
						expr: &litMatcher{
							pos:        position{line: 119, col: 4, offset: 4848},
							val:        "add",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 122, col: 7, offset: 5083},
						run: (*parser).callonOpcode179,
						expr: &seqExpr{
							pos: position{line: 122, col: 7, offset: 5083},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 122, col: 7, offset: 5083},
									val:        "push",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 122, col: 14, offset: 5090},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 122, col: 16, offset: 5092},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 122, col: 18, offset: 5094},
										name: "Value",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Timestamp",
			pos:  position{line: 125, col: 1, offset: 5145},
			expr: &actionExpr{
				pos: position{line: 125, col: 14, offset: 5158},
				run: (*parser).callonTimestamp1,
				expr: &seqExpr{
					pos: position{line: 125, col: 14, offset: 5158},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 125, col: 14, offset: 5158},
							name: "Date",
						},
						&litMatcher{
							pos:        position{line: 125, col: 19, offset: 5163},
							val:        "T",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 125, col: 23, offset: 5167},
							name: "Time",
						},
						&litMatcher{
							pos:        position{line: 125, col: 28, offset: 5172},
							val:        "Z",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Date",
			pos:  position{line: 126, col: 1, offset: 5207},
			expr: &seqExpr{
				pos: position{line: 126, col: 9, offset: 5215},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 126, col: 9, offset: 5215},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 15, offset: 5221},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 21, offset: 5227},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 27, offset: 5233},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 126, col: 33, offset: 5239},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 37, offset: 5243},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 43, offset: 5249},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 126, col: 49, offset: 5255},
						val:        "-",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 53, offset: 5259},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 126, col: 59, offset: 5265},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Time",
			pos:  position{line: 127, col: 1, offset: 5271},
			expr: &seqExpr{
				pos: position{line: 127, col: 10, offset: 5280},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 127, col: 10, offset: 5280},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 16, offset: 5286},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 127, col: 22, offset: 5292},
						val:        ":",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 26, offset: 5296},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 32, offset: 5302},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&litMatcher{
						pos:        position{line: 127, col: 38, offset: 5308},
						val:        ":",
						ignoreCase: false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 42, offset: 5312},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 48, offset: 5318},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 129, col: 1, offset: 5325},
			expr: &choiceExpr{
				pos: position{line: 130, col: 7, offset: 5339},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 130, col: 7, offset: 5339},
						name: "Integer",
					},
					&ruleRefExpr{
						pos:  position{line: 131, col: 7, offset: 5353},
						name: "ConstantRef",
					},
				},
			},
		},
		{
			name: "ConstantRef",
			pos:  position{line: 134, col: 1, offset: 5372},
			expr: &actionExpr{
				pos: position{line: 134, col: 16, offset: 5387},
				run: (*parser).callonConstantRef1,
				expr: &labeledExpr{
					pos:   position{line: 134, col: 16, offset: 5387},
					label: "k",
					expr: &ruleRefExpr{
						pos:  position{line: 134, col: 18, offset: 5389},
						name: "Constant",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 135, col: 1, offset: 5507},
			expr: &choiceExpr{
				pos: position{line: 136, col: 7, offset: 5524},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 136, col: 7, offset: 5524},
						run: (*parser).callonInteger2,
						expr: &seqExpr{
							pos: position{line: 136, col: 7, offset: 5524},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 136, col: 7, offset: 5524},
									expr: &ruleRefExpr{
										pos:  position{line: 136, col: 7, offset: 5524},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 136, col: 10, offset: 5527},
									val:        "0x",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 136, col: 15, offset: 5532},
									expr: &charClassMatcher{
										pos:        position{line: 136, col: 15, offset: 5532},
										val:        "[0-9A-Fa-f]",
										ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 137, col: 7, offset: 5624},
						run: (*parser).callonInteger9,
						expr: &seqExpr{
							pos: position{line: 137, col: 7, offset: 5624},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 137, col: 7, offset: 5624},
									expr: &ruleRefExpr{
										pos:  position{line: 137, col: 7, offset: 5624},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 137, col: 10, offset: 5627},
									val:        "addr(",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 137, col: 18, offset: 5635},
									expr: &seqExpr{
										pos: position{line: 137, col: 19, offset: 5636},
										exprs: []interface{}{
											&charClassMatcher{
												pos:        position{line: 137, col: 19, offset: 5636},
												val:        "[0-9A-Fa-f]",
												ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
												ignoreCase: false,
												inverted:   false,
											},
											&charClassMatcher{
												pos:        position{line: 137, col: 30, offset: 5647},
												val:        "[0-9A-Fa-f]",
												ranges:     []rune{'0', '9', 'A', 'F', 'a', 'f'},
												ignoreCase: false,
												inverted:   false,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 137, col: 44, offset: 5661},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 138, col: 7, offset: 5724},
						run: (*parser).callonInteger19,
						expr: &seqExpr{
							pos: position{line: 138, col: 7, offset: 5724},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 138, col: 7, offset: 5724},
									expr: &ruleRefExpr{
										pos:  position{line: 138, col: 7, offset: 5724},
										name: "_",
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 138, col: 10, offset: 5727},
									expr: &litMatcher{
										pos:        position{line: 138, col: 10, offset: 5727},
										val:        "-",
										ignoreCase: false,
									},
								},
								&oneOrMoreExpr{
									pos: position{line: 138, col: 15, offset: 5732},
									expr: &charClassMatcher{
										pos:        position{line: 138, col: 15, offset: 5732},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Bytes",
			pos:  position{line: 140, col: 1, offset: 5824},
			expr: &choiceExpr{
				pos: position{line: 141, col: 7, offset: 5839},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 141, col: 7, offset: 5839},
						run: (*parser).callonBytes2,
						expr: &labeledExpr{
							pos:   position{line: 141, col: 7, offset: 5839},
							label: "b",
							expr: &oneOrMoreExpr{
								pos: position{line: 141, col: 9, offset: 5841},
								expr: &ruleRefExpr{
									pos:  position{line: 141, col: 9, offset: 5841},
									name: "Integer",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 142, col: 7, offset: 5907},
						run: (*parser).callonBytes6,
						expr: &seqExpr{
							pos: position{line: 142, col: 7, offset: 5907},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 142, col: 7, offset: 5907},
									expr: &ruleRefExpr{
										pos:  position{line: 142, col: 7, offset: 5907},
										name: "_",
									},
								},
								&litMatcher{
									pos:        position{line: 142, col: 10, offset: 5910},
									val:        "\"",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 142, col: 14, offset: 5914},
									label: "s",
									expr: &oneOrMoreExpr{
										pos: position{line: 142, col: 16, offset: 5916},
										expr: &charClassMatcher{
											pos:        position{line: 142, col: 16, offset: 5916},
											val:        "[^\"]",
											chars:      []rune{'"'},
											ignoreCase: false,
											inverted:   true,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 142, col: 22, offset: 5922},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Address",
			pos:  position{line: 145, col: 1, offset: 5980},
			expr: &actionExpr{
				pos: position{line: 145, col: 12, offset: 5991},
				run: (*parser).callonAddress1,
				expr: &seqExpr{
					pos: position{line: 145, col: 12, offset: 5991},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 145, col: 12, offset: 5991},
							val:        "nd",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 145, col: 17, offset: 5996},
							expr: &charClassMatcher{
								pos:        position{line: 145, col: 17, offset: 5996},
								val:        "[2-9a-km-np-zA-KM-NP-Z]",
								ranges:     []rune{'2', '9', 'a', 'k', 'm', 'n', 'p', 'z', 'A', 'K', 'M', 'N', 'P', 'Z'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Constant",
			pos:  position{line: 147, col: 1, offset: 6062},
			expr: &actionExpr{
				pos: position{line: 147, col: 13, offset: 6074},
				run: (*parser).callonConstant1,
				expr: &seqExpr{
					pos: position{line: 147, col: 13, offset: 6074},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 147, col: 13, offset: 6074},
							val:        "[A-Za-z]",
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 147, col: 22, offset: 6083},
							expr: &charClassMatcher{
								pos:        position{line: 147, col: 22, offset: 6083},
								val:        "[A-Za-z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionName",
			pos:  position{line: 148, col: 1, offset: 6143},
			expr: &actionExpr{
				pos: position{line: 148, col: 17, offset: 6159},
				run: (*parser).callonFunctionName1,
				expr: &seqExpr{
					pos: position{line: 148, col: 17, offset: 6159},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 148, col: 17, offset: 6159},
							val:        "[A-Za-z]",
							ranges:     []rune{'A', 'Z', 'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
						&oneOrMoreExpr{
							pos: position{line: 148, col: 26, offset: 6168},
							expr: &charClassMatcher{
								pos:        position{line: 148, col: 26, offset: 6168},
								val:        "[A-Za-z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'A', 'Z', 'a', 'z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 151, col: 1, offset: 6226},
			expr: &oneOrMoreExpr{
				pos: position{line: 151, col: 6, offset: 6231},
				expr: &charClassMatcher{
					pos:        position{line: 151, col: 6, offset: 6231},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 153, col: 1, offset: 6239},
			expr: &seqExpr{
				pos: position{line: 153, col: 8, offset: 6246},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 153, col: 8, offset: 6246},
						expr: &ruleRefExpr{
							pos:  position{line: 153, col: 8, offset: 6246},
							name: "_",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 153, col: 11, offset: 6249},
						expr: &ruleRefExpr{
							pos:  position{line: 153, col: 11, offset: 6249},
							name: "Comment",
						},
					},
					&choiceExpr{
						pos: position{line: 153, col: 21, offset: 6259},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 153, col: 21, offset: 6259},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 153, col: 30, offset: 6268},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 153, col: 39, offset: 6277},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 153, col: 46, offset: 6284},
								val:        "\n",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 155, col: 1, offset: 6292},
			expr: &seqExpr{
				pos: position{line: 155, col: 12, offset: 6303},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 155, col: 12, offset: 6303},
						val:        ";",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 155, col: 16, offset: 6307},
						expr: &charClassMatcher{
							pos:        position{line: 155, col: 16, offset: 6307},
							val:        "[^\\r\\n]",
							chars:      []rune{'\r', '\n'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 157, col: 1, offset: 6317},
			expr: &notExpr{
				pos: position{line: 157, col: 8, offset: 6324},
				expr: &anyMatcher{
					line: 157, col: 9, offset: 6325,
				},
			},
		},
	},
}

func (c *current) onScript1(code interface{}) (interface{}, error) {
	return newScript(code, c.globalStore["functions"].(map[string]int))
}

func (p *parser) callonScript1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onScript1(stack["code"])
}

func (c *current) onCode1(rs interface{}) (interface{}, error) {
	return rs, nil
}

func (p *parser) callonCode1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCode1(stack["rs"])
}

func (c *current) onHandlerDef1(ids, s interface{}) (interface{}, error) {
	return newHandlerDef(ids.([]string), s, c.globalStore["constants"].(map[string]string))

}

func (p *parser) callonHandlerDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHandlerDef1(stack["ids"], stack["s"])
}

func (c *current) onFunctionDef1(n, s interface{}) (interface{}, error) {
	fm := c.globalStore["functions"].(map[string]int)
	name := n.(string)
	ctr := c.globalStore["functionCounter"].(int)
	fm[name] = ctr
	ctr++
	c.globalStore["functionCounter"] = ctr
	return newFunctionDef(name, s)

}

func (p *parser) callonFunctionDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionDef1(stack["n"], stack["s"])
}

func (c *current) onHandlerIDList2(v, h interface{}) (interface{}, error) {
	return append(h.([]string), v.(string)), nil
}

func (p *parser) callonHandlerIDList2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHandlerIDList2(stack["v"], stack["h"])
}

func (c *current) onHandlerIDList11(v interface{}) (interface{}, error) {
	return []string{string(c.text)}, nil
}

func (p *parser) callonHandlerIDList11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHandlerIDList11(stack["v"])
}

func (c *current) onLine2(op interface{}) (interface{}, error) {
	return op, nil
}

func (p *parser) callonLine2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine2(stack["op"])
}

func (c *current) onLine9() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonLine9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine9()
}

func (c *current) onConstDef1(k, v interface{}) (interface{}, error) {
	if key, ok := k.(string); ok {
		cm := c.globalStore["constants"].(map[string]string)
		cm[key] = v.(string)
		return v, nil
	}
	return nil, errors.New("Bad const def")

}

func (p *parser) callonConstDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstDef1(stack["k"], stack["v"])
}

func (c *current) onOpcode2() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNop)
}

func (p *parser) callonOpcode2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode2()
}

func (c *current) onOpcode4() (interface{}, error) {
	return newUnitaryOpcode(vm.OpZero)
}

func (p *parser) callonOpcode4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode4()
}

func (c *current) onOpcode6(ix interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpWChoice, ix.(string))
}

func (p *parser) callonOpcode6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode6(stack["ix"])
}

func (c *current) onOpcode12(offset interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpTuck, offset.(string))
}

func (p *parser) callonOpcode12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode12(stack["offset"])
}

func (c *current) onOpcode18() (interface{}, error) {
	return newUnitaryOpcode(vm.OpTrue)
}

func (p *parser) callonOpcode18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode18()
}

func (c *current) onOpcode20() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSwap)
}

func (p *parser) callonOpcode20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode20()
}

func (c *current) onOpcode22() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSum)
}

func (p *parser) callonOpcode22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode22()
}

func (c *current) onOpcode24() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSub)
}

func (p *parser) callonOpcode24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode24()
}

func (c *current) onOpcode26(ix interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpSort, ix.(string))
}

func (p *parser) callonOpcode26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode26(stack["ix"])
}

func (c *current) onOpcode32() (interface{}, error) {
	return newUnitaryOpcode(vm.OpSlice)
}

func (p *parser) callonOpcode32() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode32()
}

func (c *current) onOpcode34(offset interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpRoll, offset.(string))
}

func (p *parser) callonOpcode34() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode34(stack["offset"])
}

func (c *current) onOpcode40() (interface{}, error) {
	return newUnitaryOpcode(vm.OpRet)
}

func (p *parser) callonOpcode40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode40()
}

func (c *current) onOpcode42() (interface{}, error) {
	return newUnitaryOpcode(vm.OpRand)
}

func (p *parser) callonOpcode42() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode42()
}

func (c *current) onOpcode44(t interface{}) (interface{}, error) {
	return newPushTimestamp(t.(string))
}

func (p *parser) callonOpcode44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode44(stack["t"])
}

func (c *current) onOpcode50() (interface{}, error) {
	return newUnitaryOpcode(vm.OpPushL)
}

func (p *parser) callonOpcode50() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode50()
}

func (c *current) onOpcode52(ba interface{}) (interface{}, error) {
	return newPushB(ba)
}

func (p *parser) callonOpcode52() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode52(stack["ba"])
}

func (c *current) onOpcode58(a interface{}) (interface{}, error) {
	return newPushAddr(a.(string))
}

func (p *parser) callonOpcode58() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode58(stack["a"])
}

func (c *current) onOpcode64(offset interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpPick, offset.(string))
}

func (p *parser) callonOpcode64() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode64(stack["offset"])
}

func (c *current) onOpcode70() (interface{}, error) {
	return newUnitaryOpcode(vm.OpOver)
}

func (p *parser) callonOpcode70() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode70()
}

func (c *current) onOpcode72() (interface{}, error) {
	return newUnitaryOpcode(vm.OpOne)
}

func (p *parser) callonOpcode72() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode72()
}

func (c *current) onOpcode74() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNow)
}

func (p *parser) callonOpcode74() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode74()
}

func (c *current) onOpcode76() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNot)
}

func (p *parser) callonOpcode76() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode76()
}

func (c *current) onOpcode78() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNeg1)
}

func (p *parser) callonOpcode78() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode78()
}

func (c *current) onOpcode80() (interface{}, error) {
	return newUnitaryOpcode(vm.OpNeg)
}

func (p *parser) callonOpcode80() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode80()
}

func (c *current) onOpcode82() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMulDiv)
}

func (p *parser) callonOpcode82() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode82()
}

func (c *current) onOpcode84() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMul)
}

func (p *parser) callonOpcode84() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode84()
}

func (c *current) onOpcode86() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMod)
}

func (p *parser) callonOpcode86() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode86()
}

func (c *current) onOpcode88() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMin)
}

func (p *parser) callonOpcode88() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode88()
}

func (c *current) onOpcode90() (interface{}, error) {
	return newUnitaryOpcode(vm.OpMax)
}

func (p *parser) callonOpcode90() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode90()
}

func (c *current) onOpcode92() (interface{}, error) {
	return newUnitaryOpcode(vm.OpLt)
}

func (p *parser) callonOpcode92() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode92()
}

func (c *current) onOpcode94(id, count interface{}) (interface{}, error) {
	return newCallOpcode(vm.OpLookup, id.(string), count.(string))
}

func (p *parser) callonOpcode94() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode94(stack["id"], stack["count"])
}

func (c *current) onOpcode103() (interface{}, error) {
	return newUnitaryOpcode(vm.OpLen)
}

func (p *parser) callonOpcode103() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode103()
}

func (c *current) onOpcode105() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIndex)
}

func (p *parser) callonOpcode105() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode105()
}

func (c *current) onOpcode107() (interface{}, error) {
	return newUnitaryOpcode(vm.OpInc)
}

func (p *parser) callonOpcode107() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode107()
}

func (c *current) onOpcode109() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIfZ)
}

func (p *parser) callonOpcode109() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode109()
}

func (c *current) onOpcode111() (interface{}, error) {
	return newUnitaryOpcode(vm.OpIfNZ)
}

func (p *parser) callonOpcode111() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode111()
}

func (c *current) onOpcode113() (interface{}, error) {
	return newUnitaryOpcode(vm.OpGt)
}

func (p *parser) callonOpcode113() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode113()
}

func (c *current) onOpcode115(ix interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpFieldL, ix.(string))
}

func (p *parser) callonOpcode115() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode115(stack["ix"])
}

func (c *current) onOpcode121(ix interface{}) (interface{}, error) {
	return newBinaryOpcode(vm.OpField, ix.(string))
}

func (p *parser) callonOpcode121() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode121(stack["ix"])
}

func (c *current) onOpcode127() (interface{}, error) {
	return newUnitaryOpcode(vm.OpFalse)
}

func (p *parser) callonOpcode127() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode127()
}

func (c *current) onOpcode129() (interface{}, error) {
	return newUnitaryOpcode(vm.OpFail)
}

func (p *parser) callonOpcode129() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode129()
}

func (c *current) onOpcode131() (interface{}, error) {
	return newUnitaryOpcode(vm.OpExtend)
}

func (p *parser) callonOpcode131() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode131()
}

func (c *current) onOpcode133() (interface{}, error) {
	return newUnitaryOpcode(vm.OpEq)
}

func (p *parser) callonOpcode133() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode133()
}

func (c *current) onOpcode135() (interface{}, error) {
	return newUnitaryOpcode(vm.OpEndIf)
}

func (p *parser) callonOpcode135() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode135()
}

func (c *current) onOpcode137() (interface{}, error) {
	return newUnitaryOpcode(vm.OpElse)
}

func (p *parser) callonOpcode137() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode137()
}

func (c *current) onOpcode139() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDup2)
}

func (p *parser) callonOpcode139() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode139()
}

func (c *current) onOpcode141() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDup)
}

func (p *parser) callonOpcode141() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode141()
}

func (c *current) onOpcode143() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDrop2)
}

func (p *parser) callonOpcode143() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode143()
}

func (c *current) onOpcode145() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDrop)
}

func (p *parser) callonOpcode145() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode145()
}

func (c *current) onOpcode147() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDivMod)
}

func (p *parser) callonOpcode147() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode147()
}

func (c *current) onOpcode149() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDiv)
}

func (p *parser) callonOpcode149() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode149()
}

func (c *current) onOpcode151(id, count interface{}) (interface{}, error) {
	return newCallOpcode(vm.OpDeco, id.(string), count.(string))
}

func (p *parser) callonOpcode151() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode151(stack["id"], stack["count"])
}

func (c *current) onOpcode160() (interface{}, error) {
	return newUnitaryOpcode(vm.OpDec)
}

func (p *parser) callonOpcode160() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode160()
}

func (c *current) onOpcode162() (interface{}, error) {
	return newUnitaryOpcode(vm.OpChoice)
}

func (p *parser) callonOpcode162() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode162()
}

func (c *current) onOpcode164(id, count interface{}) (interface{}, error) {
	return newCallOpcode(vm.OpCall, id.(string), count.(string))
}

func (p *parser) callonOpcode164() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode164(stack["id"], stack["count"])
}

func (c *current) onOpcode173() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAvg)
}

func (p *parser) callonOpcode173() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode173()
}

func (c *current) onOpcode175() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAppend)
}

func (p *parser) callonOpcode175() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode175()
}

func (c *current) onOpcode177() (interface{}, error) {
	return newUnitaryOpcode(vm.OpAdd)
}

func (p *parser) callonOpcode177() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode177()
}

func (c *current) onOpcode179(v interface{}) (interface{}, error) {
	return newPushOpcode(v.(string))
}

func (p *parser) callonOpcode179() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpcode179(stack["v"])
}

func (c *current) onTimestamp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTimestamp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestamp1()
}

func (c *current) onConstantRef1(k interface{}) (interface{}, error) {
	cm := c.globalStore["constants"].(map[string]string)
	return cm[k.(string)], nil
}

func (p *parser) callonConstantRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstantRef1(stack["k"])
}

func (c *current) onInteger2() (interface{}, error) {
	return strings.TrimSpace(string(c.text)), nil
}

func (p *parser) callonInteger2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger2()
}

func (c *current) onInteger9() (interface{}, error) {
	return strings.TrimSpace(string(c.text)), nil
}

func (p *parser) callonInteger9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger9()
}

func (c *current) onInteger19() (interface{}, error) {
	return strings.TrimSpace(string(c.text)), nil
}

func (p *parser) callonInteger19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger19()
}

func (c *current) onBytes2(b interface{}) (interface{}, error) {
	return b, nil
}

func (p *parser) callonBytes2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBytes2(stack["b"])
}

func (c *current) onBytes6(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonBytes6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBytes6(stack["s"])
}

func (c *current) onAddress1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddress1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddress1()
}

func (c *current) onConstant1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstant1()
}

func (c *current) onFunctionName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonFunctionName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionName1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
