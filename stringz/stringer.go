package stringz

import (
	"fmt"
)

var _ fmt.Stringer = StringFunc(nil)

// StringFunc wraps a function to make it implement fmt.Stringer. It is useful when you want to delay evaluation of the function. For example:
//
//	msg := StringFunc(func() string {
//	  return fmt.Sprintf("JSON(obj)=%v", formatJSON(obj))
//	})
//	log.Debug(msg)  // 👈 only run formatJSON() if debug level is enabled
type StringFunc func() string

func (fn StringFunc) String() string { return fn() }
