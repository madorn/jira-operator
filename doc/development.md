# Development

Work in progress...

## Out-of-Cluster

Ensure dependencies are safely vendored in the project.

```
dep ensure -v
```

Start the operator out-of-cluster.

```
LOG_LEVEL="debug" operator-sdk up local --namespace <namespace-to-watch> --kubeconfig <path-to-kubeconfig>
```

## In-Cluster

Build the operator using the SDK.

```
operator-sdk build <REPO>/jira-operator
```

Push the new operator image to the remote repository.

```
docker push <REPO>/jira-operator
```
