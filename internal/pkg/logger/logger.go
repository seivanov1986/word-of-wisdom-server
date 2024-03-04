package logger

import (
	"fmt"
	"os"
)

type logger struct {
}

func New() *logger {
	return &logger{}
}

func (l *logger) Println(a ...any) {
	fmt.Println(a)
}

func (l *logger) Fatal(a ...any) {
	fmt.Println(a)
	os.Exit(1)
}
