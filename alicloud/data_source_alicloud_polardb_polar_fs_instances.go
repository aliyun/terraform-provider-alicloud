package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBPolarFsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBPolarFsInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids":                    {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"description":            {Type: schema.TypeString, Optional: true},
			"relative_db_cluster_id": {Type: schema.TypeString, Optional: true},
			"polar_fs_type":          {Type: schema.TypeString, Optional: true},
			"db_cluster_id":          {Type: schema.TypeString, Optional: true},
			"tags":                   tagsSchemaWithElements(),
			"polar_fs_instances": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id":                        {Type: schema.TypeString, Computed: true},
					"path":                      {Type: schema.TypeString, Computed: true},
					"status":                    {Type: schema.TypeString, Computed: true},
					"description":               {Type: schema.TypeString, Computed: true},
					"create_time":               {Type: schema.TypeString, Computed: true},
					"expire_time":               {Type: schema.TypeString, Computed: true},
					"pay_type":                  {Type: schema.TypeString, Computed: true},
					"region_id":                 {Type: schema.TypeString, Computed: true},
					"zone_id":                   {Type: schema.TypeString, Computed: true},
					"storage_space":             {Type: schema.TypeInt, Computed: true},
					"storage_type":              {Type: schema.TypeString, Computed: true},
					"expired":                   {Type: schema.TypeString, Computed: true},
					"bandwidth":                 {Type: schema.TypeInt, Computed: true},
					"polar_fs_type":             {Type: schema.TypeString, Computed: true},
					"vpc_id":                    {Type: schema.TypeString, Computed: true},
					"vswitch_id":                {Type: schema.TypeString, Computed: true},
					"security_group_id":         {Type: schema.TypeString, Computed: true},
					"relative_db_cluster_id":    {Type: schema.TypeString, Computed: true},
					"category":                  {Type: schema.TypeString, Computed: true},
					"accelerating_enable":       {Type: schema.TypeString, Computed: true},
					"accelerated_storage_space": {Type: schema.TypeString, Computed: true},
					"accelerate_type":           {Type: schema.TypeString, Computed: true},
					"tags":                      {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
				}},
			},
		},
	}
}

func dataSourceAlicloudPolarDBPolarFsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"RegionId": client.RegionId, "PageSize": 100, "PageNumber": 1}
	if ids, ok := d.GetOk("ids"); ok {
		request["PolarFsInstanceIds"] = strings.Join(expandStringList(ids.([]interface{})), ",")
	}
	for field, apiField := range map[string]string{
		"description": "PolarFsInstanceDescription", "relative_db_cluster_id": "RelativeDbClusterId",
		"polar_fs_type": "PolarFsType", "db_cluster_id": "DBClusterId",
	} {
		if value, ok := d.GetOk(field); ok {
			request[apiField] = value
		}
	}
	if tags, ok := d.GetOk("tags"); ok {
		request["Tag"] = expandPolarDBPolarFsTags(tags.(map[string]interface{}))
	}

	all := make([]interface{}, 0)
	for {
		response, err := polarDBRpcPost(client, "DescribePolarFs", request, 5*time.Minute)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_polar_fs_instances", "DescribePolarFs", AlibabaCloudSdkGoERROR)
		}
		items := polarDBPolarFsListItems(response["Items"])
		all = append(all, items...)
		if len(items) < 100 {
			break
		}
		request["PageNumber"] = formatInt(request["PageNumber"]) + 1
	}

	ids := make([]string, 0, len(all))
	instances := make([]map[string]interface{}, 0, len(all))
	for _, value := range all {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		id := fmt.Sprint(item["PolarFsInstanceId"])
		ids = append(ids, id)
		instances = append(instances, map[string]interface{}{
			"id": id, "path": item["PolarFsPath"], "status": item["PolarFsStatus"],
			"description": item["PolarFsInstanceDescription"], "create_time": item["CreateTime"],
			"expire_time": item["ExpireTime"], "pay_type": item["PayType"], "region_id": item["RegionId"],
			"zone_id": item["ZoneId"], "storage_space": item["StorageSpace"], "storage_type": item["StorageType"],
			"expired": item["Expired"], "bandwidth": item["Bandwidth"], "polar_fs_type": item["PolarFsType"],
			"vpc_id": item["VPCId"], "vswitch_id": item["VSwitchId"], "security_group_id": item["SecurityGroupId"],
			"relative_db_cluster_id": item["RelativeDbClusterId"], "category": item["Category"],
			"accelerating_enable": item["AcceleratingEnable"], "accelerated_storage_space": item["AcceleratedStorageSpace"],
			"accelerate_type": item["AccelerateType"], "tags": tagsToMap(polarDBPolarFsTags(item["Tags"])),
		})
	}
	if err := d.Set("polar_fs_instances", instances); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}

func polarDBPolarFsListItems(raw interface{}) []interface{} {
	if wrapper, ok := raw.(map[string]interface{}); ok {
		raw = wrapper["PolarFsPaths"]
	}
	values, _ := raw.([]interface{})
	return values
}

func polarDBPolarFsTags(raw interface{}) interface{} {
	if wrapper, ok := raw.(map[string]interface{}); ok {
		return wrapper["Tag"]
	}
	return nil
}

func expandPolarDBPolarFsTags(tags map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tags))
	for key, value := range tags {
		result = append(result, map[string]interface{}{"Key": key, "Value": value})
	}
	return result
}
