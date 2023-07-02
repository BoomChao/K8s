### 学习k8s的一些小例子

1. Creaing--Update--Delete--Deployment 是 client-go 操作 Deployment

2. Upgrade-Container 是通过命令行传递一个Deployment的名字，容器名字和新的image名字来升级 Deployment
```
如： ./main -image=nginx:1.13 -app=nginx -deployment=demo-deployment
app就表示容器名字
```

3. Kustomize 是学习 kustomize 的一个小案例

4. CreateDeployAndExposeService 是使用 client-go 创建一个deployment 并将其暴露成一个 service 

5. Access_Pod 就是访问 Pod 的两种方式，分别使用 hostPort 和 NodePort

6. Volume 是 k8s 的容器卷的使用

7. k8s_customize_controller 是自定义一个 Controller，只是实现了Controller部分，其他部分如 informer 是用   
代码生成工具直接生成的

8. helm-chart 是自己开发出的一个Ghost的chart