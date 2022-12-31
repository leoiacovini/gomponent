package gomponent

import (
	"errors"
)

type Component interface {
	Start() error
	Stop() error
}

type DepsMap = map[string]Component

type SystemComponent struct {
	Factory          func(deps DepsMap) Component
	StartedComponent Component
	Deps             []string
}

type SystemMap = map[string]*SystemComponent

// Returns the topological sorted list of components
func toposort(system SystemMap) ([]string, error) {
	visitedNodes := map[string]bool{}
	stack := []string{}

	var sortDep func(subject string) ([]string, error)
	sortDep = func(subject string) ([]string, error) {
		if visitedNodes[subject] {
			return stack, nil
		}
		visitedNodes[subject] = true
		for _, dep := range system[subject].Deps {
			_, ok := system[dep]
			if !ok {
				return nil, errors.New("invalid dependency required")
			}
			nestedStack, err := sortDep(dep)
			if err != nil {
				return nil, err
			}
			stack = nestedStack
		}
		stack = append(stack, subject)
		return stack, nil
	}

	for k := range system {
		newStack, err := sortDep(k)
		if err != nil {
			return nil, err
		}
		stack = newStack
	}

	return stack, nil
}

func GetComponent[T Component](system SystemMap, key string) (T, error) {
	component := system[key].StartedComponent
	if component == nil {
		var nilComp T
		return nilComp, errors.New("component not found")
	}
	return system[key].StartedComponent.(T), nil
}

func contains(deps []string, v string) bool {
	for _, dep := range deps {
		if dep == v {
			return true
		}
	}
	return false
}

func getDepsComponents(system SystemMap, deps []string) DepsMap {
	depsMap := map[string]Component{}
	for k, v := range system {
		if contains(deps, k) {
			depsMap[k] = v.StartedComponent
		}
	}
	return depsMap
}

func StartSystem(system SystemMap) (SystemMap, error) {
	sortedComponents, err := toposort(system)
	if err != nil {
		return nil, err
	}
	for _, k := range sortedComponents {
		currentComponent := system[k]
		// Idempotent Start
		if currentComponent.StartedComponent == nil {
			deps := getDepsComponents(system, currentComponent.Deps)
			comp := currentComponent.Factory(deps)
			err := comp.Start()
			if err != nil {
				return nil, err
			}
			currentComponent.StartedComponent = comp
		}
	}
	return system, nil
}

func reverse(list []string) []string {
	size := len(list) - 1
	reversed := []string{}
	for i := size; i >= 0; i-- {
		reversed = append(reversed, list[i])
	}
	return reversed
}

func StopSystem(system SystemMap) (SystemMap, error) {
	sortedComponents, err := toposort(system)
	if err != nil {
		return nil, err
	}
	for _, k := range reverse(sortedComponents) {
		currentComponent := system[k]
		// Idempotent Stop
		if currentComponent.StartedComponent != nil {
			err := currentComponent.StartedComponent.Stop()
			if err != nil {
				return nil, err
			}
			currentComponent.StartedComponent = nil
		}
	}
	return system, nil
}
