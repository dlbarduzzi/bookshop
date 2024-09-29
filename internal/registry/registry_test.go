package registry

import "testing"

func TestInit(t *testing.T) {
	reg, err := NewRegistry()
	if err != nil {
		t.Fatal(err)
	}

	testVar := "TESTING_ENV_VAR"
	testValue := "test-var"

	t.Setenv(testVar, testValue)
	envVar := reg.GetString(testVar)

	if envVar != testValue {
		t.Errorf("expected env var to be %s; got %s", testValue, envVar)
	}

	configPath := "test-path"
	SetConfigPath(configPath)

	if _configPath != configPath {
		t.Errorf("expected config path to be %s; got %s", configPath, _configPath)
	}

	configName := "test-name"
	SetConfigName(configName)

	if _configName != configName {
		t.Errorf("expected config name to be %s; got %s", configName, _configName)
	}

	configType := "test-path"
	SetConfigType(configType)

	if _configType != configType {
		t.Errorf("expected config type to be %s; got %s", configType, _configType)
	}
}
