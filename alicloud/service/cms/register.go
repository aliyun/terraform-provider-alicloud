package cms

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		DataSources: []conns.DataSource{
			{TypeName: "alicloud_cms_biz_traces", Factory: NewBizTracesDataSource},
		},
		Resources: []conns.Resource{
			{TypeName: "alicloud_cms_biz_trace", Factory: NewBizTraceResource},
		},
	}
}
