package logs

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	errPackage "github.com/MarlonG1/api-facturacion-sv/pkg/error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Logger *logrus.Logger
var logsLevel = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}

type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColors := map[logrus.Level]string{
		logrus.DebugLevel: "\033[37m", // Tonalidad gris para debug
		logrus.InfoLevel:  "\033[32m", // Tonalidad verde para info
		logrus.WarnLevel:  "\033[33m", // Tonalidad amarilla para warning
		logrus.ErrorLevel: "\033[31m", // Tonalidad roja para error
		logrus.FatalLevel: "\033[35m", // Tonalidad morada para fatal
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	color := levelColors[entry.Level]

	baseMessage := fmt.Sprintf("%s %s [%s]\033[0m %s",
		color,
		timestamp,
		level,
		entry.Message,
	)

	var fields string
	if len(entry.Data) > 0 {
		fields = "\n	Details:\n"
		for k, v := range entry.Data {
			fields += fmt.Sprintf("	%-10s: %v\n", k, v)
		}
	}

	return []byte(fmt.Sprintf("%s%s\n", baseMessage, fields)), nil
}

func InitLogger() error {
	Logger = logrus.New()

	formatter := &CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	}

	Logger.SetFormatter(formatter)
	Logger.SetLevel(determineLogLevel(config.Log.Level))
	Logger.SetOutput(os.Stdout)
	logDir := utils.FindProjectRoot()

	// Asegurarse de que el directorio de logs exista
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return errPackage.ErrLogDirectoryNotFound
	}
	logFilePath := filepath.Join(logDir+config.Log.Path, "dte_microservice.log")

	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return errPackage.ErrFailedToCreateLogFiles
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	Logger.SetOutput(mw)

	return nil
}

func logWithFields(level logrus.Level, msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 && fields[0] != nil {
		Logger.WithFields(fields[0]).Log(level, msg)
	} else {
		Logger.Log(level, msg)
	}
}

func Debug(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.DebugLevel, msg, fields...)
}

func Info(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.InfoLevel, msg, fields...)
}

func Warn(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.WarnLevel, msg, fields...)
}

func Error(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.ErrorLevel, msg, fields...)
}

func Fatal(msg string, fields ...map[string]interface{}) {
	logWithFields(logrus.FatalLevel, msg, fields...)
}

func determineLogLevel(logLevel string) logrus.Level {
	if level, ok := logsLevel[logLevel]; ok {
		return level
	}
	return logrus.InfoLevel
}

type WriteHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (hook *WriteHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *WriteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
