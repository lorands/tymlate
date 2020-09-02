package gen

import (
	"fmt"
	"github.com/lorands/tymlate/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type TemplateModel struct {
	SourcePath string
	TargetPath string
	ConfigPath string
	Config config.Config
	StopIfTargetNotEmpty bool
}

func (tm *TemplateModel) isExcluded(source string) bool {
	for _, exc := range tm.Config.Templates.Excludes {
		regex := regexp.MustCompile(exc)
		if regex.MatchString(source) {
			return true
		}
	}
	return false
}

func (tm *TemplateModel) isIncluded(source string) bool {
	for _, incl := range tm.Config.Templates.Includes {
		regex := regexp.MustCompile(incl)
		if regex.MatchString(source) {
			return true
		}
	}
	return false
}

func NewTemplateModel(source string, target string, configPath string, stopIfNotEmpty bool) (error, *TemplateModel) {
	if len(source) < 1 {
		return fmt.Errorf("source must be set"), nil
	}
	if len(target) < 1 {
		return fmt.Errorf("target must be set"), nil
	}

	if stopIfNotEmpty {
		if _, err := os.Stat(target); os.IsExist(err) {
			err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
				return fmt.Errorf("the target directory is not empty: %s", path)
			})
			if err != nil {
				return err, nil
			}
		}
	}

	var model config.Config
	if len(configPath) > 0 {
		readFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		if err = yaml.Unmarshal(readFile, &model); err != nil {
			return err, nil
		}
	}

	return nil, &TemplateModel{
		SourcePath:           source,
		TargetPath:           target,
		ConfigPath:           configPath,
		Config: 			  model,
		StopIfTargetNotEmpty: stopIfNotEmpty,
	}
}
