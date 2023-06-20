@description('Basename / Prefix of all resources')
param baseName string

@description('Azure Location/Region')
param location string 

@description('Id of the Container Apps Environment')
param containerAppsEnvironmentId string

@description('Container Image')
param containerImage string

// Define names
var appName = '${baseName}-cont-hello-world-app'

resource containerApp 'Microsoft.App/containerApps@2022-03-01' = {
  name: appName
  location: location
  properties: {
    managedEnvironmentId: containerAppsEnvironmentId
    configuration: {
      ingress: {
        external: true
        targetPort: 3000
      }
    }
    template: {
      containers: [
        {
          name: 'app'
          image: containerImage
          probes: [
            {
              type: 'Liveness'
              httpGet: {
                path: '/health'
                port: 3000
              }
              periodSeconds: 10
              failureThreshold: 3
              initialDelaySeconds: 20
            }
          ]
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 10
      }
    }
  }
}

output containerFqdn string = containerApp.properties.configuration.ingress.fqdn
