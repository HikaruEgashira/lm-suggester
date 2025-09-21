package suggester

import (
	"encoding/json"
	"strings"
)

type rdjson struct {
	Source      string       `json:"source"`
	Diagnostics []diagnostic `json:"diagnostics"`
}

type diagnostic struct {
	Message  string   `json:"message"`
	Location location `json:"location"`
	Severity string   `json:"severity"`
}

type location struct {
	Path  string    `json:"path"`
	Range rangeSpan `json:"range"`
}

type rangeSpan struct {
	Start position `json:"start"`
	End   position `json:"end"`
}

type position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func makeMessage(head, after string) string {
	title := head
	if title == "" {
		title = "Replace code with suggestion"
	}
	if !strings.HasSuffix(after, "\n") {
		after += "\n"
	}
	return title + "\n```suggestion\n" + after + "```"
}

func marshalRDJSON(src, path, msg string, startLine, startCol, endLine, endCol int, sev string) ([]byte, error) {
	if sev == "" {
		sev = "WARNING"
	}
	if src == "" {
		src = "reviewdog-converter"
	}
	out := rdjson{
		Source: src,
		Diagnostics: []diagnostic{
			{
				Message: msg,
				Location: location{
					Path: path,
					Range: rangeSpan{
						Start: position{Line: startLine, Column: startCol},
						End:   position{Line: endLine, Column: endCol},
					},
				},
				Severity: sev,
			},
		},
	}
	return json.MarshalIndent(out, "", "  ")
}