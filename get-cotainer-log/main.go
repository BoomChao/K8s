package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericrest "k8s.io/apiserver/pkg/registry/generic/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
	"k8s.io/client-go/util/homedir"
)

var (
	namespace      = "default"
	serviceAccount = "my-service-account"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 获取当前服务账号访问API-Servier的Token
	var expired int64 = 600
	tr, err := clientset.CoreV1().ServiceAccounts(namespace).CreateToken(ctx, serviceAccount, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences:         []string{},
			ExpirationSeconds: &expired,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	token := tr.Status.Token

	fmt.Println(token)

	location := &url.URL{
		Scheme: "https",
		Host:   net.JoinHostPort("127.0.0.1", "10250"),                                  // 访问的是node的10250端口
		Path:   fmt.Sprintf("/containerLogs/%s/%s/%s", "default", "mynginx", "mynginx"), // 获取default下的pod名为 "mynginx"中的"mynginx"容器的日志
		// RawQuery: params.Encode(),
	}

	streamer := &genericrest.LocationStreamer{
		Location: location,
		Transport: transport.NewBearerAuthRoundTripper(token, &http.Transport{

			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}),
		ContentType: "text/plain",
		Flush:       false,
		// ResponseChecker:             genericrest.NewGenericHttpResponseChecker(api.Resource("pods/log"), name),
		RedirectChecker: genericrest.PreventRedirects,
		// TLSVerificationErrorCounter: podLogsTLSFailure,
	}
	reader, _, _, err := streamer.InputStream(ctx, "not used" /*logs version is not used */, "")
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
