package cmds

import (
	"bufio"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Transpile igo files to go files",
	Run:   BuildCmd,
}

func TestIgoBuild(t *testing.T) {

	// Generating go code from igo sample file
	BuildCmd(buildCmd, []string{"./test_sample.igo"})

	// Executing Generated Command
	cmd := exec.Command("go","run", "gen_test_sample.go")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("test failed with error: %s\n", err)
	}
	t.Logf("Output from command:\n%s\n", string(out))

	// Check to make sure expected output are generated when executed
	expectedOutput := `
		civic is a valid palindrome.
		Palindrome function called!
		done processing word 'civic'
		relation is not a valid palindrome.
		done processing word 'relation'
	`

	actualOutput := string(out)
	scanner := bufio.NewScanner(strings.NewReader(actualOutput))
	for scanner.Scan() {
		if !strings.Contains(expectedOutput, scanner.Text()) {
			t.Errorf("Test failed with error: An unepected string '%s' in output", scanner.Text())
		}
	}

	// delete generated file
	t.Log("deleting Generated Test_build File...\n")
	if err := os.Remove("gen_test_sample.go"); err != nil {
		t.Errorf("Test failed with error:%s\n", err)
	}
	t.Log("deleted.\n")

}
