package utils

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// KDCClient Struct defining the KDC Client
type KDCClient struct {
	pod       *v1.Pod
	namespace string
}

const (
	POD_NAME       = "kdc"
	CONTAINER_NAME = "kdc"
)

// setNamespace Set namespace
func (k *KDCClient) SetNamespace(namespace string) {
	k.namespace = namespace
}

// deployPod Use it to deploy the kdc server
func (k *KDCClient) Deploy() {
	repoRoot, exists := os.LookupEnv("REPO_ROOT")

	if exists {
		Create(repoRoot+"/tests/suites/kafka_kerberos/resources/kdc.yaml", k.namespace)
		KClient.WaitForPod("kdc", k.namespace, 240)
	} else {
		log.Warningf("Environment variable REPO_ROOT is not set!")
	}
}

// CreateKeytabSecret Pass it string array of principals and it will create a keytab secret
func (k *KDCClient) CreateKeytabSecret(principals []string, serviceName string, secretName string) {
	//
	command := "printf \"" + strings.Join(principals, "\n") + "\n\" > /kdc/" + serviceName + "-principals.txt;" +
		"cat /kdc/" + serviceName + "-principals.txt | while read line; do /usr/sbin/kadmin -l add --use-defaults --random-password $line; done;" +
		"rm /kdc/" + serviceName + ".keytab;" +
		"cat /kdc/" + serviceName + "-principals.txt | while read line; do /usr/sbin/kadmin -l ext -k /kdc/" + serviceName + ".keytab $line; done;"

	stdout, _ := KClient.ExecInPod(k.namespace, POD_NAME, CONTAINER_NAME, []string{"/bin/sh", "-c", command})
	stdout, _ = KClient.ExecInPod(k.namespace, POD_NAME, CONTAINER_NAME, []string{"/bin/sh", "-c", "cat /kdc/" + serviceName + ".keytab | base64 -w 0"})

	KClient.createSecret(secretName, []string{"kafka.keytab", stdout}, k.namespace)
}
