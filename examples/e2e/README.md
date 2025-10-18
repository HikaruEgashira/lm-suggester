# E2E Tests with VHS

This directory contains VHS (Video Home System) tape files for visual regression testing and documentation of the lm-suggester CLI tool.

## Requirements

- [VHS](https://github.com/charmbracelet/vhs) - A tool for generating terminal GIFs

```bash
# Install VHS
brew install vhs
```

## Tape Files

### help-command.tape

Demonstrates the CLI help commands:
- `lm-suggester --help` - Main help text
- `lm-suggester version --help` - Version command help

Output: `_examples/e2e/help-command.gif`

### basic-usage.tape

Shows basic usage of the CLI tool:
- Displaying a sample input JSON file
- Converting to reviewdog JSON format
- Pretty-printing the output

Output: `_examples/e2e/basic-usage.gif`

## Generating GIFs

To regenerate the GIFs:

```bash
# From the project root directory
PATH="$PATH:$(pwd)" vhs _examples/e2e/help-command.tape
PATH="$PATH:$(pwd)" vhs _examples/e2e/basic-usage.tape
```

Or regenerate all tapes:

```bash
for tape in _examples/e2e/*.tape; do
  PATH="$PATH:$(pwd)" vhs "$tape"
done
```

## Adding New Tests

1. Create a new `.tape` file in this directory
2. Follow the VHS tape format:
   - Set `Output` to `_examples/e2e/your-test-name.gif`
   - Configure terminal settings (FontSize, Width, Height, Theme)
   - Add commands using `Type` and `Enter`
   - Use `Sleep` for appropriate timing
3. Generate the GIF: `PATH="$PATH:$(pwd)" vhs _examples/e2e/your-test-name.tape`

## VHS Tape Format

```tape
Output _examples/e2e/output.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Padding 20

Type "your-command"
Sleep 500ms
Enter
Sleep 2s
```

For more information, see the [VHS documentation](https://github.com/charmbracelet/vhs).
