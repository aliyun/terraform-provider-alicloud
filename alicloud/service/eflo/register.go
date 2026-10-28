package eflo

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "eflo", Resources: []conns.Resource{{TypeName: "alicloud_eflo_er_attachment", Factory: NewErAttachmentResource}}}
}
