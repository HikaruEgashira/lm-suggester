# VHS tape demonstrating v0.5.1 unified Convert API with auto-detection
Output e2e/v0.5.1_unified_api.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

# Show the feature being demonstrated
Type "# v0.5.1: Unified Convert API with Auto-Detection"
Enter
Sleep 500ms

# Show the input JSON
Type "cat _examples/testdata/simple_replacement.json"
Enter
Sleep 1.5s

Enter

# Convert to reviewdog format using unified API
Type "# Convert to reviewdog format using unified API"
Enter
Sleep 500ms

Type "cat _examples/testdata/simple_replacement.json | go run cmd/lm-suggester/main.go"
Enter
Sleep 3s
