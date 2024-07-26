// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sutil // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s/k8sutil"

import (
	"fmt"
	"strings"
)

// CreatePodKey concatenates namespace and podName to get a pod key
func CreatePodKey(namespace, podName string) string {
	if namespace == "" || podName == "" {
		return ""
	}
	return fmt.Sprintf("namespace:%s,podName:%s", namespace, podName)
}

// CreateContainerKey concatenates namespace, podName and containerName to get a container key
func CreateContainerKey(namespace, podName, containerName string) string {
	if namespace == "" || podName == "" || containerName == "" {
		return ""
	}
	return fmt.Sprintf("namespace:%s,podName:%s,containerName:%s", namespace, podName, containerName)
}

// ParseInstanceIdFromProviderId parses EC2 instance id from node's provider id which has format of aws:///<subnet>/<instanceId>
func ParseInstanceIdFromProviderId(providerId string) string {
	if providerId == "" || !strings.HasPrefix(providerId, "aws://") {
		return ""
	}
	return providerId[strings.LastIndex(providerId, "/")+1:]
}

type HyperPodConditionType int8

const (
	Schedulable                     HyperPodConditionType = iota
	SchedulablePreferred            HyperPodConditionType = iota
	UnschedulablePendingReplacement HyperPodConditionType = iota
	UnschedulablePendingReboot      HyperPodConditionType = iota
	Unschedulable                   HyperPodConditionType = iota
	Unknown                   		HyperPodConditionType = iota
)

// String - Creating common behavior - give the type a String function
func (ct HyperPodConditionType) String() string {
	return [...]string{"Schedulable", "SchedulablePreferred", "UnschedulablePendingReplacement", "UnschedulablePendingReboot", "Unschedulable", "Unknown"}[ct]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (ct HyperPodConditionType) EnumIndex() int {
	return int(ct)
}

var (
    HyperPodConditionTypeMap = map[string]HyperPodConditionType{
        "Schedulable":   Schedulable,
        "SchedulablePreferred": SchedulablePreferred,
        "UnschedulablePendingReplacement": UnschedulablePendingReplacement,
        "UnschedulablePendingReboot": UnschedulablePendingReboot,
        "Unschedulable":   Unschedulable,
    }
)

func ParseString(str string) (HyperPodConditionType, bool) {
    c, ok := HyperPodConditionTypeMap[str]
    return c, ok
}
