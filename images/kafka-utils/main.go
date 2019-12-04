package main

import (
	"github.com/mesosphere/kudo-kafka-operator/images/kafka-utils/pkgs/client"
	"github.com/mesosphere/kudo-kafka-operator/images/kafka-utils/pkgs/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Running Kafka utils...")

	k8sClient, err := client.GetKubernetesClient()
	if err != nil {
		log.Fatalf("Error initializing client: %+v", err)
	}
	kakfaService := service.KafkaService{
		Client: k8sClient,
		Env:    &service.EnvironmentImpl{},
	}
	kakfaService.WriteIngressToPath("/opt/kafka")
}
