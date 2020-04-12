package logger

import (
	"os"
	"sync"

	"github.com/shiena/ansicolor"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var _instance *log.Logger

var _once sync.Once

func initialLogger() {

	// arguments stuffs
	// config.SetConfigFilePath(*configPath)
	_instance = log.New()
	_instance.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	_instance.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	_instance.SetReportCaller(false)

	// Only log the warning severity or above.
	// _instance.SetLevel(log.WarnLevel)

}

func Instance() *log.Logger {
	_once.Do(initialLogger)
	return _instance
}

// Trace.Println("I have something standard to say")
// Info.Println("Special Information")
// Warning.Println("There is something you need to know about")
// Error.Println("Something has failed")

func Trace(args ...interface{}) {
	Instance().Traceln(args...)
}

func Debug(args ...interface{}) {
	Instance().Debugln(args...)
}

func Info(args ...interface{}) {
	Instance().Infoln(args...)
}

func Print(args ...interface{}) {
	Instance().Println(args...)
}

func Warn(args ...interface{}) {
	Instance().Warnln(args...)
}

func Warning(args ...interface{}) {
	Instance().Warn(args...)
}

func Error(args ...interface{}) {
	Instance().Errorln(args...)
}

func Fatal(args ...interface{}) {
	Instance().Fatalln(args...)
}

func Panic(args ...interface{}) {
	Instance().Panicln(args...)
}

func Tracef(format string, args ...interface{}) {
	Instance().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Instance().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Instance().Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	Instance().Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Instance().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	Instance().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Instance().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Instance().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Instance().Panicf(format, args...)
}
