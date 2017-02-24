# OpenCompose file reference documentation


OpenCompose file example
```yml
version: "0.1-dev"

services:
   name: foobar
   containers:
   - image: foo/bar:tag
     env:
     - foo=bar
     command: ["/bin/foobar"]
     args: 
     - "-f"
     - "/tmp/foobar"
     ports:
     - port: 8080:80

```


OpenCompose file has two main sections *version* and *services*.

## version
| type  | required |
|-------|----------| 
|string |    yes   |

Version of OpenCompose format.

## services
`services` is main section that lists all the services that this OpenCompose file describes.
`services` has to be array of [ServiceSpec](#servicespec).

### ServiceSpec
Describes one service. Service can be composed of one or multiple containers which are scheduled 
together and can communicate using localhost. (Containers in the same service share network namespace.)

Each service has name and list of the containers.

```yml
# <ServiceSpec>
  name: foo
  containers:
  - <ContainerSpec>
```

#### name
| type | required |
|------|----------|
|string|    yes   |

Name of the service.

### containers
| type                                    | required |
|-----------------------------------------|----------| 
|array of [ContainerSpec](#containerspec) |    yes   |

Each item in [containers](#containers) ([ContainerSpec](#containerspec)) defines one container.
All the containers in container array for given service will be in same Pod.


### ContainerSpec
```yml
# <ContainerSpec>
  image: foo/bar:tag
  env:
  - VARIABLE=value
  command:
    - /foo/bar
  args:
    - "-foo"
  ports:
    - <PortSpec>
```

ContainerSpec describes an image to use, its arguments, which ports should be exposed and how.

#### image
| type | required |
|------|----------|
|string|    yes   |

Name of the image that container will be started from.

#### env
| type            | required |
|-----------------|----------|
|array of strings |    no    |

List of environment variables to set in the container.
Each string has to be formated as `variable_name=value`.


#### command
**NOT YET IMPLEMENTED**

| type           | required |
|----------------|----------|
|array of strings|    no    |

Command run by the container.
It overrides default Entrypoint from Docker container image.

`command` field in in OpenCompose works in the same way as `command` filed in Kubernetes.
 You can find more about it in [Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/container-command-args/)


#### args
**NOT YET IMPLEMENTED**

| type           | required |
|----------------|----------|
|array of strings|    no    |

Arguments passed to the command.
This overrides default Command from Docker container image.

`args` field in OpenCompose works in the same way as `args` field in Kubernetes.
 You can find more about it in [Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/container-command-args/)

#### ports

| type                          | required |
|-------------------------------|----------|
|array of [PortSpec](#portspec) |  no      |

This defines what ports will be exposed for communication.


### PortSpec
```yml
# <PortSpec>
  port: 8080:80
  type: internal

```

#### port 
| type | required |
|------|----------|
|string|    yes   |

This is a string in following format: `ContainerPort:ServicePort`
- `ContainerPort` is port inside container (where service inside container accepts new connections)
- `ServicePort` is port on which this service will be accessible for others.

`ServicePort` is optional. It defaults to `ContainerPort` if not specified.


#### type 
| type        | required | possible values         |
|-------------|----------|-------------------------|
|enum(string) |    no    | `internal` or `external`|

Default value for type is `internal`.

- `internal` - port will be accessible only from inside the Kubernetes cluster
- `external` - port will be accessible from inside and outside of the Kubernetes cluster

