package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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
		return
	}
	_, err = wacUtils.Clientset.AppsV1().Deployments("default").Update(deploy)
	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
		return
	}

	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}

func (c *DeployController) ScaleDeployment() {
	deploymentName := c.Ctx.Input.Param(":deploymentName")
	var PostParams *struct {
		Replicas *int32 `json:"replicas,omitempty"`
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &PostParams)

	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
		return
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
		return
	}

	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}

func (c *DeployController) DeleteDeploy() {
	deploymentName := c.Ctx.Input.Param(":deploymentName")
	deletePolicy := metav1.DeletePropagationForeground
	err := wacUtils.Clientset.AppsV1().Deployments("default").Delete(deploymentName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
		return
	}
	c.ApiResponse(200, "操作成功", map[string]interface{}{})
}

func (c *DeployController) CreateDeploy() {
	var PostParams struct {
		AppName       string `json:"app_name,omitempty"`
		Replicas      *int32 `json:"replicas,omitempty"`
		Labels        string `json:"labels,omitempty"`
		Image         string `json:"image,omitempty"`
		ContainerName string `json:"container_name,omitempty"`
		ContainerPort int32  `json:"container_port,omitempty"`
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &PostParams)

	ImageInfo := strings.Split(PostParams.Image, ":")
	var imagesName struct {
		Repositories []string `json:"repositories,omitempty"`
	}
	httplib.Get("http://114.116.173.97:5000/v2/_catalog").ToJSON(&imagesName)

	validImage := false
	validTag := false
	for _, v := range imagesName.Repositories {
		if v == ImageInfo[0] {
			validImage = true
			break
		}
	}
	if !validImage {
		c.ApiResponse(403, "创建失败，请确认镜像是否存在", map[string]interface{}{})
		return
	}
	var tagsName struct {
		Name string   `json:"name,omitempty"`
		Tags []string `json:"tags,omitempty"`
	}
	httplib.Get("http://114.116.173.97:5000/v2/" + ImageInfo[0] + "/tags/list").ToJSON(&tagsName)
	for _, v := range tagsName.Tags {
		if v == ImageInfo[1] {
			validTag = true
			break
		}
	}
	if !validTag {
		c.ApiResponse(403, "创建失败，请确认tag是否存在", map[string]interface{}{})
		return
	}

	PostParams.Image = "192.168.13.207:5000/" + PostParams.Image
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: PostParams.AppName,
			Labels: map[string]string{
				"app": PostParams.Labels,
			},
		},
		Spec: v1.DeploymentSpec{
			Replicas: PostParams.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": PostParams.Labels,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": PostParams.Labels,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  PostParams.ContainerName,
							Image: PostParams.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: PostParams.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := wacUtils.Clientset.AppsV1().Deployments("default").Create(deployment)

	if err != nil {
		c.ApiResponse(403, "操作失败", map[string]interface{}{
			"err": err,
		})
	}

	c.ApiResponse(200, "创建成功", map[string]interface{}{
		"result": result,
	})

}
