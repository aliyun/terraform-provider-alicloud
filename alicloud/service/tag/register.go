package tag

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the TAG service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Tag",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_tag_created_by",
				Factory:  NewCreatedByResource,
			},
		},
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_tag_associated_rule", Factory: ResourceAssociatedRule},
			{TypeName: "alicloud_tag_policy", Factory: ResourcePolicy},
			{TypeName: "alicloud_tag_policy_attachment", Factory: ResourcePolicyAttachment},
			{TypeName: "alicloud_tag_meta_tag", Factory: ResourceMetaTag},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_tag_meta_tags", Factory: DataSourceMetaTags},
		},
	}
}
