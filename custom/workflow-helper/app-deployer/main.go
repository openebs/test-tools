package main

import (
	"bytes"
	"flag"
	"fmt"
	logs "log"
	"os/exec"
	"strings"
	"time"

	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/openebs/maya/pkg/util/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//AppVars maintaining all parameters require during application installation/deletion
type AppVars struct {
	namespace string
	filePath  string
	timeout   int
	operation string
	label     string
	app       string
	scope     string
}

func main() {

	//GetData is initializing required variables for app-deployer
	appVars, err := GetData()
	if err != nil {
		logs.Fatalf("err: %v", err)
	}

	config, err := getKubeConfig()
	if err != nil {
		logs.Fatalf("err: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logs.Fatalf("err: %v", err)
	}

	log.Info("[Start]: Starting App Deployer...")

	//operations for application
	//if false operation exist then default case handles it
	switch appVars.operation {
	case "apply", "create":
		if err := CreateApplication(appVars, 2, clientset); err != nil {
			logs.Fatalf("err: %v", err)
		}
		log.Infof("[Status]: %s applications has been successfully created", appVars.app)
	case "delete":
		if err := DeleteApplication(appVars, 2, clientset); err != nil {
			logs.Fatalf("err: %v", err)
		}
		log.Infof("[Status]: %s applications has been successfully deleted", appVars.app)
	default:
		logs.Fatalf("Operation '%s' not supported in app-deployer", appVars.operation)
	}
}

// GetKubeConfig function derive the kubeconfig
func getKubeConfig() (*rest.Config, error) {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	// It uses in-cluster config if kubeconfig path is not specified
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	return config, err
}

//GetData derive the application filePath and timeout
//it derive the filePath based on application scenario(week vs resilient)
func GetData() (*AppVars, error) {

	//Initialise the variables
	namespace := flag.String("namespace", "", "namespace for the application")
	filePath := flag.String("typeName", "weak", "type of the application")
	timeout := flag.Int("timeout", 300, "timeout for application status")
	operation := flag.String("operation", "apply", "type of operation for application")
	app := flag.String("app", "", "type of app for application")
	scope := flag.String("scope", "cluster", "scope of the application")
	flag.Parse()

	appVars := AppVars{
		namespace: *namespace,
		timeout:   *timeout,
		operation: *operation,
		label:     "app=" + *app,
		app:       *app,
		scope:     *scope,
	}
	//application namespace having weak and resilient filePath
	//loadtest namespace having loadtest filePath
	//sock-shop namespace having sock-shop filePath
	//podtato-head namespace having podtato-head filePath
	switch appVars.label {
	case "app=loadtest":
		appVars.filePath = "loadtest.yaml"
	case "app=sock-shop":
		appVars.filePath = *filePath + "-sock-shop.yaml"
	case "app=podtato-head":
		appVars.filePath = *filePath + "-podtato-head.yaml"
	case "app=bank-of-anthos":
		appVars.filePath = *filePath + "-bank-of-anthos.yaml"
	default:
		return &appVars, fmt.Errorf("app '%v' not supported in app-deployer", appVars.app)
	}

	return &appVars, nil
}

// CreateNamespace creates a namespace
func CreateNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	_, err := clientset.CoreV1().Namespaces().Create(nsSpec)
	return err
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	return clientset.CoreV1().Namespaces().Delete(namespaceName, &metav1.DeleteOptions{})
}

//CreateApp create the application
func CreateApp(path, ns, operation string) error {
	command := exec.Command("kubectl", operation, "-f", path, "-n", ns)
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

//DeleteApp delete the application
func DeleteApp(path, ns string) error {
	command := exec.Command("kubectl", "delete", "-f", path, "-n", ns)
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

//CreateApplication creates the application and add all corresponding resources
func CreateApplication(appVars *AppVars, delay int, clientset *kubernetes.Clientset) error {

	log.Infof("[Status]: FilePath for App Deployer is %v", appVars.filePath)
	switch strings.ToLower(appVars.scope) {
	case "cluster":
		if err := CreateNamespace(clientset, appVars.namespace); err != nil {
			if !k8serrors.IsAlreadyExists(err) {
				return err
			}
			log.Info("[Status]: Namespace already exist")
		} else {
			log.Info("[Status]: Namespace created successfully")
		}
	case "namespace":
		log.Infof("[Status]: Application is using %v namespace", appVars.namespace)
	default:
		return fmt.Errorf("Scope '%v' not supported in app-deployer", appVars.scope)
	}

	if err := CreateApp("/var/run/"+appVars.filePath, appVars.namespace, appVars.operation); err != nil {
		return fmt.Errorf("Failed to install %s", appVars.namespace)

	}
	if err := CheckApplicationStatus(appVars.namespace, appVars.label, appVars.timeout, 2, clientset); err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

//DeleteApplication deletes the application and remove all corresponding resources
func DeleteApplication(appVars *AppVars, delay int, clientset *kubernetes.Clientset) error {

	log.Infof("[Status]: FilePath for App Deployer is %v", appVars.filePath)
	log.Info("[Status]: Revert application has been started")
	if err := DeleteApp("/var/run/"+appVars.filePath, appVars.namespace); err != nil {
		return err
	}

	if err := CheckPodStatusForRevert(appVars.namespace, appVars.label, appVars.timeout, 2, clientset); err != nil {
		return err
	}
	return nil
}

// CheckApplicationStatus checks the status of the AUT
func CheckApplicationStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	// Checking whether application containers are in ready state
	log.Info("[Status]: Checking whether application containers are in ready state")
	err := CheckContainerStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	// Checking whether application pods are in running state
	log.Info("[Status]: Checking whether application pods are in running state")
	err = CheckPodStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	return nil
}

//CheckPodStatusForRevert wait for the application to terminate all pods
func CheckPodStatusForRevert(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	return retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{LabelSelector: appLabel})
			if err != nil {
				return errors.Errorf("Unable to find the pods in namespace, err: %v", err)
			}

			if len(podSpec.Items) != 0 {
				return errors.Errorf("[Status]: Pods are yet to be terminated")
			}
			return nil
		})
}

// CheckPodStatus checks the running status of the application pod
func CheckPodStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	return retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{LabelSelector: appLabel})
			if err != nil || len(podSpec.Items) == 0 {
				return errors.Errorf("Unable to find the pods with matching labels, err: %v", err)
			}
			for _, pod := range podSpec.Items {
				if string(pod.Status.Phase) != "Running" {
					return errors.Errorf("Pod is not yet in running state")
				}
				log.InfoWithValues("[Status]: The running status of Pods are as follows", logrus.Fields{
					"Pod": pod.Name, "Status": pod.Status.Phase})
			}
			return nil
		})
}

// CheckContainerStatus checks the status of the application container
func CheckContainerStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	return retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{LabelSelector: appLabel})
			if err != nil || len(podSpec.Items) == 0 {
				return errors.Errorf("Unable to find the pods with matching labels, err: %v", err)
			}
			for _, pod := range podSpec.Items {
				for _, container := range pod.Status.ContainerStatuses {
					if container.State.Terminated != nil {
						return errors.Errorf("container is in terminated state")
					}
					if container.Ready != true {
						return errors.Errorf("containers are not yet in ready state")
					}
					log.InfoWithValues("[Status]: The running status of container are as follows", logrus.Fields{
						"container": container.Name, "Pod": pod.Name, "Readiness": container.Ready})
				}
			}
			return nil
		})
}
