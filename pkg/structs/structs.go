package structs

type (
	// Config holds input parameters for the plugin
	Config struct {
		Sensitive       bool
		RoleARN         string
		Region          string
		ServerAddress   string
		K8SCert         string
		K8SToken        string
		K8SUser         string
		GoogleSA        string
		GoogleProjectID string
	}

	// Kube holds inputs for kubernetes commands and configuration
	Kube struct {
		Type           string
		Version        string
		Commands       []string
		ManifestDir    string
		ClusterName    string
		Namespace      string
		Kustomize      string
		AppVersion     string
		ImageName      string
		Rollout        string
		RolloutTimeout string
	}
)
