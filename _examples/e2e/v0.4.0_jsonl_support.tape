# NOTE: VHS 0.10.0 does not support paths starting with underscore
# Generated file should be manually moved: examples/e2e/* -> _examples/e2e/*
Output examples/e2e/v0.4.0_jsonl_support.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

# Show the JSONL input file with multiple suggestions
Type "cat _examples/testdata/example.jsonl"
Enter
Sleep 2s
Enter

# Process JSONL format via stdin (pipe)
Type "cat _examples/testdata/example.jsonl | lm-suggester"
Enter
Sleep 2s
Enter

# Process JSONL with pretty print for better readability
Type "cat _examples/testdata/example.jsonl | lm-suggester --pretty"
Enter
Sleep 2.5s
