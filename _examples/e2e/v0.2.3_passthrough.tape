# NOTE: Must be run from project root directory
Output v023_passthrough.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

Hide
Type "cd /Users/hikae/ghq/github.com/HikaruEgashira/lm-suggester"
Enter
Show

Type "# v0.2.3 Feature: Passthrough JSON Transformation"
Enter
Sleep 500ms
Type "# Input JSON with custom fields (customField, metadata)"
Enter
Sleep 500ms
Type "cat _examples/testdata/passthrough_example.json"
Enter
Sleep 2s

Enter
Type "# Transform to reviewdog format (custom fields preserved)"
Enter
Sleep 500ms
Type "cat _examples/testdata/passthrough_example.json | go run cmd/lm-suggester/main.go"
Enter
Sleep 2s
