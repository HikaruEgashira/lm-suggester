# NOTE: Must be run from project root directory
Output v040_jsonl_support.gif

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
