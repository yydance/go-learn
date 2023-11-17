package setting

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Logger = logrus.New()

// 设置日志输出方式，标准或者文件
// 设置日志级别，默认Debug

func InitLog() {
	logFormatSet()
	logModeSet()
	logLevelSet()
}
func logFormatSet() {
	var logFormat string
	fmt.Println(Viper.GetString("log-format"))
	fmt.Println(Viper.GetString("log.format"))
	if Viper.GetString("log-format") != "" {
		logFormat = Viper.GetString("log-format")
	} else {
		logFormat = Viper.GetString("log.format")
	}
	fmt.Println(logFormat)
	switch logFormat {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2023-11-15 12:00:00",
		})
	case "text":
		Logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
		})
	default:
		Logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
		})
	}
}

func logModeSet() {

	if Viper.GetBool("stdout") {
		Logger.SetOutput(os.Stdout)
		return
	}
	logFile := Viper.GetString("log.logFile")
	fmt.Println(logFile)
	logDir := filepath.Dir(logFile)
	fmt.Println(logDir)
	if _, err := pathExists(logDir); err == nil {
		if err == nil {
			writers := []io.Writer{
				//os.Stdout,
				&lumberjack.Logger{
					Filename:   logFile,
					MaxSize:    Viper.GetInt("log.maxSize"),
					MaxBackups: Viper.GetInt("log.maxBackups"),
					MaxAge:     Viper.GetInt("log.maxAge"),
					Compress:   Viper.GetBool("log.compress"),
					LocalTime:  true,
				},
			}
			fileAndStdoutWriter := io.MultiWriter(writers...)
			Logger.SetOutput(fileAndStdoutWriter)
		} else {
			Logger.Errorf("Failed to log to file: %s", logFile)
		}
	}
}

func pathExists(path string) (bool, error) {
	Logger.Infof("日志目录：%s", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Logger.Warnf("mkdir logDir [%s] failed,[%v]", path, err)
		} else {
			return true, nil
		}
	}
	return false, err
}

func logLevelSet() {
	var logLevel string
	if Viper.GetString("log-level") != "" {
		logLevel = Viper.GetString("log-level")
	} else {
		logLevel = Viper.GetString("log.level")
	}

	switch logLevel {
	case "trace", "Trace":
		Logger.SetLevel(logrus.TraceLevel)
	case "debug", "Debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info", "Info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn", "Warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error", "Error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "fatal", "Fatal":
		Logger.SetLevel(logrus.FatalLevel)
	case "panic", "Panic":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		Logger.SetLevel(logrus.WarnLevel)
	}
}
