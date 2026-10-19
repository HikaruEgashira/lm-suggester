package suggester

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// node24ActionMajors is the first major version of each action that ships a
// Node.js 24 runtime. Older majors trigger the GitHub Actions warning
// "Node.js 20 is deprecated".
var node24ActionMajors = map[string]int{
	"actions/checkout":             5,
	"actions/setup-go":             6,
	"actions/upload-artifact":      6,
	"goreleaser/goreleaser-action": 7,
}

var usesLinePattern = regexp.MustCompile(`(?m)^\s*uses:\s*([^\s#]+)(?:\s*#\s*(v[0-9][^\s]*))?`)

// TestCIWorkflowsUseNode24Actions guards against the GitHub Actions warning
// "Node.js 20 is deprecated". Workflow actions must either use a major version
// with a Node.js 24 runtime or pin a commit with a version comment.
func TestCIWorkflowsUseNode24Actions(t *testing.T) {
	pattern := filepath.Join("..", ".github", "workflows", "*.yml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob %s: %v", pattern, err)
	}
	if len(files) == 0 {
		if _, err := os.Stat(filepath.Join("..", ".github")); os.IsNotExist(err) {
			t.Skip("CI workflow directory not present")
		}
		t.Fatalf("no workflow files matched %s", pattern)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			b, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}

			for _, match := range usesLinePattern.FindAllStringSubmatch(string(b), -1) {
				ref := match[1]
				action, version, ok := strings.Cut(ref, "@")
				if !ok {
					continue
				}
				minMajor, checked := node24ActionMajors[action]
				if !checked {
					continue
				}

				major := actionMajor(version, match[2])
				if major == 0 {
					t.Errorf("%s: %s has no version comment; use v%d or newer", file, ref, minMajor)
					continue
				}
				if major < minMajor {
					t.Errorf("%s: %s targets Node.js 20; use v%d or newer", file, ref, minMajor)
				}
			}
		})
	}
}

func actionMajor(version, comment string) int {
	candidate := version
	if !strings.HasPrefix(candidate, "v") {
		candidate = comment
	}
	candidate = strings.TrimPrefix(candidate, "v")
	if end := strings.IndexFunc(candidate, func(r rune) bool {
		return r < '0' || r > '9'
	}); end >= 0 {
		candidate = candidate[:end]
	}
	major, _ := strconv.Atoi(candidate)
	return major
}
