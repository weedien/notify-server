package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/weedien/notify-server/common/config"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func Init() {
	setLevel()
	setFormatter()
	setOutput()
}

func setLevel() {
	l, err := logrus.ParseLevel(config.Config().Log.Level)
	if err != nil {
		// 默认日志级别为info
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(l)
	}
}

func setFormatter() {
	format := config.Config().Log.Format

	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else if format == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else if format == "custom" {
		logrus.SetFormatter(&CustomFormatter{})
	}
}

func setOutput() {
	outputPath := config.Config().Log.OutputPath

	var logOutput io.Writer
	if config.Config().Log.Rotate {
		logOutput = rotateLogger()
	} else if outputPath != "" {
		logOutput = openFile(outputPath)
	}

	stdout := config.Config().Log.Stdout

	if stdout && logOutput != nil {
		logOutput = io.MultiWriter(os.Stdout, logOutput)
	}

	if logOutput == nil {
		logOutput = os.Stdout
	}

	logrus.SetOutput(logOutput)
}

func openFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("Failed to log to file: %v", err)
		return nil
	}
	return file
}

func rotateLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.Config().Log.OutputPath,
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		LocalTime:  true,
	}
}
