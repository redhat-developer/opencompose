# User Guide

- [Commands](#commands)
  - [`convert`](#opencompose-convert)
  - [`completion`](#opencompose-completion)
  - [`validate`](#opencompose-validate)
  - [`version`](#opencompose-version)

## `opencompose convert`

Convert OpenCompose YAML files to Kubernetes (or OpenShift) artifacts.

### Convert to Kubernetes

```sh
opencompose convert -f hello-nginx.yaml
```

### Convert to OpenShift

```sh
opencompose convert -f hello-nginx.yaml --distro openshift
```

### Overriding the default output directory

```sh
mkdir kk8s_artifacts
opencompose convert -f hello-nginx.yaml -o k8s_artifacts
```

## `opencompose completion`

Output completion code for terminal autocompletion.

For persistance, add the below commands to your `.bashrc` and/or `.zshrc`.

### Bash

```sh
source <(kompose completion bash)
```

### Zsh

```sh
source <(kompose completion zsh)


## `opencompose validate`

Validate the OpenCompose YAML file if it is correct to specification standards. 

### Checking file validation

```sh
opencompose validate -f hello-nginx.yaml
```

## `opencompose version`

Output the current OpenCompose CLI tool version

### Version

```sh
opencompose version
```
