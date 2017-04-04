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
  host: domain.tld
  path: /admin

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

#### host 
| type        | required | possible values                                           |
|-------------|----------|-----------------------------------------------------------|
| string      |    no    | FQDN (fully qualified domain name) as defined by RFC 3986 |

Specifying host will make your application accessible outside of the cluster on domain specified as a value of `host`. (Note: this requires you to manually setup DNS records to point to your cluster's Ingress router. For development you can use services like http://nip.io/ or [http://xip.io/].)

#### path 
| type        | required | possible values                                               |
|-------------|----------|---------------------------------------------------------------|
| string      |    no    | An extended POSIX regex as defined by IEEE Std 1003.1         |
|             |          | (i.e this follows the egrep/unix syntax, not the perl syntax) |

- If you don't specify a `path` it matches all request to `host`.
- You have to specify `host` if you want to specify `path`.  

##### Path Based Ingresses
Path based ingresses specify a path component that can be compared against a URL (which requires that the traffic for the ingress be HTTP based) such that multiple ingresses can be served using the same host name, each with a different path. Ingress controller should match ingresses based on the most specific path to the least; however, this depends on the ingress controller implementation. The following table shows example ingresses and their accessibility:

| Ingress              | When Compared to     | Accessible |
|----------------------|----------------------|------------|
| www.example.com/test | www.example.com/test | Yes |
|                      | www.example.com      | No  |
| www.example.com/test and www.example.com | www.example.com/test | Yes |
|                      | www.example.com      | Yes |
| www.example.com      | www.example.com/test | Yes (Matched by the host, not the ingress) |
|                      | www.example.com      | Yes |
