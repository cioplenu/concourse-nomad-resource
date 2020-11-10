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
)

type SourceConfig struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type ParamsConfig struct {
	JobPath  string            `json:"job_path"`
	Vars     map[string]string `json:"vars"`
	VarFiles map[string]string `json:"var_files"`
}

type OutConfig struct {
	Source SourceConfig `json:"source"`
	Params ParamsConfig `json:"params"`
}

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %s\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}
	sourceDir := os.Args[1]

	var config OutConfig
	err := json.NewDecoder(os.Stdin).Decode(&config)
	check(err, "Error parsing configuration")

	templPath := filepath.Join(sourceDir, config.Params.JobPath)
	templFile, err := ioutil.ReadFile(templPath)
	check(err, "Could not read input file "+templPath)
	tmpl, err := template.New("job").Parse(string(templFile))
	check(err, "Error parsing template")

	for name, path := range config.Params.VarFiles {
		varPath := filepath.Join(sourceDir, path)
		varFile, err := ioutil.ReadFile(varPath)
		check(err, "Error reading var file")
		config.Params.Vars[name] = strings.TrimSpace(string(varFile))
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, config.Params.Vars)
	check(err, "Error executing template")

	outFile, err := os.Create(templPath)
	check(err, "Error creating output file")
	defer outFile.Close()
	_, err = outFile.Write(buf.Bytes())
	check(err, "Error writing output file")

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

	fmt.Println(out.String())
}
