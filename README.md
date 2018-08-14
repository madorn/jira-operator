# Jira Operator

A Kubernetes [Operator](https://coreos.com/operators/) to manage JIRA instances,
based on the [Operator SDK](https://github.com/operator-framework/operator-sdk).

## Overview

This goal of this operator is to manage the lifecycle of one or more instances
of [Jira Software](https://www.atlassian.com/software/jira) from
[Atlassian](https://www.atlassian.com/).

## Usage

The operator can be deployed manually into a Kubernetes cluster or through
[OLM](https://github.com/operator-framework/operator-lifecycle-manager), if that
is configured for the cluster.

Manifests for both options are provided.

#### Manual Deployment

Deploy the operator and required resources.

```
kubectl apply -f deploy/
```

After a few moments, the operator should be up and running.

```
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
jira-operator-5c558d888-swscw   1/1       Running   0          13s
```

#### OLM Deployment

If OLM is configured for the Kubernetes cluster, a manifest is available to
create the CSV for deploying Jira instances.

```
kubectl apply -f olm/jiraoperator.csv.yaml
```

#### Create Jira instance

Once the operator is running, create a new Jira instance.

```
kubectl apply -f examples/jira-minimal.yaml
```

This will create a minimal demo Jira instance, strictly for demonstration
purposes.

```
$ kubectl get jiras
NAME           AGE
jira-minimal   13s
```

Have a look at the other [examples](examples/) as well.

## Development

This operator is a work in progress, if you are interested in contributing, see
the [Development documentation](doc/development.md) to get started.

## License

JIRA Operator is under Apache 2.0 license. See the [LICENSE][license_file] file
for details.

[license_file]:./LICENSE
