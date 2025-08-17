#!/bin/bash

# Green Functions Comprehensive Test Suite
# This script tests all green functions components including performance

# Don't exit on error, let the script continue

echo "=========================================="
echo "    GREEN FUNCTIONS COMPREHENSIVE TEST"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "FAIL")
            echo -e "${RED}❌ $message${NC}"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ️  $message${NC}"
            ;;
        "WARN")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
        "HEADER")
            echo -e "${PURPLE}📋 $message${NC}"
            ;;
        "PERF")
            echo -e "${CYAN}⚡ $message${NC}"
            ;;
    esac
}

# Function to run test and capture output
run_test() {
    local test_name=$1
    local command=$2
    local description=$3
    
    echo -e "${BLUE}Running: $test_name${NC}"
    echo "Description: $description"
    echo "Command: $command"
    echo ""
    
    if eval "$command" > /tmp/green_test_output.log 2>&1; then
        print_status "PASS" "$test_name completed successfully"
        echo ""
        return 0
    else
        print_status "FAIL" "$test_name failed"
        echo "Error output:"
        cat /tmp/green_test_output.log
        echo ""
        return 1
    fi
}

# Function to run benchmark and capture results
run_benchmark() {
    local benchmark_name=$1
    local command=$2
    
    echo -e "${CYAN}Benchmarking: $benchmark_name${NC}"
    echo "Command: $command"
    echo ""
    
    if eval "$command" > /tmp/green_benchmark_output.log 2>&1; then
        print_status "PERF" "$benchmark_name completed"
        # Extract performance data
        grep "Benchmark" /tmp/green_benchmark_output.log | head -5
        echo ""
        return 0
    else
        print_status "FAIL" "$benchmark_name failed"
        cat /tmp/green_benchmark_output.log
        echo ""
        return 1
    fi
}

# Initialize counters
total_tests=0
passed_tests=0
failed_tests=0
total_benchmarks=0
passed_benchmarks=0
failed_benchmarks=0

# Start testing
print_status "HEADER" "Phase 1: Unit Tests"
echo "============================"

run_test "Base Green Function Tests" "go test ./green_functions -run TestBaseGreenFunction -v" "Test base green function functionality" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Delhommeau Tests" "go test ./green_functions -run TestDelhommeau -v" "Test Delhommeau green function implementation" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "FinGreen3D Tests" "go test ./green_functions -run TestFinGreen3D -v" "Test FinGreen3D implementation" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "LiangWuNoblesse Tests" "go test ./green_functions -run TestLiangWuNoblesse -v" "Test LiangWuNoblesse implementation" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "HAMS Tests" "go test ./green_functions -run TestHAMS -v" "Test HAMS implementation" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Integration Tests" "go test ./green_functions -run TestIntegration -v" "Test integration between different methods" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Utility Function Tests" "go test ./green_functions -run TestCompute -v" "Test utility functions" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Tabulation Cache Tests" "go test ./green_functions -run TestTabulation -v" "Test tabulation cache functionality" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

print_status "HEADER" "Phase 2: Performance Benchmarks"
echo "================================="

run_benchmark "Base Green Function Benchmarks" "go test ./green_functions -bench=BenchmarkBaseGreenFunction -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "Delhommeau Benchmarks" "go test ./green_functions -bench=BenchmarkDelhommeau -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "FinGreen3D Benchmarks" "go test ./green_functions -bench=BenchmarkFinGreen3D -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "LiangWuNoblesse Benchmarks" "go test ./green_functions -bench=BenchmarkLiangWuNoblesse -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "HAMS Benchmarks" "go test ./green_functions -bench=BenchmarkHAMS -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "Integration Benchmarks" "go test ./green_functions -bench=BenchmarkIntegration -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "Utility Function Benchmarks" "go test ./green_functions -bench=BenchmarkCompute -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

run_benchmark "Tabulation Cache Benchmarks" "go test ./green_functions -bench=BenchmarkTabulation -benchmem" && ((passed_benchmarks++)) || ((failed_benchmarks++))
((total_benchmarks++))

print_status "HEADER" "Phase 3: Comprehensive Performance Analysis"
echo "============================================="

echo "Running full benchmark suite..."
go test ./green_functions -bench=. -benchmem > /tmp/full_benchmark.log 2>&1

