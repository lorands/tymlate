package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"

	"github.com/lorands/tymlate/config"
)

type Generator struct {
	SourcePath           string
	TargetPath           string
	ConfigPath           string
	Config               config.Config
	StopIfTargetNotEmpty bool
}

func New(source, target, configPath string, stopIfNotEmpty bool) (Generator, error) {
	if len(source) == 0 {
		return Generator{}, fmt.Errorf("source must be set")
	}

	if len(target) == 0 {
		return Generator{}, fmt.Errorf("target must be set")
	}

	if len(configPath) == 0 {
		return Generator{}, fmt.Errorf("config must be set")
	}

	if stopIfNotEmpty {
		if _, err := os.Stat(target); os.IsExist(err) {
			err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
				return fmt.Errorf("the target directory is not empty: %s", path)
			})
			if err != nil {
				return Generator{}, err
			}
		}
	}

	readFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Generator{}, err
	}

	var configYaml config.Config
	if err = yaml.Unmarshal(readFile, &configYaml); err != nil {
		return Generator{}, err
	}

	return Generator{
		SourcePath:           source,
		TargetPath:           target,
		ConfigPath:           configPath,
		Config:               configYaml,
		StopIfTargetNotEmpty: stopIfNotEmpty,
	}, nil
}

func (g Generator) Generate() error {
	if err := os.MkdirAll(g.TargetPath, os.ModePerm); err != nil {
		return fmt.Errorf("generator.Generate().os.MkdirAll()")
	}

	err := filepath.Walk(g.SourcePath, func(relSource string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if g.Config.Templates.IsExcluded(relSource) && !g.Config.Templates.IsIncluded(relSource) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(g.SourcePath, relSource)
		if err != nil {
			return err
		}

		currentTarget := filepath.Join(g.TargetPath, rel)

		context, err := g.prepareContext()
		if err != nil {
			return err
		}

		realTarget, err := g.prepareTargetFilename(context, currentTarget)
		if err != nil {
			return err
		}

		currentTargetDir := filepath.Dir(realTarget)
		if _, err := os.Stat(currentTargetDir); os.IsNotExist(err) {
			if err := os.MkdirAll(currentTargetDir, os.ModePerm); err != nil {
				return err
			}
		}

		tmplSuffix := g.Config.Templates.Suffix
		if strings.HasSuffix(info.Name(), tmplSuffix) { //apply template
			err = processor{
				source:  relSource,
				target:  realTarget,
				context: context,
			}.processTemplate()
			if err != nil {
				return err
			}

			return nil
		}

		if _, err := fsCopy(relSource, realTarget); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("generator.Generate().filepath.Walk()")
	}

	return nil
}

func (g Generator) prepareTargetFilename(context map[string]any, currentTarget string) (string, error) {
	realTarget := currentTarget
	if g.Config.Templates.ProcessFilename { //try to process filename as template
		tpl := template.Must(template.New("currentTarget").Funcs(sprig.FuncMap()).Parse(currentTarget))

		destination := bytes.NewBufferString("")

		if err := tpl.Execute(destination, context); err != nil {
			return "", err
		}

		realTarget = destination.String()
	}

	suf := g.Config.Templates.Suffix
	if strings.HasSuffix(realTarget, suf) {
		realTarget = realTarget[:len(realTarget)-len(g.Config.Templates.Suffix)] //cut off extension (.tmpl) from the end
	}

	return realTarget, nil
}

func (g Generator) prepareContext() (map[string]any, error) {
	context := make(map[string]any)
	context["Env"] = envToMap()

	for k, v := range g.Config.Data {
		context[k] = v
	}

	for k, v := range g.Config.Include {
		configPath := v
		if !filepath.IsAbs(v) {
			configPath = filepath.Join(filepath.Dir(g.ConfigPath), v) //if relative, it is relative to master config
		}

		yamlConfig, err := readYamlConfig(configPath)
		if err != nil {
			return nil, err
		}
		context[k] = yamlConfig
	}

	return context, nil
}
