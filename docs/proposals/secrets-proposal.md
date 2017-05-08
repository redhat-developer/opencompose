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
    - file: /tmp/travis.password
    - data: jenkins=strongpassword@@
```

The syntax will look like -

```yaml
secrets:
  secret_name:
  - file: <file path>
  - data: <data key>=<data value>
  - data: <data key>=<base64 encoded data>
    type: base64
  - data: <data key>=<data value>
    type: plaintext (default)

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

`secrets.<secret name>.[]type` is a _mandatory_ filed with allowed keys `file` and `data`.
`file` contains the path to the file containing secret data
`data` contains plain text secret data by default along with a plain text key for the value, but a base64 encoded input can also be provided by providing a `type: base64` key with the data.
All, `file`,`type`, `data` accept `string` values.
