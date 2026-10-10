package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBLakebaseTenantToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBLakebaseTenantTokenRead,
		Schema: map[string]*schema.Schema{
			"polar_fs_instance_id": {Type: schema.TypeString, Required: true},
			"db_cluster_id":        {Type: schema.TypeString, Optional: true, Computed: true},
			"subdir":               {Type: schema.TypeString, Required: true},
			"tenant":               {Type: schema.TypeString, Optional: true, Computed: true},
			"token":                {Type: schema.TypeString, Computed: true, Sensitive: true},
			"status":               {Type: schema.TypeString, Computed: true},
		},
	}
}

// lintignore: R001
func dataSourceAlicloudPolarDBLakebaseTenantTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"PolarFsInstanceId": d.Get("polar_fs_instance_id"),
		"Subdir":            d.Get("subdir"),
	}
	if value, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = value
	}
	if value, ok := d.GetOk("tenant"); ok {
		request["Tenant"] = value
	}
	response, err := polarDBRpcPost(client, "GetLakebaseTenantToken", request, 5*time.Minute)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_lakebase_tenant_token", "GetLakebaseTenantToken", AlibabaCloudSdkGoERROR)
	}
	for field, apiField := range map[string]string{
		"token": "Token", "subdir": "Subdir", "tenant": "Tenant", "status": "Status", "db_cluster_id": "DBClusterId",
	} {
		if value, ok := response[apiField]; ok && value != nil {
			if err = d.Set(field, value); err != nil {
				return WrapError(err)
			}
		}
	}
	d.SetId(dataResourceIdHash([]string{d.Get("polar_fs_instance_id").(string), d.Get("subdir").(string), fmt.Sprint(d.Get("tenant"))}))
	return nil
}
