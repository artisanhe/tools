package unique_str

import (
	"fmt"
	"math/rand"
)

var letterMap = map[byte]byte{
	'0': 'a',
	'1': 'b',
	'2': 'c',
	'3': 'd',
	'4': 'e',
	'5': 'f',
	'6': 'g',
	'7': 'h',
	'8': 'i',
	'9': 'j',
}

var letterArr = []byte("klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateUniqueIDToStr(unique_id uint64, strLen int32) string {
	unique_str_byte := []byte(fmt.Sprintf("%d", unique_id))

	for i := 0; i < len(unique_str_byte); i = i + 2 {
		temp := unique_str_byte[i]
		unique_str_byte[i] = letterMap[temp]
	}
	for l := int32(len(unique_str_byte)); l < strLen; l++ {
		unique_str_byte = append(unique_str_byte, letterArr[rand.Intn(len(letterArr))])
	}
	return string(unique_str_byte)
}
