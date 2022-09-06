package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func readYamlConfig(yamlFilePath string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return nil, fmt.Errorf("util.readYamlConfig().ioutil.ReadFile()")
	}

	var body map[string]interface{}
	if err := yaml.Unmarshal(yamlFile, &body); err != nil {
		return nil, fmt.Errorf("util.readYamlConfig().yaml.Unmarshal()")
	}

	return body, nil
}

func fsCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("util.fsCopy(): %s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("util.fsCopy().os.Open()")
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, fmt.Errorf("util.fsCopy().os.Create()")
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return 0, fmt.Errorf("util.fsCopy().io.Copy()")
	}

	return nBytes, nil
}

func envToMap() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		part := strings.Split(v, "=")
		envMap[part[0]] = part[1]
	}

	return envMap
}
