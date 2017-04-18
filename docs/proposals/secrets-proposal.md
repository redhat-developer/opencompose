## Proposal for supporting secrets in OpenCompose

### Kubernetes Secrets

A very simple secret object in Kubernetes looks like this -

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: kubesec
type: Opaque
data:
  username: <base64 encoded value>
  password: <base64 encoded value>
```

Note that, as seen above, one secret can have multiple `data` values.

---

#### Creating a new secret:

To create a new secret, a top level `secrets` key can be used as follows:

 ```yaml
version: "0.1-dev"

services:
   name: foobar
   containers:
   - image: foo/bar:tag
     ports:
     - port: 8080:80
     env:
     - foo=bar
     - travis_creds:
         secret: ci_secret/travis
     mounts:
     - mountPath: /etc/jenkins.conf
       secret: ci_secret/jenkins

secrets:
    ci_secret:
      travis:
        file: /etc/travis.passwd
      jenkins:
        value: strongpassword@@
      circle: stillstrongpassword@@
```

Here, the short syntax treats the provided value as a secret literal.
```yaml
secrets:
  secret_name:
    secret_key: secret_value
```

Or, the long syntax will be more verbose like -

```yaml
secrets:
  secret_name:
    secret_key:
      value: secret_value
```

In order to retrieve a secret from a file, use -

```yaml
secrets:
  secret_name:
    secret_key:
      file: <file path>
```

---


To consume the above secret in OpenCompose, here is the proposed syntax.

Secrets in Kubernetes can be consumed either as an environment variable, or as a mounted volume.
 
- ##### Consuming an already created Kubernetes secret in OpenCompose as an environment variable
 
 ```yaml
version: "0.1-dev"

services:
   name: foobar
   containers:
   - image: foo/bar:tag
     ports:
     - port: 8080:80
     env:
     - foo=bar
     - db_user:
         secret: kubesec/username
     - db_pass:
         secret: kubesec/password
```

- ##### Consuming an already created Kubernetes secret in OpenCompose as a mounted volume
 
 ```yaml
version: "0.1-dev"

services:
   name: foobar
   containers:
   - image: foo/bar:tag
     ports:
     - port: 8080:80
     mounts:
     - mountPath: /var/secret_path/user_info
       secret: kubesec/username
     - mountPath: /var/secret_path/pass_info
       secret: kubesec/password
```

---

For the above mentioned syntax, the `secrets` definition at the container level will have the following fields -

- environment variables

`containers[].env[].secret` is an _optional_ field of type _string_. The string will consist of 2 parts separated by a `/`, the first being the secret name, and the second part being the key in the secret object.

- volume mounts

`containers[].mounts[].secret` is an _optional_ field of type _string_. The string will consist of 2 parts separated by a `/`, the first being the secret name, and the second part being the key in the secret object.

- top level secrets definition

`secrets.<secret name>` is a _mandatory_ field of type _string_.

`secrets.<secret name>.<secret key>` is a _mandatory_ field of type _string_. This can either by followed by the `file` or `value` directives, or be provided a _string_ containing the value of the secret key.

`secrets.<secret name>.<secret key>.file` is an _optional_ field of type _string_.

`secrets.<secret name>.<secret key>.value` is an _optional_ field of type _string_.