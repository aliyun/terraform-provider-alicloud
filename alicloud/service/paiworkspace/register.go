package paiworkspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiWorkspace (AIWorkSpace 2021-02-04) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiWorkspace",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_pai_workspace_workspace", Factory: ResourcePaiWorkspaceWorkspace},
			{TypeName: "alicloud_pai_workspace_dataset", Factory: ResourcePaiWorkspaceDataset},
			{TypeName: "alicloud_pai_workspace_datasetversion", Factory: ResourcePaiWorkspaceDatasetversion},
			{TypeName: "alicloud_pai_workspace_experiment", Factory: ResourcePaiWorkspaceExperiment},
			{TypeName: "alicloud_pai_workspace_member", Factory: ResourcePaiWorkspaceMember},
			{TypeName: "alicloud_pai_workspace_model", Factory: ResourcePaiWorkspaceModel},
			{TypeName: "alicloud_pai_workspace_model_version", Factory: ResourcePaiWorkspaceModelVersion},
			{TypeName: "alicloud_pai_workspace_run", Factory: ResourcePaiWorkspaceRun},
			{TypeName: "alicloud_pai_workspace_user_config", Factory: ResourcePaiWorkspaceUserConfig},
			{TypeName: "alicloud_pai_workspace_code_source", Factory: ResourcePaiWorkspaceCodeSource},
		},
	}
}
