package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"goselect/parser"
	"goselect/parser/context"
	"goselect/parser/executor"
	"goselect/parser/source"
	"goselect/parser/writer"
	"os"
	"strings"
)

func newExecuteCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "execute",
		Aliases: []string{"ex"},
		Short:   "Execute a select query",
		Long:    `Execute a select query. Select query syntax: select <attributes> from <source directory> [where <condition>] [order by] [limit]`,
		Example: `
1. goselect execute -q='select filename, absolutepath from .'
2. goselect ex -q='select name, size, extension from . where like(name, results.*) order by 2'
3. goselect ex -q='select name, size, extension from . where or(like(name, results.*), gt(size, 2048)) order by 2 limit 5'
`,
		Run: func(cmd *cobra.Command, args []string) {
			errorColor := "\033[31m"

			buildOptions := func() *executor.Options {
				nestedTraversal, _ := cmd.Flags().GetBool("nestedTraversal")
				ignoreTraversal, _ := cmd.Flags().GetStringSlice("skipDirectoryTraversal")

				options := executor.NewDefaultOptions()
				if nestedTraversal {
					options.EnableNestedTraversal()
				} else {
					options.DisableNestedTraversal()
				}
				options.DirectoriesToIgnoreTraversal(ignoreTraversal)
				return options
			}
			executeQuery := func(cmd *cobra.Command) (*executor.EvaluatingRows, *parser.SelectQuery, error) {
				rawQuery, _ := cmd.Flags().GetString("query")
				newContext := context.NewContext(context.NewFunctions(), context.NewAttributes())
				newParser, err := parser.NewParser(rawQuery, newContext)
				if err != nil {
					return nil, nil, err
				}
				query, err := newParser.Parse()
				if err != nil {
					return nil, nil, err
				}
				rows, err := executor.NewSelectQueryExecutor(query, newContext, buildOptions()).Execute()
				if err != nil {
					return nil, nil, err
				}
				return rows, query, nil
			}
			formatter := func(cmd *cobra.Command) (writer.Formatter, string, error) {
				exportFormat, _ := cmd.Flags().GetString("format")
				switch strings.ToLower(exportFormat) {
				case "json":
					return writer.NewJsonFormatter(), strings.ToLower(exportFormat), nil
				case "html":
					return writer.NewHtmlFormatter(), strings.ToLower(exportFormat), nil
				case "table":
					minWidth, _ := cmd.Flags().GetUint16("minWidth")
					maxWidth, _ := cmd.Flags().GetUint16("maxWidth")

					if minWidth == 0 && maxWidth == 0 {
						return writer.NewTableFormatter(), strings.ToLower(exportFormat), nil
					}
					if minWidth == 0 && maxWidth != 0 {
						return writer.NewTableFormatterWithWidthOptions(writer.NewAttributeWidthOptions(
							writer.UnspecifiedMinWidth,
							int(maxWidth),
						)), strings.ToLower(exportFormat), nil
					}
					if minWidth != 0 && maxWidth == 0 {
						return writer.NewTableFormatterWithWidthOptions(writer.NewAttributeWidthOptions(
							int(minWidth),
							writer.UnspecifiedMaxWidth,
						)), strings.ToLower(exportFormat), nil
					}
					return writer.NewTableFormatterWithWidthOptions(
						writer.NewAttributeWidthOptions(int(minWidth), int(maxWidth)),
					), strings.ToLower(exportFormat), nil
				default:
					return nil, "", fmt.Errorf(ErrorMessageInvalidExportFormat, SupportedExportFormats())
				}
			}
			writeToConsole := func(formattedResult string) {
				buffer := new(bytes.Buffer)
				consoleWriter := writer.NewWriter(buffer)
				_ = consoleWriter.Write(formattedResult)
				cmd.Print(buffer.String())
			}
			writeToFile := func(format, formattedResult string) error {
				directoryPath, _ := cmd.Flags().GetString("path")
				if strings.EqualFold(format, "table") {
					return errors.New(ErrorMessageAttemptedToExportTableToFile)
				}
				directoryPath, err := source.ExpandDirectoryPath(directoryPath)
				if err != nil {
					return err
				}
				if filePath, err := os.Stat(directoryPath); err != nil {
					return err
				} else {
					if !filePath.IsDir() {
						return errors.New(ErrorMessageExpectedFilePathToBeADirectory)
					}
					pathSeparator := string(os.PathSeparator)
					filePath := directoryPath + pathSeparator + fmt.Sprintf("results.%v", format)
					if strings.HasSuffix(directoryPath, pathSeparator) {
						filePath = directoryPath + fmt.Sprintf("results.%v", format)
					}
					fileWriter, err := writer.NewFileWriter(filePath)
					if err != nil {
						return err
					}
					return fileWriter.Write(formattedResult)
				}
			}
			write := func(format string, formattedResult string) error {
				directoryPath, _ := cmd.Flags().GetString("path")
				if len(directoryPath) == 0 {
					writeToConsole(formattedResult)
					return nil
				}
				return writeToFile(format, formattedResult)
			}
			run := func() {
				rows, query, err := executeQuery(cmd)
				if err != nil {
					cmd.Println(errorColor, err)
					return
				}
				exportFormatter, format, err := formatter(cmd)
				if err != nil {
					cmd.Println(errorColor, err)
					return
				}
				err = write(format, exportFormatter.Format(query.Projections, rows))
				if err != nil {
					cmd.Println(errorColor, err)
					return
				}
			}
			run()
		},
	}
}

func SupportedExportFormats() []string {
	return []string{"json", "html", "table"}
}

func init() {
	executeCmd := newExecuteCommand()
	rootCmd.AddCommand(executeCmd)
	executeCmd.PersistentFlags().StringP(
		"query",
		"q",
		"",
		"specify the query. Use --query=<query> or -q=<query>",
	)
	executeCmd.PersistentFlags().BoolP(
		"nestedTraversal",
		"n",
		true,
		"specify if nested directories should be traversed. Use --nestedTraversal=<true/false> or -n=<true/false>",
	)
	executeCmd.PersistentFlags().StringSliceP(
		"skipDirectoryTraversal",
		"s",
		[]string{".git", ".github"},
		"specify the directory names that should not be traversed. Use --skipDirectoryTraversal=<directory> or -s=<directory>. Multiple directory names can be passed by using --skipDirectoryTraversal=.git --skipDirectoryTraversal=.github",
	)
	executeCmd.PersistentFlags().StringP(
		"format",
		"f",
		"table",
		"specify the export format. Supported values include: json, html and table. Use --format=<format>",
	)
	executeCmd.PersistentFlags().StringP(
		"path",
		"p",
		"",
		"specify the directory path to export the results. Use --path=<directoryPath>",
	)
	executeCmd.PersistentFlags().Uint16P(
		"minWidth",
		"m",
		0,
		"specify the minimum character length to be used for each attribute/column. This flag is relevant only for the table format and will be needed only if the default formatting breaks. For the best results, use minWidth and maxWidth together. Use --minWidth=<value greater than zero>",
	)
	executeCmd.PersistentFlags().Uint16P(
		"maxWidth",
		"x",
		0,
		"specify the maximum character length to be used for each attribute/column. This flag is relevant only for the table format and will be needed only if the default formatting breaks. For the best results, use minWidth and maxWidth together. Use --maxWidth=<value greater than zero>",
	)
}
