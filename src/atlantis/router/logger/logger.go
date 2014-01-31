/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package logger

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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
	signal.Notify(logger.sigCh, syscall.SIGHUP)

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
		log.Printf("[nil] [ERROR] "+fmt, args...)
	}
}

func Printf(fmt string, args ...interface{}) {
	if pkgLogger != nil {
		pkgLogger.Printf(fmt, args...)
	} else {
		log.Printf("[nil] [PRINTF] "+fmt, args...)
	}
}

func Debugf(fmt string, args ...interface{}) {
	if pkgLogger != nil {
		pkgLogger.Debugf(fmt, args...)
	} else {
		log.Printf("[nil] [DEBUG] "+fmt, args...)
	}
}
