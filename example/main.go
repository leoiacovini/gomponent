package main

import (
	"fmt"

	"github.com/leoiacovini/gomponent"
)

type Config struct {
	started bool
}

func (config *Config) Start() error {
	config.started = true
	return nil
}

func (config *Config) Stop() error {
	config.started = false
	return nil
}

func ConfigFactory(system gomponent.DepsMap) gomponent.Component {
	return &Config{}
}

func main() {

	system := gomponent.SystemMap{
		"config": &gomponent.SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{},
		},
	}

	startedSystem, err := gomponent.StartSystem(system)
	if err != nil {
		fmt.Printf("Error starting system, %v\n", err)
	}
	fmt.Println(startedSystem)

	stoppedSystem, err := gomponent.StopSystem(startedSystem)
	if err != nil {
		fmt.Printf("Error stopping system, %v", err)
	}
	fmt.Println(stoppedSystem)
}
