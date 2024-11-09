package logs

import "github.com/fatih/color"

var Verbose bool

func Info(f string, a ...any) {
	color.Green(f, a)
}

func Warn(f string, a ...any) {
	color.Yellow(f, a)
}

func Error(f string, a ...any) {
	color.Red(f, a)
}

func Debug(f string, a ...any) {
	if Verbose {
		color.Blue(f, a)
	}
}
