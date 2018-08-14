# Jira Operator

A Kubernetes operator to manage JIRA instances.

## Overview

This operator will deploy Jira...

## Usage

Deploy the operator and required resources.

```
kubectl apply -f deploy
```

Deploy a new Jira instance.

```
kubectl apply -f examples/jira-minimal.yaml
```

## Development

See the [Development documentation](doc/development.md).

## License

JIRA Operator is under Apache 2.0 license. See the [LICENSE][license_file] file for details.

[license_file]:./LICENSE
