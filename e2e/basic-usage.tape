Output e2e/basic-usage.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

# Show the input JSON file
Type "cat _examples/testdata/simple_replacement.json"
Sleep 500ms
Enter
Sleep 2s
Enter
Sleep 500ms

# Convert to reviewdog JSON format
Type "lm-suggester -i _examples/testdata/simple_replacement.json"
Sleep 500ms
Enter
Sleep 2s

# Convert with pretty print
Type "lm-suggester -i _examples/testdata/simple_replacement.json --pretty"
Sleep 500ms
Enter
Sleep 2.5s
