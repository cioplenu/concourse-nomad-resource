package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	resource "github.com/cioplenu/concourse-nomad-resource"
)

func Check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %s\n", err)
		os.Exit(1)
	}
}

func GetHistory(source resource.Source) resource.History {
	cmd := exec.Command(
		"nomad",
		"job",
		"history",
		"-json",
		"-address="+source.URL,
		"-token="+source.Token,
		source.Name,
	)
	var histResp bytes.Buffer
	cmd.Stdout = &histResp
	err := cmd.Run()
	Check(err, "Error checking versions")

	var history []resource.JobVersion
	json.Unmarshal(histResp.Bytes(), &history)

	return history
}
