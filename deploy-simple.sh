#!/bin/bash
set -e

echo "🚀 Starting deployment to EC2..."

# Environment variables should be set by GitHub Actions
ECR_REGISTRY="${ECR_REGISTRY:-649024131095.dkr.ecr.us-east-1.amazonaws.com/consistency-service}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo "📋 Deployment Configuration:"
echo "  ECR Registry: $ECR_REGISTRY"
echo "  Image Tag: $IMAGE_TAG"

echo "🔐 Logging into ECR..."
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 649024131095.dkr.ecr.us-east-1.amazonaws.com

echo "🐳 Pulling latest Docker image from ECR..."
docker pull $ECR_REGISTRY:$IMAGE_TAG

echo "🛑 Stopping existing container (if any)..."
docker stop consistency-api || true
docker rm consistency-api || true

echo "🚀 Starting new container..."
docker run -d \
  --name consistency-api \
  --restart unless-stopped \
  -p 80:8080 \
  -e PORT=8080 \
  -e ENV=production \
  -e LOG_LEVEL=info \
  -e LOG_PRETTY=false \
  -e DB_HOST=consistency-prod-db.cs5c86m8c7jh.us-east-1.rds.amazonaws.com \
  -e DB_PORT=5432 \
  -e DB_NAME=consistency_service \
  -e DB_USER=postgres \
  -e DB_PASSWORD=consistency1july \
  -e DB_SSL_MODE=require \
  $ECR_REGISTRY:$IMAGE_TAG

echo "🔍 Verifying deployment..."
sleep 15

# Check if container is running
if docker ps | grep -q consistency-api; then
  echo "✅ Container is running successfully!"
  
  # Test health endpoint
  echo "🏥 Testing health endpoint..."
  for i in {1..6}; do
    if curl -f http://localhost/health > /dev/null 2>&1; then
      echo "✅ Health check passed!"
      break
    else
      echo "⏳ Waiting for application to start... ($i/6)"
      sleep 10
    fi
  done
else
  echo "❌ Container failed to start!"
  echo "📋 Container logs:"
  docker logs consistency-api || true
  exit 1
fi

echo "🧹 Cleaning up old images..."
docker image prune -f

echo "✅ Deployment completed successfully!"
echo "🌐 Your API is available at: http://$(curl -s http://checkip.amazonaws.com)" 