package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"wac/wacUtils"
)

// 获取资源使用情况,cpu，内存
func GetClusterUsed(clientset *kubernetes.Clientset) ClusterResources {
	clusterResouces := ClusterResources{}
	resources := GetPods("", clientset)
	var cpu int64
	var memory int64
	for _, item := range resources {
		containers := item.Spec.Containers
		for _, container := range containers {
			cpu += container.Resources.Limits.Cpu().Value()
			memory += container.Resources.Limits.Memory().Value()
		}
	}
	clusterResouces.Services = GetServiceNumber(clientset, "")
	clusterResouces.UsedCpu = cpu
	clusterResouces.UsedMem = memory / 1024 / 1024 / 1024
	return clusterResouces
}

// 获取集群数据,首页使用
func GetClusterData() CloudClusterDetail {
	details := CloudClusterDetail{}
	GoGetCluseterDetail(&details)

	if details.ClusterCpu > 0 && details.ClusterMem > 0 {
		floatCpu := (float64(details.UsedCpu) / float64(details.ClusterCpu)) * 100
		floatMem := (float64(details.UsedMem) / float64(details.ClusterMem)) * 100
		cp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatCpu), 64)
		mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatMem), 64)
		details.CpuUsePercent = cp
		details.MemUsePercent = mp
		details.MemFree = details.ClusterMem - details.UsedMem
		details.CpuFree = details.ClusterCpu - details.UsedCpu
	}
	return details
}

// 通过多线程添加和获取数据
func GoGetCluseterDetail(details *CloudClusterDetail) {
	detail := GetClusterDetailData()
	details.ClusterCpu = details.ClusterCpu + detail.ClusterCpu
	details.ClusterMem = details.ClusterMem + detail.ClusterMem
	details.ClusterNode = details.ClusterNode + detail.ClusterNode
	details.ClusterPods = details.ClusterPods + detail.ClusterPods
	details.Services = details.Services + detail.Services
	details.UsedMem = details.UsedMem + detail.UsedMem
	details.UsedCpu = details.UsedCpu + detail.UsedCpu
	details.Couters = details.Couters + 1
}

// 缓存集群详情到redis里
// 获取集群资源使用详细情况
func GetClusterDetailData() CloudClusterDetail {
	data := CloudCluster{}
	detail := CloudClusterDetail{}
	c := wacUtils.Clientset
	detail.ClusterId = data.ClusterId
	clusterStatus := GetNodeFromCluster(c)
	detail.ClusterCpu = clusterStatus.CpuNum
	detail.ClusterMem = clusterStatus.MemSize
	detail.ClusterNode = clusterStatus.Nodes
	detail.ClusterPods = clusterStatus.PodNum
	detail.OsVersion = clusterStatus.OsVersion
	if detail.ClusterCpu > 0 && detail.ClusterMem > 0 {
		used := GetClusterUsed(c)
		detail.UsedMem = used.UsedMem
		detail.UsedCpu = used.UsedCpu
		floatCpu := (float64(detail.UsedCpu) / float64(detail.ClusterCpu)) * 100
		floatMem := (float64(detail.UsedMem) / float64(detail.ClusterMem)) * 100
		cp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatCpu), 64)
		mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatMem), 64)
		detail.CpuUsePercent = cp
		detail.MemUsePercent = mp
		detail.MemFree = detail.ClusterMem - detail.UsedMem
		detail.CpuFree = detail.ClusterCpu - detail.UsedCpu
		detail.Services = used.Services
	}
	return detail

}
