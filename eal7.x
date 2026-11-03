#!/data/data/com.termux/files/usr/bin/bash
# eal7.x - FLOP QUADRUPLE Port 8038 - oenLytics.x IANA 1460172
# Port 38 or 80 as flop quadruple = 8038 = 80*100 + 38
echo "eal7.x FLOP QUADRUPLE - Port 8038"
echo "<script src=\"http://localhost:8038/oenlytics.js\"></script>"
echo "<script>oenLytics.track('page_view')</script>"
# Run Go with QUAD precision FLOP counter
PORT=8038 go run cmd/oenlytics/main.go
