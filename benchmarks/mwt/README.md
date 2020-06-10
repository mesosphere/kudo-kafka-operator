

## Prerequisites

- Kubernetes cluster up
- KUDO CLI installed on node kuttl is running from
- KUDO manager installed in the cluster

# TestSuites

- **Setup:** used to setup a MWT test use the following command

` KAFKA_PARAMETER_FILE=params-mwt kubectl kuttl test setup/ --parallel 1 --skip-delete`

or

` KAFKA_PARAMETER_FILE=params-mwt kubectl kuttl test --config kuttl-setup.yaml`


- **Teardown:** used to remove zookeeper/kafka (and verify it has been removed)

`kubectl kuttl test --config kuttl-teardown.yaml`

## Tests under **Setup**

- 00-kudo-check: confirms

  1. commands locally are present
  1. kudo manager is ready at the server

- 01-zookeeper-install
  1. first creates and asserts "kafka-mwt" namespace exists
  1. installs the zookeeper operator CRDs and asserts they exist
  1. installs zookeeper and asserts that the deployment plan is
     complete
     
- 02-kafka-install
  1. installs the kafka operator CRDs and asserts they exist
  1. installs kafka and asserts that the deployment plan is
     complete
     
- 03-dashboard-install
  1. installs the grafana configmap dashboard, that is loaded by the grafana operator.
   Test expects the operator to reading the kubeaddons namespace
     
## Verification

run the next command to setup the utils pod we will use for verification

```
kubectl kuttl verify-workload/
```

Now that kuttl confirms that utils pod is up and ready, we can exec inside the pod to run verifications

```
kubectl exec -ti deploy/utils-pod -n kafka-mwt -- sh
```

messages in:
```
curl "prometheus-kubeaddons-prom-prometheus.kubeaddons.svc.cluster.local:9090/api/v1/query?query=kafka_server_BrokerTopicMetrics_MessagesIn_total\{namespace='kafka-mwt',service='kafka-instance-svc',topic=''\}" | jq -r .data.result[].value[1] | awk '{ sum += $1 } END { print sum }'
```

message rate:

```
curl "prometheus-kubeaddons-prom-prometheus.kubeaddons.svc.cluster.local:9090/api/v1/query?query=sum(rate(kafka_server_BrokerTopicMetrics_MessagesIn_total\{service='kafka-instance-svc',namespace='kafka-mwt',topic=''\}\[1m\]))" | jq -r .data.result[].value[1]
```

Alternative approach is to monitor the dashboard installed with the `setup`

**NOTES:**

1. The MWT parameter values need to be changed. This was built and tested on
   konvoy but not on a MWT env
2. The timeouts may need to change. In particular the wait for namespace, the
   wait for deploy to finish and the wait for deletes.
3. This was tested with konvoy with config setting established with
   `./konvoy apply kubeconfig --force-overwrite`
