Guestbook Example
-----------------

### Example Reference:
https://github.com/kubernetes/kubernetes/tree/master/examples/guestbook

```
$ opencompose convert -f guestbook.yml
created file "redis-master-service.yaml"
created file "redis-master-deployment.yaml"
created file "redis-slave-service.yaml"
created file "redis-slave-deployment.yaml"
created file "frontend-service.yaml"
created file "frontend-deployment.yaml"
```

### Deploy the Guestbook on kubernetes

```
$ kubectl create -f .
deployment "frontend" created
service "frontend" created
deployment "redis-master" created
service "redis-master" created
deployment "redis-slave" created
service "redis-slave" created
```

```
$ kubectl get deployments
NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
frontend       1         1         1            1           1m
redis-master   1         1         1            1           1m
redis-slave    1         1         1            1           1m
```

```
$ kubectl get services
NAME           CLUSTER-IP   EXTERNAL-IP   PORT(S)        AGE
frontend       10.0.0.51    <pending>     80:30482/TCP   5m
kubernetes     10.0.0.1     <none>        443/TCP        5m
redis-master   10.0.0.104   <none>        6379/TCP       5m
redis-slave    10.0.0.66    <none>        6379/TCP       5m
```

Once it's exposed to external IP, visit the IP at `http://<EXTERNAL-IP:<PORT>`, you should see a webpage with guestbook.