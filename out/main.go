package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	resource "github.com/cioplenu/concourse-nomad-resource"
	"github.com/cioplenu/concourse-nomad-resource/common"
)

type Params struct {
	JobPath  string            `json:"job_path"`
	Vars     map[string]string `json:"vars"`
	VarFiles map[string]string `json:"var_files"`
}

type OutConfig struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Response struct {
	Version  resource.Version  `json:"version"`
	Metadata resource.Metadata `json:"metadata"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}
	sourceDir := os.Args[1]

	var config OutConfig
	err := json.NewDecoder(os.Stdin).Decode(&config)
	common.Check(err, "Error parsing configuration")

	templPath := filepath.Join(sourceDir, config.Params.JobPath)
	templFile, err := ioutil.ReadFile(templPath)
	common.Check(err, "Could not read input file "+templPath)
	tmpl, err := template.New("job").Parse(string(templFile))
	common.Check(err, "Error parsing template")

	for name, path := range config.Params.VarFiles {
		varPath := filepath.Join(sourceDir, path)
		varFile, err := ioutil.ReadFile(varPath)
		common.Check(err, "Error reading var file")
		config.Params.Vars[name] = strings.TrimSpace(string(varFile))
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, config.Params.Vars)
	common.Check(err, "Error executing template")

	outFile, err := os.Create(templPath)
	common.Check(err, "Error creating output file")
	defer outFile.Close()
	_, err = outFile.Write(buf.Bytes())
	common.Check(err, "Error writing output file")

	cmd := exec.Command(
		"nomad",
		"job",
		"run",
		"-address="+config.Source.URL,
		"-token="+config.Source.Token,
		templPath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing nomad: %s\n", err)
		fmt.Fprint(os.Stderr, out.String())
		os.Exit(1)
	}

	fmt.Fprint(os.Stderr, out.String())

	history := common.GetHistory(config.Source)

	response := Response{
		Version: resource.Version{
			Version: history[0].Version,
		},
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
