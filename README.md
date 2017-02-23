# OpenCompose

[![Build Status](https://travis-ci.org/redhat-developer/opencompose.svg?branch=master)](https://travis-ci.org/redhat-developer/opencompose)

The goal of OpenCompose is to make it easier to on-board developers to Kubernetes.
It is a curated higher level abstraction for specific Kubernetes resources.
Very simple idea, isn't it? Developer shouldn't have to learn various Kubernetes concepts just to test and deploy their applications on Kubernetes.
The developer is generally concerned with the application that is being developed.

We are at a very nascent stage of this project and we have listed some of our ideas as [issues](https://github.com/redhat-developer/opencompose/issues)
and some of our ideas have taken shape in form of [examples](https://github.com/redhat-developer/opencompose/blob/master/examples/) 
but we are definitely looking for suggestions and contributions from the Kubernetes community.
Please send PRs to improve this project, file issues and RFEs against this repository.

### Installation
We don't have a release yet, but you can easily build OpenCompose tool from source.
All you need is [Go installed](https://golang.org/dl/) on your system. Than simply run:
```sh
go get -u https://github.com/redhat-developer/opencompose
```

### Example
1) Create or download [hello-nginx.yaml](https://github.com/redhat-developer/opencompose/blob/master/examples/hello-nginx.yaml) file.

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
This will create two files in current directory - *deployment-helloworld.yaml* and *service-helloworld.yaml*.
First file is a Kubernetes Deployment object and second is a Kubernetes Service.

To deploy your application to Kubernetes cluster just run `kubectl create -f service-helloworld.yaml -f deployment-helloworld.yaml`


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
