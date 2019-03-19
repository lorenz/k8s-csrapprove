package main

import (
	"crypto/x509"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"k8s.io/klog"

	"git.dolansoft.org/dolansoft/k8s-csrapprove/approver"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	klog.InitFlags(nil)
	klog.SetOutput(os.Stdout)

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	validationCAPool := x509.NewCertPool()
	if len(config.CAData) > 0 {
		if !validationCAPool.AppendCertsFromPEM(config.CAData) {
			klog.Errorf("Failed to load K8s root cert for validation")
		}
	}
	if config.CAFile != "" {
		caData, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			klog.Errorf("Failed to load K8s root cert for validation: %v", err)
		}
		if !validationCAPool.AppendCertsFromPEM(caData) {
			klog.Errorf("Failed to load K8s root cert for validation")
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	csrInformer := informerFactory.Certificates().V1beta1().CertificateSigningRequests()
	go csrInformer.Informer().Run(stopper)
	nodeInformer := informerFactory.Core().V1().Nodes()
	go nodeInformer.Informer().Run(stopper)
	csrApprover := approver.NewCSRApprovingController(clientset, csrInformer, nodeInformer, validationCAPool)
	csrApprover.Run(1, stopper)
}
