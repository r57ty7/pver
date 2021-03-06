package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func Test_npmCmd_Execute(t *testing.T) {
	setUp(t)

	mockFvm := &mockFvm{
		version: "1.0.0",
	}
	npmFvm = mockFvm

	cmd := newNpmCmd()

	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	_ = cmd.Execute()
	got, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	want := "Version: 1.0.0\n"

	if string(got) != want {
		t.Fatalf("want '%s', got '%s'", want, got)
	}

}

func Test_npmCmd_Execute_Update(t *testing.T) {
	setUp(t)

	mockFvm := &mockFvm{
		version: "1.0.0",
	}
	npmFvm = mockFvm

	cmd := newNpmCmd()

	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"-u", "1.2.3"})
	_ = cmd.Execute()
	got, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}

	wantMessage := "Version: 1.0.0\nUpdated to => 1.2.3\n"
	wantVersion := "1.2.3"

	if string(got) != wantMessage {
		t.Fatalf("want '%s', got '%s'", wantMessage, got)
	}

	if mockFvm.version != wantVersion {
		t.Fatalf("want '%s', got '%s'", wantVersion, mockFvm.version)
	}

}
