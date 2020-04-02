package internal

import (
	"encoding/json"
	"net/http"
	"runtime"
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
		ilog.Log.Fatalln(err)
	}

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		ilog.Log.Fatalln(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	info := new(versionInfo)

	// marshal response body into versionInfo struct
	_ = json.NewDecoder(res.Body).Decode(info)

	if version != info.Version {
		// update is available
		return true
	}

	// no update found
	return false
}
