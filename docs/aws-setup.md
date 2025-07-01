# AWS Setup Guide for Consistency Service

This guide provides detailed instructions to set up AWS infrastructure for staging and production environments.

## Prerequisites

- AWS CLI installed and configured (`aws configure`)
- AWS account with appropriate permissions
- GitHub repository with admin access

## 1. Create ECR Repository

```bash
# Create ECR repository
aws ecr create-repository \
    --repository-name consistency-service \
    --region us-east-1

# Get login token for Docker
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
```

**Note the repository URI:** `<account-id>.dkr.ecr.us-east-1.amazonaws.com/consistency-service`

## 2. Create VPC and Networking

```bash
# Create VPC
aws ec2 create-vpc --cidr-block 10.0.0.0/16 --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=consistency-vpc}]'

# Create Internet Gateway
aws ec2 create-internet-gateway --tag-specifications 'ResourceType=internet-gateway,Tags=[{Key=Name,Value=consistency-igw}]'

# Attach Internet Gateway to VPC
aws ec2 attach-internet-gateway --vpc-id <vpc-id> --internet-gateway-id <igw-id>

# Create public subnets (for ALB)
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.1.0/24 --availability-zone us-east-1a --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=consistency-public-1a}]'
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.2.0/24 --availability-zone us-east-1b --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=consistency-public-1b}]'

# Create private subnets (for ECS and RDS)
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.3.0/24 --availability-zone us-east-1a --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=consistency-private-1a}]'
aws ec2 create-subnet --vpc-id <vpc-id> --cidr-block 10.0.4.0/24 --availability-zone us-east-1b --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=consistency-private-1b}]'
```

## 3. Create Security Groups

```bash
# Security group for ALB
aws ec2 create-security-group \
    --group-name consistency-alb-sg \
    --description "Security group for ALB" \
    --vpc-id <vpc-id>

aws ec2 authorize-security-group-ingress \
    --group-id <alb-sg-id> \
    --protocol tcp \
    --port 80 \
    --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
    --group-id <alb-sg-id> \
    --protocol tcp \
    --port 443 \
    --cidr 0.0.0.0/0

# Security group for ECS tasks
aws ec2 create-security-group \
    --group-name consistency-ecs-sg \
    --description "Security group for ECS tasks" \
    --vpc-id <vpc-id>

aws ec2 authorize-security-group-ingress \
    --group-id <ecs-sg-id> \
    --protocol tcp \
    --port 8080 \
    --source-group <alb-sg-id>

# Security group for RDS
aws ec2 create-security-group \
    --group-name consistency-rds-sg \
    --description "Security group for RDS" \
    --vpc-id <vpc-id>

aws ec2 authorize-security-group-ingress \
    --group-id <rds-sg-id> \
    --protocol tcp \
    --port 5432 \
    --source-group <ecs-sg-id>
```

## 4. Create RDS Subnet Group

```bash
aws rds create-db-subnet-group \
    --db-subnet-group-name consistency-subnet-group \
    --db-subnet-group-description "Subnet group for Consistency service" \
    --subnet-ids <private-subnet-1a-id> <private-subnet-1b-id>
```

## 5. Create RDS Instances

### Staging Database
```bash
aws rds create-db-instance \
    --db-instance-identifier consistency-staging-db \
    --db-instance-class db.t3.micro \
    --engine postgres \
    --engine-version 15.4 \
    --allocated-storage 20 \
    --master-username postgres \
    --master-user-password 'your-secure-password' \
    --db-name consistency_service \
    --vpc-security-group-ids <rds-sg-id> \
    --db-subnet-group-name consistency-subnet-group \
    --backup-retention-period 7 \
    --storage-encrypted \
    --tags Key=Environment,Value=staging
```

### Production Database
```bash
aws rds create-db-instance \
    --db-instance-identifier consistency-prod-db \
    --db-instance-class db.t3.small \
    --engine postgres \
    --engine-version 15.4 \
    --allocated-storage 20 \
    --master-username postgres \
    --master-user-password 'your-secure-password' \
    --db-name consistency_service \
    --vpc-security-group-ids <rds-sg-id> \
    --db-subnet-group-name consistency-subnet-group \
    --backup-retention-period 30 \
    --storage-encrypted \
    --multi-az \
    --tags Key=Environment,Value=production
```

