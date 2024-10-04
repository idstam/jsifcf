package internal

import (
	"fmt"
	"log"
	"os"
	"time"
	//"github.com/mailgun/mailgun-go"
)

type LogLevel int

type Logger struct {
	Level    LogLevel
	FileName string
	FileChan chan string
}

func (me *Logger) Init(level LogLevel) {
	now := time.Now().UTC().Format("2006-01-02")

	me.Level = LogLevel(level)
	me.FileName = "jsifcf-" + now + ".log"
	me.FileChan = make(chan string)

	go me.writeToFile()
}

const (
	Trace       LogLevel = iota // Trace = 0
	Debug                       // Debug = 1
	Info                        // Info = 2
	Warn                        // Warn = 3
	Error                       // Error = 4
	Fatal                       // Fatal = 5
	Transaction                 // Transaktion som ska skrivas till logg oavsett inställning av log-level och som inte ska vidaer till ops
)

func (me *Logger) Transaction(where string, msg string) {
	me.writeLog(where, Trace, msg)
}

func (me *Logger) Trace(where string, msg string) {
	if me.Level <= Trace {
		me.writeLog(where, Transaction, msg)
	}
}
func (me *Logger) Debug(where string, msg string) {
	if me.Level <= Debug {
		me.writeLog(where, Debug, msg)
	}
}
func (me *Logger) Info(where string, msg string) {
	if me.Level <= Info {
		me.writeLog(where, Info, msg)
	}
}
func (me *Logger) Warn(where string, msg string) {
	if me.Level <= Warn {
		me.writeLog(where, Warn, msg)
	}
}
func (me *Logger) Error(where string, msg string) {
	if me.Level <= Error {
		me.writeLog(where, Error, msg)
	}
}
func (me *Logger) Fatal(where string, msg string) {
	if me.Level <= Fatal {
		me.writeLog(where, Fatal, msg)
		log.Fatal(where, msg)
	}
}
func (me *Logger) writeLog(where string, level LogLevel, msg string) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	line := fmt.Sprintf("%s\t%s\t%s\t%s", now, LogLevelName(level), where, msg)

	if me.FileName != "" {
		me.FileChan <- line
	}
	log.Println(line)
}

func (me *Logger) writeToFile() {
	f, err := os.OpenFile(me.FileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for {
		line := <-me.FileChan
		if _, err := f.WriteString(line + "\n"); err != nil {
			log.Fatal("Failed to write log: " + err.Error())
		}
	}

}

func LogLevelName(level LogLevel) string {
	switch level {
	case 0:
		return "Trace" // Trace = 0
	case 1:
		return "Debug" // Debug = 1
	case 2:
		return "Info" // Info = 2
	case 3:
		return "Warn" // Warn = 3
	case 4:
		return "Error" // Error = 4
	case 5:
		return "Fatal" // Fatal = 5
	default:
		return "Unknown"
	}
}
