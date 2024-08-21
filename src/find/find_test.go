package find

import (
	"bytes"
	"context"
	"testing"
)

func TestRun(t *testing.T) {
	for name, tc := range map[string]struct {
		Dir            string
		Pattern        string
		Separator      string
		ExpectedStdout string
		ExpectedStderr string
		ShouldErr      bool
	}{
		"match pattern": {
			Dir:            "testdata",
			Pattern:        `alice`,
			ExpectedStdout: "testdata/alice\n",
		},

		"match pattern with ./ dir prefix": {
			Dir:            "./testdata",
			Pattern:        `alice`,
			ExpectedStdout: "./testdata/alice\n",
		},

		"null separator": {
			Dir:            "testdata",
			Pattern:        `alice`,
			Separator:      "\x00",
			ExpectedStdout: "testdata/alice\x00",
		},

		"No pattern walks all files": {
			Dir:     "testdata",
			Pattern: "",
			ExpectedStdout: `testdata
testdata/alice
testdata/bob
testdata/carol
`,
		},

		"missing file returns error": {
			Dir:            "unknown",
			Pattern:        "",
			ExpectedStderr: "\"unknown\" does not exist\n",
		},

		"invalid regex": {
			Pattern:   `(?!alice)`,
			ShouldErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}

			f := Find{
				Stdout:    &stdoutBuf,
				Stderr:    &stderrBuf,
				Dir:       tc.Dir,
				Pattern:   tc.Pattern,
				Separator: tc.Separator,
			}

			err := f.Run(context.Background())
			if (err != nil) != tc.ShouldErr {
				t.Fatal(err)
			}

			if got, want := stdoutBuf.String(), tc.ExpectedStdout; got != want {
				t.Errorf("stdout: got %s, want %s", got, want)
			}

			if got, want := stderrBuf.String(), tc.ExpectedStderr; got != want {
				t.Errorf("stderr: got %s, want %s", got, want)
			}
		})
	}

}
