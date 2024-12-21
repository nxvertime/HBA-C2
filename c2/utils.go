package main

import (
	"fmt"
	"math/rand"
)

// ////////// OTHERS ////////
// ///// ARGS TO MAP
func argsToMap(args []string) map[string]interface{} {
	argsMap := make(map[string]interface{})
	for i, arg := range args {
		key := fmt.Sprintf("arg%d", i+1)
		argsMap[key] = arg
	}
	return argsMap
}

// /////// OTHER THINGS ///////////
const CHARS = "abcdefghijklmnopqrstuvwxyz123456789"

func createSID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHARS[rand.Intn(len(CHARS))]

	}
	return string(b)
}
