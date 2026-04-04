#!/bin/sh

echo "🔍 Running pre-commit checks..."

# Format code
echo "👉 Formatting code..."
go fmt ./...

# Vet code
echo "👉 Running go vet..."
go vet ./...
if [ $? -ne 0 ]; then
  echo "❌ go vet failed"
  exit 1
fi

# Lint
echo "👉 Running golangci-lint..."
golangci-lint run
if [ $? -ne 0 ]; then
  echo "❌ Lint failed"
  exit 1
fi

# Check for secrets (basic)
echo "👉 Checking for secrets..."
if git diff --cached | grep -E "AKIA|SECRET|PASSWORD"; then
  echo "❌ संभावित secret मिला! Commit रोक दिया गया"
  exit 1
fi

echo "✅ Pre-commit checks passed"
