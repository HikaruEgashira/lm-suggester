# lm-suggester Examples

This directory contains examples demonstrating the usage of the lm-suggester library.

## Usage

```bash
cat testdata/simple_replacement.json | go run ./simple
```

## Example Scenarios

### 1. Simple Code Replacement (`simple_replacement.json`)
Basic single-line code replacement with suggestion.

### 2. Japanese Comment Translation (`japanese_comment.json`)
Handles UTF-8 multibyte characters correctly when translating Japanese comments.

### 3. Multiline Function Refactoring (`multiline_function.json`)
Demonstrates replacing multiple lines with a simplified version.

### 4. Full File Update (`full_file_update.json`)
Shows diff calculation when LMBefore is not provided (full file replacement).

### 5. Emoji in Code (`emoji_in_code.json`)
Tests proper handling of emoji characters in code.
