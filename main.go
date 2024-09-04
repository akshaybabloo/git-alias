package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var version = "dev"

var Sort bool
var ConfigPath string

// RootCmd is the root command for the application.
func RootCmd() *cobra.Command {

	var command = &cobra.Command{
		Use:   "git-alias [flags]",
		Short: "Over-engineered Git alias pretty printer",
		RunE: func(cmd *cobra.Command, args []string) error {

			if strings.HasPrefix(ConfigPath, "~/") {
				homeDir, _ := os.UserHomeDir()
				ConfigPath = filepath.Join(homeDir, ConfigPath[2:])
			}
			_, err := os.Stat(ConfigPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("error: %s file does not exist", ConfigPath)
				}
				return fmt.Errorf("error: could not stat file")
			}

			config, err := ini.LoadSources(ini.LoadOptions{
				SpaceBeforeInlineComment: true,
			}, ConfigPath)
			if err != nil {
				return fmt.Errorf("error: could not load file")
			}

			section, err := config.GetSection("alias")
			if err != nil {
				return fmt.Errorf("error: could not find alias section")
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"#", "Alias", "Command", "Description"})
			if len(args) > 0 {
				c := color.New(color.BgCyan, color.Bold)
				index := 1
				searchString := args[0]
				for _, key := range section.Keys() {
					if strings.Contains(key.Value(), searchString) || strings.Contains(key.Name(), searchString) {
						valueIndex := strings.Index(key.Value(), searchString)
						if valueIndex != -1 {
							t.AppendRow(table.Row{index, key.Name(), key.Value()[0:valueIndex] + c.Sprint(key.Value()[valueIndex:valueIndex+len(searchString)]) + key.Value()[valueIndex+len(searchString):], key.Comment})
							t.AppendSeparator()
							index++
						}
						keyIndex := strings.Index(key.Name(), searchString)
						if keyIndex != -1 {
							t.AppendRow(table.Row{index, key.Name()[0:keyIndex] + c.Sprint(key.Name()[keyIndex:keyIndex+len(searchString)]) + key.Name()[keyIndex+len(searchString):], key.Value(), key.Comment})
							t.AppendSeparator()
							index++
						}
					}
				}
				t.SetTitle(fmt.Sprintf("Found %d aliases", t.Length()))
			} else {
				for i, key := range section.Keys() {
					t.AppendRow(table.Row{i + 1, key.Name(), key.Value(), key.Comment})
					t.AppendSeparator()
				}
			}
			t.SetColumnConfigs([]table.ColumnConfig{
				{Name: "Command", WidthMax: 60},
			})
			if Sort {
				t.SortBy([]table.SortBy{{Name: "Alias", Mode: table.Asc}})
			}

			t.Render()

			return nil
		},
	}

	command.PersistentFlags().StringVar(&ConfigPath, "config", "~/.gitconfig", "Path to git config file")
	command.PersistentFlags().BoolVar(&Sort, "sort", false, "Sort aliases by alias name")

	command.Version = version

	return command
}

func main() {
	rootCmd := RootCmd()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
