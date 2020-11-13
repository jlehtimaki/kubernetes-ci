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
	"log"
	"os/exec"
)

type KubeConfig struct {
	Apiversion     	string              `yaml:"apiversion"`
	Kind           	string              `yaml:"kind"`
	CurrentContext 	string              `yaml:"current-context"`
	Contexts       	[]Context			`yaml:"contexts"`
	Users          	[]User				`yaml:"users"`
	Clusters       	struct {
		Name    	string `yaml:"name"`
		Cluster 	struct {
			Server 	string `yaml:"server"`
			CAData 	string `yaml:"certificate-authority-data"`
		}
	}
}

type Context struct {
	Context 	struct{
		Cluster	string 	`yaml:"cluster"`
		User	string	`yaml:"user"`
	}	`yaml:"context"`
	Name		string	`yaml:"name"`
}

type User struct {
	Name	string	`yaml:"name"`
	User	struct{
		AuthProvider	struct{
			Name	string	`yaml:"name"`
		}	`yaml:"auth-provider"`
	}
}
//kubeUsers := fmt.Sprintf("[{name: user-1, user: {auth-provider: {name: gcp}}}]")


type GoogleBackend struct {
	backend.BaseBackend
	Config	s.Config
	Kube	s.Kube
}

func NewGoogleBackend(config s.Config, kube s.Kube) (*GoogleBackend, error) {
	backend := &GoogleBackend{Config: config, Kube: kube}
	if backend.Config.GoogleSA == "" || backend.Config.GoogleProjectID == "" {
		return nil, fmt.Errorf("missing configuration parameters for GKE")
	}
	return backend, nil
}

func (b *GoogleBackend) Login() []*exec.Cmd {
	var commands []*exec.Cmd
	explicitLogin(b.Config, b.Kube)
	return commands
}

// implicit uses Application Default Credentials to authenticate.
func explicitLogin(c s.Config, k s.Kube) {
	ctx := context.Background()
	jsonByte :=  []byte(c.GoogleSA)

	client, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(jsonByte))
	if err != nil {
		log.Fatal(err)
	}
	req := &containerpb.ListClustersRequest{
		ProjectId: c.GoogleProjectID,
		Zone:      c.Region,
	}

	clusters, err := client.ListClusters(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	for _, cluster := range clusters.Clusters {
		if cluster.Name == k.ClusterName {
			writeConfig(cluster)
		}
	}
}

func writeConfig(cluster *containerpb.Cluster){
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
	kubeUser.User.AuthProvider.Name = "gcp"
	kubeConfig.Contexts = append(kubeConfig.Contexts, *kubeContext)
	kubeConfig.Users = append(kubeConfig.Users, *kubeUser)
	kubeConfig.Clusters.Name = cluster.Name
	kubeConfig.Clusters.Cluster.Server = fmt.Sprintf("https://%s", cluster.Endpoint)
	kubeConfig.Clusters.Cluster.CAData = cluster.MasterAuth.ClusterCaCertificate

	d, err := yaml.Marshal(&kubeConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}