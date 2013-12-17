package logger

import (
	"log"
	"os"
	"os/signal"
)

type Logger struct {
	verbose bool
	sigCh   chan os.Signal
	killCh  chan bool
}

func New() *Logger {
	logger := &Logger{
		verbose: false,
		sigCh:   make(chan os.Signal),
		killCh:  make(chan bool),
	}
	signal.Notify(logger.sigCh, os.Interrupt)

	go logger.monitor()

	return logger
}

func (l *Logger) monitor() {
	for {
		select {
		case <-l.sigCh:
			l.verbose = !l.verbose
			if l.verbose {
				l.Printf("verbose logging enabled")
			} else {
				l.Printf("verbose logging disabled")
			}
		case <-l.killCh:
			return
		}
	}
}

func (l *Logger) Shutdown() {
	l.killCh <- true
}

func (l *Logger) Errorf(fmt string, args ...interface{}) {
	log.Printf("[ERROR] "+fmt, args...)
}

func (l *Logger) Printf(fmt string, args ...interface{}) {
	log.Printf("[INFO] "+fmt, args...)
}

func (l *Logger) Debugf(fmt string, args ...interface{}) {
	if l.verbose {
		log.Printf("[DEBUG] "+fmt, args...)
	}
}

// Global Logger

var pkgLogger *Logger

func InitPkgLogger() {
	if pkgLogger == nil {
		pkgLogger = New()
	}
}

func Errorf(fmt string, args ...interface{}) {
	if pkgLogger != nil {
		pkgLogger.Errorf(fmt, args...)
	} else {
		log.Printf("[ERROR] dropping logs, pkgLogger is nil")
	}
}

func Printf(fmt string, args ...interface{}) {
	if pkgLogger != nil {
		pkgLogger.Printf(fmt, args...)
	} else {
		log.Printf("[ERROR] dropping logs, pkgLogger is nil")
	}
}

func Debugf(fmt string, args ...interface{}) {
	if pkgLogger != nil {
		pkgLogger.Debugf(fmt, args...)
	} else {
		log.Printf("[ERROR] dropping logs, pkgLogger is nil")
	}
}
