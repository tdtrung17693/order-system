package utils

import (
	"log"
	"os"
)

var errorLogger = log.New(os.Stderr, "[error] ", 0)

//  := log.New(os.Stderr, "", 0)
// l.Println("log msg")
func LogErrorLn(v ...any) {
	errorLogger.Println(v...)
}
