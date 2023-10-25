package version

import "time"

var V = "0.0.1"
var Commit = ""
var Updated = time.Now()

var Version = struct {
	Version string    `json:"version"`
	Commit  string    `json:"commit"`
	Updated time.Time `json:"updated"`
}{
	Version: V,
	Commit:  Commit,
	Updated: Updated,
}