## 6. Store Database Credentials in Secrets Manager

```bash
# Staging credentials
aws secretsmanager create-secret \
    --name "consistency/staging/db" \
    --description "Database credentials for staging environment" \
    --secret-string '{"username":"postgres","password":"your-secure-password","host":"<staging-db-endpoint>","port":"5432","dbname":"consistency_service"}'

# Production credentials
aws secretsmanager create-secret \
    --name "consistency/production/db" \
    --description "Database credentials for production environment" \
    --secret-string '{"username":"postgres","password":"your-secure-password","host":"<prod-db-endpoint>","port":"5432","dbname":"consistency_service"}'
```

## 7. Create ECS Cluster

```bash
aws ecs create-cluster \
    --cluster-name consistency-cluster \
    --capacity-providers FARGATE \
    --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1
```

## 8. Create Application Load Balancer

```bash
# Create ALB
aws elbv2 create-load-balancer \
    --name consistency-alb \
    --subnets <public-subnet-1a-id> <public-subnet-1b-id> \
    --security-groups <alb-sg-id>

# Create target groups
aws elbv2 create-target-group \
    --name consistency-staging-tg \
    --protocol HTTP \
    --port 8080 \
    --vpc-id <vpc-id> \
    --target-type ip \
    --health-check-path /health \
    --health-check-interval-seconds 30 \
    --health-check-timeout-seconds 5 \
    --healthy-threshold-count 2 \
    --unhealthy-threshold-count 5

aws elbv2 create-target-group \
    --name consistency-prod-tg \
    --protocol HTTP \
    --port 8080 \
    --vpc-id <vpc-id> \
    --target-type ip \
    --health-check-path /health \
    --health-check-interval-seconds 30 \
    --health-check-timeout-seconds 5 \
    --healthy-threshold-count 2 \
    --unhealthy-threshold-count 5

# Create listeners
aws elbv2 create-listener \
    --load-balancer-arn <alb-arn> \
    --protocol HTTP \
    --port 80 \
    --default-actions Type=forward,TargetGroupArn=<staging-tg-arn>
```

## 9. Create IAM Roles

### ECS Task Execution Role
```bash
aws iam create-role \
    --role-name consistency-ecs-execution-role \
    --assume-role-policy-document '{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "Service": "ecs-tasks.amazonaws.com"
                },
                "Action": "sts:AssumeRole"
            }
        ]
    }'

aws iam attach-role-policy \
    --role-name consistency-ecs-execution-role \
    --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy

# Create custom policy for Secrets Manager access
aws iam create-policy \
    --policy-name consistency-secrets-policy \
    --policy-document '{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "secretsmanager:GetSecretValue"
                ],
                "Resource": [
                    "arn:aws:secretsmanager:us-east-1:*:secret:consistency/staging/db*",
                    "arn:aws:secretsmanager:us-east-1:*:secret:consistency/production/db*"
                ]
            }
        ]
    }'

aws iam attach-role-policy \
    --role-name consistency-ecs-execution-role \
    --policy-arn arn:aws:iam::<account-id>:policy/consistency-secrets-policy
```

### ECS Task Role
```bash
aws iam create-role \
    --role-name consistency-ecs-task-role \
    --assume-role-policy-document '{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "Service": "ecs-tasks.amazonaws.com"
                },
                "Action": "sts:AssumeRole"
            }
        ]
    }'
```

## 10. Create GitHub Actions IAM User

