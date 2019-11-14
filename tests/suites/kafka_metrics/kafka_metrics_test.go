package kafka_metrics

import (
	"fmt"
	"testing"

	. "github.com/mesosphere/kudo-kafka-operator/tests/suites"

	"github.com/mesosphere/kudo-kafka-operator/tests/utils"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var (
	customNamespace = "metrics-ns"
)

var _ = Describe("KafkaTest", func() {
	Describe("[Kafka Metrics Checks]", func() {
		Context("metrics-ns installation", func() {
			It("Kafka and Zookeeper statefulset should have 3 replicas with status READY", func() {
				err := utils.KClient.WaitForStatefulSetReadyReplicasCount(DefaultZkStatefulSetName, customNamespace, 3, 240)
				Expect(err).To(BeNil())
				err = utils.KClient.WaitForStatefulSetReadyReplicasCount(DefaultKafkaStatefulSetName, customNamespace, 3, 300)
				Expect(err).To(BeNil())
				Expect(utils.KClient.GetStatefulSetCount(DefaultKafkaStatefulSetName, customNamespace)).To(Equal(3))
			})
			It("Check for metrics endpoint", func() {
				kafkaClient := utils.NewKafkaClient(utils.KClient, &utils.KafkaClientConfiguration{
					Namespace:       utils.String(customNamespace),
				})
				out, err := kafkaClient.GetMetricsEndpointOutput(GetBrokerPodName(0), DefaultContainerName)
				Expect(err).To(BeNil())
				Expect(out).To(ContainSubstring("kafka_controller_ControllerStats_LeaderElectionRateAndTimeMs"))
			})
		})
	})
})

var _ = BeforeSuite(func() {
	utils.TearDown(customNamespace)
	Expect(utils.DeletePVCs("data-dir")).To(BeNil())
	utils.KClient.CreateNamespace(customNamespace, false)
	utils.Setup(customNamespace)
})

var _ = AfterSuite(func() {
	utils.TearDown(customNamespace)
	Expect(utils.DeletePVCs("data-dir")).To(BeNil())
	utils.KClient.DeleteNamespace(customNamespace)
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter(fmt.Sprintf("%s-junit.xml", "kafka-metrics"))
	RunSpecsWithDefaultAndCustomReporters(t, "KafkaMetrics Suite", []Reporter{junitReporter})
}
