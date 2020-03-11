// Package tools implements some Functions which are needed to run this
// Application
package tools

import _log "gitlab.com/l0nax/changelog-go/internal/log"

// log is used like the regular `log` package but with hooks for sentry
var log = _log.Log
