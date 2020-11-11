package main

import (
	"encoding/json"
	"os"

	resource "github.com/cioplenu/concourse-nomad-resource"
	"github.com/cioplenu/concourse-nomad-resource/common"
)

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version,omitempty"`
}

func main() {
	var request Request
	err := json.NewDecoder(os.Stdin).Decode(&request)
	common.Check(err, "Error parsing request")

	lastVersion := request.Version.Version

	history := common.GetHistory(request.Source)
	versions := make([]resource.Version, 0)

	for _, jobVersion := range history {
		if lastVersion != 0 && lastVersion > jobVersion.Version {
			continue
		}
		versions = append(versions, resource.Version{jobVersion.Version})
		if lastVersion == 0 {
			break
		}
	}

	json.NewEncoder(os.Stdout).Encode(versions)

}
