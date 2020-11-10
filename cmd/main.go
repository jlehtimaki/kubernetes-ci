package main

import (
	kci "github.com/jlehtimaki/kubernetes-ci"
s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
"github.com/sirupsen/logrus"
"github.com/urfave/cli/v2"
"os"
)


func main() {
	app := cli.NewApp()
	app.Name = "kubernetes plugin"
	app.Usage = "kubernetes plugin"
	app.Action = run
	app.Flags = []cli.Flag{

		//
		// plugin args
		//

		&cli.StringSliceFlag{
			Name:   "actions",
			Usage:  "a list of actions to have kubectl perform",
			EnvVars: []string{"PLUGIN_ACTIONS","ACTIONS"},
			Required: true,
		},
		&cli.StringFlag{
			Name:   "type",
			Usage:  "A Type of Kubernetes deployment. eg. EKS, GKE, Baremetal",
			EnvVars: []string{"PLUGIN_TYPE", "TYPE"},
			Value:  "Baremetal",
			Required: true,
		},
		&cli.StringFlag{
			Name:   "k8s_ca",
			Usage:  "CA Certificate to Kubernetes",
			EnvVars: []string{"PLUGIN_CA", "CA"},
		},
		&cli.StringFlag{
			Name:   "k8s_token",
			Usage:  "Token to Kubernetes",
			EnvVars: []string{"PLUGIN_TOKEN", "TOKEN"},
		},
		&cli.StringFlag{
			Name:   "k8s_user",
			Usage:  "Kubernetes user to authenticate",
			EnvVars: []string{"PLUGIN_K8S_USER", "K8S_USER"},
			Value:  "default",
		},
		&cli.StringFlag{
			Name:   "k8s_server",
			Usage:  "Kubernetes server address",
			EnvVars: []string{"PLUGIN_K8S_SERVER", "K8S_SERVER"},
		},
		&cli.StringFlag{
			Name:   "assume_role",
			Usage:  "A role to assume before running the awscli commands",
			EnvVars: []string{"PLUGIN_ASSUME_ROLE", "ASSUME_ROLE"},
		},
		&cli.StringFlag{
			Name:   "kubectl_version",
			Usage:  "kubectl version number",
			EnvVars: []string{"PLUGIN_KUBECTL_VERSION", "KUBECTL_VERSION"},
		},
		&cli.StringFlag{
			Name:   "cluster_name",
			Usage:  "EKS Cluster Name",
			EnvVars: []string{"PLUGIN_CLUSTER_NAME", "CLUSTER_NAME"},
			Value:  "EKS-Cluster",
		},
		&cli.StringFlag{
			Name:   "manifest_dir",
			Usage:  "Directory that holds manifests",
			EnvVars: []string{"PLUGIN_MANIFEST_DIR", "MANIFEST_DIR"},
			Value:  "./",
		},
		&cli.StringFlag{
			Name:   "kubernetes_namespace",
			Usage:  "Namespace for Kubernetes",
			EnvVars: []string{"PLUGIN_NAMESPACE", "NAMESPACE"},
		},
		&cli.StringFlag{
			Name:   "aws_region",
			Usage:  "AWS Region to use",
			EnvVars: []string{"AWS_REGION", "AWS_REGION"},
			Value:  "eu-west-1",
		},
		&cli.StringFlag{
			Name:   "kustomize",
			Usage:  "To use kustomize",
			EnvVars: []string{"PLUGIN_KUSTOMIZE", "KUSTOMIZE"},
			Value:  "false",
		},
		&cli.StringFlag{
			Name:   "image.version",
			Usage:  "Version to be deployed",
			EnvVars: []string{"PLUGIN_IMAGE_VERSION", "IMAGE_VERSION"},
		},
		&cli.StringFlag{
			Name:   "image.name",
			Usage:  "Image name to be changed",
			EnvVars: []string{"PLUGIN_IMAGE", "IMAGE"},
		},
		&cli.StringFlag{
			Name:   "rolloutCheck",
			Usage:  "Checking rollout status",
			EnvVars: []string{"PLUGIN_ROLLOUT_CHECK", "ROLLOUT_CHECK"},
			Value:  "true",
		},
		&cli.StringFlag{
			Name:   "rolloutTimeout",
			Usage:  "Timeout of rollout",
			EnvVars: []string{"PLUGIN_ROLLOUT_TIMEOUT", "ROLLOUT_TIMEOUT"},
			Value:  "1m",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := kci.Plugin{
		Config: s.Config{
			RoleARN:       c.String("assume_role"),
			Region:        c.String("aws_region"),
			ServerAddress: c.String("k8s_server"),
			K8SUser:       c.String("k8s_user"),
			K8SCert:       c.String("k8s_ca"),
			K8SToken:      c.String("k8s_token"),
		},
		Kube: s.Kube{
			Type:           c.String("type"),
			Version:        c.String("kubectl_version"),
			Commands:       c.StringSlice("actions"),
			ClusterName:    c.String("cluster_name"),
			ManifestDir:    c.String("manifest_dir"),
			Namespace:      c.String("kubernetes_namespace"),
			AppVersion:     c.String("image.version"),
			Kustomize:      c.String("kustomize"),
			ImageName:      c.String("image.name"),
			Rollout:        c.String("rolloutCheck"),
			RolloutTimeout: c.String("rolloutTimeout"),
		},
	}

	return plugin.Exec()
}
