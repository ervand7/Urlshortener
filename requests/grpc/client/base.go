package main

import (
	"fmt"
	"math/rand"
	"time"
)

var userID = "89510fb9-ec3c-4393-937c-3f97ef8b0fab"

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
