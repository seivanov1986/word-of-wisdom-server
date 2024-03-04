package logger

type Logger interface {
	Println(v ...any)
	Fatal(a ...any)
}
