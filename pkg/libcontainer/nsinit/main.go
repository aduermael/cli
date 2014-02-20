package main

import (
	"encoding/json"
	"github.com/dotcloud/docker/pkg/libcontainer"
	"log"
	"os"
)

func main() {
	container, err := loadContainer()
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "exec":
		exitCode, err := execCommand(container)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(exitCode)
	case "init":
		if err := initCommand(container, os.Args[2]); err != nil {
			log.Fatal(err)
		}
	}
}

func loadContainer() (*libcontainer.Container, error) {
	f, err := os.Open("container.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var container *libcontainer.Container
	if err := json.NewDecoder(f).Decode(&container); err != nil {
		return nil, err
	}
	return container, nil
}
