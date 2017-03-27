# OpenCompose file reference documentation


OpenCompose file example
```yml
version: "0.1-dev"

services:
- name: foobar
  replicas: 3
  labels:
    foo_label: bar_label
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
    mounts:
    - volumeName: db
      mountPath: /app/store
      volumeSubPath: foo/bar
      readOnly: true
  emptyDirVolumes:
  - name: temp

volumes:
- name: db
  size: 5GiB
  accessMode: ReadWriteMany
  storageClass: fast
```


OpenCompose file has three main sections *version*, *services* and *volumes*.

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

Each service has name and list of the containers, while `replicas` can be optional.

```yml
# <ServiceSpec>
  name: foo
  replicas: 4
  labels:
    foo_label: bar_label
  containers:
  - <ContainerSpec>
  emptyDirVolumes:
  - <EmptyDirVolumeSpec>
```

#### name
| type | required |
|------|----------|
|string|    yes   |

Name of the service.

#### replicas
| type    | required |
|---------|----------|
| integer |    no    |

Number of desired pods of this particluar service. This is an optional field. The valid value can only be a positive number.

#### labels
| type    | required |
|---------|----------|
| map with string keys and string values |    no    |

Desired labels to be applied to the resulting Kubernetes objects from the service.

### containers
| type                                    | required |
|-----------------------------------------|----------|
|array of [ContainerSpec](#containerspec) |    yes   |

Each item in [containers](#containers) ([ContainerSpec](#containerspec)) defines one container.
All the containers in container array for given service will be in same Pod.

### emptyDirVolumes
| type                                              | required |
|---------------------------------------------------|----------|
|array of [EmptyDirVolumeSpec](#emptydirvolumespec) |    no    |

Each item in [emptyDirVolumes](#emptydirvolumes) ([EmptyDirVolumeSpec](#emptydirvolumespec)) defines a EmptyDir volume.
These volumes will be shared among the containers in the service.

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
  mounts:
    - <MountSpec>
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

#### mounts

| type                            | required |
|---------------------------------|----------|
|array of [MountSpec](#mountspec) |  no      |

This defines what volumes will be mounted inside the container.


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

### MountSpec
```yml
# <MountSpec>
  volumeName: db
  mountPath: /app/store
  volumeSubPath: foo/bar
  readOnly: true
```

#### volumeName
| type | required | possible values                                                                 |
|------|----------|---------------------------------------------------------------------------------|
|string|    yes   | should conform to the definition of a subdomain in DNS (RFC 1123), [details](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md). |

Should match the name of volume from root level [`volumes`](#volumes) directive or `service` level [`emptyDirVolumes`](#emptydirvolumes) directive.

#### mountPath
| type | required |
|------|----------|
|string|    yes   |

Absolute path within the container at which the volume should be mounted. Must not contain ':'.

#### volumeSubPath
| type | required | Default value      |
|------|----------|--------------------|
|string|    no    | "" (volume's root) |

Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).

#### readOnly
| type | required | Default value |
|------|----------|---------------|
|bool  |    no    | `false`       |

Volume is mounted read-only if `true`, read-write otherwise.


### EmptyDirVolumeSpec
Describes one EmptyDir volume. This can be referenced from the [container mounts](#mountspec).

```yml
# <EmptyDirVolumeSpec>
  name: tmp
```

#### name
| type | required |
|------|----------|
|string|    yes   |

Name of the EmptyDir volume. 


## volumes
`volumes` is main section that lists all the volumes that this OpenCompose file describes.
`volumes` has to be array of [VolumeSpec](#volumespec).

### VolumeSpec
Describes one volume. This can be referenced from the [container mounts](#mountspec).

Each volume has size defined, the accessmodes defined and storage class from where this should be read.

```yml
# <VolumeSpec>
  name: db
  size: 5GiB
  accessMode: ReadWriteMany
  storageClass: fast
```

#### name
| type | required | possible values                                                                 |
|------|----------|---------------------------------------------------------------------------------|
|string|    yes   | should conform to the definition of a subdomain in DNS (RFC 1123), [details](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md). |

Name of the volume.

#### size
| type | required | possible values                                                                   |
|------|----------|-----------------------------------------------------------------------------------|
|string|    yes   | must match the regular expression `'^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$` |

Size of the volume created.

#### accessMode

| type | required | possible values                                      |
|------|----------|------------------------------------------------------|
|string|    yes   | `ReadWriteOnce` or `ReadOnlyMany` or `ReadWriteMany` |

AccessMode contains the desired access mode the volume should have.

#### storageClass

| type | required |
|------|----------|
|string|    no    | 

Name of the StorageClass that will back this volume.