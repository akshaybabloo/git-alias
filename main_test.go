package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func execute(args string) string {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	if args != "" {
		RootCmd.SetArgs(strings.Split(args, " "))
	}
	err := RootCmd.Execute()
	if err != nil {
		return err.Error()
	}
	return actual.String()
}

func Test_Main_Fail(t *testing.T) {
	actual := execute("blah")
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

	actual := execute(file.Name())
	expected := "co"
	assert.Contains(t, expected, actual)
}

func Test_Main_Search(t *testing.T) {
	SearchString = "co"
	file, err := os.CreateTemp("", "local-gitconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	config := `[alias]
	co = checkout
`
	_, err = file.WriteString(config)

	actual := execute(file.Name())
	expected := "co"
	assert.Contains(t, expected, actual)
}
