### 学习k8s的一些小例子

1. Creaing--Update--Delete--Deployment 是 client-go 操作 Deployment

2. Upgrade-Container 是通过命令行传递一个Deployment的名字，容器名字和新的image名字来升级 Deployment

```
如： ./main -image=nginx:1.13 -app=nginx -deployment=demo-deployment
app就表示容器名字
```
