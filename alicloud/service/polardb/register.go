package polardb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "polardb", Resources: []conns.Resource{{TypeName: "alicloud_polardb_db_node", Factory: NewDbNodeResource}}}
}
