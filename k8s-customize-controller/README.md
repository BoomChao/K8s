### 开发自定义的 controller 的整个过程   

前置:需要提交去[这里下](https://github.com/kubernetes/code-generator/tree/v0.27.0 )载好对应的code-generator的工具
使用的go语言是1.18就去下载对应的go.mod也是1.18的版本

执行脚本自动生成文件:./generate-groups.sh all \\   
k8s_customize_controller/pkg/client \\   
k8s_customize_controller/pkg/apis \\   
bolingcavalry:v1

1.创建 CRD(Custom Resources Definition)，令 k8s 明白我们自定的 API 对象

2.编写代码，将 CRD 的情况写入对应的代码中,然后通过代码生成工具生成除controller之外的informer,client等内容   
较为固定的代码

3.编写controller，在里面判断实际情况是否达到了API对象的声明情况，如果未达到，就要进行实际业务处理,   
而这也是controller的通用做法   

4.go build 生成可执行文件，注意这里用了 controller.go 中的文件，所以不能直接 go run main.go   
./k8s_customize_controller --alsologtostderr
   
5.这里日志打印选用的是 glog,而不是 golang 自带的 log
 
glog 是 Google 开源日志库 C++ glog 的Go语言简洁版,
