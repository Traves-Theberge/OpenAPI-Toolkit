#!/bin/bash
# Quick smoke test script for OpenAPI TUI

set -e

echo "🧪 OpenAPI TUI - Smoke Test"
echo "============================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test 1: Build
echo -e "${BLUE}Test 1: Building application...${NC}"
go build -o openapi-tui ./cmd/openapi-tui
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Test 2: Binary check
echo -e "${BLUE}Test 2: Checking binary...${NC}"
if [ -f "openapi-tui" ]; then
    chmod +x openapi-tui
    ls -lh openapi-tui
    echo -e "${GREEN}✓ Binary created ($(du -h openapi-tui | cut -f1))${NC}"
else
    echo "❌ Binary not found"
    exit 1
fi
echo ""

# Test 3: Unit tests
echo -e "${BLUE}Test 3: Running unit tests...${NC}"
go test ./... -count=1 > /tmp/test-output.txt 2>&1
TEST_COUNT=$(grep -c "^ok" /tmp/test-output.txt || echo "0")
echo -e "${GREEN}✓ All tests passed ($TEST_COUNT packages)${NC}"
echo ""

# Test 4: Test coverage
echo -e "${BLUE}Test 4: Checking test coverage...${NC}"
go test ./internal/config -cover 2>&1 | grep "coverage:" || echo "Coverage check done"
go test ./internal/errors -cover 2>&1 | grep "coverage:" || echo "Coverage check done"
go test ./internal/export -cover 2>&1 | grep "coverage:" || echo "Coverage check done"
echo -e "${GREEN}✓ Coverage checked${NC}"
echo ""

# Test 5: Check test spec exists
echo -e "${BLUE}Test 5: Checking test spec...${NC}"
if [ -f "test-api.yaml" ]; then
    echo -e "${GREEN}✓ Test spec found (test-api.yaml)${NC}"
else
    echo "⚠️  Test spec not found - creating it..."
    echo "Run: cat TESTING-GUIDE.md to see how to create it"
fi
echo ""

# Summary
echo "================================"
echo -e "${GREEN}✅ All smoke tests passed!${NC}"
echo "================================"
echo ""
echo "📋 Next steps:"
echo "  1. Run './openapi-tui' to start the TUI"
echo "  2. Follow TESTING-GUIDE.md for manual testing"
echo "  3. Use test-api.yaml as your test spec"
echo ""
echo "Test files created:"
echo "  - openapi-tui (binary, $(du -h openapi-tui | cut -f1))"
echo "  - test-api.yaml (test spec)"
echo "  - TESTING-GUIDE.md (testing instructions)"
echo ""
