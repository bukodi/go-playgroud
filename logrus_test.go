package playground

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNegotiateHTTPServer(t *testing.T) {
	var myHook = MyHook{}
	//logrus.AddHook(&myHook)
	x := logrus.New()
	x.AddHook(&myHook)
	logger := x.WithField("component", "test")
	logger.Info("Hello 1")
	logrus.Info("Hello 2")
	logger.Info("Hello 3")
	internalFunctionLevel1(logger)
}

func internalFunctionLevel1(l *logrus.Entry) {
	l.Info("hello 4")
}

type MyHook struct {
}

func (myhook MyHook) Fire(entry *logrus.Entry) error {
	fmt.Printf("%+v", entry)
	return nil
}

func (myhook MyHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}
