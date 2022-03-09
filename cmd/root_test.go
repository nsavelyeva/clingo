package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory:
	// in IDE click on green double triangles to execute package tests, or
	// in terminal navigate to repository root and execute 'go test ./...'
	tmpDir, e1 := ioutil.TempDir("..", "tmp-test-clingo")
	eventsFileName := filepath.Join(tmpDir, "events.json")
	require.NoError(t, e1, fmt.Sprintf("error creating a temporary test folder %s", tmpDir))

	defer func() {
		e2 := os.Remove(eventsFileName)
		require.NoError(t, e2, fmt.Sprintf("error removing %s", eventsFileName))
		e2 = os.Remove(tmpDir)
		require.NoError(t, e2, fmt.Sprintf("error removing temporary test folder %s", tmpDir))
	}()

	testDir, e3 := os.Getwd()
	require.NoError(t, e3, "error getting the current working directory")

	defer func(dir string) {
		e4 := os.Chdir(dir)
		require.NoError(t, e4, fmt.Sprintf("error changing working directory to %s", dir))
	}(testDir)

	e5 := os.Chdir(tmpDir)
	require.NoError(t, e5, fmt.Sprintf("error changing to the temporary test directory %s", tmpDir))

	f, e6 := os.Create(eventsFileName)
	require.NoError(t, e6, fmt.Sprintf("failed to create file %s due to %s", eventsFileName, e6))

	defer func(f *os.File) {
		e7 := f.Close()
		if e7 != nil {
			require.NoError(t, e7, fmt.Sprintf("failed to close file %s due to %s", eventsFileName, e7))
		}
	}(f)
	today := time.Now()
	s := fmt.Sprintf("{\"%02d-%02d\": {\"year\": 2000, \"remind\": 3, \"type\": \"birthday\", \"event\": \"Someone's birthday\"},\"%02d-%02d\": {\"year\": 2000, \"remind\": 3, \"type\": \"anniversary\", \"event\": \"Someone's aniversary\"}}",
		today.Month(), today.Day(), today.Month(), today.Day()+1)
	_, e8 := f.WriteString(s)
	require.NoError(t, e8, fmt.Sprintf("failed to write file %s due to %s", eventsFileName, e8))

	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		readPath := filepath.Join(testDir, "..", "clingo-conf.toml")
		configB, e9 := ioutil.ReadFile(readPath)
		require.NoError(t, e9, fmt.Sprintf("error reading test config file %s", readPath))

		writePath := filepath.Join(tmpDir, "clingo-conf.toml")
		e10 := ioutil.WriteFile(writePath, configB, 0644)
		require.NoError(t, e10, fmt.Sprintf("error writing test config file %s", writePath))

		defer func(name string) {
			e11 := os.Remove(name)
			require.NoError(t, e11, fmt.Sprintf("error removing test config file %s", name))
		}(filepath.Join(tmpDir, "clingo-conf.toml"))

		// Run ./clingo
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		err := cmd.Execute()
		require.NoError(t, err, fmt.Sprintf("error executing cli command %s", err))

		gotOutput := output.String()
		wantOutput := fmt.Sprintf("Today is %d %s %d: Someone's birthday [%d year(s)]\nIn 1 day(s) will be %d-%02d-%02d: Someone's aniversary [%d year(s)]\n",
			today.Day(), today.Month(), today.Year(), today.Year()-2000, today.Year(), today.Month(), today.Day()+1, today.Year()-2000)
		assert.Equal(t, wantOutput, gotOutput, "expected the 'events' option from the config file and the 'filter' from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		// Run CLINGO=purple ./clingo
		e9 := os.Setenv("CLINGO_FILTER", "anniversary")
		require.NoError(t, e9, "error setting CLINGO_FILTER env variable")

		defer func() {
			e10 := os.Unsetenv("CLINGO_FILTER")
			require.NoError(t, e10, "error unsetting CLINGO_FILTER env variable")
		}()

		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		err := cmd.Execute()
		require.NoError(t, err, "error executing cli command")

		gotOutput := output.String()
		wantOutput := fmt.Sprintf("In 1 day(s) will be %d-%02d-%02d: Someone's aniversary [%d year(s)]\n",
			today.Year(), today.Month(), today.Day()+1, today.Year()-2000)
		assert.Equal(t, wantOutput, gotOutput, "expected the 'filter' option to use the environment variable value and the 'events' option to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./clingo --number 2
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--filter", "birthday"})
		err := cmd.Execute()
		require.NoError(t, err, "error executing cli command")

		gotOutput := output.String()
		wantOutput := fmt.Sprintf("Today is %d %s %d: Someone's birthday [%d year(s)]\n",
			today.Day(), today.Month(), today.Year(), today.Year()-2000)
		assert.Equal(t, wantOutput, gotOutput, "expected the 'filter' option to use the flag value and 'events' option to use the flag default")
	})
}
