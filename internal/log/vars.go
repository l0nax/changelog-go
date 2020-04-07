// Package log contains the initialized and configured logrus logger
package log

import (
	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	"github.com/sirupsen/logrus"

	"gitlab.com/l0nax/changelog-go/pkg/version"
)

// Log is the global "logger" which can be used by any other package
//
// Example (global variable):
//	import ilog gitlab.com/l0nax/changelog-go/internal/log
//
//	var log = ilog.Log
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
				Log.Debugln("Discard sending crash to sentry server because Environment == 'develop'")
				return nil
			}

			return event
		},
	})
	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
	}

	Log = logrus.New()
	// Overwrite the default exit function to flush the sentry
	// events
	Log.ExitFunc = func(e int) {
		// sentry flush
		fmt.Printf("Sentry flushed: %t\n", sentry.Flush(time.Second*2))

		os.Exit(1)
	}

	// add logrus sentry hook
	Log.AddHook(sentryhook.New([]logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}))
}
