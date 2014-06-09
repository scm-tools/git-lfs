package commands

import (
	"bytes"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSmudge(t *testing.T) {
	repo := NewRepository(t, "empty")
	defer repo.Test()

	prePushHookFile := filepath.Join(repo.Path, ".git", "hooks", "pre-push")

	cmd := repo.Command("smudge")
	cmd.Input = bytes.NewBufferString("# git-media\nSOMEOID")
	cmd.Output = "whatever"

	cmd.Before(func() {
		path := filepath.Join(repo.Path, ".git", "media", "SO", "ME")
		file := filepath.Join(path, "SOMEOID")
		assert.Equal(t, nil, os.MkdirAll(path, 0755))
		assert.Equal(t, nil, ioutil.WriteFile(file, []byte("whatever\n"), 0755))
	})

	cmd.After(func() {
		// assert hooks
		stat, err := os.Stat(prePushHookFile)
		assert.Equal(t, nil, err)
		assert.Equal(t, false, stat.IsDir())
	})

	cmd = repo.Command("smudge")
	cmd.Input = bytes.NewBufferString("# git-media\nSOMEOID")
	cmd.Output = "whatever"
	customHook := []byte("echo 'yo'")

	cmd.Before(func() {
		path := filepath.Join(repo.Path, ".git", "media", "SO", "ME")
		file := filepath.Join(path, "SOMEOID")
		assert.Equal(t, nil, os.MkdirAll(path, 0755))
		assert.Equal(t, nil, ioutil.WriteFile(file, []byte("whatever\n"), 0755))
		assert.Equal(t, nil, ioutil.WriteFile(prePushHookFile, customHook, 0755))
	})

	cmd.After(func() {
		by, err := ioutil.ReadFile(prePushHookFile)
		assert.Equal(t, nil, err)
		assert.Equal(t, string(customHook), string(by))
	})
}
