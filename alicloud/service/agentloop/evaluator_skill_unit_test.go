package agentloop

import (
	"testing"
)

func TestUnitAliCloudAgentloopEvaluatorSkillEnableDefault(t *testing.T) {
	got := ResourceEvaluatorSkill().Schema["enable"].Default
	if got != false {
		t.Fatalf("enable default = %v, want false", got)
	}
}
