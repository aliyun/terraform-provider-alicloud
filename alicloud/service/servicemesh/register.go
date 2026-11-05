package servicemesh

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Servicemesh service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Servicemesh",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_service_mesh_service_mesh", Factory: ResourceServiceMeshServiceMesh},
			{TypeName: "alicloud_service_mesh_user_permission", Factory: ResourceServiceMeshUserPermission},
			{TypeName: "alicloud_service_mesh_extension_provider", Factory: ResourceServiceMeshExtensionProvider},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_service_mesh_service_meshes", Factory: DataSourceServiceMeshServiceMeshes},
			{TypeName: "alicloud_service_mesh_versions", Factory: DataSourceServiceMeshVersions},
			{TypeName: "alicloud_service_mesh_extension_providers", Factory: DataSourceServiceMeshExtensionProviders},
		},
	}
}
