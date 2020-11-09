package main

import (
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	kci "github.com/jlehtimaki/kubernetes-ci"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var revision string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "kubernetes plugin"
	app.Usage = "kubernetes plugin"
	app.Action = run
	app.Version = revision
	app.Flags = []cli.Flag{

		//
		// plugin args
		//

		cli.StringSliceFlag{
			Name:   "actions",
			Usage:  "a list of actions to have kubectl perform",
			EnvVar: "PLUGIN_ACTIONS",
			Value:  &cli.StringSlice{"diff"},
		},
		cli.StringFlag{
			Name:   "type",
			Usage:  "A Type of Kubernetes deployment. eg. EKS, GKE, Baremetal",
			EnvVar: "PLUGIN_TYPE",
			Value:  "Baremetal",
		},
		cli.StringFlag{
			Name:   "k8s_ca",
			Usage:  "CA Certificate to Kubernetes",
			EnvVar: "PLUGIN_CA",
		},
		cli.StringFlag{
			Name:   "k8s_token",
			Usage:  "Token to Kubernetes",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "k8s_user",
			Usage:  "Kubernetes user to authenticate",
			EnvVar: "PLUGIN_K8S_USER",
			Value:  "default",
		},
		cli.StringFlag{
			Name:   "k8s_server",
			Usage:  "Kubernetes server address",
			EnvVar: "PLUGIN_K8S_SERVER",
		},
		cli.StringFlag{
			Name:   "assume_role",
			Usage:  "A role to assume before running the awscli commands",
			EnvVar: "PLUGIN_ASSUME_ROLE",
		},
		cli.StringFlag{
			Name:   "kubectl_version",
			Usage:  "kubectl version number",
			EnvVar: "PLUGIN_KUBECTL_VERSION",
		},
		cli.StringFlag{
			Name:   "cluster_name",
			Usage:  "EKS Cluster Name",
			EnvVar: "PLUGIN_CLUSTER_NAME",
			Value:  "EKS-Cluster",
		},
		cli.StringFlag{
			Name:   "manifest_dir",
			Usage:  "Directory that holds manifests",
			EnvVar: "PLUGIN_MANIFEST_DIR",
			Value:  "./",
		},
		cli.StringFlag{
			Name:   "kubernetes_namespace",
			Usage:  "Namespace for Kubernetes",
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "aws_region",
			Usage:  "AWS Region to use",
			EnvVar: "AWS_REGION",
			Value:  "eu-west-1",
		},
		cli.StringFlag{
			Name:   "kustomize",
			Usage:  "To use kustomize",
			EnvVar: "PLUGIN_KUSTOMIZE",
			Value:  "false",
		},
		cli.StringFlag{
			Name:   "image.version",
			Usage:  "Version to be deployed",
			EnvVar: "PLUGIN_IMAGE_VERSION",
		},
		cli.StringFlag{
			Name:   "image.name",
			Usage:  "Image name to be changed",
			EnvVar: "PLUGIN_IMAGE",
		},
		cli.StringFlag{
			Name:   "rolloutCheck",
			Usage:  "Checking rollout status",
			EnvVar: "PLUGIN_ROLLOUT_CHECK",
			Value:  "true",
		},
		cli.StringFlag{
			Name:   "rolloutTimeout",
			Usage:  "Timeout of rollout",
			EnvVar: "PLUGIN_ROLLOUT_TIMEOUT",
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