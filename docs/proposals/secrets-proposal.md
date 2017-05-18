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
 - name: foobar
   containers:
   - image: foo/bar:tag
     ports:
     - port: 8080:80
     env:
     - name: foo
       value: bar
     - name: travis_creds
       secretRef: ci_secret/travis
     mounts:
     - mountPath: /etc/jenkins.conf
       secretRef: ci_secret/jenkins

secrets:
  - name: ci_secret
    fromFile:
    - /tmp/travis.password
    fromLiteral:
    - jenkins=strongpassword@@
```

The syntax will look like -

```yaml
secrets:
  - name: <secret name>
    fromFile:
    - <file path>
    - <file path>
    fromLiteral:
    - <data key>=<base64 encoded data>
    - <data key>=<base64 encoded data>
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
     - name: foo
       value: bar
     - name: db_user
       secretRef: kubesec/username
     - name: db_pass
       secretRef: kubesec/password
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
       secretRef: kubesec/username
     - mountPath: /var/secret_path/pass_info
       secretRef: kubesec/password
```

---

For the above mentioned syntax, the `secrets` definition at the container level will have the following fields -

- environment variables

`containers.env.secretRef` is an _optional_ field of type _string_. The string will consist of 2 parts separated by a `/`, the first being the secret name, and the second part being the key in the secret object.

- volume mounts

`containers.mounts.secretRef` is an _optional_ field of type _string_. The string will consist of 2 parts separated by a `/`, the first being the secret name, and the second part being the key in the secret object.

- top level secrets definition

`secrets.name` is a _mandatory_ field of type _string_, which sets the name of the secret.

`secrets.fromFile` is an _optional_ field which takes an _array_ of file paths that contain the secret data.

`secrets.fromLiteral` is an _optional_ field which takes an _array_ of strings, of the form _plain text key=base64 encoded secret data_