# Volumes for opencompose

### Example

```yaml
services:
- name: db
  containers:
  - image: mysql
    mounts:
    - volumeName: db
      mountPath: /var/lib/mysql

- name: backup
  containers:
  - image: backup
    mounts:
    - volumeName: db
      mountPath: /app/store
      volumeSubPath: foo/bar
      readOnly: true

- name: process
  containers:
  - image: process
    mounts:
    - volumeName: temp
      mountPath: /app/data  
  emptyDirVolumes:
  - name: temp

volumes:
- name: db
  size: 5GiB
  accessMode: ReadWriteMany
  storageClass: fast
```

### Application explanation

In above example, we have services `db`, `backup` and `process`.

- `db` is database which needs storage and is writable.
- `backup` is a service which periodically sends backup somewhere of database, so it needs access to the data of database but only as readonly.
- `process` is a more of a job which is dataprocessing application, which does not need persistent storage but needs storage until it does it's job.


### Root level `volumes` directive:

We generate `pvc` for each entry in the volumes directive.

- `name`:          (required field, type: *string*) name of volume, which is referenced from containers for mounting. `pvc` will be created with this name.
- `size`:          (required field, type: *string*) In `pvc` this will be put in `pvc.spec.resources.requests`
- `accessMode`:    (required field, type: *string*) In `pvc` this will turn into `pvc.spec.accessModes`. The possible values of this field could be: `ReadOnlyMany`, `ReadWriteOnce` or `ReadWriteMany`.
- `storageClass`:  (optional field) In `pvc` this will be into `pvc.metadata.annotations` as
  ```json
  "annotations": {
    "volume.beta.kubernetes.io/storage-class": "user-value"
  }
  ```
  src: http://blog.kubernetes.io/2016/10/dynamic-provisioning-and-storage-in-kubernetes.html

### Service level `emptyDirVolumes` directive:

- `name`: (required field, type: *string*) name of `emptyDir` to be created. Containers in same service/pod can share one or more `emptyDir`. This will be mapped in `deployment` to `deployment.spec.template.spec.volumes.emptyDir`.

### Container level `mounts` directive:

- `volumeName`: (required field, type: *string*) name of volume from root level `volumes` directive or `service` level `emptyDirVolumes` directive.
- `mountPath`:  (required field, type: *string*) This will be mapped to `pod.spec.containers.volumeMounts.mountPath` in `deployment`.
- `volumeSubPath`: (optional field, type: *string*) This will be mapped to `pod.spec.containers.volumeMounts.subPath` in `deployment`.
- `readOnly`: (optional field, type: *bool*, default: false) This will be mapped to `pod.spec.containers.volumeMounts.readOnly` in `deployment`.
