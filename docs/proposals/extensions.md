## Proposal for extending an OpenCompose file with another OpenCompose file

##### Scenario:

Let's say, we have a sample OpenCompose file defining a Wordpress application -

_wordpress-opencompose.yml_ -
```yaml
version: '0.1-dev'

services:
- name: database
  containers:
  - image: mariadb:10
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: rootpasswd
    - name: MYSQL_DATABASE
      value: wordpress
    - name: MYSQL_USER
      value: wordpress
    - name: MYSQL_PASSWORD
      value: wordpress
    ports:
    - port: 3306

- name: web
  containers:
  - image: wordpress:4
    env:
    - name: WORDPRESS_DB_HOST
      value: database:3306
    - name: WORDPRESS_DB_PASSWORD
      value: wordpress
    - name: WORDPRESS_DB_USER
      value: wordpress
    - name: WORDPRESS_DB_NAME
      value: wordpress
    ports:
    - port: 80
      type: external
```

This file works fine for a developer, but since there is no persistent volumes defined here, we need to alter this for production use cases; but changing the original file disrupts the developer workflow.

So, the production folks create another OpenCompose file which will extend _wordpress-opencompose.yml_ which will look like -

_wordpress.extension.yml_

```yaml
version: '0.1-dev'
type: extension

services:
- name: database
  env:
  - name: MYSQL_ROOT_PASSWORD
    secretRef: dbcreds/rootpassword
  mounts:
  - volumeRef: database
    mountPath: /var/lib/mysql

secrets:
- name: dbcreds
  data:
  - key: rootpassword
    base64: dmVyeXN0cm9uZ3Bhc3N3b3Jk

volumes:
- name: database
  size: 100Mi
  accessMode: ReadWriteOnce
```

With this file,
- the OpenCompose definitions which already exist in _wordpress-opencompose.yml_ are updated/overwritten, which in this case is the _database_ service, which gets a volume mount and a secret exposed.
- the OpenCompose definitions that do not exist in _wordpress-opencompose.yml_, are created, which in this case are the root level secrets and root level volumes.
- the OpenCompose definitions can be deleted using the immutable keys for different fields -

| field                        | immutable keys |
|--------------------------|-----------------------|
| services                 | name                 |
| containers              | name                |
| env                         | name               |
| ports                      | servicePort       |
| mounts                  | mountPath        |
| emptyDirVolumes  | name                |
| volumes                 | name                |

Support, we want to delete the environment variable "foo", then out extension file with look like -

```yaml
version: '0.1-dev'
type: extension

services:
- name: database
  env:
  - name: foo
    $operation: delete
```

Now the command, `opencompose -f wordpress-opencompose.yml,wordpress.extension.yml convert` is run, the magic happens and wordpress runs happily in the wild.

---

The extension file -
- does not have to be a legal OpenCompose file. It can be a complete OpenCompose file, in which case nothing will be overridden and only appended to the original file.
- specifies a root level `type: extension` field which makes it an extension file, which means this file is not expected to have legal OpenCompose syntax and mandatory fields can be missing. The validation will be carried out once the file is merged with the other files.
- _everything_ mentioned in the extension file takes precedence over the file it is extending, whenever a conflicting field appears.
- is passed normally like any other OpenCompose file besides the normal file with the `-f` directive, e.g. `opencompose -f original.yml,extension.yml convert`
