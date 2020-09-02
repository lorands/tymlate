package gen

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type genConfig struct {
	source string
	target string
	context map[string]map[string]interface{}
}

func (tm *TemplateModel) Generate() error {
	var excludeMode = true

	if len(tm.Config.Templates.Excludes) > 0 && len(tm.Config.Templates.Includes)>0 {
		println("Both excludes and includes are defined in config file, so only includes will be used.")
		excludeMode = false
	}

	osErr:=os.MkdirAll(tm.TargetPath, os.ModePerm)
	if osErr != nil {
		return osErr
	}

	err := filepath.Walk(tm.SourcePath, func(relSource string, info os.FileInfo, err error) error {
		var skip = false
		if excludeMode {
			skip = tm.isExcluded(relSource)
		} else { //include
			skip = !tm.isIncluded(relSource)
		}

		if ! skip {
			rel, err := filepath.Rel(tm.SourcePath, relSource)
			if err != nil {
				return err
			}
			currentTarget := filepath.Join(tm.TargetPath, rel)
			if ! info.IsDir() {
				currentTargetDir := filepath.Dir(currentTarget)
				if _, err := os.Stat(currentTargetDir); os.IsNotExist(err) {
					if err := os.MkdirAll(currentTargetDir, os.ModePerm); err != nil {
						return err
					}
				}

				tmplSuffix := tm.Config.Templates.Suffix
				if strings.HasSuffix(info.Name(), tmplSuffix) { //apply template
					realTarget := currentTarget[:len(currentTarget)-len(tmplSuffix)] //cut off .tmpl from the end
					context, cErr := tm.prepareContext()
					if cErr != nil {
						return cErr
					}
					cErr = genConfig{
						source:  relSource,
						target:  realTarget,
						context: context,
					}.processTemplate()
					if cErr != nil {
						return cErr
					}
				} else { //simple copy
					if _, err := fsCopy(relSource, currentTarget); err != nil {
						return err
					}
				}
			}
		}

		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (tm *TemplateModel) prepareContext() (map[string]map[string]interface{}, error) {
	//env
	context := make(map[string]map[string]interface{})
	//context := tm.Config.Data
	envMap, _ := envToMap()
	envIMap := make(map[string]interface{})
	for k,v := range *envMap {
		envIMap[k] = v
	}
	context["Env"] = envIMap
	for k,v := range tm.Config.Data {
		context[k] = v
	}


	for k,v := range tm.Config.Include {
		var cPath string
		if filepath.IsAbs(v) {
			cPath = v
		} else {
			cPath = filepath.Join(filepath.Dir(tm.ConfigPath), v) //if relative, it is relative to master config
		}
		if yamlConfig, err := readYamlConfig(cPath); err != nil {
			return nil, err
		} else {
			context[k] = yamlConfig
		}
	}

	return context, nil
}

func readYamlConfig(yamlFilePath string) (map[string]interface{}, error)  {
	var body map[string]interface{}

	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}


func (c genConfig) processTemplate() error {
	fmt.Printf("%s --> %s\n", c.source, c.target)

	tpl := template.Must(
		template.New(filepath.Base(c.source)).Funcs(sprig.FuncMap()).ParseFiles(c.source))

	destination, err := os.Create(c.target)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err := tpl.Execute(destination, c.context); err != nil {
		return err
	}

	return nil
}

func fsCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func envToMap() (*map[string]string, error) {
	envMap := make(map[string]string)
	var err error

	for _, v := range os.Environ() {
		part := strings.Split(v, "=")
		envMap[part[0]] = part[1]
	}

	return &envMap, err
}