package suggester

import (
	"encoding/json"
	"strings"
	"testing"
)

type rdOutput struct {
	Source      string `json:"source"`
	Diagnostics []struct {
		Message  string `json:"message"`
		Location struct {
			Path  string `json:"path"`
			Range struct {
				Start struct {
					Line   int `json:"line"`
					Column int `json:"column"`
				} `json:"start"`
				End struct {
					Line   int `json:"line"`
					Column int `json:"column"`
				} `json:"end"`
			} `json:"range"`
		} `json:"location"`
		Severity string `json:"severity"`
	} `json:"diagnostics"`
}

func parseRDJSON(t *testing.T, b []byte) rdOutput {
	t.Helper()
	var out rdOutput
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal rdjson: %v", err)
	}
	return out
}

func TestBuildRDJSON_WithBefore_SingleLine(t *testing.T) {
	in := Input{
		FilePath:  "main.go",
		BaseText:  "a\nb\nc\n",
		LLMBefore: "b\n",
		LLMAfter:  "B!\n",
	}
	out, err := BuildRDJSON(in)
	if err != nil {
		t.Fatalf("BuildRDJSON error: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("suggestion block not found: %s", string(out))
	}
	parsed := parseRDJSON(t, out)
	if parsed.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", parsed.Diagnostics[0].Location.Range.Start.Line)
	}
	if parsed.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("end line = %d, want 3", parsed.Diagnostics[0].Location.Range.End.Line)
	}
	if parsed.Source != "llm-suggester" {
		t.Fatalf("source = %s, want llm-suggester", parsed.Source)
	}
}

func TestBuildRDJSON_WithBefore_MultiLine(t *testing.T) {
	in := Input{
		FilePath:  "f.txt",
		BaseText:  "x\ny\nz\n",
		LLMBefore: "y\nz\n",
		LLMAfter:  "Y\nZ\n",
	}
	out, err := BuildRDJSON(in)
	if err != nil {
		t.Fatalf("BuildRDJSON error: %v", err)
	}
	parsed := parseRDJSON(t, out)
	if parsed.Diagnostics[0].Location.Range.Start.Line != 2 || parsed.Diagnostics[0].Location.Range.End.Line != 4 {
		t.Fatalf("unexpected range lines: %d..%d", parsed.Diagnostics[0].Location.Range.Start.Line, parsed.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestBuildRDJSON_NoBefore_Diff(t *testing.T) {
	in := Input{
		FilePath: "f.txt",
		BaseText: "a\nb\nc\n",
		LLMAfter: "a\nB!\nc\n",
	}
	out, err := BuildRDJSON(in)
	if err != nil {
		t.Fatalf("BuildRDJSON error: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("suggestion block missing: %s", string(out))
	}
	parsed := parseRDJSON(t, out)
	if parsed.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", parsed.Diagnostics[0].Location.Range.Start.Line)
	}
	if parsed.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("end line = %d, want 3", parsed.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestBuildRDJSON_CRLF(t *testing.T) {
	in := Input{
		FilePath:  "f.txt",
		BaseText:  "a\r\nb\r\nc\r\n",
		LLMBefore: "b\r\n",
		LLMAfter:  "B!\r\n",
	}
	out, err := BuildRDJSON(in)
	if err != nil {
		t.Fatalf("BuildRDJSON error: %v", err)
	}
	parsed := parseRDJSON(t, out)
	if parsed.Diagnostics[0].Location.Range.Start.Line != 2 || parsed.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("unexpected range for CRLF input")
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("normalized suggestion missing: %s", string(out))
	}
}

func TestBuildRDJSON_Duplicate(t *testing.T) {
	in := Input{
		FilePath:  "f.txt",
		BaseText:  "foo\nbar\nbar\n",
		LLMBefore: "bar\n",
		LLMAfter:  "BAR\n",
	}
	out, err := BuildRDJSON(in)
	if err != nil {
		t.Fatalf("BuildRDJSON error: %v", err)
	}
	parsed := parseRDJSON(t, out)
	if parsed.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("should choose first match, got line %d", parsed.Diagnostics[0].Location.Range.Start.Line)
	}
}

func TestBuildRDJSON_EmptyAfter(t *testing.T) {
	in := Input{
		FilePath:  "f.txt",
		BaseText:  "x\n",
		LLMBefore: "x\n",
	}
	if _, err := BuildRDJSON(in); err == nil || err != ErrEmptyAfter {
		t.Fatalf("want ErrEmptyAfter, got %v", err)
	}
}

func TestMinimalRange_NoChange(t *testing.T) {
	if _, _, _, err := minimalRangeFromFullAfter("a\n", "a\n"); err != ErrNoChange {
		t.Fatalf("want ErrNoChange, got %v", err)
	}
}
