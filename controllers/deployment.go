package controllers

import (
	"encoding/json"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wac/wacUtils"
)

type DeployController struct {
	BaseController
}

func (c *DeployController) GetDeployments() {
	deployments, _ := wacUtils.Clientset.AppsV1().Deployments("default").List(metav1.ListOptions{})

	returnData := make([]map[string]interface{}, len(deployments.Items))

	for k, value := range deployments.Items {
		returnData[k] = map[string]interface{}{
			"name":         value.Name,
			"replicas":     value.Spec.Replicas,
			"runningCount": value.Status.ReadyReplicas,
			"labels":       value.ObjectMeta.Labels,
			"created_at":   value.GetCreationTimestamp().Format("2006-01-02 15:04:05"),
			"containers":   value.Spec.Template.Spec.Containers,
			"value":        value,
		}
	}

	c.ApiResponse(200, "获取成功", map[string]interface{}{
		"deployments": returnData,
	})
}

func (c *DeployController) UpdateDeployment() {
	var deploy *v1.Deployment

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &deploy)
	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
	}
	_, err = wacUtils.Clientset.AppsV1().Deployments("default").Update(deploy)
	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
	}

	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}

func (c *DeployController) ScaleDeployment() {
	deploymentName := c.Ctx.Input.Param(":deploymentName")
	var PostParams *struct {
		Replicas *int32 `json:"replicas,omitempty"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &PostParams)
	//println(c.Ctx.Input.RequestBody)

	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
	}

	var deployment *v1.Deployment
	deployment, err = wacUtils.Clientset.AppsV1().Deployments("default").Get(deploymentName, metav1.GetOptions{})

	println(deployment.Spec.Replicas)
	deployment.Spec.Replicas = PostParams.Replicas
	println(PostParams.Replicas)

	_, err = wacUtils.Clientset.AppsV1().Deployments("default").Update(deployment)
	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
	}

	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}
