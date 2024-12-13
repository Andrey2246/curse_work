package encryption

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

func RuneToBin(input rune) []int {
	inNum := int(input)
	result := make([]int, 16)
	for i := 15; inNum > 0; i-- {
		result[i] = inNum % 2
		inNum = inNum / 2
	}
	return result
}

func AddDESPadding(a []byte) []byte {
	prop := byte(len(a) % 8)
	for i := 0; len(a)%8 != 0; i++ {
		a = append(a, prop)
	}
	return a
}

func StrToBin(a string) []int {
	result := make([]int, 0)
	for _, c := range a {
		bin := RuneToBin(c)
		if len(bin) == 8 {
			result = append(result, bin...)
		}
		result = append(result, bin...)
	}

	return result
}

func removeDESPadding(a []byte) []byte {
	prop := a[len(a)-1]
	realSize := len(a) - 1
	for i := 0; i <= int(prop); i++ {
		if a[len(a)-i-1] != prop {
			return a
		}
		realSize -= 1
	}
	return a[:realSize]
}

func Permute(input []int, table []int) []int {
	output := make([]int, len(table))
	for i, position := range table {
		output[i] = input[position-1]
	}
	return output
}

func LeftShift(input []int, shifts int) []int {
	return append(input[shifts:], input[:shifts]...)
}

func GenerateRoundKeys(key []int) [][]int {
	PermutedKey := Permute(key, PC1)

	left, right := PermutedKey[:28], PermutedKey[28:]

	roundKeys := make([][]int, 16)
	for i := 0; i < 16; i++ {
		left = LeftShift(left, LeftShifts[i])
		right = LeftShift(right, LeftShifts[i])

		combinedKey := append(left, right...)
		roundKeys[i] = Permute(combinedKey, PC2)
	}

	return roundKeys
}

func Expand(input []int) []int { //расширение части по таблице
	return Permute(input, ExpansionTable)
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

func Xor(a, b []int) []int {
	result := make([]int, len(a))
	for i := range a {
		if a[i] == 0 && b[i] == 0 || a[i] == 1 && b[i] == 1 {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}
	return result
}

func Feistel(right []int, roundKey []int) []int {
	expanded := Expand(right)
	xored := Xor(expanded, roundKey)
	substituted := substitute(xored)
	return Permute(substituted, PermutationTable)
}

func mainDES(block []int, roundKeys [][]int) []int {
	PermutedBlock := Permute(block, IPTable) // начальная перестановка

	left, right := PermutedBlock[:32], PermutedBlock[32:]

	for i := 0; i < 16; i++ {
		temp := right
		right = Xor(left, Feistel(right, roundKeys[i]))
		left = temp
	}

	combined := append(right, left...)
	return Permute(combined, FPTable)
}

func ByteToBinarySlice(b byte) []int {
	binary := make([]int, 8) // A byte has 8 bits
	for i := 0; i < 8; i++ {
		// Extract each bit, starting from the most significant bit
		binary[7-i] = int((b >> i) & 1)
	}
	return binary
}

func ByteSToBinS(bytes []byte) []int {
	binary := make([]int, 0)
	for _, b := range bytes {
		binary = append(binary, ByteToBinarySlice(b)...)
	}
	return binary
}

func BinarySliceToByte(binary []int) byte {
	var b byte
	for i, bit := range binary {
		b |= byte(bit) << (7 - i)
	}
	return b
}

func BinSToByteS(bins []int) []byte {
	bytes := make([]byte, 0)
	for i := 0; i < len(bins); i += 8 {
		bytes = append(bytes, BinarySliceToByte(bins[i:i+8]))
	}
	return bytes
}

func EncryptDES(plaintext []byte, key string) []byte {
	text := ByteSToBinS(AddDESPadding(plaintext))

	roundkeys := GenerateRoundKeys(StrToBin(key))

	cyphertext := make([]int, 0)
	for i := 0; i < len(text); i += 64 {
		cyphertext = append(cyphertext, mainDES(text[i:i+64], roundkeys)...)
	}
	return BinSToByteS(cyphertext)
}

func DecryptDES(ciphertext []byte, key string) []byte {
	text := ByteSToBinS(ciphertext)

	roundkeys := GenerateRoundKeys(StrToBin(key))

	reversedKeys := make([][]int, 16)
	for i := 0; i < 16; i++ {
		reversedKeys[i] = roundkeys[15-i]
	}

	cyphertext := make([]int, 0)
	for i := 0; i < len(text); i += 64 {
		cyphertext = append(cyphertext, mainDES(text[i:i+64], reversedKeys)...)
	}
	return removeDESPadding(BinSToByteS(cyphertext))
}
