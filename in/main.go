package main

import (
	"encoding/json"
	"os"

	resource "github.com/cioplenu/concourse-nomad-resource"
	"github.com/cioplenu/concourse-nomad-resource/common"
)

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
}

type Response struct {
	Version  resource.Version  `json:"version"`
	Metadata resource.Metadata `json:"metadata"`
}

func main() {
	var request Request
	err := json.NewDecoder(os.Stdin).Decode(&request)
	common.Check(err, "Error parsing request")

	response := Response{
		Version: request.Version,
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
