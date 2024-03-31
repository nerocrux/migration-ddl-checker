package main

import (
	"log"
	"os"
	"strings"

	"github.com/nerocrux/migration-ddl-checker/analyzer"
	"github.com/nerocrux/migration-ddl-checker/ddl"
	"github.com/urfave/cli/v2"
)

type Analyzer interface {
	Analyze(contents string) (bool, error)
	IsHazardousDDL(d ddl.DDL) bool
}

func main() {
	var (
		targetFiles   cli.StringSlice
		syntax        string
		hazardousDDLs cli.StringSlice
	)
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "target-files",
				Usage:       "changes in these files will be checked",
				Destination: &targetFiles,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "syntax",
				Usage:       "available values: mysql, postgres, spanner",
				Destination: &syntax,
				Required:    false,
			},
			&cli.StringSliceFlag{
				Name:        "hazardous-ddl",
				Usage:       "specify hazardous ddl with comma separated values",
				Destination: &hazardousDDLs,
				Required:    false,
				DefaultText: "ALTER,DROP",
			},
		},
		Action: func(ctx *cli.Context) error {
			ddls := ddl.FromConfig(hazardousDDLs.Value())

			var a Analyzer
			switch syntax {
			case "mysql":
				a = analyzer.NewMysqlAnalyzer(ddls)
			case "postgresql":
				a = analyzer.NewPostgresqlAnalyzer(ddls)
			case "spanner":
				a = analyzer.NewSpannerAnalyzer(ddls)
			default:
				// do nothing
			}

			hazardousFiles, err := analyze(a, targetFiles.Value())
			if err != nil {
				return err
			}
			os.Stdout.WriteString(strings.Join(hazardousFiles, ",") + "\n")
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
