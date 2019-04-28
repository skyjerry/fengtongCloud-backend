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
			"name":       value.Name,
			"nodeName":   value.Spec.NodeName,
			"status":     value.Status.Phase,
			"image":      value.Spec.Containers[0].Image,
			"created_at": value.GetCreationTimestamp().Format("2006-01-02 15:04:05"),
			"value":      value,
		}
	}

	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"pods": returnData,
	})
}

func (c *PodController) DeletePod() {
	podName := c.Ctx.Input.Param(":podName")

	err := wacUtils.Clientset.CoreV1().Pods("default").Delete(podName, &metav1.DeleteOptions{})

	c.ApiResponse(200, "删除成功", map[string]interface{}{
		"err": err,
	})
}
