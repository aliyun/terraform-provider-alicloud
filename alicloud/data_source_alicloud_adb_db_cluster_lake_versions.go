package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAdbDbClusterLakeVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAdbDbClusterLakeVersionsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Preparing", "Creating", "Restoring", "Running", "Deleting", "ClassChanging", "NetAddressCreating", "NetAddressDeleting"}, false),
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  30,
			},
			"versions": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"commodity_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"compute_resource": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"connection_string": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"db_cluster_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"db_cluster_version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"engine": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"engine_version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"expire_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"expired": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"lock_mode": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"lock_reason": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"payment_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"storage_resource": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vswitch_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zone_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAdbDbClusterLakeVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["DBClusterStatus"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []interface{}
	var response map[string]interface{}
	var err error

	for {
		action := "DescribeDBClusters"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_adb_db_cluster_lake_versions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBCluster", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBClusterId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	adbService := AdbService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                 fmt.Sprint(object["DBClusterId"]),
			"commodity_code":     object["CommodityCode"],
			"compute_resource":   object["ComputeResource"],
			"connection_string":  object["ConnectionString"],
			"create_time":        object["CreateTime"],
			"db_cluster_id":      object["DBClusterId"],
			"db_cluster_version": object["DBVersion"],
			"engine":             object["Engine"],
			"expire_time":        object["ExpireTime"],
			"expired":            object["Expired"],
			"lock_mode":          object["LockMode"],
			"lock_reason":        object["LockReason"],
			"port":               object["Port"],
			"resource_group_id":  object["ResourceGroupId"],
			"status":             object["DBClusterStatus"],
			"storage_resource":   object["StorageResource"],
			"vpc_id":             object["VPCId"],
			"vswitch_id":         object["VSwitchId"],
			"zone_id":            object["ZoneId"],
			"payment_type":       convertAdbDbClusterLakeVersionPaymentTypeResponse(object["PayType"]),
		}

		ids = append(ids, fmt.Sprint(object["DBClusterId"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DBClusterId"])
		object, err = adbService.DescribeAdbDbClusterLakeVersion(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["engine_version"] = object["EngineVersion"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}
	return nil
}
