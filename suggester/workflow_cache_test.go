package suggester

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCIWorkflowsDisableSetupGoCache guards against the self-hosted runner
// cache-service warning "Failed to restore: Cache service responded with 400".
// actions/setup-go must not touch the GitHub-hosted cache, and no workflow
// may use actions/cache directly.
func TestCIWorkflowsDisableSetupGoCache(t *testing.T) {
	dir := filepath.Join("..", ".github")
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		t.Skip("CI workflow directory not present")
	}

	pattern := filepath.Join(dir, "workflows", "*.yml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob %s: %v", pattern, err)
	}
	if len(files) == 0 {
		t.Fatalf("no workflow files matched %s", pattern)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			b, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}
			src := string(b)

			if strings.Contains(src, "actions/cache@") {
				t.Errorf("%s uses actions/cache, which triggers the self-hosted cache-service warning", file)
			}

			if !strings.Contains(src, "actions/setup-go@") {
				return
			}
			if strings.Contains(src, "cache: true") {
				t.Errorf("%s enables the actions/setup-go cache", file)
			}
			if !strings.Contains(src, "cache: false") {
				t.Errorf("%s must set \"cache: false\" on actions/setup-go", file)
			}
		})
	}
}
