1. 先执行base目录下的yml文件，kubectl apply -k ./base,  这就相当于部署服务
访问service映射的端口发现字体为红色   

2. kustomize ./overlays/production 就相当于   
kubectl appyly -k ./overlays/production 这会通过 kubectl 使用 kustomize 并检查指定目录下的 kustomization.yml, 这就相当于管理服务，再次访问service映射的端口发现字体为蓝色   
   
3. kustomize 可以理解为是打补丁，对已经部署的应用做一些小的修改，这个修改当然也可以通过直接修改原来的yml源文件来进行，但是当服务很多时使用直接 yml 管理就很繁琐，这时 kustomize 的作用就凸显出来了
