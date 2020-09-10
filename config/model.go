package config

//import "gopkg.in/dealancer/validate.v2"
//
//type Include struct {
//	Data map[string]string
//}

type Templates struct {
	Includes []string `yaml:"includes,omitempty"`
	Excludes []string `yaml:"excludes,omitempty"`
	Suffix string `default:".tmpl" yaml:"suffix,omitempty"`
}

type Config struct {
	Data map[string]map[string]interface{} `yaml:"data,omitempty"`
	Include map[string]string
	//Include Include `yaml:"include,omitempty"`
	Templates Templates `yaml:"templates,omitempty"`
}
