package print

import (
	"github.com/fatih/color"
)

func Info(a string)  {
	color.New(color.FgGreen).Print(a)
}


func Infoln(a string)  {
	c := color.New(color.FgGreen)
	c.Println(a)
}

func Error(a ...interface{})  {
	color.New(color.FgRed).Print(a)
}

func Errorln(a string)  {
	color.New(color.FgRed).Println(a)
}