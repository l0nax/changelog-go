package tools

import (
	"crypto/rand"
)

const randStrLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	randData := make([]byte, length)
	_, err := rand.Read(randData)
	if err != nil {
		log.Fatal(err)
	}

	for i := range randData {
		randData[i] = randStrLetters[randData[i]%byte(len(randStrLetters))]
	}

	return string(randData)
}
