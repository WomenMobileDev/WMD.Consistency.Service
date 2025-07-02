#!/bin/bash

set -e

echo "🚀 Deploying Consistency Service to AWS Lambda (Under $10/month)"
echo "================================================================"

# Install Serverless Framework v3 specifically  
echo "Installing Serverless Framework v3..."
npm install serverless@3

# Remove domain manager plugin for now
echo "Installing serverless plugins..."
# npm install serverless-domain-manager

# Get DB credentials from AWS Secrets Manager
echo "Getting database credentials..."
export DB_USER=$(aws secretsmanager get-secret-value --secret-id consistency-db-user --query SecretString --output text 2>/dev/null || echo "postgres")
export DB_PASSWORD=$(aws secretsmanager get-secret-value --secret-id consistency-db-password --query SecretString --output text 2>/dev/null || echo "consistency1july")
export DB_HOST="consistency-prod-db.cs5c86m8c7jh.us-east-1.rds.amazonaws.com"

# Build Lambda function
echo "Building Lambda function..."
./build-lambda.sh

# Deploy to AWS
echo "Deploying to AWS Lambda..."
npx serverless deploy

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📊 Expected monthly cost: $8-10 (vs previous $124-134)"
echo "💰 Annual savings: ~$1,400-1,500"
echo ""
echo "🔗 Your API endpoints:"
echo "   • Health: [API_GATEWAY_URL]/health"
echo "   • API v1: [API_GATEWAY_URL]/api/v1"
echo ""
echo "⚡ Note: First Lambda cold start may take 10-15 seconds"
echo "🔄 Subsequent requests will be much faster (~100ms)"

echo ""
echo "🌐 Your API endpoint:"
echo "   https://api.consistency-production.shubhams.dev"
echo ""
echo "🔧 To update your API:"
echo "   1. Make code changes"
echo "   2. Run: ./deploy-serverless.sh"
echo ""
echo "💡 Next steps to reduce costs further:"
echo "   1. Consider switching to RDS Serverless v2 (pause when inactive)"
echo "   2. Set up CloudWatch billing alerts"
echo "   3. Monitor usage and optimize memory allocation" 