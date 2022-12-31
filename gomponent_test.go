package gomponent

import (
	"fmt"
	"testing"
)

// Fixtures
type Config struct {
	started bool
}

func ConfigFactory(system DepsMap) Component {
	return &Config{}
}

func (config *Config) Start() error {
	config.started = true
	return nil
}

func (config *Config) Stop() error {
	config.started = false
	return nil
}

func newSystem() SystemMap {
	return SystemMap{
		"testComponent1": &SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{},
		},
		"testComponent2": &SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{"testComponent1"},
		},
		"testComponent3": &SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{"testComponent2"},
		},
		"testComponent4": &SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{"testComponent1", "testComponent5"},
		},
		"testComponent5": &SystemComponent{
			Factory: ConfigFactory,
			Deps:    []string{"testComponent3"},
		},
	}
}

func compareSlice(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		if v1 != s2[i] {
			return false
		}
	}
	return true
}

func TestReverse(t *testing.T) {
	actual := reverse([]string{"1", "2", "3", "4"})
	expected := []string{"4", "3", "2", "1"}
	if !compareSlice(actual, expected) {
		t.Errorf("Slice wrongly reversed. Actual %v, Expected: %v", actual, expected)
	}
}

// Tests
func TestToposort(t *testing.T) {
	t.Run("successful toposort run", func(t *testing.T) {
		testSystem := newSystem()
		sorted, err := toposort(testSystem)
		expected := []string{"testComponent1", "testComponent2", "testComponent3", "testComponent5", "testComponent4"}
		if err != nil {
			t.Errorf("Error when executing toposort: %v", err)
			return
		}
		if !compareSlice(sorted, expected) {
			t.Errorf("Sorted is not correct\nexpected: %v\nactual: %v", expected, sorted)
		}
	})
	t.Run("error on toposort - invalid dependency", func(t *testing.T) {
		brokenSystem := SystemMap{
			"comp1": &SystemComponent{
				Factory: ConfigFactory,
				Deps:    []string{"notFound"},
			},
		}
		sorted, err := toposort(brokenSystem)
		if err == nil {
			t.Errorf("Exptect Error, found sorted list: %v", sorted)
		}
	})
}

func TestStartSystem(t *testing.T) {
	testSystem := newSystem()
	started, err := StartSystem(testSystem)
	if err != nil {
		t.Errorf("Error starting system: %v", err)
	}
	config, err := GetComponent[*Config](started, "testComponent3")
	if err != nil {
		t.Errorf("Erro getting component: %v", err)
	}
	fmt.Printf("Started System: %v", config)
}

func TestStopSystem(t *testing.T) {
	testSystem := newSystem()
	StartSystem(testSystem)
	config1, _ := GetComponent[*Config](testSystem, "testComponent3")
	stopped, err := StopSystem(testSystem)
	if err != nil {
		t.Errorf("Error stopping system: %v", err)
	}
	config2, err := GetComponent[*Config](stopped, "testComponent3")
	if err == nil {
		t.Errorf("Found component %v, when should not: %v", config2, err)
	}
	if config1.started {
		t.Errorf("Config component is started after system stop")
	}
}
