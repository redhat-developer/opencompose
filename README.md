# OpenCompose

[![Build Status Widget]][Build Status] [![GoDoc Widget]][GoDoc] [![Slack Widget]][Slack]

OpenCompose is a declarative higher level abstraction for specific Kubernetes resources.

A developer shouldn't have to learn various Kubernetes concepts just to test and deploy their applications.
Focus on the application that is being developed.

The goal of OpenCompose is to make it easier for developers to on-board to Kubernetes.

We are at a very evolving stage of this project and we have listed some of our ideas as [issues](https://github.com/redhat-developer/opencompose/issues)
and [examples](https://github.com/redhat-developer/opencompose/blob/master/examples/).
We are open to suggestions and contributions from the Kubernetes community as our project grows.
Please send any PRs, issues or RFCs to improve this project.

## Use Case

Go from a simple [hello-nginx.yaml](https://github.com/redhat-developer/opencompose/blob/master/examples/hello-nginx.yaml) example to a full Kubernetes environment:

Create (or download) `hello-nginx.yaml`

```yaml
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: nginx
    ports:
    - port: 80:8080
      type: external
```

Convert the file using `opencompose`

```sh
opencompose convert -f hello-nginx.yaml
# Alternatively, you can pass the URL of the remote file
opencompose convert -f https://raw.githubusercontent.com/redhat-developer/opencompose/master/examples/hello-nginx.yaml
```

Deploy your generate artifacts to Kubernetes with `kubectl`

```sh
kubectl create -f helloworld-service.yaml -f helloworld-deployment.yaml
```

## Installation

#### Binary installation

The easiest way to install OpenCompose is through our binary on our [GitHub release page](https://github.com/redhat-developer/opencompose/releases).

```sh
# Linux 
curl -L https://github.com/redhat-developer/opencompose/releases/download/v0.2.0/opencompose-linux-amd64 -o opencompose

# macOS
curl -L https://github.com/redhat-developer/opencompose/releases/download/v0.2.0/opencompose-darwin-amd64 -o opencompose

# Windows
curl -L https://github.com/redhat-developer/opencompose/releases/download/v0.2.0/opencompose-windows-amd64.exe -o opencompose.exe

chmod +x ./opencompose
sudo mv ./opencompose /usr/local/bin/opencompose
```

#### Go (from source)

To test the latest changes (as our project is still in it's infancy), a simple `go get` is all you need to get the latest source.

```sh
go get -u github.com/redhat-developer/opencompose
```

Although the binary is installed via `go get`. In order to create a properly signed build (ex. `opencompose version`), you will have to build it with `make bin`:

```sh
go get -u github.com/redhat-developer/opencompose
make bin
```

## Shell autocompletion

We support both Bash and Zsh autocompletion.

```sh
# Bash (add to .bashrc for persistence)
source <(opencompose completion bash)

# Zsh (add to .zshrc for persistence)
source <(opencompose completion zsh)
```

## Documentation

 - [OpenCompose file reference documentation](https://github.com/redhat-developer/opencompose/blob/master/docs/file-reference.md)
 - You can find more in [doc/](https://github.com/redhat-developer/opencompose/tree/master/docs) folder

## Community

We always welcome your feedback and thoughts on the project! Come and join our mailing list - [opencompose [at] googlegroups.com](https://groups.google.com/forum/#!forum/opencompose). We also hangout on [slack.k8s.io](http://slack.k8s.io/) ([#sig-apps](https://kubernetes.slack.com/messages/sig-apps/)).

[Build Status]: https://travis-ci.org/redhat-developer/opencompose
[Build Status Widget]: https://travis-ci.org/redhat-developer/opencompose.svg?branch=master
[GoDoc]: https://godoc.org/github.com/redhat-developer/opencompose
[GoDoc Widget]: https://godoc.org/github.com/redhat-developer/opencompose?status.svg
[Slack]: http://slack.kubernetes.io#sig-apps
[Slack Widget]: https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png
