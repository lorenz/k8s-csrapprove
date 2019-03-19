package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/klog"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/controller/certificates/approver"
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

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	csrInformer := informerFactory.Certificates().V1beta1().CertificateSigningRequests()
	go csrInformer.Informer().Run(make(<-chan struct{}))
	csrApprover := approver.NewCSRApprovingController(clientset, csrInformer)
	csrApprover.Run(1, make(<-chan struct{}))
}
