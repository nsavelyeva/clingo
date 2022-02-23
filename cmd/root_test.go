package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory
	tmpDir, e1 := ioutil.TempDir("..", "tmp-test-clingo")
	require.NoError(t, e1, fmt.Sprintf("error creating a temporary test folder %s", tmpDir))

	defer func(name string) {
		e2 := os.Remove(name)
		require.NoError(t, e2, fmt.Sprintf("error removing temporary test folder %s", name))
	}(tmpDir)

	testDir, e3 := os.Getwd()
	require.NoError(t, e3, "error getting the current working directory")

	defer func(dir string) {
		e4 := os.Chdir(dir)
		require.NoError(t, e4, fmt.Sprintf("error changing working directory to %s", dir))
	}(testDir)

	e5 := os.Chdir(tmpDir)
	require.NoError(t, e5, fmt.Sprintf("error changing to the temporary test directory %s", tmpDir))

	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		readPath := filepath.Join(testDir, "..", "clingo-conf.toml")
		configB, e6 := ioutil.ReadFile(readPath)
		require.NoError(t, e6, fmt.Sprintf("error reading test config file %s", readPath))

		writePath := filepath.Join(tmpDir, "clingo-conf.toml")
		e7 := ioutil.WriteFile(writePath, configB, 0644)
		require.NoError(t, e7, fmt.Sprintf("error writing test config file %s", writePath))

		defer func(name string) {
			e8 := os.Remove(name)
			require.NoError(t, e8, fmt.Sprintf("error removing test config file %s", name))
		}(filepath.Join(tmpDir, "clingo-conf.toml"))

		// Run ./clingo
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		err := cmd.Execute()
		require.NoError(t, err, "error executing cli command")

		gotOutput := output.String()
		wantOutput := `Your favorite color is: blue
The magic number is: 7
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color from the config file and the number from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		// Run CLINGO=purple ./clingo
		e6 := os.Setenv("CLINGO_FAVORITE_COLOR", "purple")
		require.NoError(t, e6, "error setting CLINGO_FAVORITE_COLOR env variable")

		defer func() {
			e7 := os.Unsetenv("CLINGO_FAVORITE_COLOR")
			require.NoError(t, e7, "error unsetting CLINGO_FAVORITE_COLOR env variable")
		}()

		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		err := cmd.Execute()
		require.NoError(t, err, "error executing cli command")

		gotOutput := output.String()
		wantOutput := `Your favorite color is: purple
The magic number is: 7
`
		assert.Equal(t, wantOutput, gotOutput, "expected the color to use the environment variable value and the number to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./clingo --number 2
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"--number", "2"})
		err := cmd.Execute()
		require.NoError(t, err, "error executing cli command")

		gotOutput := output.String()
		wantOutput := `Your favorite color is: red
The magic number is: 2
`
		assert.Equal(t, wantOutput, gotOutput, "expected the number to use the flag value and the color to use the flag default")
	})
}
