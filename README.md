# Gomponent

A Go module inspired by Alessandra Sierra's Clojure [Component Library](https://github.com/stuartsierra/component)

# Usage

## Create a Component

```go

import (
    "github.com/leoiacovini/gomponent"
)

// Define your component as a struct
type Config struct {
	started bool
}

// Implement both Start and Stop method from the Component Interface
// if something went wrong during any of those steps, it is possible
// to return an Error that will interrupt the rest of the execution
func (config *Config) Start() error {
	config.started = true
	return nil
}

func (config *Config) Stop() error {
	config.started = false
	return nil
}

// Factory function that should return a new instance reference of your
// component. Notice that the signature should be the following.
func ConfigFactory(system gomponent.DepsMap) gomponent.Component {
	return &Config{}
}
```

## Building a System

```go

import (
    "github.com/leoiacovini/gomponent"
)

func main() {
    // SystemMap just like a regular map, that expect string as keys, and 
    // and gomponent.SystemComponent pointers as values
    system := gomponent.SystemMap{
		"config": &gomponent.SystemComponent{
            // use your previously defined Factory function here (or you can also pass an
            // anonymous/lambda function instead)
			Factory: ConfigFactory,
            // your dependencies list, should be the name of other components present in this
            // SystemMap. Components are initialized orderly, and the Factory function is provided
            // with all the Deps specified here in this list
			Deps:    []string{},
		},
	}

    // Starts the provided SystemMap, if something goes wrong started will be nil, and
    // there will be an Error in the `err` variable
    started, err := gomponent.StartSystem(system)

    // To get any specific component from the started system. If the component if not found
    // err will be filled.
    config, err := gomponent.GetComponent[*Config](system, "config")

    //... do whatever you need to do with your system

    // Stops a started system in the correct order (inverse of the started one). Similarly,
    // if something goes wrong during this proccess `err` will contain its details
    // *Warning*: if something fails during the system stop, it maybe in an inconsistent
    // state (partially stopped, and is not suitable to be continued used)
    stopped, err := gomponent.StopSystem(started)
}

```