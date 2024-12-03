package DES

import (
	"fmt"
)

// Permutation Table
var IPTable = []int{
	58, 50, 42, 34, 26, 18, 10, 2,
	60, 52, 44, 36, 28, 20, 12, 4,
	62, 54, 46, 38, 30, 22, 14, 6,
	64, 56, 48, 40, 32, 24, 16, 8,
	57, 49, 41, 33, 25, 17, 9, 1,
	59, 51, 43, 35, 27, 19, 11, 3,
	61, 53, 45, 37, 29, 21, 13, 5,
	63, 55, 47, 39, 31, 23, 15, 7,
}

var IPInvTable = []int{
	40, 8, 48, 16, 56, 24, 64, 32,
	39, 7, 47, 15, 55, 23, 63, 31,
	38, 6, 46, 14, 54, 22, 62, 30,
	37, 5, 45, 13, 53, 21, 61, 29,
	36, 4, 44, 12, 52, 20, 60, 28,
	35, 3, 43, 11, 51, 19, 59, 27,
	34, 2, 42, 10, 50, 18, 58, 26,
	33, 1, 41, 9, 49, 17, 57, 25,
}

var ExpansionTable = []int{
	32, 1, 2, 3, 4, 5,
	4, 5, 6, 7, 8, 9,
	8, 9, 10, 11, 12, 13,
	12, 13, 14, 15, 16, 17,
	16, 17, 18, 19, 20, 21,
	20, 21, 22, 23, 24, 25,
	24, 25, 26, 27, 28, 29,
	28, 29, 30, 31, 32, 1,
}

var PermutationTable = []int{
	16, 7, 20, 21, 29, 12, 28, 17,
	1, 15, 23, 26, 5, 18, 31, 10,
	2, 8, 24, 14, 32, 27, 3, 9,
	19, 13, 30, 6, 22, 11, 4, 25,
}

