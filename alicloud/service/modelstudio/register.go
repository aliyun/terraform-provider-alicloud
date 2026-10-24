package modelstudio

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "modelstudio", Resources: []conns.Resource{{TypeName: "alicloud_model_studio_workspace", Factory: NewWorkspaceResource}}}
}
