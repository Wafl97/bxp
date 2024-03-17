package util

import "github.com/fatih/color"

const (
	Colour_reset  = "\033[0m"
	Colour_red    = "\033[31m"
	Colour_green  = "\033[32m"
	Colour_yellow = "\033[33m"
	Colour_blue   = "\033[34m"
	Colour_purple = "\033[35m"
	Colour_cyan   = "\033[36m"
	Colour_gray   = "\033[37m"
	Colour_white  = "\033[97m"
)

var (
	red   = color.New(color.FgRed)
	green = color.New(color.FgGreen)
	blue  = color.New(color.FgHiBlue)
)

func Green(str string) string {
	return green.Sprint(str)
}

func Red(str string) string {
	return red.Sprint(str)
}

func Blue(str string) string {
	return blue.Sprint(str)
}
