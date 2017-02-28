# OpenCompose

[![Build Status](https://travis-ci.org/redhat-developer/opencompose.svg?branch=master)](https://travis-ci.org/redhat-developer/opencompose)

The goal of OpenCompose is to make it easier for developers to on-board to Kubernetes.
OpenCompose is a declarative higher level abstraction for specific Kubernetes resources.

Very simple idea, isn't it? 
A developer shouldn't have to learn various Kubernetes concepts just to test and deploy their applications.
Focus on the application that is being developed.

We are at a very evolving stage of this project and we have listed some of our ideas as [issues](https://github.com/redhat-developer/opencompose/issues)
and [examples](https://github.com/redhat-developer/opencompose/blob/master/examples/).
We are open to suggestions and contributions from the Kubernetes community as our project grows.
Please send any PRs, issues or RFCs to improve this project.

### Installation
#### From source
You can easily build OpenCompose tool from source. All you need is [Go](https://golang.org/dl/):
```sh
go get -u github.com/redhat-developer/opencompose
```

#### Binaries
You can retrieve binaries for Linux, macOS and Windows on our [GitHub release page](https://github.com/redhat-developer/opencompose/releases).

##### Linux
If you prefer to do it straight from CLI here's a one-liner for you:
```bash
curl -L https://github.com/redhat-developer/opencompose/releases/download/v0.1.0-alpha.0/opencompose-v0.1.0-alpha.0-d0edfd9-linux-64bit.tar.xz | tar -xJf - -C ${HOME}/bin ./opencompose
```

### Example
1) Create or download [hello-nginx.yaml](https://github.com/redhat-developer/opencompose/blob/master/examples/hello-nginx.yaml).

```yaml
version: 0.1-dev
services:
- name: helloworld
  containers:
  - image: tomaskral/nonroot-nginx
    ports:
    - port: 8080
```

2) Convert OpenCompose file to Kubernetes objects

```sh
opencompose convert -f hello-nginx.yaml
```

This will create two Kubernetes files in current directory - `deployment-helloworld.yaml` and `service-helloworld.yaml`.

To deploy your application to Kubernetes run:

```sh
kubectl create -f service-helloworld.yaml -f deployment-helloworld.yaml
```


### Command-line Completions
#### Bash
For Bash auto completion run the following command:

```bash
source <(opencompose completion bash)
```

To make it permanent add this line to your `~/.bashrc`.

#### Zsh
For Zsh auto completion run the following command:

```zsh
source <(opencompose completion zsh)
```

To make it permanent add this line to your `~/.zshrc`.

### Documentation
 - [OpenCompose file reference documentation](https://github.com/redhat-developer/opencompose/blob/master/docs/file-reference.md)
 - You can find more in [doc/](https://github.com/redhat-developer/opencompose/tree/master/docs) folder
