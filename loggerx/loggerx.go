package loggerx

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

func Init() {

	logFilePath := filepath.Join(".", "loggerx", "logfile.txt")

	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file: ", err)
	}
}

func Info(message string) {
	logMessage := fmt.Sprintf("[INFO][%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)

	if _, err := logFile.WriteString(logMessage); err != nil {
		fmt.Println("Failed to write to log file: ", err)
	}
}

func Error(message string) {
	logMessage := fmt.Sprintf("[ERROR][%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)

	if _, err := logFile.WriteString(logMessage); err != nil {
		fmt.Println("Failed to write to log file: ", err)
	}
}

func SimpleHttpGet(url string) {

	_, err := http.Get(url)
	if err != nil {
		Info("This is an info log message.")
	} else {
		Error("This is an error log message.")
	}
}
