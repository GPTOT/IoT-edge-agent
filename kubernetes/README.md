
# Edge Agent Kubernetes Example

The [deploy](deploy.yaml) file will deploy a single-replica statefulset with a container running the redpanda edge agent. A configmap is included in this file, and will need to be modified for your deployment.

The instructions below assume you don't already have source and destination clusters available, and will walk you through the steps for both deploying these two clusters and deploying the redpanda edge agent. If you already have source and destination clusters, then you can skip [the steps](#deploy-source-and-destinaton-clusters) for deploying these clusters.

## Prerequisites
