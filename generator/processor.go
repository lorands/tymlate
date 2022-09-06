package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig"
)

type processor struct {
	source  string
	target  string
	context map[string]any
}

func (c processor) processTemplate() error {
	fmt.Printf("%s --> %s\n", c.source, c.target)

	tpl := template.Must(template.New(filepath.Base(c.source)).Funcs(sprig.FuncMap()).ParseFiles(c.source))

	destination, err := os.Create(c.target)
	if err != nil {
		return fmt.Errorf("processor.processTemplate().os.Create()")
	}
	defer destination.Close()

	if err := tpl.Execute(destination, c.context); err != nil {
		return fmt.Errorf("processor.processTemplate().tpl.Execute()")
	}

	return nil
}