print_status "PERF" "Extracting performance metrics..."

# Extract and display key performance metrics
echo ""
echo "=== TOP PERFORMANCE METRICS ==="
echo ""

# Fastest operations
echo "🏃 FASTEST OPERATIONS:"
grep "Benchmark" /tmp/full_benchmark.log | sort -k3 -n | head -5

echo ""
echo "🐌 SLOWEST OPERATIONS:"
grep "Benchmark" /tmp/full_benchmark.log | sort -k3 -nr | head -5

echo ""
echo "💾 MEMORY EFFICIENT:"
grep "Benchmark" /tmp/full_benchmark.log | grep "0 B/op" | head -5

echo ""
echo "🔢 HIGH THROUGHPUT:"
grep "Benchmark" /tmp/full_benchmark.log | sort -k2 -nr | head -5

print_status "HEADER" "Phase 4: Method Comparison"
echo "============================="

echo "Comparing different Green Function methods..."
echo ""

# Compare Delhommeau vs LiangWuNoblesse vs HAMS
echo "📊 METHOD COMPARISON (Small Problems):"
grep "BenchmarkIntegration_MethodComparison" /tmp/full_benchmark.log | head -3

echo ""
echo "📊 METHOD COMPARISON (Medium Problems):"
grep "BenchmarkHAMS_vs_LiangWuNoblesse" /tmp/full_benchmark.log | head -2

print_status "HEADER" "Phase 5: Code Quality Checks"
echo "==============================="

run_test "Go Vet Analysis" "go vet ./green_functions" "Static analysis of green functions code" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Code Formatting" "test -z \"\$(go fmt ./green_functions)\"" "Check code formatting" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

run_test "Build Verification" "go build ./green_functions" "Verify green functions package builds" && ((passed_tests++)) || ((failed_tests++))
((total_tests++))

print_status "HEADER" "Phase 6: Detailed Performance Report"
echo "====================================="

# Generate detailed performance report
echo "Generating detailed performance report..."
echo ""

# Create performance summary
echo "=== PERFORMANCE SUMMARY ===" > /tmp/green_performance_report.txt
echo "Generated: $(date)" >> /tmp/green_performance_report.txt
echo "" >> /tmp/green_performance_report.txt

echo "=== ALL BENCHMARK RESULTS ===" >> /tmp/green_performance_report.txt
cat /tmp/full_benchmark.log >> /tmp/green_performance_report.txt

echo ""
echo "📄 Detailed performance report saved to: /tmp/green_performance_report.txt"

print_status "HEADER" "Test Summary"
echo "============="

print_status "INFO" "Unit Tests: $passed_tests/$total_tests passed"
print_status "INFO" "Benchmarks: $passed_benchmarks/$total_benchmarks passed"

if [ $failed_tests -eq 0 ] && [ $failed_benchmarks -eq 0 ]; then
    print_status "PASS" "ALL GREEN FUNCTIONS TESTS PASSED!"
else
    print_status "FAIL" "Some tests failed. Please review the output above."
fi

echo ""
echo "=== SYSTEM INFORMATION ==="
echo "Go Version: $(go version)"
echo "OS: $(uname -s) $(uname -r)"
echo "Architecture: $(uname -m)"
echo "CPU: $(grep 'model name' /proc/cpuinfo | head -1 | cut -d: -f2 | xargs)"
echo "Memory: $(free -h | grep Mem | awk '{print $2}')"

echo ""
echo "=== GREEN FUNCTIONS TEST COMPLETE ==="

# Calculate success rates
if [ $total_tests -gt 0 ]; then
    test_success_rate=$(( (passed_tests * 100) / total_tests ))
else
    test_success_rate=0
fi

if [ $total_benchmarks -gt 0 ]; then
    benchmark_success_rate=$(( (passed_benchmarks * 100) / total_benchmarks ))
else
    benchmark_success_rate=0
fi

echo "Test Success Rate: ${test_success_rate}%"
echo "Benchmark Success Rate: ${benchmark_success_rate}%"

# Clean up temporary files
rm -f /tmp/green_test_output.log
rm -f /tmp/green_benchmark_output.log

if [ $failed_tests -eq 0 ] && [ $failed_benchmarks -eq 0 ]; then
    exit 0
else
    exit 1
fi 