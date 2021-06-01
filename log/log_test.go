package log_test

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/artisanhe/tools/log"
	"github.com/artisanhe/tools/log/context"
)

var logger = log.Log{
	Name:  "test",
	Level: "Debug",
	//Path:  "./logs/test.log",
}

func init() {
	logger.Init()
}

func TestLog(t *testing.T) {
	context.SetLogID("1123123")

	logrus.Info("Info")
	logrus.Warning("Warn")
	logrus.Error("Error")
	logrus.WithField("test2", 2).Info("test")
}
