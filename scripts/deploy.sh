#!/bin/bash

# Exit on any error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment...${NC}"

# Load environment variables from .env file if it exists
if [ -f ".env" ]; then
    echo -e "${YELLOW}Loading environment variables from .env file...${NC}"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check required environment variables
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}Error: DATABASE_URL environment variable is required${NC}"
    echo "Set it with: export DATABASE_URL='your-turso-database-url'"
    exit 1
fi

if [ -z "$JWT_KEY" ]; then
    echo -e "${RED}Error: JWT_KEY environment variable is required${NC}"
    echo "Set it with: export JWT_KEY='your-secret-key'"
    exit 1
fi

# Deploy frontend
echo -e "${YELLOW}Deploying frontend...${NC}"
cd frontend
yarn run build && yarn run deploy

# Deploy backend
echo -e "${YELLOW}Deploying backend...${NC}"
cd ../backend/

# Build Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t administratum-backend .

# Tag for Artifact Registry
echo -e "${YELLOW}Tagging image...${NC}"
docker tag administratum-backend europe-west1-docker.pkg.dev/administratum-468510/administratum-repo/backend:latest

# Push to Artifact Registry
echo -e "${YELLOW}Pushing to Artifact Registry...${NC}"
docker push europe-west1-docker.pkg.dev/administratum-468510/administratum-repo/backend:latest

# Deploy to Cloud Run with environment variables
echo -e "${YELLOW}Deploying to Cloud Run...${NC}"
gcloud run deploy backend \
  --image europe-west1-docker.pkg.dev/administratum-468510/administratum-repo/backend:latest \
  --region europe-west1 \
  --allow-unauthenticated \
  --set-env-vars "PLATFORM=production,DATABASE_URL=${DATABASE_URL},JWT_KEY=${JWT_KEY}"

echo -e "${GREEN}Deployment completed successfully!${NC}"
