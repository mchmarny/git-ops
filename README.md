# Continuous deployment demo using GitHub actions

> this demo is still being developed, don't use it!

> cluster setup 

## Demo 

### Edit it

Start by editing the `staticMessage` variable in [app/main.go](app/main.go) to simulate developer making code changes:

> Make sure to save your changes

```go
const (
	staticMessage = "hello PDX"
)
```

Then increment the version number variable (`APP_VERSION`) in the [app/Makefile](app/Makefile):

```shell
APP_VERSION ?=v0.1.4
```

### Tag it

When ready to make a release, tag it and push the tag to GitHub:

```shell
make tag
```

This will `git tag` it and `git push origin` your version tag to trigger to pipeline

### View it

Navigate to the cluster where the app is deployed to get the current release:

https://gitops.thingz.io/

You can also monitor the GitOps pipeline to see when you are ready to refresh the app in the browser:


## Setup 

### Deploy

To setup the demo, first create the namespace: 

```shell
kubectl apply -f k8s/ns.yaml
```

If you have TLS certs for this the demo domain create a TLS secret 

```shell
kubectl create secret tls tls-secret -n gitops --key cert-pk.pem --cert cert-ca.pem
```

Than applying the rest:

```shell
kubectl apply -f k8s/
```

Check on the status: 

```shell
kubectl get pods -n gitops
```

The response should include the `gitops` pod in status `Running` with container ready state `2/2`:

```shell
NAME                      READY   STATUS    RESTARTS   AGE
gitops-5fb4d4d6f9-6m74l   2/2     Running   0          25s
```

Also, check on the ingress: 

```shell
kubectl get ingress -n gitops
```

Should include `gitops` host as well as the cluster IP mapped in your DNS:

```shell
NAME                   HOSTS              ADDRESS    PORTS   AGE
gitops-ingress-rules   gitops.thingz.io   x.x.x.x    80      19s
```

If everything went well, you should be able to navigate now to: 

https://gitops.thingz.io

## GitHub

To configure the GitHub action so it can apply new builds to your cluster, first you'll need to get your service principal. For AKS you can run:

```shell
az ad sp create-for-rbac --sdk-auth
```

The resulting file will look something like this:

```json
{
  "clientId": "...",
  "clientSecret": "...",
  "subscriptionId": "...",
  "tenantId": "...",
  "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
  "resourceManagerEndpointUrl": "https://management.azure.com/",
  "activeDirectoryGraphResourceId": "https://graph.windows.net/",
  "sqlManagementEndpointUrl": "https://management.core.windows.net:8443/",
  "galleryEndpointUrl": "https://gallery.azure.com/",
  "managementEndpointUrl": "https://management.core.windows.net/"
}
```

Copy that JSON and create following secrets in your GitHub repo where the action will run:

* `AZURE_CREDENTIALS` - with the content of the above file 
* `AZURE_CLUSTER_NAME` - with the name of your cluster 
* `AZURE_RESOURCE_GROUP` - with the name of your Azure resource group 


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
