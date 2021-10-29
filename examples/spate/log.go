package main

import (
	"fmt"

	"github.com/fatih/color"
)

func Error(format string, a ...interface{}) {
	error_label := color.New(color.FgRed, color.Bold)
	error_label.Print("error ")
	color.Red(format, a...)
}

func Warn(format string, a ...interface{}) {
	warn_label := color.New(color.FgYellow, color.Bold)
	warn_label.Print("warn  ")
	color.Yellow(format, a...)
}

func Info(format string, a ...interface{}) {
	info_label := color.New(color.FgBlue)
	info_label.Print("info  ")
	fmt.Printf(format+"\n", a...)
}
