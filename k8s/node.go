package k8s

import (
	"github.com/astaxie/beego/logs"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNodeFromCluster(clientSet *kubernetes.Clientset) ClusterStatus {
	nodes := GetNodes(clientSet, "")
	clusterStatus := ClusterStatus{}
	for _, item := range nodes {
		if clusterStatus.MemSize == 0 && clusterStatus.CpuNum == 0 {
			clusterStatus.CpuNum = item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = 1
		} else {
			clusterStatus.CpuNum = clusterStatus.CpuNum + item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = clusterStatus.MemSize + item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = clusterStatus.Nodes + 1
		}
		clusterStatus.OsVersion = item.Status.NodeInfo.OSImage
	}
	clusterStatus.PodNum = GetPodsNumber("", clientSet)
	clusterStatus.Services = GetServiceNumber(clientSet, "")
	clusterStatus.MemSize = clusterStatus.MemSize / 1024 / 1024 / 1024

	return clusterStatus
}

// 获取nodes
func GetNodes(clientset *kubernetes.Clientset, labels string) []v1.Node {
	opt := metav1.ListOptions{}
	if labels != "" {
		opt.LabelSelector = labels
	}
	nodes, err := clientset.CoreV1().Nodes().List(opt)
	if err != nil {
		logs.Error("获取Nodes错误", err.Error())
		return make([]v1.Node, 0)
	}
	return nodes.Items
}
