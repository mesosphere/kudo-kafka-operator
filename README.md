![kudo-kafka](./docs/resources/images/kudo-kafka.png)

# KUDO Kafka Operator

The KUDO Kafka operator creates, configures and manages [Apache Kafka](https://kafka.apache.org/) clusters running on Kubernetes.

### Overview

KUDO Kafka is a Kubernetes operator built on [KUDO](kudo.dev) to manage Apache Kafka in a scalable, repeatable, and standardized way over Kubernetes. Currently KUDO Kafka supports:

- Securing the cluster in various ways: TLS encryption, Kerberos authentication, Kafka AuthZ
- Prometheus metrics right out of the box with example Grafana dashboards
- Kerberos support
- Graceful rolling updates for any cluster configuration changes
- Graceful rolling upgrades when upgrading the operator version
- External access through LB/Nodeports
- Mirror-maker integration
- Cruise Control integration
- Connect integration

To get more information around KUDO Kafka architecture please take a look on the [KUDO Kafka concepts](./docs/concepts.md) document.

## Getting started

The latest stable version of Kafka operator is `1.3.3`
For more details, please see the [docs](./docs/) folder.

## Testing Features

You can also test some of key features of KUDO Kafka present in [features runbook](./docs/features-runbooks.md).


## Releases

| KUDO Kafka | Apache Kafka | Minimum KUDO Version |
| ---------- | ------------ | -------------------- |
| 1.2.1      | 2.4.1        | 0.11.0               |
| 1.3.0      | 2.5.0        | 0.11.0               |
| 1.3.1      | 2.5.0        | 0.13.0               |
| 1.3.2      | 2.5.1        | 0.14.0               |
| **1.3.3**  | **2.5.1**    | **0.14.0**           |


## Unreleased version

| KUDO Kafka | Apache Kafka | Minimum KUDO Version |
| ---------- | ------------ | -------------------- |
| 1.3.4      | 2.6.0        | 0.14.0               |
