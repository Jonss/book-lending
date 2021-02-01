package logger

import (
	"fmt"
	"log"
)

func Info(message string) {
	log.Println(fmt.Sprintf("INFO %s", message))
}

func Debug(message string) {
	log.Println(fmt.Sprintf("DEBUG %s", message))
}

func Error(message string) {
	log.Println(fmt.Sprintf("ERROR %s", message))
}

func Warn(message string) {
	log.Println(fmt.Sprintf("WARN %s", message))
}
