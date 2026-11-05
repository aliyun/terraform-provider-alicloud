package openapiexplorer

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage wires the openapiexplorer product into the provider registry.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "OpenAPIExplorer",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_open_api_explorer_api_mcp_server", Factory: resourceAliCloudOpenApiExplorerApiMcpServer},
		},
	}
}
