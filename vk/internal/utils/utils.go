package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomStrings(count int, minLen int, maxLen int) []string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		length := rand.Intn(maxLen-minLen+1) + minLen
		s := make([]rune, length)
		for j := range s {
			s[j] = letters[rand.Intn(len(letters))]
		}
		result = append(result, string(s))
	}
	return result
}
