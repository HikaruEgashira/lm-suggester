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
	title := strings.TrimSpace(head)
	if title == "" {
		title = "Replace code with suggestion"
	}
	body := after
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return title + "\n```suggestion\n" + body + "```"
}

func marshalRDJSON(src, path, msg string, startLine, startCol, endLine, endCol int, severity string) ([]byte, error) {
	if severity == "" {
		severity = "WARNING"
	}
	if src == "" {
		src = "llm-suggester"
	}
	payload := rdjson{
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
				Severity: severity,
			},
		},
	}
	return json.MarshalIndent(payload, "", "  ")
}
