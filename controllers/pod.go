package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wac/wacUtils"
)

type PodController struct {
	BaseController
}

func (c *PodController) GetPods() {
	pods, _ := wacUtils.Clientset.CoreV1().Pods("default").List(metav1.ListOptions{})

	returnData := make([]map[string]interface{}, len(pods.Items))

	for k, value := range pods.Items {
		returnData[k] = map[string]interface{}{
			"value": value,
		}
	}

	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"pods": returnData,
	})
}