```bash
aws iam create-user --user-name github-actions-consistency

aws iam create-policy \
    --policy-name github-actions-consistency-policy \
    --policy-document '{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "ecr:GetAuthorizationToken",
                    "ecr:BatchCheckLayerAvailability",
                    "ecr:GetDownloadUrlForLayer",
                    "ecr:BatchGetImage",
                    "ecr:InitiateLayerUpload",
                    "ecr:UploadLayerPart",
                    "ecr:CompleteLayerUpload",
                    "ecr:PutImage"
                ],
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "ecs:RegisterTaskDefinition",
                    "ecs:UpdateService",
                    "ecs:DescribeServices",
                    "ecs:DescribeTaskDefinition"
                ],
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "iam:PassRole"
                ],
                "Resource": [
                    "arn:aws:iam::<account-id>:role/consistency-ecs-execution-role",
                    "arn:aws:iam::<account-id>:role/consistency-ecs-task-role"
                ]
            }
        ]
    }'

aws iam attach-user-policy \
    --user-name github-actions-consistency \
    --policy-arn arn:aws:iam::<account-id>:policy/github-actions-consistency-policy

aws iam create-access-key --user-name github-actions-consistency
```

## 11. Update ECS Task Definition

Update `ecs-task-def.json` with your actual values:

```json
{
  "family": "consistency-service-task",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::<account-id>:role/consistency-ecs-execution-role",
  "taskRoleArn": "arn:aws:iam::<account-id>:role/consistency-ecs-task-role",
  "containerDefinitions": [
    {
      "name": "app",
      "image": "<IMAGE_URI>",
      "portMappings": [
        { "containerPort": 8080, "hostPort": 8080 }
      ],
      "environment": [
        { "name": "ENV", "value": "<ENVIRONMENT>" },
        { "name": "LOG_LEVEL", "value": "info" },
        { "name": "LOG_PRETTY", "value": "false" }
      ],
      "secrets": [
        {
          "name": "DB_HOST",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:<account-id>:secret:consistency/<environment>/db:host::"
        },
        {
          "name": "DB_PORT",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:<account-id>:secret:consistency/<environment>/db:port::"
        },
        {
          "name": "DB_NAME",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:<account-id>:secret:consistency/<environment>/db:dbname::"
        },
        {
          "name": "DB_USER",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:<account-id>:secret:consistency/<environment>/db:username::"
        },
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:<account-id>:secret:consistency/<environment>/db:password::"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/consistency-service",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "essential": true
    }
  ]
}
```

## 12. Create CloudWatch Log Group

```bash
aws logs create-log-group --log-group-name /ecs/consistency-service
```

## 13. Create ECS Services

### Staging Service
```bash
aws ecs create-service \
    --cluster consistency-cluster \
    --service-name consistency-staging \
    --task-definition consistency-service-task:1 \
    --desired-count 1 \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[<private-subnet-1a-id>,<private-subnet-1b-id>],securityGroups=[<ecs-sg-id>],assignPublicIp=DISABLED}" \
    --load-balancers "targetGroupArn=<staging-tg-arn>,containerName=app,containerPort=8080"
```

### Production Service
```bash
aws ecs create-service \
    --cluster consistency-cluster \
    --service-name consistency-prod \
    --task-definition consistency-service-task:1 \
    --desired-count 2 \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[<private-subnet-1a-id>,<private-subnet-1b-id>],securityGroups=[<ecs-sg-id>],assignPublicIp=DISABLED}" \
    --load-balancers "targetGroupArn=<prod-tg-arn>,containerName=app,containerPort=8080"
```

## 14. Configure GitHub Secrets

In your GitHub repository, go to Settings > Secrets and Variables > Actions, and add:

- `AWS_ACCESS_KEY_ID`: From step 10
- `AWS_SECRET_ACCESS_KEY`: From step 10
- `AWS_ACCOUNT_ID`: Your AWS account ID
- `ECR_REPOSITORY`: consistency-service
- `ECS_CLUSTER`: consistency-cluster
- `ECS_EXECUTION_ROLE_ARN`: arn:aws:iam::<account-id>:role/consistency-ecs-execution-role
- `ECS_TASK_ROLE_ARN`: arn:aws:iam::<account-id>:role/consistency-ecs-task-role

## 15. Deploy

After completing the setup:

1. Commit and push your changes to `staging` or `main` branch
2. GitHub Actions will automatically build and deploy your application
3. Monitor the deployment in the ECS console

## URLs

- **Staging**: `http://<alb-dns-name>` (with staging listener rules)
- **Production**: `http://<alb-dns-name>` (with production listener rules)

## Cleanup (Optional)

To delete all resources, run the cleanup commands in reverse order of creation. 