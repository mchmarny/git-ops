# Continuous deployment using GitHub actions

![git-ops release on tag](https://github.com/mchmarny/git-ops/workflows/git-ops%20release%20on%20tag/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/git-ops)](https://goreportcard.com/report/github.com/mchmarny/git-ops) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/git-ops)

Simple pipeline for test, image, and deploy app onto Kubernetes cluster using GitHub action. 

![](image/diagram.png)

This pipeline is triggered by the creation of a release tag (e.g. `v.1.2.3`). It sets up the test environment and runs unit tests. If tests are successful, the pipeline creates configuration data, builds and publishes the image to the GitHub container registry, as well as deploys the resulting image using the simple deployment strategy. 

> This demo assumes yoo already have a Kubernetes cluster configured with Dapr. If not, consider the opinionated Dapr install in [dapr-demos/setup](https://github.com/mchmarny/dapr-demos/tree/master/setup).

## Demo

To walk-through the demo steps, start by navigating to the the already deployed app and take a note of the release version and the deployment time (in UTC):

https://gitops.thingz.io/

### Edit code

Next, edit the `staticMessage` variable in [app/main.go](app/main.go) to simulate developer making code changes:

> Make sure to save your changes

```go
const (
  greetingsMessage = "hello PDX"
)
```

Now, increment the version number variable (`APP_VERSION`) in the [app/Makefile](app/Makefile):

```shell
APP_VERSION ?=v0.1.5 # was v0.1.4
```

### Sync changes

Add, commit, and push your local changes upstream:

```shell
make sync
```

This will `git add`, `git commit`, and `git push` your changes to GitHub

### Create a release tag

When ready to make a release, tag it and push the tag to GitHub:

```shell
make tag
```

This will `git tag` it and `git push origin` your version tag to trigger to pipeline

> Note, the GitHub pipeline takes about ~2 min from the time you tag it to when the new app is deployed. To monitor the results either check the [GitHub notifications](https://github.com/notifications) or watch the action execute the [individual steps](https://github.com/mchmarny/git-ops/actions?query=workflow%3A%22git-ops+release+on+tag%22) although that link will depend on your GitHUb username (e.g. `mchmarny` above).

### View it

Once the pipeline is finished, you navigate again to the app. 

https://gitops.thingz.io/

If everything went well, the new release version should reflect the change you made to the variable (`APP_VERSION`) in the [app/Makefile](app/Makefile) and the deployment time (in UTC) should be also updated. 

If the changes are not there, check the [GitHub Action](https://github.com/mchmarny/git-ops/actions?query=workflow%3A%22git-ops+release+on+tag%22) to check on the status. 

## Setup Demo

### Initial deployment

To setup the demo, first create the namespace: 

```shell
kubectl apply -f config/space.yaml
```

If you have certs for the demo domain create a TLS secret:

```shell
kubectl create secret tls tls-secret -n gitops --key cert-pk.pem --cert cert-ca.pem
```

Now apply the Dapr component and the rest of [deployment](config/):

```shell
kubectl apply -f component/
kubectl apply -f config/
```

When the command completed, check on the status: 

```shell
kubectl get pods -n gitops
```

The response should include the `gitops` pod in status `Running` with container ready state `2/2`:

```shell
NAME                      READY   STATUS    RESTARTS   AGE
gitops-5fb4d4d6f9-6m74l   2/2     Running   0          25s
```

One last check on ingress: 

```shell
kubectl get ingress -n gitops
```

It should include `gitops` host as well as the cluster IP that's mapped in your DNS. If any of this sounds confusing, the the cluster setup instructions [here](https://github.com/mchmarny/dapr-demos/tree/master/setup).

```shell
NAME                   HOSTS              ADDRESS    PORTS   AGE
gitops-ingress-rules   gitops.thingz.io   x.x.x.x    80      19s
```

If everything went well, you should be able to navigate now to: 

https://gitops.thingz.io

### Kubernetes config

To enable GitHub action to deploy the built images to your cluster you'll first need to configure its context. If you already have authenticated to that cluster you can find this info in the `.kube` folder in your home directory. To ensure that the config has only information for that one cluster, it may be easier to simply export it from your managed Kubernetes provider.

> Warning, the exported file has sensitive information, make sure to delete it after

For AKS for example:

```shell
az aks get-credentials --name demo --file sa.json
```

next, create GitHub secret (named `KUBE_CONFIG`) with the content of that file on the repo where the action will run.

That concludes the setup. You can navigate to the top of this readme and run the [demo](#demo).

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
