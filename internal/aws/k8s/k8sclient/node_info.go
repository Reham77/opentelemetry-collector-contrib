// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sclient // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s/k8sclient"

import (
	v1 "k8s.io/api/core/v1"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s/k8sutil"
)

type Label string

const (
	SageMakerNodeHealthStatus Label = "sagemaker.amazonaws.com/node-health-status"
	SageMakerNodeHealthStatusSC Label = "n"
)

type NodeInfo struct {
	Name           string
	Conditions     []*NodeCondition
	Capacity       v1.ResourceList
	Allocatable    v1.ResourceList
	ProviderId     string
	InstanceType   string
	HyperPodLabels map[Label]k8sutil.HyperPodConditionType
}

type NodeCondition struct {
	Type   v1.NodeConditionType
	Status v1.ConditionStatus
}
