// Copyright 2018-19 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"math/rand"
	// "time"
)

func init() {
	// rand.Seed(time.Now().UnixNano()) // Make variables deterministic
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq generates random identifiers
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
