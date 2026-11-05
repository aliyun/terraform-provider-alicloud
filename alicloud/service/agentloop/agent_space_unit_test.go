package agentloop

import (
	"testing"
)

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
