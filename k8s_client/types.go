package k8s_client

type DeploymentResource struct {
	Name  string `json:"name"`
	Limit string `json:"limit"`
}

type DeploymentConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeploymentSecret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CustomDeploymentRequest struct {
	Name           string               `json:"name"`
	Replicas       int32                `json:"replicas"`
	ImageAddress   string               `json:"image_address"`
	ImageTag       string               `json:"image_tag"`
	DomainAddress  string               `json:"domain_address"`
	ServicePort    int32                `json:"service_port"`
	Resources      []DeploymentResource `json:"resources"`
	Configs        []DeploymentConfig   `json:"configs"`
	Secrets        []DeploymentSecret   `json:"secrets"`
	ExternalAccess bool                 `json:"external_access"`
	Monitor        bool                 `json:"monitor"`
}

type PostgresDeploymentRequest struct {
	Name           string               `json:"name"`
	Resources      []DeploymentResource `json:"resources"`
	ExternalAccess bool                 `json:"external_access"`
}

type PodStatus struct {
	Name      string `json:"name"`
	Phase     string `json:"phase"`
	HostIP    string `json:"host_ip"`
	PodIP     string `json:"pod_ip"`
	StartTime string `json:"start_time"`
}

type DeploymentStatus struct {
	Name          string      `json:"name"`
	Replicas      int32       `json:"replicas"`
	ReadyReplicas int32       `json:"ready_replicas"`
	PodStatuses   []PodStatus `json:"pod_statuses"`
}
