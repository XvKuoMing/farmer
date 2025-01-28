package internals

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func GetLogger(from *string) *Logger {
	if from == nil {
		return nil // reason: https://github.com/golang/go/issues/47164
	}
	var path string = *from
	var _name []string = strings.Split(strings.TrimSuffix(path, "/"), "/")
	var name string = _name[len(_name)-1]

	logger := Logger{}
	flags := log.Ldate | log.Ltime
	var perm fs.FileMode = 0666
	log_file := path + ".log"
	err_file := path + ".err"
	os.MkdirAll(filepath.Dir(log_file), perm)

	logFile, _ := os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	logErrFile, _ := os.OpenFile(err_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	logger.info = log.New(logFile, name+" INFO: ", flags)
	logger.warn = log.New(logFile, name+" WARNING: ", flags)
	logger.err = log.New(logErrFile, name+" ERROR: ", flags)

	return &logger
}

func (logger *Logger) log(f func()) {
	if logger == nil {
		return
	}
	f()
}

func (logger *Logger) Info(v ...any) {
	logger.log(func() { logger.info.Println(v...) })
}

func (logger *Logger) Warn(v ...any) {
	logger.log(func() { logger.warn.Println(v...) })
}

func (logger *Logger) Err(v ...any) {
	logger.log(func() { logger.err.Println(v...) })
}

func (logger *Logger) FatalErr(v ...any) {
	logger.log(func() { logger.err.Fatalln(v...) })
}
