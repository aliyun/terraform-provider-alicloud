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

func resourceAliCloudClickHouseEnterpriseDBCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDBClusterCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDBClusterRead,
		Update: resourceAliCloudClickHouseEnterpriseDBClusterUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDBClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
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
			"region_id": {
				Type:     schema.TypeString,
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

func resourceAliCloudClickHouseEnterpriseDBClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("scale_min"); ok {
		request["ScaleMin"] = v
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
	if v, ok := d.GetOk("multi_zones"); ok {
		multiZoneMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VSwitchIds"] = dataLoopTmp["vswitch_ids"].(*schema.Set).List()
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			multiZoneMapsArray = append(multiZoneMapsArray, dataLoopMap)
		}
		multiZoneMapsJson, err := json.Marshal(multiZoneMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["MultiZone"] = string(multiZoneMapsJson)
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
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION", "ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDBClusterStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudClickHouseEnterpriseDBClusterRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDBClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDBCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster DescribeClickHouseEnterpriseDBCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("scale_max", objectRaw["ScaleMax"])
	d.Set("scale_min", objectRaw["ScaleMin"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("zone_id", objectRaw["ZoneId"])

	multiZonesRaw := objectRaw["MultiZones"]
	multiZonesMaps := make([]map[string]interface{}, 0)
	if multiZonesRaw != nil {
		for _, multiZonesChildRaw := range multiZonesRaw.([]interface{}) {
			multiZonesMap := make(map[string]interface{})
			multiZonesChildRaw := multiZonesChildRaw.(map[string]interface{})
			multiZonesMap["zone_id"] = multiZonesChildRaw["ZoneId"]

			vSwitchIdsRaw := make([]interface{}, 0)
			if multiZonesChildRaw["VSwitchIds"] != nil {
				vSwitchIdsRaw = multiZonesChildRaw["VSwitchIds"].([]interface{})
			}

			multiZonesMap["vswitch_ids"] = vSwitchIdsRaw
			multiZonesMaps = append(multiZonesMaps, multiZonesMap)
		}
	}
	if err := d.Set("multi_zones", multiZonesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudClickHouseEnterpriseDBClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyDBInstanceClass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("scale_min") {
		update = true
		request["ScaleMin"] = d.Get("scale_min")
	}

	if d.HasChange("scale_max") {
		update = true
		request["ScaleMax"] = d.Get("scale_max")
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
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickHouseServiceV2.ClickHouseEnterpriseDBClusterStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudClickHouseEnterpriseDBClusterRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDBClusterDelete(d *schema.ResourceData, meta interface{}) error {

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
