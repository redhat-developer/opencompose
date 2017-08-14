# HealthChecks

- This will be container level, optional field, called `health`.
- Under this field will be two sections `livenessProbe` or `readinessProbe`.
- Fields under both of these sections are same.
- The fields are all similar to how kubernetes specifies probes(`livenessProbe`/`readinessProbe`).

## Sample Examples

```yaml
containers:
- image: centos/httpd
  health:
    livenessProbe:
      httpGet:
        port: 8080
        path: /v1
      initialDelaySeconds: 10
      timeoutSeconds: 3
      failureThreshold: 5
```

**OR**

```yaml
containers:
- image: centos/httpd
  health:
    readinessProbe:
      exec:
        command:
        - cat
        - /foo/bar
```

The detailed information of this can be found at following:

- [API reference](https://kubernetes.io/docs/api-reference/v1.6/#probe-v1-core)
- [Configuring Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)

