package main

import (
	"log"
	"os"
	"strings"

	"github.com/nerocrux/migration-ddl-checker/analyzer"
	"github.com/urfave/cli/v2"
)

type Analyzer interface {
	Analyze(contents string) (bool, error)
}

func main() {
	var (
		targetFiles cli.StringSlice
		syntax      string
	)
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "target-files",
				Usage:       "changes in these files will be tolerated",
				Destination: &targetFiles,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "syntax",
				Usage:       "available values: mysql, postgres, spanner",
				Destination: &syntax,
				Required:    false,
			},
		},
		Action: func(ctx *cli.Context) error {
			var a Analyzer
			switch syntax {
			case "mysql":
				a = analyzer.NewMysqlAnalyzer()
			case "postgresql":
				a = analyzer.NewPostgresqlAnalyzer()
			case "spanner":
				a = analyzer.NewSpannerAnalyzer()
			default:
				// do nothing
			}

			hazardousFiles, err := analyze(a, targetFiles.Value())
			if err != nil {
				return err
			}
			os.Stdout.WriteString(strings.Join(hazardousFiles, "\n") + "\n")
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func analyze(a Analyzer, targetFiles []string) ([]string, error) {
	var hazardousFiles []string
	if a == nil {
		// If analyzer not exist, treat all target files as hazardous because we cannot analyze them.
		return targetFiles, nil
	}

	for _, file := range targetFiles {
		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		isHazardous, err := a.Analyze(string(contents))
		if err != nil {
			return nil, err
		}
		if isHazardous {
			hazardousFiles = append(hazardousFiles, file)
		}
	}
	return hazardousFiles, nil
}
