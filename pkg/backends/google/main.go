package google

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	backend "github.com/jlehtimaki/kubernetes-ci/pkg/backends"
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type KubeConfig struct {
	Apiversion     string     `yaml:"apiVersion"`
	Kind           string     `yaml:"kind"`
	CurrentContext string     `yaml:"current-context"`
	Contexts       []Context  `yaml:"contexts"`
	Users          []User     `yaml:"users"`
	Clusters       []Clusters `yaml:"clusters"`
}

type Clusters struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server string `yaml:"server"`
		CAData string `yaml:"certificate-authority-data"`
	} `yaml:"cluster"`
}

type Context struct {
	Context struct {
		Cluster string `yaml:"cluster"`
		User    string `yaml:"user"`
	} `yaml:"context"`
	Name string `yaml:"name"`
}

type User struct {
	Name string `yaml:"name"`
	User struct {
		AuthProvider struct {
			Name string `yaml:"name"`
		} `yaml:"auth-provider"`
	} `yaml:"user"`
}

type GoogleBackend struct {
	backend.BaseBackend
	Config s.Config
	Kube   s.Kube
}

var (
	kubeConfigDir  = ".kube"
	kubeConfigFile = fmt.Sprintf("%s/config", kubeConfigDir)
	kubeAuthFile   = fmt.Sprintf("%s/auth.json", kubeConfigDir)
)

func NewGoogleBackend(config s.Config, kube s.Kube) (*GoogleBackend, error) {
	googleBackend := &GoogleBackend{Config: config, Kube: kube}
	if googleBackend.Config.GoogleSA == "" || googleBackend.Config.GoogleProjectID == "" {
		return nil, fmt.Errorf("missing configuration parameters for GKE")
	}
	return googleBackend, nil
}

func (b *GoogleBackend) Login() []*exec.Cmd {
	explicitLogin(b.Config, b.Kube)
	return nil
}

// implicit uses Application Default Credentials to authenticate.
func explicitLogin(c s.Config, k s.Kube) {
	ctx := context.Background()

	// Convert ServiceAccount to JSON
	jsonByte := []byte(c.GoogleSA)

	// Create new ClusterClient
	client, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(jsonByte))
	if err != nil {
		log.Fatal(err)
	}

	// List clusters and find correct one
	req := &containerpb.ListClustersRequest{
		ProjectId: c.GoogleProjectID,
		Zone:      c.Region,
	}

	clusters, err := client.ListClusters(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	// Write configuration files to .kube folder
	err = os.Mkdir(kubeConfigDir, 0744)
	if err != nil {
		log.Fatal(err)
	}
	for _, cluster := range clusters.Clusters {
		if cluster.Name == k.ClusterName {
			writeConfig(cluster)
			writeAuthentication(c.GoogleSA)
		}
	}
}

func writeAuthentication(googleSA string) {
	d := []byte(googleSA)
	err := ioutil.WriteFile(kubeAuthFile, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", kubeAuthFile)
	if err != nil {
		log.Fatal(err)
	}
}

func writeConfig(cluster *containerpb.Cluster) {
	// Set kubeconfig params
	kubeConfig := &KubeConfig{
		Apiversion:     "v1",
		Kind:           "Config",
		CurrentContext: "my-cluster",
	}
	kubeContext := &Context{
		Context: struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		}{
			cluster.Name,
			"user-1",
		},
		Name: "my-cluster",
	}
	kubeUser := &User{
		Name: "user-1",
	}
	kubeCluster := &Clusters{
		Name: cluster.Name,
		Cluster: struct {
			Server string `yaml:"server"`
			CAData string `yaml:"certificate-authority-data"`
		}{
			fmt.Sprintf("https://%s", cluster.Endpoint),
			cluster.MasterAuth.ClusterCaCertificate,
		},
	}
	kubeUser.User.AuthProvider.Name = "gcp"
	kubeConfig.Contexts = append(kubeConfig.Contexts, *kubeContext)
	kubeConfig.Users = append(kubeConfig.Users, *kubeUser)
	kubeConfig.Clusters = append(kubeConfig.Clusters, *kubeCluster)

	// Convert struct to YAML
	d, err := yaml.Marshal(&kubeConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Write kubeconfig
	err = ioutil.WriteFile(kubeConfigFile, d, 0644)
	if err != nil {
		log.Fatalf("write file: %s", err)
	}

	// Set kubeconfig env variable
	err = os.Setenv("KUBECONFIG", kubeConfigFile)
	if err != nil {
		log.Fatalf("os.Setenv: %s", err)
	}

}
