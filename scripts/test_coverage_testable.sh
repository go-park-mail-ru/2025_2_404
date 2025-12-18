#!/bin/bash

# Script to calculate test coverage for testable code only
# Excludes: cmd, protos, delivery, connections, repository, storage, external

echo "Running tests on testable packages (excluding repository/storage/external)..."

# Get coverage for testable packages only
go test \
  ./pkg/convertImage \
  ./pkg/utils \
  ./pkg/logger \
  ./pkg/ReadYooKassaIP \
  ./internal/service/ad/config \
  ./internal/service/ad/domain/user \
  ./internal/service/ad/usecase/ad \
  ./internal/service/ad/usecase/budget \
  ./internal/service/auth/domain \
  ./internal/service/auth/service \
  ./internal/service/profile/config \
  ./internal/service/profile/service \
  ./internal/service/slot/config \
  ./internal/service/slot/usecase/slot \
  ./internal/service/slot/usecase/metric \
  ./internal/service/storage/config \
  ./internal/service/storage/usecase/filestorage \
  -coverprofile=coverage_testable.out \
  2>&1 | grep -E "coverage:|ok |FAIL"

echo ""
echo "=== Coverage Summary (Testable Code Only) ==="
go tool cover -func=coverage_testable.out | tail -1
echo ""
echo "Note: Excludes repository, storage, external, connections layers (require DB/HTTP mocks)"

# Generate HTML report if requested
if [ "$1" == "html" ]; then
  echo "Generating HTML coverage report..."
  go tool cover -html=coverage_testable.out -o coverage_testable.html
  echo "Report saved to coverage_testable.html"
fi
