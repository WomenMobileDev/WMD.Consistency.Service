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
# Extract registry URL from ECR_REGISTRY (remove repository name if present)
REGISTRY_URL=$(echo $ECR_REGISTRY | cut -d'/' -f1)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $REGISTRY_URL

echo "🐳 Pulling latest Docker image from ECR..."
echo "Attempting to pull: $ECR_REGISTRY:$IMAGE_TAG"

# Try to pull the specific commit image first, fall back to latest
if ! docker pull $ECR_REGISTRY:$IMAGE_TAG; then
  echo "⚠️  Failed to pull $IMAGE_TAG, trying 'latest' tag..."
  if ! docker pull $ECR_REGISTRY:latest; then
    echo "❌ Failed to pull both $IMAGE_TAG and latest tags"
    echo "This might be the first deployment or build step failed"
    echo "Check GitHub Actions build logs for Docker build/push errors"
    exit 1
  else
    echo "✅ Using latest tag instead"
    IMAGE_TAG="latest"
  fi
else
  echo "✅ Successfully pulled $IMAGE_TAG"
fi

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