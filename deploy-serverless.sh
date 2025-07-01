#!/bin/bash

set -e

echo "🚀 Deploying Consistency Service to AWS Lambda (Under $10/month)"
echo "================================================================"

# Check if serverless is installed
if ! command -v serverless &> /dev/null; then
    echo "Installing Serverless Framework..."
    npm install -g serverless
fi

# Install serverless plugins
if [ ! -d "node_modules" ]; then
    echo "Installing serverless plugins..."
    npm init -y
    npm install serverless-domain-manager
fi

# Get database credentials from AWS Secrets Manager
echo "Getting database credentials..."
export DB_USER=$(aws secretsmanager get-secret-value --secret-id consistency-db-user-LUUpbF --query SecretString --output text --region us-east-1)
export DB_PASSWORD=$(aws secretsmanager get-secret-value --secret-id consistency-db-password-LUUpbF --query SecretString --output text --region us-east-1)
export DB_HOST=consistency-prod-db.cs5c86m8c7jh.us-east-1.rds.amazonaws.com

echo "Building Lambda function..."
./build-lambda.sh

echo "Deploying to AWS Lambda..."
serverless deploy

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📊 Cost comparison:"
echo "   Before (ECS): ~$50-70/month"
echo "   After (Lambda): ~$3-8/month"
echo "   💰 Savings: ~$42-62/month ($500-750/year)"
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