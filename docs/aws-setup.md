# AWS Setup Guide for Consistency Service

This guide helps you set up AWS resources for two environments: `staging` and `production` (main branch).

## 1. ECR (Elastic Container Registry)
- Create one ECR repository: `consistency-service`

## 2. RDS (PostgreSQL)
- Create two RDS PostgreSQL instances:
  - `consistency-staging-db`
  - `consistency-prod-db`
- Set DB name: `consistency_service`
- Set username and password (store in Secrets Manager)

## 3. Secrets Manager
- Store DB credentials for each environment:
  - `consistency-staging-db-user`
  - `consistency-staging-db-password`
  - `consistency-prod-db-user`
  - `consistency-prod-db-password`

## 4. ECS (Fargate)
- Create one ECS cluster: `consistency-cluster`
- Create two ECS services:
  - `consistency-staging`
  - `consistency-prod`
- Use the same task definition family: `consistency-service-task`
- Set environment variables and secrets in the ECS task definition (see `ecs-task-def.json`)

## 5. Networking
- Use a VPC with public/private subnets
- Attach security groups to allow traffic to app and DB

## 6. IAM
- Create a user/role for GitHub Actions with permissions for ECR, ECS, and Secrets Manager

## 7. GitHub Secrets
- Add the following to your repo settings:
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY`

## 8. Domain & SSL (Optional)
- Use Route 53 for DNS
- Use ACM for SSL certificates
- Attach an Application Load Balancer to ECS services

---

**After setup:**
- Update `ecs-task-def.json` with your actual RDS endpoints and Secrets ARNs.
- Push to `main` or `staging` to trigger deployment. 