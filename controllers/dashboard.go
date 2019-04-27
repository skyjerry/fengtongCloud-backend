package controllers

import (
	"wac/k8s"
	"wac/wacUtils"
)

type DashboardController struct {
	BaseController
}

func (c *DashboardController) GetDashboardInfo() {
	//data := k8s.GetClusterData()
	pods := k8s.GetPods("", wacUtils.Clientset)
	var cpu int64
	var memory int64
	for _, item := range pods {
		containers := item.Spec.Containers
		for _, container := range containers {
			cpu += container.Resources.Limits.Cpu().Value()
			memory += container.Resources.Limits.Memory().Value()
		}
	}
	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"data":   pods,
		"cpu":    cpu,
		"memory": memory,
	})
}
