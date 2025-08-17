#!/bin/bash

# Simple Green Functions Test Runner
# Usage: ./run_green_tests.sh [option]
# Options:
#   unit     - Run only unit tests
#   bench    - Run only benchmarks
#   perf     - Run performance analysis
#   all      - Run all tests (default)
#   report   - Generate detailed report

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}    GREEN FUNCTIONS TEST RUNNER${NC}"
echo -e "${BLUE}==========================================${NC}"
echo ""

# Function to run unit tests
run_unit_tests() {
    echo -e "${CYAN}Running Unit Tests...${NC}"
    go test . -v
    echo ""
}

# Function to run benchmarks
run_benchmarks() {
    echo -e "${CYAN}Running Benchmarks...${NC}"
    go test . -bench=. -benchmem
    echo ""
}

# Function to run performance analysis
run_performance() {
    echo -e "${CYAN}Running Performance Analysis...${NC}"
    go test . -bench=. -benchmem > /tmp/green_benchmark.log 2>&1
    
    echo "=== TOP PERFORMANCE METRICS ==="
    echo ""
    
    echo "🏃 FASTEST OPERATIONS:"
    grep "Benchmark" /tmp/green_benchmark.log | sort -k3 -n | head -5
    
    echo ""
    echo "💾 MEMORY EFFICIENT:"
    grep "Benchmark" /tmp/green_benchmark.log | grep "0 B/op" | head -5
    
    echo ""
    echo "🔢 HIGH THROUGHPUT:"
    grep "Benchmark" /tmp/green_benchmark.log | sort -k2 -nr | head -5
    echo ""
}

# Function to generate report
generate_report() {
    echo -e "${CYAN}Generating Detailed Report...${NC}"
    ./test_green_functions.sh
    echo ""
    echo -e "${GREEN}Report generated! Check GREEN_FUNCTIONS_REPORT.md${NC}"
}

# Main execution
case "${1:-all}" in
    "unit")
        run_unit_tests
        ;;
    "bench")
        run_benchmarks
        ;;
    "perf")
        run_performance
        ;;
    "report")
        generate_report
        ;;
    "all")
        run_unit_tests
        run_benchmarks
        run_performance
        ;;
    *)
        echo "Usage: $0 [unit|bench|perf|all|report]"
        echo "  unit   - Run only unit tests"
        echo "  bench  - Run only benchmarks"
        echo "  perf   - Run performance analysis"
        echo "  all    - Run all tests (default)"
        echo "  report - Generate detailed report"
        exit 1
        ;;
esac

echo -e "${GREEN}✅ Green Functions tests completed!${NC}" 