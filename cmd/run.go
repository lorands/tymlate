package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lorands/tymlate/generator"
)

var (
	source   string
	target   string
	confFile string

	rootCmd = &cobra.Command{
		Use:   "tymlate",
		Short: "tymlate is a directory structure-aware file generator with go templating",
		Long:  `A long long long introduction...`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if confFile == "" {
				//most probably we have .tymlate.yml
				confFile = filepath.Join(source, ".tymlate.yml")
				info, err := os.Stat(confFile)
				if os.IsNotExist(err) {
					return err
				}

				if info.IsDir() {
					return fmt.Errorf("provide path to config file, not a directory")
				}
			}

			generatorModel, err := generator.New(source, target, confFile, false)
			if err != nil {
				return err
			}

			additionalDataSourcesIncluded, err := cmd.Flags().GetStringSlice("datasource")
			if err != nil {
				return err
			}

			if len(additionalDataSourcesIncluded) > 0 {
				if generatorModel.Config.Include == nil {
					generatorModel.Config.Include = make(map[string]string)
				}

				for _, v := range additionalDataSourcesIncluded {
					parts := strings.Split(v, "=")
					if len(parts) != 2 {
						return fmt.Errorf("datasource must be given in a form of name=pathToFile")
					}

					generatorModel.Config.Include[parts[0]] = parts[1]
				}
			}

			if err := generatorModel.Generate(); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "path to template source folder")
	_ = rootCmd.MarkFlagRequired("source")

	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "path to target folder")
	_ = rootCmd.MarkFlagRequired("target")

	rootCmd.PersistentFlags().StringVarP(&confFile, "configuration", "c", "", "path to configuration file")
	rootCmd.PersistentFlags().StringSliceP("datasource", "d", nil, "Datasource in name=file format")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
