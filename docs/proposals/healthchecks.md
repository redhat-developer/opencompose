# HealthChecks


- This will be container level, optional field, called `health`
- Liveness and Readiness Probes will be generated out of this in the `deployment.spec.template.spec.containers`.

## Sample example

```yaml
containers:
- image: centos/httpd
  health:
  - ready: true
    get: "http://localhost:8080/v1"
    delay: 10
    timeout: 3
    tries: 5


  - ready: true
    exec: "cat /foo/bar"

  - ready: true
    get: "http://localhost:8080/healthz"

  - live: true
    exec: "cat /foo/bar"

  - live: true
    get: "http://localhost:8080/healthz"
```

## Field description

| Field name | Type   | Optional | Conflicting field | Details                                                                                                       | Converted to (append to `deployment.spec.template.spec.containers`)                                                                          |
|------------|--------|----------|-------------------|---------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| `ready`    | bool   | no       | `live`            | Generate Readiness Probe config.                                                                              | `readinessProbe`                                                               |
| `live`     | bool   | no       | `ready`           | Generate Liveness Probe config.                                                                               | `livenessProbe`                                                                |
| `get`      | string | no       | `exec`            | A URL to do HTTP request on. Detailed information of what field in URL gets converted to what is given below. | `readinessProbe.httpGet` **OR** `livenessProbe.httpGet`                         |
| `exec`     | string | no       | `get`             | A command to run inside the container.                                                                        | `readinessProbe.exec` **OR** `livenessProbe.exec`                               |
| `delay`    | uint32 | yes      | None              | Number of seconds after the container has started before probes are initiated.                                | `livenessProbe.initialDelaySeconds` **OR** `readinessProbe.initialDelaySeconds` |
| `timeout`  | uint32 | yes      | None              | Number of seconds after which the probe times out.                                                            | `livenessProbe.timeoutSeconds` **OR** `readinessProbe.timeoutSeconds`           |
| `tries`    | uint32 | yes      | None              | Minimum consecutive failures for the probe to be considered failed after, having succeeded.                   | `livenessProbe.failureThreshold` **OR** `readinessProbe.failureThreshold`       |


## URL explanation

A URL gets deconstructed as follows

```
http://localhost:8080/healthz
```

| Field       | Optional | Details                                                                                                                                                  | Converted to (append to `deployment.spec.template.spec.containers`)                                                            |
|-------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| `http`      | Yes      | Scheme to use for connecting to the host.                                                                                                                | `readinessProbe.httpGet.scheme` |
| `localhost` | No       | If the value is `localhost` or `127.0.0.1` then it is not assigned to anything after conversion. But if the value is some other valid IP it will be set. | `readinessProbe.httpGet.host`   |
| `8080`      | No       | Number of the port to access on the container. Number must be in the range 1 to 65535.                                                                   | `readinessProbe.httpGet.port`   |
| `healthz`   | Yes      | Path to access on the HTTP server.                                                                                                                       | `readinessProbe.httpGet.path`   |

**NOTE**: In above converted to field for each URL field if `livenessProbe` is enabled then those respective fields will be created.

## Conflicting fields

- `ready` and `live` are conflicting fields and cannot co-exist but it is also important that either of them should be present at any point of time.
- `get` and `exec` are conflicting fields and cannot co-exist but it is also important that either of them should be present at any point of time.