var SBox = [8][4][16]int{
	// S-Box 1
	{
		{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
		{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
		{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
		{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13},
	},
	// S-Box 2
	{
		{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
		{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
		{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
		{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
	},
	// S-Box 3
	{
		{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
		{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
		{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
		{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
	},
	// S-Box 4
	{
		{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
		{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
		{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
		{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14},
	},
	// S-Box 5
	{
		{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
		{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
		{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
		{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
	},
	// S-Box 6
	{
		{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
		{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
		{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
		{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
	},
	// S-Box 7
	{
		{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
		{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
		{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
		{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12},
	},
	// S-Box 8
	{
		{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
		{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
		{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
		{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11},
	},
}

var PC1 = []int{
	57, 49, 41, 33, 25, 17, 9,
	1, 58, 50, 42, 34, 26, 18,
	10, 2, 59, 51, 43, 35, 27,
	19, 11, 3, 60, 52, 44, 36,
	63, 55, 47, 39, 31, 23, 15,
	7, 62, 54, 46, 38, 30, 22,
	14, 6, 61, 53, 45, 37, 29,
	21, 13, 5, 28, 20, 12, 4,
}

// Permuted Choice 2 (PC-2) - reduces 56-bit key to 48 bits for each subkey
var PC2 = []int{
	14, 17, 11, 24, 1, 5, 3, 28,
	15, 6, 21, 10, 23, 19, 12, 4,
	26, 8, 16, 7, 27, 20, 13, 2,
	41, 52, 31, 37, 47, 55, 30, 40,
	51, 45, 33, 48, 44, 49, 39, 56,
	34, 53, 46, 42, 50, 36, 29, 32,
}

var FPTable = []int{
	40, 8, 48, 16, 56, 24, 64, 32,
	39, 7, 47, 15, 55, 23, 63, 31,
	38, 6, 46, 14, 54, 22, 62, 30,
	37, 5, 45, 13, 53, 21, 61, 29,
	36, 4, 44, 12, 52, 20, 60, 28,
	35, 3, 43, 11, 51, 19, 59, 27,
	34, 2, 42, 10, 50, 18, 58, 26,
	33, 1, 41, 9, 49, 17, 57, 25,
}

// Number of left shifts for each round
var LeftShifts = [16]int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}

func addDESPadding(a []int) []int {
	if len(a)%64 != 0 { // iso padding
		addByte := make([]int, 8)
		need := len(a) % 64
		for i := 0; need > 0; i++ {
			addByte[8-i] = need % 2
			need /= 2
		}
		for i := 0; i < (len(a)%64)/8; i++ {
			a = append(a, addByte...)
		}
	}
	return a
}

func StrToBin(a string) []int {
	result := make([]int, len(a)*16)
	for _, c := range a {
		bin := RuneToBin(c)
		result = append(result, bin...)
	}

	return result
}

func pow(x, y int) int {
	result := 1
	for i := 0; i < x; i++ {
		result *= y
	}
	return result
}

func removeDESPadding(a []int) []int {
	for i := 0; i < 8; i++ {
		if a[len(a)-i] != a[len(a)-8-i] {
			return a //no padding
		}
	}
	n := 0
	for i := 0; i < 8; i++ {
		n += a[len(a)-8-i] * int(pow(2, i))
	}
	a = a[:n]
	return a
}

func BinToStr(a []int) string {
	result := ""
	for i := 0; i < len(a); i += 16 {
		result += string(BinToRune(a[i : i+16]))
	}
	return result
}

func permute(input []int, table []int) []int {
	output := make([]int, len(table))
	for i, position := range table {
		output[i] = input[position-1]
	}
	return output
}

func leftShift(input []int, shifts int) []int {
	return append(input[shifts:], input[:shifts]...)
}

func generateRoundKeys(key []int) [][]int {
	permutedKey := permute(key, PC1)

	left, right := permutedKey[:28], permutedKey[28:]

	roundKeys := make([][]int, 16)
	for i := 0; i < 16; i++ {
		left = leftShift(left, LeftShifts[i])
		right = leftShift(right, LeftShifts[i])

		combinedKey := append(left, right...)
		roundKeys[i] = permute(combinedKey, PC2)
	}

	return roundKeys
}

func expand(input []int) []int { //расширение части по таблице
	return permute(input, ExpansionTable)
}

func substitute(input []int) []int { //s-box подстановка
	output := []int{}
	for i := 0; i < 8; i++ {
		block := input[i*6 : (i+1)*6]
		row := block[0]*2 + block[5]
		col := block[1]*8 + block[2]*4 + block[3]*2 + block[4]
		value := SBox[i][row][col]
		binary := fmt.Sprintf("%04b", value)
		for _, bit := range binary {
			output = append(output, int(bit-'0'))
		}
	}
	return output
}

func xor(a, b []int) []int {
	result := make([]int, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func feistel(right []int, roundKey []int) []int {
	expanded := expand(right)
	xored := xor(expanded, roundKey)
	substituted := substitute(xored)
	return permute(substituted, PermutationTable)
}

func mainDES(text string, roundKeys [][]int) []int {
	blocks := StrToBin(text)
	cyphertext := make([]int, len(blocks))
	for i := 0; i < len(blocks); i += 64 {
		block := blocks[i : i+64]
		permutedBlock := permute(block, IPTable) // начальная перестановка

		left, right := permutedBlock[:32], permutedBlock[32:]
		left = feistel(left, roundKeys[0])

		for i := 0; i < 16; i++ {
			temp := right
			right = xor(left, feistel(right, roundKeys[i]))
			left = temp
		}

		combined := append(right, left...)
		combined = permute(combined, FPTable)
		cyphertext = append(cyphertext, combined...)
	}
	return cyphertext
}

func EncryptDES(text, key string) []int {
	roundKeys := generateRoundKeys(StrToBin(key))
	return mainDES(text, roundKeys)
}

func DecryptDES(text, key string) []int {
	roundKeys := generateRoundKeys(StrToBin(key))

	reversedKeys := make([][]int, 16)
	for i := 0; i < 16; i++ {
		reversedKeys[i] = roundKeys[15-i]
	}
	return mainDES(text, roundKeys)
}
