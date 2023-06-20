######################################
## Define Variables - UPDATE VALUES
$BASE_NAME = "todoaca101"
$LOCATION = "westeurope"
######################################

## Resource Group & Deployment
$RESOURCE_GROUP_NAME = "$BASE_NAME-rg"
# $DEPLOYMENT_NAME = "$BASE_NAME-deployment-$(Get-Date -Format 'yyyyMMddHHmmss')"

## Register Providers
az provider register --wait --namespace Microsoft.App
az provider register --wait --namespace Microsoft.ContainerService
az provider register --wait --namespace Microsoft.Cdn

## Create Resource Group
az group create `
    --name $RESOURCE_GROUP_NAME `
    --location $LOCATION
Write-Host "...Created Resource Group"

$REGISTRY_NAME="acr$BASE_NAME"
$ACR_SERVER="$REGISTRY_NAME.azurecr.io"

# Create the Container Registry
az acr create --resource-group $RESOURCE_GROUP_NAME --name $REGISTRY_NAME --sku Basic
az acr login --name $REGISTRY_NAME
Write-Host "...Created container registry"

# Build docker images
Write-Host "Building docker images"
docker build -t todo/frontend:latest ../frontend/.
docker build -t todo/api:latest ../api/.

docker tag todo/frontend:latest "$ACR_SERVER/todo-frontend:latest"
docker tag todo/api:latest "$ACR_SERVER/todo-api:latest"

docker push "$ACR_SERVER/todo-frontend:latest"
docker push "$ACR_SERVER/todo-api:latest"

# Verify and list your images in the Registry
az acr repository list --name $REGISTRY_NAME --output table

## Deploy Template
 $RESULT = az deployment group create `
     --resource-group $RESOURCE_GROUP_NAME `
     --name $DEPLOYMENT_NAME `
     --template-file main.bicep `
     --parameters baseName=$BASE_NAME `
     --query properties.outputs.result

## Output Result
 $PRIVATE_LINK_ENDPOINT_CONNECTION_ID = $($RESULT | jq -r '.value.privateLinkEndpointConnectionId')
 $FQDN = $($RESULT | jq -r '.value.fqdn')
 $PRIVATE_LINK_SERVICE_ID = $($RESULT | jq -r '.value.privateLinkServiceId')

# FALLBACK: Private Link Service approval
# if ([string]::IsNullOrEmpty($PRIVATE_LINK_ENDPOINT_CONNECTION_ID)) {
#     Write-Host "Failed to get privateLinkEndpointConnectionId"
#     while ([string]::IsNullOrEmpty($PRIVATE_LINK_ENDPOINT_CONNECTION_ID)) {
#         Write-Host "- retrying..."
#         $PRIVATE_LINK_ENDPOINT_CONNECTION_ID = $(az network private-endpoint-connection list --id $PRIVATE_LINK_SERVICE_ID --query "[0].id" -o tsv)
#         Start-Sleep -Seconds 5
#     }
# }

## Approve Private Link Service
 Write-Host "Private link endpoint connection ID: $PRIVATE_LINK_ENDPOINT_CONNECTION_ID"
 az network private-endpoint-connection approve --id $PRIVATE_LINK_ENDPOINT_CONNECTION_ID --description "(Frontdoor) Approved by CI/CD"

 Write-Host "...Deployment FINISHED!"
 Write-Host "Please wait a few minutes until endpoint is established..."
 Write-Host "--- FrontDoor FQDN: https://$FQDN ---"
