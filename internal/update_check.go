package internal

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"time"

	ilog "gitlab.com/l0nax/changelog-go/internal/log"
)

type versionInfo struct {
	Version string `json:"Version"`
	Sha256  string `json:"Sha256"`
}

// CheckUpdate will check if there is a new update available.
// If so it will return 'true'.
func CheckUpdate(version, dataURL string) bool {
	// initialize HTTP Client
	client := http.Client{
		// max 500ms timeout, to not slow the application too much down
		Timeout: time.Millisecond * 500,
	}

	// prepare API URL
	url := dataURL + "/" + runtime.GOOS + "-" + runtime.GOARCH + ".json"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// NOTE(l0nax): Do not print an error when the update check fails.
		// ilog.Log.WithError(err).Errorln("could not initialize http requester")
		return false
	}

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		// NOTE(l0nax): Do not print an error when the update check fails.
		// ilog.Log.Errorln(err)
		return false
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	info := new(versionInfo)

	// marshal response body into versionInfo struct
	_ = json.NewDecoder(res.Body).Decode(info)

	// remove 'v' prefix from any versions
	info.Version = strings.TrimPrefix(info.Version, "v")
	version = strings.TrimPrefix(version, "v")

	if version != info.Version {
		// update is available
		return true
	}

	// no update found
	return false
}
