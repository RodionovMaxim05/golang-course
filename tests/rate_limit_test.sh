#!/bin/bash

# Target Settings
TARGET_URL="http://localhost:28080/api/ping"
RATE=15
DURATION="2s"

echo "Running the Rate Limiter test via Vegeta..."
echo "Target: $TARGET_URL, Rate: $RATE requests/sec, Duration: $DURATION"
echo "------------------------------------------------------------"

# Launch an attack through Vegeta and save a binary report
echo "GET $TARGET_URL" | vegeta attack -rate=$RATE -duration=$DURATION > result.bin

# Outputting a text report to the console
echo "Results:"
vegeta report result.bin

# Perform a deep analysis of status codes for automated verification (CI/CD)
echo "------------------------------------------------------------"
echo "Checking status code distribution..."

# Count how many times the code 429 has been encountered
STATUS_429=$(vegeta report -type=json result.bin | jq '."status_codes"."429" // 0')
# Count how many times the code 200 has been encountered
STATUS_200=$(vegeta report -type=json result.bin | jq '."status_codes"."200" // 0')

echo "Successful requests (200 OK): ${STATUS_200:-0}"
echo "Blocked requests (429 Too Many Requests): ${STATUS_429:-0}"

rm result.bin

# Logic for test completion: if no blocked requests are found, the limiter did not work correctly
if [ "$STATUS_429" -eq 0 ]; then
    echo "❌ TEST FAILED: The limiter allowed all requests. Status code 429 not received."
    exit 1
else
    echo "✅ TEST PASSED: The limiter correctly filtered out excessive load."
    exit 0
fi
