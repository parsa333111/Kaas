package k8s_client

import (
	"context"
	"fmt"
	"log"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateCustomDeployment(deploymentRequest CustomDeploymentRequest) bool {
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	ok := createConfig(clientSet, deploymentRequest.Name, deploymentRequest.Configs)
	if !ok {
		return false
	}

	ok = createSecret(clientSet, deploymentRequest.Name, deploymentRequest.Secrets)
	if !ok {
		return false
	}

	resourceList, ok := assertResourceList(deploymentRequest.Resources)
	if !ok {
		return false
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentRequest.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentRequest.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     deploymentRequest.Name,
					"monitor": strconv.FormatBool(deploymentRequest.Monitor),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     deploymentRequest.Name,
						"monitor": strconv.FormatBool(deploymentRequest.Monitor),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deploymentRequest.Name,
							Image: fmt.Sprintf("%s:%s", deploymentRequest.ImageAddress, deploymentRequest.ImageTag),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: deploymentRequest.ServicePort,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-config", deploymentRequest.Name),
										},
									},
								},
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-secret", deploymentRequest.Name),
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   resourceList,
								Requests: resourceList,
							},
						},
					},
				},
			},
		},
	}

	_, err = clientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	ok = createService(
		clientSet, deploymentRequest.Name, deploymentRequest.ServicePort, deploymentRequest.Monitor, false)
	if !ok {
		return false
	}

	if deploymentRequest.ExternalAccess {
		ok = createIngress(
			clientSet, deploymentRequest.Name, deploymentRequest.ServicePort, deploymentRequest.DomainAddress)
		if !ok {
			return false
		}
	}

	return true
}

func CreatePostgresDeployment(deploymentRequest PostgresDeploymentRequest) bool {
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	ok := createConfig(clientSet, deploymentRequest.Name,
		[]DeploymentConfig{
			{
				Key:   "POSTGRES_USER",
				Value: postgresUser,
			},
			{
				Key:   "POSTGRES_DB",
				Value: postgresDB,
			},
		},
	)
	if !ok {
		return false
	}

	ok = createSecret(clientSet, deploymentRequest.Name,
		[]DeploymentSecret{
			{
				Key:   "POSTGRES_PASSWORD",
				Value: postgresPassword,
			},
		},
	)
	if !ok {
		return false
	}

	resourceList, ok := assertResourceList(deploymentRequest.Resources)
	if !ok {
		return false
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentRequest.Name,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: deploymentRequest.Name,
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     deploymentRequest.Name,
					"monitor": "false",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     deploymentRequest.Name,
						"monitor": "false",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deploymentRequest.Name,
							Image: postgresImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: postgresPort,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-config", deploymentRequest.Name),
										},
									},
								},
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("%s-secret", deploymentRequest.Name),
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   resourceList,
								Requests: resourceList,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      fmt.Sprintf("%s-data", deploymentRequest.Name),
									MountPath: fmt.Sprintf("/mnt/data/%s-data", deploymentRequest.Name),
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: fmt.Sprintf("%s-data", deploymentRequest.Name),
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(postgresStorage),
							},
						},
					},
				},
			},
		},
	}

	_, err = clientSet.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	ok = createService(
		clientSet, deploymentRequest.Name, postgresPort, false, deploymentRequest.ExternalAccess)

	return ok
}

func GetDeploymentStatus(deploymentName string) (DeploymentStatus, bool) {
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Println(err.Error())
		return DeploymentStatus{}, false
	}

	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
		return DeploymentStatus{}, false
	}

	deploymentStatus, ok := extarctDeploymentStatus(clientSet, deployment)
	if !ok {
		return DeploymentStatus{}, false
	}

	return deploymentStatus, true
}

func GetDeploymentsStatus() ([]DeploymentStatus, bool) {
	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Println(err.Error())
		return []DeploymentStatus{}, false
	}

	deployments, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return []DeploymentStatus{}, false
	}

	var deploymentStatuses []DeploymentStatus
	for _, deployment := range deployments.Items {
		deploymentStatus, ok := extarctDeploymentStatus(clientSet, &deployment)
		if !ok {
			return []DeploymentStatus{}, false
		}

		deploymentStatuses = append(deploymentStatuses, deploymentStatus)
	}

	return deploymentStatuses, true
}
