package controllers

import (
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// NodeAllocatedResources describes node allocated resources.

// Operations about Users
type MainController struct {
	BaseController
}

var (
	clientset *kubernetes.Clientset
)

func (m *MainController) GetServices() {
	services, err := clientset.CoreV1().Services("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	returnData := make(map[string]interface{})

	for _, value := range services.Items {
		returnData[value.Name] = map[string]interface{}{
			"Value": value,
			//"Status": value.Status,
		}
	}
	m.Data["json"] = returnData
	m.ServeJSON()

}

func (m *MainController) ScaleService() {
	//serviceId := m.GetString("id")
	var kubeconfig string
	home := homeDir()
	kubeconfig = filepath.Join(home, ".kube", "config")

	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		m.Data["json"] = map[string]int{"error": 500}
		m.ServeJSON()
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	result, getErr := deploymentsClient.Get("kubernetes-bootcamp", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	result.Spec.Replicas = int32Ptr(1)
	_, updateErr := deploymentsClient.Update(result)
	//returnData := make(map[string]interface{})
	//
	//for _, value := range list.Items {
	//	returnData[value.Name] = map[string]interface{}{
	//		"Name":   value.Name,
	//		"Status": value.Status,
	//	}
	//}
	res := "失败"
	if updateErr == nil {
		res = "成功"
	}
	m.Data["json"] = res

	m.ServeJSON()
}

func (m *MainController) GetDeployments() {
	deployments, _ := clientset.AppsV1().Deployments("default").List(metav1.ListOptions{})

	returnData := make(map[string]interface{})
	for _, value := range deployments.Items {
		returnData[value.Name] = map[string]interface{}{
			"Value": value,
			//"Status": value.Status,
		}
	}
	m.Data["json"] = returnData
	m.ServeJSON()
}

func (m *MainController) StopNode() {
	nodeClient := clientset.CoreV1().Nodes()

	node, _ := nodeClient.Get("huawei1", metav1.GetOptions{})

	node.Spec.Unschedulable = true
	//node.Spec.Taints = nil
	_, updateErr := nodeClient.Update(node)
	result := "success"
	if updateErr != nil {
		result = "fail"
	}
	m.Data["json"] = result
	m.ServeJSON()
}

func (m *MainController) Ping() {
	m.ApiResponse(200, "success", map[string]interface{}{})
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func int32Ptr(i int32) *int32 { return &i }
