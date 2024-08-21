package cat

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	for name, tc := range map[string]struct {
		Paths          []string
		Stdin          string
		ExpectedStdout string
		ExpectedStderr string
	}{
		"stdin only": {
			Stdin:          "a b c",
			ExpectedStdout: "a b c",
		},

		"file only": {
			Paths:          []string{"./testdata/foo.txt"},
			ExpectedStdout: "foo\n",
		},

		"file and stdin": {
			Stdin:          "bar",
			Paths:          []string{"./testdata/foo.txt", "-"},
			ExpectedStdout: "foo\nbar",
		},

		"missing file": {
			Paths:          []string{"missing.txt"},
			ExpectedStderr: "open missing.txt: no such file or directory\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			stdinBuf := bytes.NewBufferString(tc.Stdin)
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			cat := Cat{
				Stdin:  stdinBuf,
				Stdout: &stdoutBuf,
				Stderr: &stderrBuf,
			}

			cat.Run(tc.Paths)

			if got, want := stdoutBuf.String(), tc.ExpectedStdout; got != want {
				t.Errorf("got %s, want %s", got, want)
			}

			if got, want := stderrBuf.String(), tc.ExpectedStderr; got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	}

}
