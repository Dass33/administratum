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

if [ -z "$GCP_PROJECT_ID" ]; then
    echo -e "${RED}Error: GCP_PROJECT_ID environment variable is required${NC}"
    echo "Set it with: export GCP_PROJECT_ID='your-gcp-project-id'"
    exit 1
fi

if [ -z "$GCP_REGION" ]; then
    echo -e "${RED}Error: GCP_REGION environment variable is required${NC}"
    echo "Set it with: export GCP_REGION='your-gcp-region'"
    exit 1
fi

if [ -z "$GCP_ARTIFACT_REGISTRY_REPO" ]; then
    echo -e "${RED}Error: GCP_ARTIFACT_REGISTRY_REPO environment variable is required${NC}"
    echo "Set it with: export GCP_ARTIFACT_REGISTRY_REPO='your-artifact-registry-repo'"
    exit 1
fi

# Deploy frontend
echo -e "${YELLOW}Deploying frontend...${NC}"
cd frontend
yarn run build && yarn run deploy

# Check if Docker is running, start if needed
echo -e "${YELLOW}Checking Docker status...${NC}"
if ! docker info >/dev/null 2>&1; then
    echo -e "${YELLOW}Docker is not running. Starting Docker...${NC}"
    sudo systemctl start docker
    
    # Wait for Docker to start
    echo -e "${YELLOW}Waiting for Docker to start...${NC}"
    timeout=30
    while [ $timeout -gt 0 ] && ! docker info >/dev/null 2>&1; do
        sleep 1
        ((timeout--))
    done
    
    if ! docker info >/dev/null 2>&1; then
        echo -e "${RED}Error: Failed to start Docker after 30 seconds${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Docker started successfully${NC}"
else
    echo -e "${GREEN}Docker is already running${NC}"
fi

# Deploy backend
echo -e "${YELLOW}Deploying backend...${NC}"
cd ../backend/

# Build Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t administratum-backend .

# Tag for Artifact Registry
echo -e "${YELLOW}Tagging image...${NC}"
docker tag administratum-backend ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/${GCP_ARTIFACT_REGISTRY_REPO}/backend:latest

# Push to Artifact Registry
echo -e "${YELLOW}Pushing to Artifact Registry...${NC}"
docker push ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/${GCP_ARTIFACT_REGISTRY_REPO}/backend:latest

# Deploy to Cloud Run with environment variables
echo -e "${YELLOW}Deploying to Cloud Run...${NC}"
gcloud run deploy backend \
  --image ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT_ID}/${GCP_ARTIFACT_REGISTRY_REPO}/backend:latest \
  --region ${GCP_REGION} \
  --allow-unauthenticated \
  --set-env-vars "PLATFORM=production,DATABASE_URL=${DATABASE_URL},JWT_KEY=${JWT_KEY}"

echo -e "${GREEN}Deployment completed successfully!${NC}"
