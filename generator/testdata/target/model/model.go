package config

type Model struct {
    Strings []string `yaml:"includes,omitempty"`
    Data map[string]map[string]interface{} `yaml:"data,omitempty"`
}
