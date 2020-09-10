package cmd

import (
	"fmt"
	gen "github.com/lorands/tymlate/generator"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	source string
	target string
	confFile string
	rootCmd = &cobra.Command{
		Use:   "tymlate",
		Short: "tymlate is a directory structure-aware file generator with go templating",
		Long: `A long long long introduction...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if confFile == "" {
				//most probably we have .tymlate.yml
				confFile = filepath.Join(source, ".tymlate.yml")
				info, err := os.Stat(confFile)
				if os.IsNotExist(err) {
					return err
				}
				if info.IsDir() {
					return fmt.Errorf("provide path to config file not to directory")
				}
			}

			err, templateModel := gen.NewTemplateModel(source, target, confFile, false)
			if err != nil {
				return err
			}

			//read up config to config
			dss, _ := cmd.Flags().GetStringSlice("datasource")

			if len(dss) > 0 {
				if templateModel.Config.Include == nil {
					templateModel.Config.Include = make(map[string]string)
				}
				for _, v := range dss {
					parts := strings.Split(v, "=")
					if len(parts) != 2 {
						return fmt.Errorf("datasource must be given in a form of name=pathToFile")
					}
					templateModel.Config.Include[parts[0]] = parts[1]
				}
			}

			return templateModel.Generate()

		},
	}
)


func init() {
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "path to template source folder")
	rootCmd.MarkFlagRequired("source")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "path to target folder")
	rootCmd.MarkFlagRequired("target")
	rootCmd.PersistentFlags().StringVarP(&confFile, "configuration", "c", "", "path to configuration file")
	rootCmd.PersistentFlags().StringSliceP("datasource", "d", nil, "Datasource in name=file format")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}