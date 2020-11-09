package kubernetes

import (
	"bytes"
	"fmt"
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

const kubeExe = "kubectl"
const kustomizeExe = "kustomize"

var (
	path      = "/bin/kubectl"
	namespace = "default"
)

type DeploymentFile struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	}
}

func KubeCommand(kube s.Kube, command string) *exec.Cmd {
	var args []string
	if kube.Namespace != "" {
		args = append(args, "--namespace", kube.Namespace)
	}
	if kube.Kustomize == "true" {
		args = append(args, command, "-f", "-")
	} else {
		args = append(args, command, "-f", kube.ManifestDir)
	}
	return exec.Command(kubeExe, args...)
}

func KustomizeSetVersion(kube s.Kube) *exec.Cmd {
	imageName := fmt.Sprintf("%s:%s", kube.ImageName, kube.AppVersion)
	return exec.Command(kustomizeExe, "edit", "set", "image", imageName)
}

func InstallKubectl(version string) error {
	arch := os.Getenv("GOARCH")
	if arch == "" {
		arch = "amd64"
	}
	downloadUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubectl", version, arch)
	logrus.Info("Installing Kubectl version ", version)
	err := downloadFile(path, downloadUrl)
	if err != nil {
		return err
	}
	err = addExecRights()
	if err != nil {
		return err
	}
	return nil
}

func addExecRights() error {
	err := os.Chmod("/bin/kubectl", 0777)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	//Get the response bytes from the url
	logrus.Info("Downloading file ", url)
	response, err := http.Get(url)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	// Check that file exists
	if !checkFileExists(filepath) {
		return fmt.Errorf("kubectl file not found")
	}
	return nil
}

func checkFileExists(filepath string) bool {
	// Returns true if file exists
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return true
}

// Check the rollout statuses of each Deployment, Statefulset & DaemonSet
func CheckRolloutStatus(kube s.Kube) []*exec.Cmd {
	// Init commands
	var commands []*exec.Cmd

	// find all yaml/yml files from the manifest directory
	yamlFiles := findYAMLFiles(kube.ManifestDir)

	// loop through the files and parse the yaml data
	for _, file := range yamlFiles {
		deploymentFile := DeploymentFile{}
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			logrus.Errorf("could not read file %s", file)
		}

		// Set the correct namespace
		if kube.Namespace != "" {
			namespace = kube.Namespace
		}

		reader := bytes.NewReader(yamlFile)

		decoder := yaml.NewDecoder(reader)

		for decoder.Decode(&deploymentFile) == nil {
			// Set the correct namespace
			if deploymentFile.Metadata.Namespace != "" {
				namespace = deploymentFile.Metadata.Namespace
			}
			// Set the commands rollout status commands if deployment kind matches
			if deploymentFile.Kind == "Deployment" || deploymentFile.Kind == "StatefulSet" || deploymentFile.Kind == "DaemonSet" {
				commands = append(
					commands, exec.Command(
						kubeExe,
						"-n",
						namespace,
						"rollout", "status",
						deploymentFile.Kind,
						deploymentFile.Metadata.Name,
						"--timeout", kube.RolloutTimeout,
					),
				)
			}
		}
	}
	return commands
}

func findYAMLFiles(path string) []string {
	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".yml", f.Name())
			if err == nil && r {
				files = append(files, path)
			}
			r, err = regexp.MatchString(".yaml", f.Name())
			if err == nil && r {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}
