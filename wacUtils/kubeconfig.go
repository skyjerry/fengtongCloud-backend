package wacUtils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var KubeConfig, _ = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube-online", "config"))
var Clientset, _ = kubernetes.NewForConfig(KubeConfig)
