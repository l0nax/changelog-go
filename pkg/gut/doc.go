// gut (git utils) implements some Helper Functions around the `git` Library to easier
// interact with the Repository.
package gut

import _log "gitlab.com/l0nax/changelog-go/internal/log"

// log is used like the regular `log` package but with hooks for sentry
var log = _log.Log
