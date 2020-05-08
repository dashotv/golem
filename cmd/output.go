package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

var au aurora.Aurora

func Info(msg string) {
	fmt.Println(au.Cyan(msg))
}

func Error(msg string) {
	fmt.Println(au.Red(msg))
}

func Warn(msg string) {
	fmt.Println(au.Yellow(msg))
}
