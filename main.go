// Developed by Rui Lopes (ruilopes.com). Released under the LGPLv3 license.
//
// See https://bitbucket.org/rgl/docker-native-inspect

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"text/template"
)

var format = flag.String("format", "", "Format the output using the given Go template; e.g. {{.network_state.veth_host}}")
var dockerBinary = flag.String("docker", "docker", "Path to the Docker binary")
var dockerHome = flag.String("docker-home", "/var/lib/docker", "Path to the Docker home")

func main() {
	log.SetFlags(0)

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		log.Fatalf("\nERROR You MUST pass exactly two positional arguments: [state|container] [partial container id or name]\n")
	}

	containerId, err := getContainerId(flag.Arg(1))

	if err != nil {
		log.Fatalf("ERROR %s\n", err)
	}

	filePath := path.Join(*dockerHome, "execdriver/native", containerId, flag.Arg(0)+".json")

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatalf("ERROR Failed to read from %s: %v\n", filePath, err)
	}

	if len(*format) > 0 {
		formatTemplate, err := template.New("").Parse(*format)
		if err != nil {
			log.Fatalf("ERROR Failed to parse template: %s\n", err)
		}

		var value interface{}

		if err := json.Unmarshal(data, &value); err != nil {
			log.Fatalf("ERROR Failed to parse JSON file %s: %v\n", filePath, err)
		}

		if err := formatTemplate.Execute(os.Stdout, value); err != nil {
			log.Fatalf("ERROR Failed to execute the template: %v\n", err)
		}
	} else {
		var formatted bytes.Buffer

		if err := json.Indent(&formatted, data, "", "  "); err != nil {
			log.Fatalf("ERROR Failed to Indent JSON %s\n", err)
		}

		formatted.WriteTo(os.Stdout)
	}
}

func getContainerId(partialNameOrId string) (string, error) {
	var stdErr, stdOut bytes.Buffer

	c := exec.Command(*dockerBinary, "inspect", "--format='{{.Id}}'", partialNameOrId)
	c.Stderr = &stdErr
	c.Stdout = &stdOut

	err := c.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus := exitError.Sys().(syscall.WaitStatus); ok {
				return "", fmt.Errorf("Failed to get docker container id. exitCode=%v stdOut=%v stdErr=%v", waitStatus.ExitStatus(), stdOut.String(), stdErr.String())
			}
		}
		return "", fmt.Errorf("Failed to get docker container id. err=%v stdOut=%v stdErr=%v", err, stdOut.String(), stdErr.String())
	}

	return strings.TrimSpace(stdOut.String()), nil
}
