## OpenCompose file reference

_A brief example of the OpenCompose specification:_

```yaml
version: "0.1-dev"

services:
- name: foobar
  replicas: 3
  labels:
    foo_label: bar_label
  containers:
  - name: baz
    image: foo/bar:tag
    env:
    - name: foo
      value: bar
    command: ["/bin/foobar"]
    args:
    - "-f"
    - "/tmp/foobar"
    ports:
    - port: 8080:80
      type: internal
      host: domain.tld
      path: /admin
    mounts:
    - volumeRef: db
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


The OpenCompose format has three main sections: *version*, *services* and *volumes*.

## Section: Version

```yaml
version: "0.1-dev"
```

| Type  | Required |
|-------|----------|
|string |    yes   |

The current version of the OpenCompose format

## Section: Services

```yaml
services:
- <Service>
  ...

- <Service>
  ...
```

`services` is main section that lists all the services. 

`services` is an array of [Service](#service).

### Service

Describes one service. Service can be composed of one or multiple containers which are scheduled
together and can communicate using localhost. (Containers in the same service share network namespace.)

Each service has name and list of the containers, while `replicas` can be optional.

_Example:_

```yaml
services:
- name: foobar
  replicas: 3
  labels:
    foo_label: bar_label
  containers:
  - <Container>
  emptyDirVolumes:
  - <EmptyDirVolume>
```

#### name

| Type | Required |
|------|----------|
|string|    yes   |

Name of the service.

#### replicas

| Type    | Required |
|---------|----------|
| integer |    no    |

Number of desired pods of this particluar service. This is an optional field. The valid value can only be a positive number.

#### labels

| Type    | Required |
|---------|----------|
| map with string keys and string values |    no    |

Desired labels to be applied to the resulting Kubernetes objects from the service.

#### containers

| Type                                    | Required |
|-----------------------------------------|----------|
|array of [Container](#container) |    yes   |

Each item in [containers](#containers) ([Container](#container)) defines one container.
All the containers in container array for given service will be in same Pod.

#### emptyDirVolumes

| Type                                              | Required |
|---------------------------------------------------|----------|
|array of [EmptyDirVolumeSpec](#emptydirvolumespec) |    no    |

Each item in [emptyDirVolumes](#emptydirvolumes) ([EmptyDirVolumeSpec](#emptydirvolumespec)) defines a EmptyDir volume.
These volumes will be shared among the containers in the service.

### Container

```yaml
services:
- name: foobar
  ...
  containers:
  - name: baz
    image: foo/bar:tag
    env:
    - name: foo
      value: bar
    command: ["/bin/foobar"]
    args:
    - "-f"
    - "/tmp/foobar"
    ports:
    - <Port>
    mounts:
    - <Mount>
  ...
```

Container describes an image to use, its arguments, which ports should be exposed and how.

#### name

| Type | Required |
|------|----------|
|string|    yes   |

Name of the container.

#### image

| Type | Required |
|------|----------|
|string|    yes   |

Name of the image that container will be started from.

#### env

| Type            | Required |
|-----------------|----------|
|array of [EnvVariables](#envVariables) |    no    |

List of environment variables to set in the container.
Each string has to be formated as `variable_name=value`.


#### command
**NOT YET IMPLEMENTED**

| Type           | Required |
|----------------|----------|
|array of strings|    no    |

Command ran by the container. Overrides default Entrypoint from Docker container image.

`command` field in in OpenCompose works in the same way as `command` filed in Kubernetes.

You can find more information in the [Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/container-command-args/)


#### args
**NOT YET IMPLEMENTED**

| Type           | Required |
|----------------|----------|
|array of strings|    no    |

Arguments passed to the command.
This overrides default Command from Docker container image.

`args` field in OpenCompose works in the same way as `args` field in Kubernetes.
 You can find more about it in [Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/container-command-args/)

#### ports

| Type                          | Required |
|-------------------------------|----------|
|array of [Port](#port) |  no      |

This defines what ports will be exposed for communication.

#### mounts

| Type                            | Required |
|---------------------------------|----------|
|array of [Mount](#mount) |  no      |

This defines what volumes will be mounted inside the container.


### envVariables

```yaml
services:
- name: foobar
  ...
  containers:
    - ...
      env:
      - name: foo
        value: bar
    ...
  ...
```

Inside the containers, the environment variables will generally be visible as `foo=bar`.

#### name

| Type | Required |
|------|----------|
|string|    yes   |

This is a string which defines the name of the environment variable being set.

#### value

| Type | Required |
|------|----------|
|string|    yes   |

This is string which defines the value of the environment variable being set.

### Port

```yaml
services:
- name: foobar
  ...
  containers:
    - ...
      ports:
        - port: 8080:80
          type: internal
          host: domain.tld
          path: /admin
    ...
  ...
