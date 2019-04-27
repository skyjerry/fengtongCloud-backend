package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wac/wacUtils"
)

type NodeController struct {
	BaseController
}

func (c *NodeController) GetNodes() {
	nodes, _ := wacUtils.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})

	returnData := make([]map[string]interface{}, len(nodes.Items))

	for k, value := range nodes.Items {
		returnData[k] = map[string]interface{}{
			"name":          value.Name,
			"role":          value.Labels["node-role.kubernetes.io/worker"],
			"status":        value.Status.Conditions,
			"unschedulable": value.Spec.Unschedulable,
			"created_at":    value.GetCreationTimestamp().Format("2006-01-02 15:04:05"),
			"addresses":     value.Status.Addresses,
			"value":         value,
		}
	}

	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"nodes": returnData,
	})
}

func (c *NodeController) GetNode() {
	nodeName := c.Ctx.Input.Param(":nodeName")
	if len(nodeName) == 0 {
		c.ApiResponse(404, "获取失败", map[string]interface{}{})
	}
	node, _ := wacUtils.Clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})

	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"nodeInfo": node,
	})
}

func (c *NodeController) StartNode() {
	nodeName := c.Ctx.Input.Param(":nodeName")
	if len(nodeName) == 0 {
		c.ServeJSON()
	}
	node, _ := wacUtils.Clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})

	node.Spec.Unschedulable = false
	_, err := wacUtils.Clientset.CoreV1().Nodes().Update(node)
	if err != nil {
		c.ApiResponse(500, "设置失败", map[string]interface{}{
			"errInfo": err,
		})
	}

	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}

func (c *NodeController) StopNode() {
	nodeName := c.Ctx.Input.Param(":nodeName")
	if len(nodeName) == 0 {
		c.ServeJSON()
	}
	node, _ := wacUtils.Clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})

	node.Spec.Unschedulable = true
	_, err := wacUtils.Clientset.CoreV1().Nodes().Update(node)
	if err != nil {
		c.ApiResponse(500, "设置失败", map[string]interface{}{
			"errInfo": err,
		})
	}
	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}
