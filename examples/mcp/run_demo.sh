#!/bin/bash
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}'
  sleep 0.2
  echo '{"jsonrpc":"2.0","method":"notifications/initialized"}'
  sleep 0.2
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"suggest","arguments":{"file_path":"/tmp/test.go","base_text":"package main\n\nfunc main() {\n\tprint(\"Hello\")\n}","lm_before":"print(\"Hello\")","lm_after":"fmt.Println(\"Hello\")","message":"Use fmt.Println","reporter":"local"}}}'
  sleep 1
) | lm-suggester mcp 2>&1 &
PID=$!
sleep 2
kill $PID 2>/dev/null || true
wait $PID 2>/dev/null || true
