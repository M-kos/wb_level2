#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

test_count=0
passed_count=0
failed_count=0

run_test() {
    local test_name="$1"
    local command="$2"
    local expected_exit="$3"

    ((test_count++))
    echo -n "Test $test_count: $test_name ... "

    eval "$command" > /dev/null 2>&1
    local exit_code=$?

    if [ "$exit_code" -eq "$expected_exit" ]; then
        echo -e "${GREEN}PASSED${NC}"
        ((passed_count++))
    else
        echo -e "${RED}FAILED${NC} (expected exit code $expected_exit, got $exit_code)"
        ((failed_count++))
    fi
}

# Test built-in commands
echo "=== Built-in Commands Tests ==="
run_test "echo command" "echo 'hello world'" 0
run_test "pwd command" "pwd" 0
run_test "cd to /tmp" "cd /tmp && pwd | grep -q tmp" 0

# Test logic operators
echo
echo "=== Logic Operators Tests ==="
run_test "AND operator - success" "true && echo 'success'" 0
run_test "OR operator - fallback" "false || echo 'success'" 0

# Test pipelines
echo
echo "=== Pipeline Tests ==="
run_test "Simple pipeline" "echo 'hello' | grep 'hello'" 0
run_test "Multiple pipes" "echo 'hello world' | grep 'hello' | grep 'world'" 0

echo
echo "=== Test Summary ==="
echo -e "Total tests: $test_count"
echo -e "Passed: ${GREEN}$passed_count${NC}"
echo -e "Failed: ${RED}$failed_count${NC}"

if [ "$failed_count" -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
