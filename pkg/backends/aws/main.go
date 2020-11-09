package aws

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	backend "github.com/jlehtimaki/kubernetes-ci/pkg/backends"
	s "github.com/jlehtimaki/kubernetes-ci/pkg/structs"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

const awsCliExe = "aws"

type AWSBackend struct {
	backend.BaseBackend
	Config	s.Config
	Kube	s.Kube
}

func NewAWSBackend(config s.Config, kube s.Kube) *AWSBackend {
	backend := &AWSBackend{Config: config, Kube: kube}
	return backend
}

func (b *AWSBackend) Login() []*exec.Cmd {
	var commands []*exec.Cmd
	commands = append(commands, exec.Command(awsCliExe, "--version"))
	if b.Config.RoleARN != "" {
		assumeRole(b.Config.RoleARN)
	}
	commands = append(commands, awsGetKubeConfig(b.Kube.ClusterName, b.Config.Region))

	return commands
}

func awsGetKubeConfig(clusterName string, region string) *exec.Cmd {
	return exec.Command(awsCliExe, "eks", "--region", region, "update-kubeconfig", "--name", clusterName)
}

func assumeRole(roleArn string) {
	logrus.Infof("assuming role %s", roleArn)
	client := sts.New(session.New())
	duration := time.Hour * 1
	stsProvider := &stscreds.AssumeRoleProvider{
		Client:          client,
		Duration:        duration,
		RoleARN:         roleArn,
		RoleSessionName: "drone",
	}

	value, err := credentials.NewCredentials(stsProvider).Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error assuming role!")
	}
	os.Setenv("AWS_ACCESS_KEY_ID", value.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", value.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", value.SessionToken)
}