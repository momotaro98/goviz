package main

import (
	"os"
	"testing"
)

func TestOption(t *testing.T) {
	os.Args = []string{
		"goviz",
		"-i",
		"github.com/momotaro98/goviz",
	}
	option, err := getOptions()
	if err != nil {
		t.Errorf("option error %s", err)
	}
	if option.InputDir != "github.com/momotaro98/goviz" {
		t.Error("incorrect parse")
	}
}

func TestOption2(t *testing.T) {
	os.Args = []string{
		"goviz",
		"-i",
		"gitdaichi/goviz",
	}
	if p := process(); p != 1 {
		t.Errorf("exit status expect 1 ,but %d", p)
	}
}
