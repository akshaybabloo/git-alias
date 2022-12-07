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

var SearchString string
var Sort bool
var ConfigPath string

var RootCmd = &cobra.Command{
	Use:   "git-alias [file] (default: ~/.gitconfig)",
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

		config, err := ini.Load(ConfigPath)
		if err != nil {
			return fmt.Errorf("error: could not load file")
		}

		section, err := config.GetSection("alias")
		if err != nil {
			return fmt.Errorf("error: could not find alias section")
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Alias", "Command"})
		if SearchString != "" {
			c := color.New(color.BgCyan, color.Bold)
			index := 1
			for _, key := range section.Keys() {
				if strings.Contains(key.Value(), SearchString) {
					valueIndex := strings.Index(key.Value(), SearchString)
					if valueIndex != -1 {
						t.AppendRow(table.Row{index, key.Name(), key.Value()[0:valueIndex] + c.Sprint(key.Value()[valueIndex:valueIndex+len(SearchString)]) + key.Value()[valueIndex+len(SearchString):]})
						t.AppendSeparator()
						index++
					}
				}
			}
			t.SetTitle(fmt.Sprintf("Found %d aliases", t.Length()))
		} else {
			for i, key := range section.Keys() {
				t.AppendRow(table.Row{i + 1, key.Name(), key.Value()})
				t.AppendSeparator()
			}
		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{Name: "Command", WidthMax: 80},
		})
		if Sort {
			t.SortBy([]table.SortBy{{Name: "Alias", Mode: table.Asc}})
		}

		t.Render()

		return nil
	},
}

func main() {

	RootCmd.PersistentFlags().StringVarP(&SearchString, "search", "s", "", "Search for aliases containing the given string")
	RootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "~/.gitconfig", "Path to git config file")
	RootCmd.PersistentFlags().BoolVar(&Sort, "sort", false, "Sort aliases by alias name")

	RootCmd.Version = version
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
