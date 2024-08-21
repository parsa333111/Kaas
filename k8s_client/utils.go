package k8s_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 { return &i }

func extarctDeploymentStatus(clientSet *kubernetes.Clientset, deployment *v1.Deployment) (DeploymentStatus, bool) {
	labelSelector := fields.SelectorFromSet(deployment.Spec.Selector.MatchLabels).String()

	pods, err := clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Println(err.Error())
		return DeploymentStatus{}, false
	}

	var podStatuses []PodStatus
	for _, pod := range pods.Items {
		podStatuses = append(podStatuses, PodStatus{
			Name:      pod.Name,
			Phase:     string(pod.Status.Phase),
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
			StartTime: pod.Status.StartTime.Format(time.RFC3339),
		})
	}

	deploymentStatus := DeploymentStatus{
		Name:          deployment.Name,
		Replicas:      *deployment.Spec.Replicas,
		ReadyReplicas: deployment.Status.ReadyReplicas,
		PodStatuses:   podStatuses,
	}

	return deploymentStatus, true
}

func assertResourceList(deploymentResources []DeploymentResource) (corev1.ResourceList, bool) {
	resourceList := make(corev1.ResourceList)

	for _, deploymentResource := range deploymentResources {
		resourceName := corev1.ResourceName(strings.ToLower(deploymentResource.Name))
		quantity, err := resource.ParseQuantity(deploymentResource.Limit)
		if err != nil {
			log.Println(err)
			return corev1.ResourceList{}, false
		}
		resourceList[resourceName] = quantity
	}

	return resourceList, true
}

func createConfig(clientSet *kubernetes.Clientset, deploymentName string, deploymentConfigs []DeploymentConfig) bool {
	data := make(map[string]string)

	for _, deploymentConfig := range deploymentConfigs {
		data[deploymentConfig.Key] = deploymentConfig.Value
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-config", deploymentName),
		},
		Data: data,
	}

	_, err := clientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func createSecret(clientSet *kubernetes.Clientset, deploymentName string, deploymentSecrets []DeploymentSecret) bool {
	data := make(map[string][]byte)

	for _, deploymentSecret := range deploymentSecrets {
		data[deploymentSecret.Key] = []byte(base64.StdEncoding.EncodeToString([]byte(deploymentSecret.Value)))
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-secret", deploymentName),
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}

	_, err := clientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func createService(clientSet *kubernetes.Clientset, deploymentName string, servicePort int32, monitor bool, externalAccess bool) bool {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":     deploymentName,
				"monitor": strconv.FormatBool(monitor),
			},
			Ports: []corev1.ServicePort{
				{
					Port:       servicePort,
					TargetPort: intstr.FromInt(int(servicePort)),
				},
			},
		},
	}

	if externalAccess {
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
	}

	_, err := clientSet.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func createIngress(clientSet *kubernetes.Clientset, deploymentName string, servicePort int32, domainName string) bool {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-ingress", deploymentName),
			Namespace: namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: domainName,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: deploymentName,
											Port: networkingv1.ServiceBackendPort{
												Number: servicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientSet.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
