// logrus contains the initialized and configured logrus logger
package logrus

import (
	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	// intialize sentry before starting
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://e87c950e1b544af385a198035234b248@fabmation.info/3",
	})
	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
	}

	Log = logrus.New()

	// add logrus sentry hook
	Log.AddHook(sentryhook.New([]logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}))
}
