package paiplugin

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the PAIPlugin service.
//
// PAIPlugin (POP product PaiPlugin, version 2022-01-12) exposes a REST/ROA-style
// gateway: every operation maps to an HTTP verb against an /api/v2/* path, with
// request parameters carried in the JSON body (POST/PUT) or the query string
// (GET). The connectivity client's Roa* methods encode that contract, so unlike
// the RPC services the call sites pass a path rather than an Action name.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiPlugin",
		Resources: []conns.Resource{
			{TypeName: "alicloud_pai_plugin_content_generation_job_config", Factory: NewContentGenerationJobConfigResource},
			{TypeName: "alicloud_pai_plugin_content_generation_job", Factory: NewContentGenerationJobResource},
			{TypeName: "alicloud_pai_plugin_message", Factory: NewMessageResource},
		},
		DataSources: []conns.DataSource{
			{TypeName: "alicloud_pai_plugin_algorithm", Factory: NewAlgorithmDataSource},
			{TypeName: "alicloud_pai_plugin_user", Factory: NewUserDataSource},
			{TypeName: "alicloud_pai_plugin_message_config", Factory: NewMessageConfigDataSource},
			{TypeName: "alicloud_pai_plugin_message_metrics", Factory: NewMessageMetricsDataSource},
		},
	}
}
