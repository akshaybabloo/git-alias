package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func execute(args string) string {
	actual := new(bytes.Buffer)
	rootCmd := RootCmd()
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	if args != "" {
		rootCmd.SetArgs(strings.Split(args, " "))
	}
	err := rootCmd.Execute()
	if err != nil {
		return err.Error()
	}
	return actual.String()
}

func Test_Main_Fail(t *testing.T) {
	actual := execute("--config blah")
	expected := "error: blah file does not exist"
	assert.Equal(t, expected, actual)
}

func Test_Main_Pass(t *testing.T) {
	file, err := os.CreateTemp("", "local-gitconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	config := `[alias]
	co = checkout
`
	_, err = file.WriteString(config)

	actual := execute(fmt.Sprintf("--config %s", file.Name()))
	expected := "co"
	assert.Contains(t, expected, actual)
}

func Test_Main_Search_Pass(t *testing.T) {
	file, err := os.CreateTemp("", "local-gitconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	config := `[alias]
	co = checkout
`
	_, err = file.WriteString(config)

	actual := execute(fmt.Sprintf("--config %s -s checkout", file.Name()))
	expected := "checkout"
	assert.Contains(t, expected, actual)
}
