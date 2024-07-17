// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"

	ci "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s/k8sclient"
)

func TestSetGetCPUCapacity(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeInfo.setCPUCapacity(int(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int32(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int64(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint32(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint64(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	// with invalid type
	nodeInfo.setCPUCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())

	// with negative value
	nodeInfo.setCPUCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int(-3))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int32(-4))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
}

func TestSetGetMemCapacity(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeInfo.setMemCapacity(int(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int32(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int64(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint32(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint64(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	// with invalid type
	nodeInfo.setMemCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())

	// with negative value
	nodeInfo.setMemCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())
}

func TestGetNodeStatusCapacityPods(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityPods, valid := nodeInfo.getNodeStatusCapacityPods()
	assert.True(t, valid)
	assert.Equal(t, uint64(5), nodeStatusCapacityPods)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityPods, valid = nodeInfo.getNodeStatusCapacityPods()
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCapacityPods)
	assert.False(t, nodeInfo.isHyperPodNode())
}

func TestGetNodeStatusAllocatablePods(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusAllocatablePods, valid := nodeInfo.getNodeStatusAllocatablePods()
	assert.True(t, valid)
	assert.Equal(t, uint64(15), nodeStatusAllocatablePods)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusAllocatablePods, valid = nodeInfo.getNodeStatusAllocatablePods()
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusAllocatablePods)
	assert.False(t, nodeInfo.isHyperPodNode())
}

func TestGetNodeStatusCondition(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid := nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeDiskPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeMemoryPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodePIDPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeNetworkUnavailable)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeNetworkUnavailable)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())
}

func TestGetNodeConditionUnknown(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid := nodeInfo.getNodeConditionUnknown()
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeConditionUnknown()
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.False(t, nodeInfo.isHyperPodNode())
}

func TestGetLabelValueUnknown(t *testing.T) {
	nodeInfo := newNodeInfo("hyperpod-testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid := nodeInfo.getLabelValueUnknown(k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	assert.True(t, nodeInfo.isHyperPodNode())

	nodeInfo = newNodeInfo("hyperpod-testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getLabelValueUnknown(k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)
	assert.True(t, nodeInfo.isHyperPodNode())
}

func TestGetLabelValue(t *testing.T) {
	nodeInfo := newNodeInfo("hyperpod-testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	assert.True(t, nodeInfo.isHyperPodNode())

	nodeStatusCondition, valid := nodeInfo.getLabelValue(ci.UnschedulablePendingReplacement, k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	nodeStatusCondition, valid = nodeInfo.getLabelValue(ci.Schedulable, k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)

	nodeStatusCondition, valid = nodeInfo.getLabelValue(ci.SchedulablePreferred, k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	nodeStatusCondition, valid = nodeInfo.getLabelValue(ci.UnschedulablePendingReboot, k8sclient.SageMakerNodeHealthStatus)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	const TestingLabel k8sclient.Label = "TestingLabel"
	nodeStatusCondition, valid = nodeInfo.getLabelValue(ci.UnschedulablePendingReboot, TestingLabel)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
}
