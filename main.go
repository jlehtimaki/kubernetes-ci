package kubernetes_ci

import (
	"fmt"
	backend "github.com/jlehtimaki/kubernetes-ci/pkg/backends"
	"github.com/jlehtimaki/kubernetes-ci/pkg/backends/aws"
	"github.com/jlehtimaki/kubernetes-ci/pkg/backends/baremetal"
	"github.com/jlehtimaki/kubernetes-ci/pkg/backends/google"
	k "github.com/jlehtimaki/kubernetes-ci/pkg/kubernetes"
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"strings"
)

type (
	// Plugin represents the plugin instance to be executed
	Plugin struct {
		Config s.Config
		Kube   s.Kube
	}
)

var (
	allowedCommands = []string{"apply", "delete", "diff"}
	kubeExe         = "kubectl"
	kustomizeExe    = "kustomize"
)

func allowedCommand(command string) bool {
	for _, com := range allowedCommands {
		if com == command {
			return true
		}
	}
	return false
}

// Exec executes the plugin
func (p Plugin) Exec() error {
	// Install specified version of kubectl
	if p.Kube.Version != "" {
		err := k.InstallKubectl(p.Kube.Version)
		if err != nil {
			return err
		}
	}

	// Initialize commands
	var commands []*exec.Cmd

	// Print Kubectl version
	commands = append(commands, exec.Command(kubeExe, "version", "--client=true"))

	// Set the backend provider
	var b backend.Backend
	var err error
	switch p.Kube.Type {
	case "EKS":
		fmt.Println("Using EKS type of Kubernetes")
		b, err = aws.NewAWSBackend(p.Config, p.Kube)
	case "Baremetal":
		fmt.Println("Using Baremetal type of Kubernetes")
		b, err = baremetal.NewBaremetalBackend(p.Config, p.Kube)
	case "GKE":
		fmt.Println("Using GKE type of Kubernetes")
		b, err = google.NewGoogleBackend(p.Config, p.Kube)
	default:
		log.Fatalf("unknown backend: %s", p.Kube.Type)
	}

	if err != nil {
		return err
	}

	// Login to backend
	commands = append(commands, b.Login()...)

	// Set version with Kustomize
	if p.Kube.AppVersion != "" {
		commands = append(commands, k.KustomizeSetVersion(p.Kube))
	}

	// Add commands listed in actions
	for _, action := range p.Kube.Commands {
		if allowedCommand(action) {
			commands = append(commands, k.KubeCommand(p.Kube, action))
		} else {
			return fmt.Errorf("valid actions are: apply, destroy.  You provided %s", action)
		}
	}

	if p.Kube.Rollout == "true" {
		commands = append(commands, k.CheckRolloutStatus(p.Kube)...)
	}

	// Run commands
	for _, c := range commands {
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if !p.Config.Sensitive {
			trace(c)
		}

		if strings.Contains(c.String(), "edit") {
			c.Dir = p.Kube.ManifestDir
		}

		if p.Kube.Kustomize == "true" {
			// Pipeline the kustomize build command with kubectl command
			c1 := exec.Command(kustomizeExe, "build", p.Kube.ManifestDir)
			c2 := c

			// initialize error
			var err error

			// pipe the commands
			c2.Stdin, err = c1.StdoutPipe()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to pipeline commands")
			}
			c2.Stdout = os.Stdout

			// run the commands
			err = c2.Start()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute kubectl command")
			}
			err = c1.Run()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute kustomize command")
			}
			// wait for the first command to finish
			err = c2.Wait()
			if err != nil && !strings.Contains(c.String(), "diff") {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to wait kustomize command")
			}
		} else {
			err := c.Run()
			// If kubectl command is diff ignore exit code since diff returns exit 1 if the is changes
			if err != nil && !strings.Contains(c.String(), "diff") {
				logrus.Info(c.String())
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute a command")
			}
		}

		logrus.Debug("Command completed successfully")
	}

	return nil
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
