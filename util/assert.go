package util

import "fmt"

// Assert
//   - param condition :     bool
//   - param messageFormat : type of string (@See fmt.Sprintf)
//   - param args :          slice of string
func Assert(condition bool, messageFormat string, args ...string) {
	if condition {
		panic(fmt.Sprintf(messageFormat, args))
	}
}
