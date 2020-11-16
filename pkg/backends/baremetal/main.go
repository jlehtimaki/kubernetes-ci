package baremetal

import (
	"fmt"
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	"io/ioutil"
	"os/exec"
)

var (
	kubeExe = "kubectl"
)

type BaremetalBackend struct {
	Config s.Config
	Kube   s.Kube
}

func NewBaremetalBackend(config s.Config, kube s.Kube) (*BaremetalBackend, error) {
	backend := &BaremetalBackend{Config: config, Kube: kube}
	if backend.Config.K8SUser == "" || backend.Config.K8SCert == "" || backend.Config.K8SToken == "" || backend.Config.ServerAddress == "" {
		return nil, fmt.Errorf("missing configuration parameters for baremetal")
	}
	return backend, nil
}

func (b *BaremetalBackend) Login() []*exec.Cmd {
	return baremetalKubeConfig(b.Config.K8SToken, b.Config.K8SCert, b.Config.ServerAddress, b.Config.K8SUser)
}

func baremetalKubeConfig(token string, cert string, server string, user string) []*exec.Cmd {
	fmt.Println("Setting up Baremetal Kubernetes configuration")
	if cert != "" {
		// Write certificate file
		writeCertToFile(cert)
	}

	// Assign all needed kubernetes commands and return them
	var commands []*exec.Cmd
	tokenString := fmt.Sprintf("--token=%s", token)
	serverString := fmt.Sprintf("--server=%s", server)
	userString := fmt.Sprintf("--user=%s", user)
	commands = append(commands, exec.Command(kubeExe, "config", "set-credentials", "default", tokenString))
	if cert != "" {
		commands = append(commands, exec.Command(kubeExe, "config", "set-cluster", "default", serverString, "--certificate-authority=ca.crt"))
	} else {
		commands = append(commands, exec.Command(kubeExe, "config", "set-cluster", "default", serverString, "--insecure-skip-tls-verify=true"))
	}
	commands = append(commands, exec.Command(kubeExe, "config", "set-context", "default", "--cluster=default", userString))
	commands = append(commands, exec.Command(kubeExe, "config", "use-context", "default"))
	return commands
}

func writeCertToFile(cert string) {
	err := ioutil.WriteFile("ca.crt", []byte(cert), 0644)
	if err != nil {
		fmt.Printf("Could not write certificate file: %s", err.Error())
	}
}
