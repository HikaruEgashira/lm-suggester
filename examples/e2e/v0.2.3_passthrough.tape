# VHS tape demonstrating v0.2.3 passthrough transformation
Output v0.2.3_passthrough.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

Type "# v0.2.3 Feature: Passthrough JSON Transformation"
Enter
Sleep 500ms
Type "# Input JSON with custom fields (customField, metadata)"
Enter
Sleep 500ms
Type "cat ../testdata/passthrough_example.json"
Enter
Sleep 2s

Enter
Type "# Transform to reviewdog format (custom fields preserved)"
Enter
Sleep 500ms
Type "cat ../testdata/passthrough_example.json | go run cmd/lm-suggester/main.go"
Enter
Sleep 2s
