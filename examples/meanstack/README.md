Meanstack Example
-----------------

### Example Reference:
https://kubernetes.io/docs/getting-started-guides/meanstack/

##### Note: This example is tested on minikube, For creation of PV, refer http://suraj.io/post/quick-pv-for-local-k8s/

```
$ opencompose convert -f meanstack.yml
created file "mongo-service.yaml"
created file "mongo-deployment.yaml"
created file "web-service.yaml"
created file "web-deployment.yaml"
created file "mongo-persistent-storage-persistentvolumeclaim.yaml"
```

### Deploy the example on kubernetes

```
$ kubectl create -f .
deployment "mongo" created
persistentvolumeclaim "mongo-persistent-storage" created
service "mongo" created
deployment "web" created
service "web" created
```

```
$ kubectl get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
mongo     1         1         1            1           5m
web       2         2         2            2           5m
```

```
$ kubectl get services
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)           AGE
kubernetes   10.0.0.1     <none>        443/TCP           20d
mongo        10.0.0.133   <pending>     27017:31916/TCP   58s
web          10.0.0.26    <pending>     80:32583/TCP      58s
```

Once it's exposed to external IP, visit the IP at `http://<MINIKUBE-IP>:<PORT>`, you should see a webapp.

Note: Persistent volume should be present in order to run this example.
