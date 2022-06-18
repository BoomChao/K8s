1.先执行base目录下的yml文件，kubectl apply -k ./base,  这就相当于部署服务
访问service映射的端口发现字体为红色
<br >  

2.kustomize ./overlays/production 就相当于 
kubectl appyly -k ./overlays/production 这会通过kubectl使用kustomize并检查
指定目录下的 kustomization.yml, 这就相当于管理服务，再次访问service映射的端口发现字体为蓝色
<br >
   
3.kustomize可以理解为是打补丁，对已经部署的应用做一些小的修改，这个修改当然也可以通过
直接修改原来的yml源文件来进行，但是当服务很多时使用直接yml管理就很繁琐，这时kustomize的作用就凸显出来了
