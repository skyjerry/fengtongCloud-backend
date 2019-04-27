package k8s

import (
	"github.com/astaxie/beego/logs"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(namespace string, clientSet *kubernetes.Clientset) []v1.Pod {
	opt := metav1.ListOptions{}
	pods, err := clientSet.CoreV1().Pods(namespace).List(opt)
	if err != nil {
		logs.Error("获取Pods错误", err.Error())
		return make([]v1.Pod, 0)
	}
	return pods.Items
}

// 获取pods数量
func GetPodsNumber(namespace string, clientSet *kubernetes.Clientset) int {
	opt := metav1.ListOptions{}

	pods, err := clientSet.CoreV1().Pods(namespace).List(opt)
	if err != nil {
		logs.Error("获取k8s Pods失败", err.Error())
		return 0
	}
	return len(pods.Items)
}
