param location string = az.resourceGroup().location

@description('Name of the Azure container registry')
param registry string = 'todoapp101'

var tags = {
    ENV: 'TEST'
    APP: 'TODO101'
}

resource acr 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing  = {
  name: registry
}


resource vnet 'Microsoft.Network/virtualNetworks@2022-11-01' = {
  name:  'Virtual-Network-${uniqueString(az.resourceGroup().id)}'
  location: location
  tags: tags
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.3.0.0/16'
      ]
    }
    subnets: [
    {
      name: 'subnet-app'
      properties: {
          addressPrefix: '10.3.0.0/23'
          privateEndpointNetworkPolicies: 'Disabled'
          privateLinkServiceNetworkPolicies: 'Disabled'
      }
    }
  ]
  }
}

resource acaEnv 'Microsoft.App/managedEnvironments@2022-11-01-preview' = {
  name: 'ACA-Env-${uniqueString(az.resourceGroup().id)}'
  location: location
  tags: tags
  properties: {
    vnetConfiguration: {
      internal: true
      infrastructureSubnetId: vnet.properties.subnets[0].id
    }
  }
}

resource acaAPI 'Microsoft.App/containerApps@2022-11-01-preview' = {
  name: 'api'
  location: location
  tags: tags
  properties: {
    configuration: {
      ingress: {
        allowInsecure: false
        external: false
        stickySessions: {
          affinity: 'sticky'
        }
        targetPort: 8080
        transport: 'auto'
      }
      registries: [
        {
          passwordSecretRef: 'acr-password'
          server: acr.properties.loginServer
          username: acr.listCredentials().username
        }
      ]
      secrets: [{
        name: 'acr-password'
        value: acr.listCredentials().passwords[0].value
      }]
    }
    managedEnvironmentId: acaEnv.id
    template: {
      containers: [
        {
          env: [
            {
              name: 'DATABASE_URI'
              value: ''
            }
          ]
          image: '${registry}.azurecr.io/todo-api'
          name: 'api'
          resources: {
            cpu: 1
            memory: '2Gi'
          }
        }
      ]
      scale: {
        maxReplicas: 4
        minReplicas: 1
      }
    }
  }
}

output vnetName string = vnet.name
output apiFqdn string = acaAPI.properties.configuration.ingress.fqdn
