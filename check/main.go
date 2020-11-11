package main

import (
	"encoding/json"
	"os"
	"sort"

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
	sort.Sort(history) // Nomad provides history from newest to oldest
	versions := make([]resource.Version, 0)

	for i, jobVersion := range history {
		// Only return the current version and newer ones
		if lastVersion != 0 && lastVersion > jobVersion.Version {
			continue
		}
		// If no version was provided return only the latest
		if lastVersion == 0 && i < len(history)-1 {
			continue
		}
		versions = append(versions, resource.Version{jobVersion.Version})
	}

	json.NewEncoder(os.Stdout).Encode(versions)

}
