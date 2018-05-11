// Code generated by "stringer -trimprefix Op -type Opcode opcodes.go"; DO NOT EDIT.

package vm

import "strconv"

const _Opcode_name = "NopDropDrop2DupDup2SwapOverPickRollRetFailZeroPush64OneNeg1PushTNowRandPushLAddSubMulDivModNotNegIncDecIndexLenAppendExtendSliceFieldFieldLIfzIfnzElseEndSumAvgMaxMinChoiceWChoiceSortLookup"

var _Opcode_map = map[Opcode]string{
	0:   _Opcode_name[0:3],
	1:   _Opcode_name[3:7],
	2:   _Opcode_name[7:12],
	5:   _Opcode_name[12:15],
	6:   _Opcode_name[15:19],
	9:   _Opcode_name[19:23],
	13:  _Opcode_name[23:27],
	14:  _Opcode_name[27:31],
	15:  _Opcode_name[31:35],
	16:  _Opcode_name[35:38],
	17:  _Opcode_name[38:42],
	32:  _Opcode_name[42:46],
	41:  _Opcode_name[46:52],
	42:  _Opcode_name[52:55],
	43:  _Opcode_name[55:59],
	44:  _Opcode_name[59:64],
	45:  _Opcode_name[64:67],
	47:  _Opcode_name[67:71],
	48:  _Opcode_name[71:76],
	64:  _Opcode_name[76:79],
	65:  _Opcode_name[79:82],
	66:  _Opcode_name[82:85],
	67:  _Opcode_name[85:88],
	68:  _Opcode_name[88:91],
	69:  _Opcode_name[91:94],
	70:  _Opcode_name[94:97],
	71:  _Opcode_name[97:100],
	72:  _Opcode_name[100:103],
	80:  _Opcode_name[103:108],
	81:  _Opcode_name[108:111],
	82:  _Opcode_name[111:117],
	83:  _Opcode_name[117:123],
	84:  _Opcode_name[123:128],
	96:  _Opcode_name[128:133],
	112: _Opcode_name[133:139],
	128: _Opcode_name[139:142],
	129: _Opcode_name[142:146],
	130: _Opcode_name[146:150],
	136: _Opcode_name[150:153],
	144: _Opcode_name[153:156],
	145: _Opcode_name[156:159],
	146: _Opcode_name[159:162],
	147: _Opcode_name[162:165],
	148: _Opcode_name[165:171],
	149: _Opcode_name[171:178],
	150: _Opcode_name[178:182],
	151: _Opcode_name[182:188],
}

func (i Opcode) String() string {
	if str, ok := _Opcode_map[i]; ok {
		return str
	}
	return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
}