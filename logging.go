package main

import "fmt"

func LogMsg(format string, a ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(fmt.Sprintf("%+v\n", format), a...)
}
