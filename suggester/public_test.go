package suggester

import (
	"encoding/json"
	"strings"
	"testing"
)

type rdOut struct {
	Source struct {
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	} `json:"source"`
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

func mustParse(t *testing.T, b []byte) rdOut {
	t.Helper()
	var o rdOut
	if err := json.Unmarshal(b, &o); err != nil {
		t.Fatalf("json: %v", err)
	}
	return o
}

func TestWithBefore_SingleLine(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "main.go",
		"BaseText": "a\nb\nc\n",
		"LMBefore": "b\n",
		"LMAfter":  "B!\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("suggestion not found:\n%s", string(out))
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", rd.Diagnostics[0].Location.Range.Start.Line)
	}
	if rd.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("end line = %d, want 3", rd.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestWithBefore_MultiLine(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "f.txt",
		"BaseText": "x\ny\nz\n",
		"LMBefore": "y\nz\n",
		"LMAfter":  "Y\nZ\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 || rd.Diagnostics[0].Location.Range.End.Line != 4 {
		t.Fatalf("range lines got %d..%d, want 2..4",
			rd.Diagnostics[0].Location.Range.Start.Line,
			rd.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestNoBefore_FullAfterDiff(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "f.txt",
		"BaseText": "a\nb\nc\n",
		"LMAfter":  "a\nB!\nc\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("suggestion missing:\n%s", string(out))
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 || rd.Diagnostics[0].Location.Range.End.Line != 2 {
		t.Fatalf("range lines got %d..%d, want 2..2",
			rd.Diagnostics[0].Location.Range.Start.Line,
			rd.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestCRLF_Normalize(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "f.txt",
		"BaseText": "a\r\nb\r\nc\r\n",
		"LMBefore": "b\r\n",
		"LMAfter":  "B!\r\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 || rd.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("CRLF range wrong")
	}
	if !strings.Contains(string(out), "```suggestion\\nB!\\n```") {
		t.Fatalf("normalized suggestion missing")
	}
}

func TestDuplicate_FirstMatchPreferred(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "f.txt",
		"BaseText": "foo\nbar\nbar\n",
		"LMBefore": "bar\n",
		"LMAfter":  "BAR\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("should select first 'bar' at line 2, got %d", rd.Diagnostics[0].Location.Range.Start.Line)
	}
}

func TestEmptyAfter_Error(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "f.txt",
		"BaseText": "x\n",
		"LMBefore": "x\n",
		"LMAfter":  "",
	}
	inputJSON, _ := json.Marshal(input)
	_, err := Convert(inputJSON, "reviewdog")
	if err == nil {
		t.Fatalf("want ErrEmptyAfter, got %v", err)
	}
}

func TestUTF8_Japanese(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "main.go",
		"BaseText": "package main\n// こんにちは世界\nfunc main() {}\n",
		"LMBefore": "// こんにちは世界\n",
		"LMAfter":  "// Hello, World!\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\n// Hello, World!\\n```") {
		t.Fatalf("suggestion not found:\n%s", string(out))
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", rd.Diagnostics[0].Location.Range.Start.Line)
	}
	if rd.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("end line = %d, want 3", rd.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestUTF8_ChineseEmoji(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "test.txt",
		"BaseText": "第一行\n需要修改的行 🚀\n第三行\n",
		"LMBefore": "需要修改的行 🚀\n",
		"LMAfter":  "已修改的行 ✅\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\n已修改的行 ✅\\n```") {
		t.Fatalf("suggestion not found:\n%s", string(out))
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", rd.Diagnostics[0].Location.Range.Start.Line)
	}
	if rd.Diagnostics[0].Location.Range.End.Line != 3 {
		t.Fatalf("end line = %d, want 3", rd.Diagnostics[0].Location.Range.End.Line)
	}
}

func TestUTF8_MixedContent(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "mixed.go",
		"BaseText": "// English comment\n// 日本語コメント\n// 中文注释\n// Emoji 🎉\n",
		"LMBefore": "// 日本語コメント\n",
		"LMAfter":  "// Japanese comment\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	rd := mustParse(t, out)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 {
		t.Fatalf("start line = %d, want 2", rd.Diagnostics[0].Location.Range.Start.Line)
	}
}

func TestUTF8_NoBefore_FullAfter(t *testing.T) {
	input := map[string]interface{}{
		"FilePath": "unicode.txt",
		"BaseText": "こんにちは\n世界\nWorld\n",
		"LMAfter":  "こんにちは\nせかい\nWorld\n",
	}
	inputJSON, _ := json.Marshal(input)
	out, err := Convert(inputJSON, "reviewdog")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(string(out), "```suggestion\\nせかい\\n```") {
		t.Fatalf("suggestion not found:\n%s", string(out))
	}
	rd := mustParse(t, out)
	// Should replace "世界\n" with "せかい\n" (line 2)
	if rd.Diagnostics[0].Location.Range.Start.Line != 2 || rd.Diagnostics[0].Location.Range.End.Line != 2 {
		t.Fatalf("range lines got %d..%d, want 2..2",
			rd.Diagnostics[0].Location.Range.Start.Line,
			rd.Diagnostics[0].Location.Range.End.Line)
	}
}