// Package log contains the initialized and configured logrus logger
package log

import (
	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	"github.com/sirupsen/logrus"

	"gitlab.com/l0nax/changelog-go/pkg/version"
)

var Log *logrus.Logger

func init() {
	// initialize sentry before starting
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://e87c950e1b544af385a198035234b248@fabmation.info/3",
		Release:     "changelog-go@" + version.Version,
		Environment: version.Environment,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// discard the event if Environment is 'develop'
			if event.Environment == "develop" {
				return nil
			}

			return event
		},
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
	}))
}
