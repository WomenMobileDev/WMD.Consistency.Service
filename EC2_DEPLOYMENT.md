# EC2 Deployment Guide

This guide explains how to deploy the Consistency API to AWS EC2 using Docker for cost-effective hosting.

## Cost Breakdown

**Monthly Costs (approximately $8-12/month):**
- EC2 t3.micro (free tier): $0-8/month
- RDS db.t4g.micro: $8/month  
- Data transfer: ~$1/month
- **Total: $8-12/month** (vs $124-134/month with ECS Fargate)

## Prerequisites

1. **AWS Account** with free tier eligibility
2. **EC2 Instance** (t3.micro recommended for free tier)
3. **RDS Database** (existing production database)
4. **GitHub Secrets** configured for deployment

## EC2 Instance Setup

### 1. Launch EC2 Instance

```bash
# Launch a t3.micro instance with Amazon Linux 2
# Security Group: Allow HTTP (80), HTTPS (443), SSH (22)
# Storage: 8GB (free tier)
```

### 2. Install Docker and Docker Compose

```bash
# Connect to your EC2 instance
ssh -i your-key.pem ec2-user@your-ec2-ip

# Update system
sudo yum update -y

# Install Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install AWS CLI v2
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Log out and back in for group changes to take effect
exit
```

### 3. Configure AWS CLI

```bash
# Configure AWS credentials for ECR access
aws configure
# Enter your AWS Access Key ID, Secret Key, and region (us-east-1)
```

## GitHub Actions Setup

### Required GitHub Secrets

Add these secrets to your GitHub repository:

```
EC2_SSH_PRIVATE_KEY  # Your EC2 private key (.pem file content)
EC2_HOST            # Your EC2 public IP or domain
EC2_USER            # ec2-user (for Amazon Linux)
AWS_ACCESS_KEY_ID   # AWS access key with ECR permissions
AWS_SECRET_ACCESS_KEY # AWS secret key
```

### ECR Repository

Ensure your ECR repository exists:
```bash
aws ecr create-repository --repository-name consistency-service --region us-east-1
```

## Deployment Process

1. **Push to main branch** triggers automatic deployment
2. **GitHub Actions** builds Docker image and pushes to ECR
3. **Deployment script** is copied to EC2 and executed
4. **Docker Compose** pulls and runs the latest image

## Manual Deployment

If you need to deploy manually:

```bash
# On your local machine
docker build -t your-account.dkr.ecr.us-east-1.amazonaws.com/consistency-service:latest .
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin your-account.dkr.ecr.us-east-1.amazonaws.com
docker push your-account.dkr.ecr.us-east-1.amazonaws.com/consistency-service:latest

# On EC2 instance
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin your-account.dkr.ecr.us-east-1.amazonaws.com
docker pull your-account.dkr.ecr.us-east-1.amazonaws.com/consistency-service:latest
docker run -d --name consistency-api -p 80:8080 \
  -e ENV=production \
  -e DB_HOST=your-rds-endpoint \
  -e DB_PORT=5432 \
  -e DB_NAME=consistency_service \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your-password \
  your-account.dkr.ecr.us-east-1.amazonaws.com/consistency-service:latest
```

## Database Configuration

Update your RDS security group to allow connections from your EC2 instance:

```bash
# Add EC2 security group or specific IP to RDS security group
# Port: 5432 (PostgreSQL)
```

## Monitoring and Logs

```bash
# Check application status
curl http://your-ec2-ip/health

# View logs
docker logs consistency-api

# Check running containers
docker ps

# Restart application
docker-compose -f docker-compose.production.yml restart
```

## SSL/HTTPS Setup (Optional)

For production, consider setting up a reverse proxy with SSL:

```bash
# Install nginx
sudo yum install -y nginx

# Configure nginx as reverse proxy with Let's Encrypt SSL
# This will enable HTTPS access to your API
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Ensure ec2-user is in docker group
2. **ECR login fails**: Check AWS credentials and permissions
3. **Database connection**: Verify RDS security group allows EC2 access
4. **Health check fails**: Check application logs and database connectivity

### Useful Commands

```bash
# Check application health
curl http://localhost/health

# View application logs
docker logs consistency-api -f

# Restart application
docker-compose down && docker-compose up -d

# Check disk space
df -h

# Clean up old Docker images
docker image prune -f
```

## Security Considerations

1. **SSH Key Management**: Keep your EC2 private key secure
2. **Database Credentials**: Consider using AWS Secrets Manager
3. **Security Groups**: Restrict access to necessary ports only
4. **Regular Updates**: Keep EC2 instance and Docker updated

## Scaling

For higher traffic, consider:
- Upgrading to larger EC2 instance types
- Using Application Load Balancer with multiple instances
- Implementing auto-scaling groups
- Using RDS read replicas

## Backup Strategy

1. **Database**: RDS automated backups (enabled by default)
2. **Application**: Docker images stored in ECR
3. **Configuration**: All configs in Git repository 