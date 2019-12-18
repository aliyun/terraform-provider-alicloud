package alicloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestApplyDefaultArgs(t *testing.T) {
	args := make([]string, 0)
	args = applyDefaultArgs(args)
	if len(args) != 0 {
		t.Log("pass TestApplyDefaultArgs")
		return
	}
	t.Error("TestApplyDefaultArgs failed to apply default args")
}

func TestCreateScalingGroupTags(t *testing.T) {
	validLabels := "a=b,c=d"
	validTaints := "e=f:NoSchedule"
	tags := createScalingGroupTags(validLabels, validTaints)

	validLabelsArr := strings.Split(validLabels, ",")

	validTaintsArr := strings.Split(validTaints, ",")

	for _, label := range validLabelsArr {
		labelKeyValue := strings.Split(label, "=")
		if ok := strings.Contains(tags, fmt.Sprintf("%s%s", LabelPattern, labelKeyValue[0])); ok != true {
			t.Error("failed to pass TestCreateScalingGroupTags,because convert labels failure")
		}
	}

	for _, taint := range validTaintsArr {
		taintKeyValue := strings.Split(taint, "=")
		if ok := strings.Contains(tags, fmt.Sprintf("%s%s", TaintPattern, taintKeyValue[0])); ok != true {
			t.Error("failed to pass TestCreateScalingGroupTags,because convert taints failure")
		}
	}
	t.Log("pass TestCreateScalingGroupTags")
}
