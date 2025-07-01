#!/bin/bash

echo "Building Go application for AWS Lambda..."

# Build for Lambda (Linux AMD64)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/lambda/main.go

# Make it executable
chmod +x bootstrap

echo "Build complete! Ready for serverless deployment."
echo ""
echo "To deploy:"
echo "1. npm install -g serverless"
echo "2. npm install serverless-domain-manager"  
echo "3. Get DB credentials from AWS Secrets Manager:"
echo "   export DB_USER=\$(aws secretsmanager get-secret-value --secret-id consistency-db-user-LUUpbF --query SecretString --output text)"
echo "   export DB_PASSWORD=\$(aws secretsmanager get-secret-value --secret-id consistency-db-password-LUUpbF --query SecretString --output text)"
echo "   export DB_HOST=consistency-prod-db.cs5c86m8c7jh.us-east-1.rds.amazonaws.com"
echo "4. serverless deploy"
echo ""
echo "Expected monthly cost: \$3-8 (vs current \$50-70)" 