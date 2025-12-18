#!/bin/bash

# Script to calculate test coverage for relevant packages only
# Excludes: cmd, protos, delivery, connections

echo "Running tests on relevant packages..."

# Get coverage for relevant packages
go test \
  ./pkg/... \
  ./internal/service/ad/config \
  ./internal/service/ad/domain/user \
  ./internal/service/ad/usecase/ad \
  ./internal/service/ad/usecase/budget \
  ./internal/service/ad/repository/postgres/ad \
  ./internal/service/ad/repository/postgres/budget \
  ./internal/service/auth/config \
  ./internal/service/auth/domain \
  ./internal/service/auth/service \
  ./internal/service/auth/storage/postgres \
  ./internal/service/profile/config \
  ./internal/service/profile/service \
  ./internal/service/profile/storage/postgres \
  ./internal/service/profile/external/http \
  ./internal/service/slot/config \
  ./internal/service/slot/domain/slot \
  ./internal/service/slot/usecase/slot \
  ./internal/service/slot/usecase/metric \
  ./internal/service/slot/repository/postgres/slot \
  ./internal/service/slot/repository/postgres/metric \
  ./internal/service/storage/config \
  ./internal/service/storage/usecase/filestorage \
  -coverprofile=coverage_relevant.out \
  2>&1 | grep -E "coverage:|ok |FAIL"

echo ""
echo "=== Coverage Summary ==="
go tool cover -func=coverage_relevant.out | tail -1
echo ""

# Generate HTML report if requested
if [ "$1" == "html" ]; then
  echo "Generating HTML coverage report..."
  go tool cover -html=coverage_relevant.out -o coverage_relevant.html
  echo "Report saved to coverage_relevant.html"
fi
