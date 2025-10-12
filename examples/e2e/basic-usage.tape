Output basic-usage.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

# Show the input JSON file
Type "cat ../testdata/simple_replacement.json"
Enter
Sleep 1s
Enter

# Convert to reviewdog JSON format
Type "lm-suggester -i ../testdata/simple_replacement.json"
Enter
Sleep 1s

# Convert with pretty print
Type "lm-suggester -i ../testdata/simple_replacement.json --pretty"
Enter
Sleep 1.5s
