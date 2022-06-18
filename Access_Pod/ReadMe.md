#### 外部访问k8s中的pod服务端口
方式一：使用 hostPort 方式(就是将容器端口做一个端口映射),hostpath.yml 就是使用这种方式    
这样做有个缺点，因为 Pod 重新调度的时候该Pod被调度到的宿主机可能会变动，这样就变化了，用户必须自己维护一个 Pod 与所在宿主机的对应关系   
<br>

方式二：使用 NodePort 的方式(也就是使用Service)
nodePort 在 kubenretes 里是一个广泛应用的服务暴露方式     
Kubernetes 中的 service 默认情况下都是使用的 ClusterIP 这种类型，这样的 service 会产生一个 ClusterIP    
这个IP只能在集群内部访问，要想让外部能够直接访问 service，需要将 service type 修改为 NodePort     
httpdpod.yml 和 rclserver.yml 就是使用这种方式