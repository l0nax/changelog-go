// changelog implements the parser and writer for changelog entries and the
// CHANGELOG.md generator.
package changelog

import _log "gitlab.com/l0nax/changelog-go/internal/log"

// log is used like the regular `log` package but with hooks for sentry
var log = _log.Log
