// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudClickHouseEnterpriseDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDbClusterCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDbClusterRead,
		Update: resourceAliCloudClickHouseEnterpriseDbClusterUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"computing_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disabled_ports": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"computing_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"engine_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multi_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"node_scale_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"node_scale_min": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"object_store_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scale_max": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scale_min": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_quota": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDbClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("node_scale_max"); ok {
		request["NodeScaleMax"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("scale_min"); ok {
		request["ScaleMin"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("node_count"); ok {
		request["NodeCount"] = v
	}
	if v, ok := d.GetOk("multi_zones"); ok {
		multiZoneMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["VSwitchIds"] = convertToInterfaceArray(dataLoop1Tmp["vswitch_ids"])
			dataLoop1Map["ZoneId"] = dataLoop1Tmp["zone_id"]
			multiZoneMapsArray = append(multiZoneMapsArray, dataLoop1Map)
		}
		multiZoneMapsJson, err := json.Marshal(multiZoneMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["MultiZone"] = string(multiZoneMapsJson)
	}

	if v, ok := d.GetOkExists("node_scale_min"); ok {
		request["NodeScaleMin"] = v
	}
	if v, ok := d.GetOk("scale_max"); ok {
		request["ScaleMax"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VswitchId"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.DBInstanceId", response)
	d.SetId(fmt.Sprint(id))

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION", "ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudClickHouseEnterpriseDbClusterUpdate(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDbCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster DescribeClickHouseEnterpriseDbCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("category", objectRaw["Category"])
	d.Set("charge_type", objectRaw["ChargeType"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("disabled_ports", objectRaw["DisabledPorts"])
	d.Set("engine_minor_version", objectRaw["EngineMinorVersion"])
	d.Set("engine_version", objectRaw["EngineVersion"])
	d.Set("maintain_end_time", objectRaw["MaintainEndTime"])
	d.Set("maintain_start_time", objectRaw["MaintainStartTime"])
	d.Set("node_count", formatInt(objectRaw["NodeCount"]))
	d.Set("node_scale_max", formatInt(objectRaw["NodeScaleMax"]))
	d.Set("node_scale_min", formatInt(objectRaw["NodeScaleMin"]))
	d.Set("object_store_size", objectRaw["ObjectStoreSize"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scale_max", objectRaw["ScaleMax"])
	d.Set("scale_min", objectRaw["ScaleMin"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_quota", objectRaw["StorageQuota"])
	d.Set("storage_size", objectRaw["StorageSize"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	multiZonesRaw := objectRaw["MultiZones"]
	multiZonesMaps := make([]map[string]interface{}, 0)
	if multiZonesRaw != nil {
		for _, multiZonesChildRaw := range convertToInterfaceArray(multiZonesRaw) {
			multiZonesMap := make(map[string]interface{})
			multiZonesChildRaw := multiZonesChildRaw.(map[string]interface{})
			multiZonesMap["zone_id"] = multiZonesChildRaw["ZoneId"]

			vSwitchIdsRaw := make([]interface{}, 0)
			if multiZonesChildRaw["VSwitchIds"] != nil {
				vSwitchIdsRaw = convertToInterfaceArray(multiZonesChildRaw["VSwitchIds"])
			}

			multiZonesMap["vswitch_ids"] = vSwitchIdsRaw
			multiZonesMaps = append(multiZonesMaps, multiZonesMap)
		}
	}
	if err := d.Set("multi_zones", multiZonesMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = clickHouseServiceV2.DescribeEnterpriseDbClusterDescribeEndpoints(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("instance_network_type", objectRaw["InstanceNetworkType"])

	endpointsRaw := objectRaw["Endpoints"]
	endpointsMaps := make([]map[string]interface{}, 0)
	if endpointsRaw != nil {
		for _, endpointsChildRaw := range convertToInterfaceArray(endpointsRaw) {
			endpointsMap := make(map[string]interface{})
			endpointsChildRaw := endpointsChildRaw.(map[string]interface{})
			endpointsMap["computing_group_id"] = endpointsChildRaw["ComputingGroupId"]
			endpointsMap["connection_string"] = endpointsChildRaw["ConnectionString"]
			endpointsMap["endpoint_name"] = endpointsChildRaw["EndpointName"]
			endpointsMap["ip_address"] = endpointsChildRaw["IPAddress"]
			endpointsMap["net_type"] = endpointsChildRaw["NetType"]
			endpointsMap["status"] = endpointsChildRaw["Status"]
			endpointsMap["vswitch_id"] = endpointsChildRaw["VSwitchId"]
			endpointsMap["vpc_id"] = endpointsChildRaw["VpcId"]
			endpointsMap["vpc_instance_id"] = endpointsChildRaw["VpcInstanceId"]

			portsRaw := endpointsChildRaw["Ports"]
			portsMaps := make([]map[string]interface{}, 0)
			if portsRaw != nil {
				for _, portsChildRaw := range convertToInterfaceArray(portsRaw) {
					portsMap := make(map[string]interface{})
					portsChildRaw := portsChildRaw.(map[string]interface{})
					portsMap["port"] = portsChildRaw["Port"]
					portsMap["protocol"] = portsChildRaw["Protocol"]

					portsMaps = append(portsMaps, portsMap)
				}
			}
			endpointsMap["ports"] = portsMaps
			endpointsMaps = append(endpointsMaps, endpointsMap)
		}
	}
	if err := d.Set("endpoints", endpointsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudClickHouseEnterpriseDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyDBInstanceClass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("node_scale_max") {
		update = true
		request["NodeScaleMax"] = d.Get("node_scale_max")
	}

	if !d.IsNewResource() && d.HasChange("node_scale_min") {
		update = true
		request["NodeScaleMin"] = d.Get("node_scale_min")
	}

	if !d.IsNewResource() && d.HasChange("scale_min") {
		update = true
		request["ScaleMin"] = d.Get("scale_min")
	}

	if !d.IsNewResource() && d.HasChange("scale_max") {
		update = true
		request["ScaleMax"] = d.Get("scale_max")
	}

	if !d.IsNewResource() && d.HasChange("node_count") {
		update = true
		request["NodeCount"] = d.Get("node_count")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "EnterpriseDBCluster"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDbClusterStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDBInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
	}
	request["AttributeValue"] = d.Get("description")
	request["AttributeType"] = "DBInstanceDescription"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("description"))}, d.Timeout(schema.TimeoutUpdate), 0, clickHouseServiceV2.ClickHouseEnterpriseDbClusterStateRefreshFunc(d.Id(), "Description", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		if err := clickHouseServiceV2.SetResourceTags(d, "EnterpriseDBCluster"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudClickHouseEnterpriseDbClusterRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDbClusterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
