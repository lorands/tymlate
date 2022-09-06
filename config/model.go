package config

import "regexp"

type Config struct {
	Data      map[string]map[string]any `yaml:"data,omitempty"`
	Include   map[string]string         `yaml:"include,omitempty"`
	Templates Templates                 `yaml:"templates,omitempty"`
}

type Templates struct {
	Includes        []string `yaml:"includes,omitempty"`
	Excludes        []string `yaml:"excludes,omitempty"`
	Suffix          string   `default:".tmpl" yaml:"suffix,omitempty"`
	ProcessFilename bool     `yaml:"processFilename,omitempty"`
}

func (t Templates) IsExcluded(source string) bool {
	for _, exclude := range t.Excludes {
		regex := regexp.MustCompile(exclude)
		if regex.MatchString(source) {
			return true
		}
	}

	return false
}

func (t Templates) IsIncluded(source string) bool {
	for _, include := range t.Includes {
		regex := regexp.MustCompile(include)
		if regex.MatchString(source) {
			return true
		}
	}

	return false
}