```

#### port

| Type | Required |
|------|----------|
|string|    yes   |

This is a string in following format: `ContainerPort:ServicePort`
- `ContainerPort` is port inside container (where service inside container accepts new connections)
- `ServicePort` is port on which this service will be accessible for others.

`ServicePort` is optional. It defaults to `ContainerPort` if not specified.


#### type

| Type        | Required | possible values         |
|-------------|----------|-------------------------|
|enum(string) |    no    | `internal` or `external`|

Default value for Type is `internal`.

- `internal` - port will be accessible only from inside the Kubernetes cluster
- `external` - port will be accessible from inside and outside of the Kubernetes cluster

#### host

| Type        | Required | possible values                                           |
|-------------|----------|-----------------------------------------------------------|
| string      |    no    | FQDN (fully qualified domain name) as defined by RFC 3986 |

Specifying host will make your application accessible outside of the cluster on domain specified as a value of `host`. (Note: this requires you to manually setup DNS records to point to your cluster's Ingress router. For development you can use services like http://nip.io/ or [http://xip.io/].)

#### path

| Type        | Required | possible values                                               |
|-------------|----------|---------------------------------------------------------------|
| string      |    no    | An extended POSIX regex as defined by IEEE Std 1003.1         |
|             |          | (i.e this follows the egrep/unix syntax, not the perl syntax) |

- If you don't specify a `path` it matches all request to `host`.
- You have to specify `host` if you want to specify `path`.

##### Path Based Ingresses

Path based ingresses specify a path component that can be compared against a URL (which requires that the traffic for the ingress be HTTP based) such that multiple ingresses can be served using the same host name, each with a different path. Ingress controller should match ingresses based on the most specific path to the least; however, this depends on the ingress controller implementation. The following table shows example ingresses and their accessibility:

| Ingress              | When Compared to     | Accessible |
|----------------------|----------------------|------------|
| example.com/test     | example.com/test | Yes |
|                      | example.com      | No  |
| example.com/test and example.com | example.com/test | Yes |
|                      | example.com      | Yes |
| example.com      | example.com/test | Yes (Matched by the host, not the ingress) |
|                      | example.com      | Yes |

### Mount

```yaml
...
services:
- name: foobar
  ...
  containers:
    - ...
      mounts:
      - volumeRef: db
        mountPath: /app/store
        volumeSubPath: foo/bar
        readOnly: true
    ...
...
```

#### volumeRef

| Type | Required | possible values                                                                 |
|------|----------|---------------------------------------------------------------------------------|
|string|    yes   | should conform to the definition of a subdomain in DNS (RFC 1123), [details](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md). |

Should match the name of volume from root level [`volumes`](#volumes) directive or `service` level [`emptyDirVolumes`](#emptydirvolumes) directive.

#### mountPath

| Type | Required |
|------|----------|
|string|    yes   |

Absolute path within the container at which the volume should be mounted. Must not contain ':'.

#### volumeSubPath

| Type | Required | Default value      |
|------|----------|--------------------|
|string|    no    | "" (volume's root) |

Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).

#### readOnly

| Type | Required | Default value |
|------|----------|---------------|
|bool  |    no    | `false`       |

Volume is mounted read-only if `true`, read-write otherwise.


### EmptyDirVolume

Describes one EmptyDir volume. This can be referenced from the [container mounts](#mount).

```yaml
services:
- name: foobar
  ...
  emptyDirVolumes:
  - name: temp
```

#### name

| Type | Required |
|------|----------|
|string|    yes   |

Name of the EmptyDir volume. 


## Section: Volumes

`volumes` is main section that lists all the volumes that this OpenCompose file describes.
`volumes` has to be array of [VolumeSpec](#volumespec).

### Volume

Describes one volume. This can be referenced from the [container mounts](#mount).

Each volume has size defined, the accessmodes defined and storage class from where this should be read.

```yaml
volumes:
- name: db
  size: 5GiB
  accessMode: ReadWriteMany
  storageClass: fast
```

#### name

| Type | Required | possible values                                                                 |
|------|----------|---------------------------------------------------------------------------------|
|string|    yes   | should conform to the definition of a subdomain in DNS (RFC 1123), [details](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md). |

Name of the volume.

#### size

| Type | Required | possible values                                                                   |
|------|----------|-----------------------------------------------------------------------------------|
|string|    yes   | must match the regular expression `'^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$` |

Size of the volume created.

#### accessMode

| Type | Required | possible values                                      |
|------|----------|------------------------------------------------------|
|string|    yes   | `ReadWriteOnce` or `ReadOnlyMany` or `ReadWriteMany` |

AccessMode contains the desired access mode the volume should have.

#### storageClass

| Type | Required |
|------|----------|
|string|    no    | 

Name of the StorageClass that will back this volume.
