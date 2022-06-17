package utils

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// InitLog Init Log
func InitLog(logLevel, path, filename string) {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("init logger error. %+v", errors.WithStack(err))
	}

	log.SetLevel(level)
	ConfigLocalFilesystemLogger(path, filename, 7*time.Hour*24, time.Second*20)
}

// ConfigLocalFilesystemLogger Rotate Log
func ConfigLocalFilesystemLogger(logPath, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	writer := NewRotatelogsWriter(path.Join(logPath, logFileName), maxAge, rotationTime)
	debugWriter := NewRotatelogsWriter(path.Join(logPath, "debug-"+logFileName), maxAge, rotationTime)
	errorWriter := NewRotatelogsWriter(path.Join(logPath, "error-"+logFileName), maxAge, rotationTime)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: debugWriter,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: errorWriter,
		log.FatalLevel: errorWriter,
		log.PanicLevel: errorWriter,
	}, &log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	log.AddHook(lfHook)
}

func NewRotatelogsWriter(baseLogPath string, maxAge, rotationTime time.Duration) *rotatelogs.RotateLogs {
	debugWriter, err := rotatelogs.New(
		baseLogPath+".%Y%m%d",
		// rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)

	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	return debugWriter
}
