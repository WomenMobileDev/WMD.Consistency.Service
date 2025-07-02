#!/bin/bash
set -e

echo "🚀 Starting deployment to EC2..."

# Environment variables (will be set by GitHub Actions)
INSTANCE_ID="${INSTANCE_ID:-}"
PUBLIC_IP="${PUBLIC_IP:-}"
ECR_REGISTRY="${ECR_REGISTRY:-}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo "📋 Deployment Configuration:"
echo "  Instance ID: $INSTANCE_ID"
echo "  Public IP: $PUBLIC_IP"
echo "  ECR Registry: $ECR_REGISTRY"
echo "  Image Tag: $IMAGE_TAG"

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
  -e DATABASE_URL="$DATABASE_URL" \
  -e JWT_SECRET="$JWT_SECRET" \
  -e ENV="production" \
  $ECR_REGISTRY:$IMAGE_TAG

echo "🔍 Verifying deployment..."
sleep 10

# Check if container is running
if docker ps | grep -q consistency-api; then
  echo "✅ Container is running successfully!"
  
  # Test health endpoint
  echo "🏥 Testing health endpoint..."
  if curl -f http://localhost/health; then
    echo ""
    echo "✅ Health check passed!"
  else
    echo ""
    echo "⚠️  Health check failed, but container is running"
  fi
else
  echo "❌ Container failed to start!"
  echo "📋 Container logs:"
  docker logs consistency-api || true
  exit 1
fi

echo "🧹 Cleaning up old images..."
docker image prune -f

echo "✅ Deployment completed successfully!"
echo "🌐 Your API is available at: http://$PUBLIC_IP" 