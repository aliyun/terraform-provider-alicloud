package agentloop

import (
	"testing"
)

func TestUnitAliCloudAgentloopServicePackageResources(t *testing.T) {
	want := map[string]bool{
		"alicloud_agentloop_agent_space":        true,
		"alicloud_agentloop_dataset":            true,
		"alicloud_agentloop_endpoint_connector": true,
		"alicloud_agentloop_evaluation_task":    true,
		"alicloud_agentloop_evaluator":          true,
		"alicloud_agentloop_evaluator_skill":    true,
		"alicloud_agentloop_pipeline":           true,
	}
	servicePackage := ServicePackage()
	if servicePackage.Name != "AgentLoop" || len(servicePackage.SDKResources) != len(want) {
		t.Fatalf("unexpected service package: name=%q resources=%d", servicePackage.Name, len(servicePackage.SDKResources))
	}
	for _, resource := range servicePackage.SDKResources {
		if !want[resource.TypeName] || resource.Factory == nil {
			t.Fatalf("unexpected or invalid resource registration: %s", resource.TypeName)
		}
		delete(want, resource.TypeName)
	}
	if len(want) != 0 {
		t.Fatalf("missing resource registrations: %v", want)
	}
}

func TestUnitAliCloudAgentloopAgentSpaceTrajectoryStoreSchema(t *testing.T) {
	trajectoryStore := ResourceAgentSpace().Schema["trajectory_store_enabled"]
	if !trajectoryStore.Optional || !trajectoryStore.ForceNew {
		t.Fatalf("trajectory_store_enabled must remain Optional and ForceNew")
	}
	// The API does not return this field on GET. A schema Default would add an
	// unverifiable value during import; Create sends the bool zero value instead.
	if trajectoryStore.Default != nil {
		t.Fatalf("trajectory_store_enabled default = %v, want nil", trajectoryStore.Default)
	}
}
