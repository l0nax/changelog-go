// Package tools implements some Functions which are needed to run this Application
package tools

import (
	"crypto/rand"
	log "github.com/sirupsen/logrus"
)

const randStrLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	randData := make([]byte, length)
	_, err := rand.Read(randData)
	if err != nil {
		log.Fatal(err)
	}

	for i, rByte := range randData {
		randData[i] = randStrLetters[randData[i]%byte(len(randStrLetters))]
	}

	return string(randData)
}
