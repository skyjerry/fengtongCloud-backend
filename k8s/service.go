package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 获取某个集群的服务信息
func GetServices(clientset *kubernetes.Clientset, namespace string) ([]v1.Service, error) {
	opt := metav1.ListOptions{}
	data, err := clientset.CoreV1().Services(namespace).List(opt)
	if err != nil {
		logs.Error("获取service失败啦", err)
		return make([]v1.Service, 0), err
	}
	json.Marshal(data.Items)
	return data.Items, err
}

func GetServiceNumber(clientset *kubernetes.Clientset, namespace string) int {
	data, _ := GetServices(clientset, namespace)
	return len(data)
}
