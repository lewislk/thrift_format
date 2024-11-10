package logs

import "github.com/fatih/color"

var Verbose bool

func Info(msg string) {
	color.Green(msg)
}

func InfoF(f string, a ...any) {
	color.Green(f, a)
}

func Warn(msg string) {
	color.Yellow(msg)
}

func WarnF(f string, a ...any) {
	color.Yellow(f, a)
}

func Error(msg string) {
	color.Red(msg)
}

func ErrorF(f string, a ...any) {
	color.Red(f, a)
}

func DebugF(f string, a ...any) {
	if Verbose {
		color.Blue(f, a)
	}
}
