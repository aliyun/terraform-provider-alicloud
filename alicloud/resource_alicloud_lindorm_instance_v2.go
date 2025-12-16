// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudLindormInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudLindormInstanceV2Create,
		Read:   resourceAliCloudLindormInstanceV2Read,
		Update: resourceAliCloudLindormInstanceV2Update,
		Delete: resourceAliCloudLindormInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(101 * time.Minute),
			Update: schema.DefaultTimeout(1001 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arbiter_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arbiter_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arch_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_renew_duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renewal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cloud_storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cloud_storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"engine_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"connect_address_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"is_last_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"node_group": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_disk_size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"enable_attach_local_disk": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu_core_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"node_count": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"node_spec": {
										Type:     schema.TypeString,
										Required: true,
									},
									"node_disk_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"memory_size_gi_b": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"spec_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_group_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"instance_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"standby_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"white_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_list": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudLindormInstanceV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLindormV2Instance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("engine_list"); ok {
		engineListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["node_group"]
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["ResourceGroupName"] = dataLoop1Tmp["resource_group_name"]
				dataLoop1Map["NodeSpec"] = dataLoop1Tmp["node_spec"]
				dataLoop1Map["NodeDiskSize"] = dataLoop1Tmp["node_disk_size"]
				dataLoop1Map["NodeCount"] = dataLoop1Tmp["node_count"]
				if dataLoop1Tmp["node_disk_type"] != "" {
					dataLoop1Map["NodeDiskType"] = dataLoop1Tmp["node_disk_type"]
				}
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["NodeGroupList"] = localMaps
			dataLoopMap["EngineType"] = dataLoopTmp["engine_type"]
			engineListMapsArray = append(engineListMapsArray, dataLoopMap)
		}
		request["EngineList"] = engineListMapsArray
	}

	if v, ok := d.GetOk("arbiter_zone_id"); ok {
		request["ArbiterZoneId"] = v
	}
	if v, ok := d.GetOk("standby_vswitch_id"); ok {
		request["StandbyVSwitchId"] = v
	}
	if v, ok := d.GetOk("standby_zone_id"); ok {
		request["StandbyZoneId"] = v
	}
	if v, ok := d.GetOk("cloud_storage_type"); ok {
		request["CloudStorageType"] = v
	}
	request["PayType"] = d.Get("payment_type")
	request["ArchVersion"] = d.Get("arch_version")
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	request["InstanceAlias"] = d.Get("instance_alias")
	if v, ok := d.GetOk("arbiter_vswitch_id"); ok {
		request["ArbiterVSwitchId"] = v
	}
	if v, ok := d.GetOkExists("cloud_storage_size"); ok {
		request["CloudStorageSize"] = v
	}
	if v, ok := d.GetOk("primary_zone_id"); ok {
		request["PrimaryZoneId"] = v
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}
	if v, ok := d.GetOkExists("auto_renewal"); ok {
		request["AutoRenewal"] = v
	}
	request["VPCId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("primary_vswitch_id"); ok {
		request["PrimaryVSwitchId"] = v
	}
	request["VSwitchId"] = d.Get("vswitch_id")
	request["ZoneId"] = d.Get("zone_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_lindorm_instance_v2", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	lindormServiceV2 := LindormServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, lindormServiceV2.LindormInstanceV2StateRefreshFunc(d.Id(), "$.InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudLindormInstanceV2Update(d, meta)
}

func resourceAliCloudLindormInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	lindormServiceV2 := LindormServiceV2{client}

	objectRaw, err := lindormServiceV2.DescribeLindormInstanceV2(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_lindorm_instance_v2 DescribeLindormInstanceV2 Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("arbiter_vswitch_id", objectRaw["ArbiterVSwitchId"])
	d.Set("arbiter_zone_id", objectRaw["ArbiterZoneId"])
	d.Set("cloud_storage_size", objectRaw["CloudStorageSize"])
	d.Set("cloud_storage_type", objectRaw["DiskCategory"])
	d.Set("deletion_protection", objectRaw["DeletionProtection"])
	d.Set("instance_alias", objectRaw["InstanceAlias"])
	d.Set("payment_type", objectRaw["PayType"])
	d.Set("primary_vswitch_id", objectRaw["PrimaryVSwitchId"])
	d.Set("primary_zone_id", objectRaw["PrimaryZoneId"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("standby_vswitch_id", objectRaw["StandbyVSwitchId"])
	d.Set("standby_zone_id", objectRaw["StandbyZoneId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("vswitch_id", objectRaw["VswitchId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("auto_renewal", objectRaw["AutoRenew"])

	engineListRaw := objectRaw["EngineList"]
	engineListMaps := make([]map[string]interface{}, 0)
	if engineListRaw != nil {
		for _, engineListChildRaw := range convertToInterfaceArray(engineListRaw) {
			engineListMap := make(map[string]interface{})
			engineListChildRaw := engineListChildRaw.(map[string]interface{})
			engineListMap["engine_type"] = engineListChildRaw["Engine"]
			engineListMap["is_last_version"] = engineListChildRaw["IsLastVersion"]
			engineListMap["latest_version"] = engineListChildRaw["LatestVersion"]
			engineListMap["version"] = engineListChildRaw["Version"]

			connectAddressListRaw := engineListChildRaw["ConnectAddressList"]
			connectAddressListMaps := make([]map[string]interface{}, 0)
			if connectAddressListRaw != nil {
				for _, connectAddressListChildRaw := range convertToInterfaceArray(connectAddressListRaw) {
					connectAddressListMap := make(map[string]interface{})
					connectAddressListChildRaw := connectAddressListChildRaw.(map[string]interface{})
					connectAddressListMap["address"] = connectAddressListChildRaw["Address"]
					connectAddressListMap["port"] = connectAddressListChildRaw["Port"]
					connectAddressListMap["type"] = connectAddressListChildRaw["Type"]

					connectAddressListMaps = append(connectAddressListMaps, connectAddressListMap)
				}
			}
			engineListMap["connect_address_list"] = connectAddressListMaps
			nodeGroupRaw := engineListChildRaw["NodeGroup"]
			nodeGroupMaps := make([]map[string]interface{}, 0)
			if nodeGroupRaw != nil {
				for _, nodeGroupChildRaw := range convertToInterfaceArray(nodeGroupRaw) {
					nodeGroupMap := make(map[string]interface{})
					nodeGroupChildRaw := nodeGroupChildRaw.(map[string]interface{})
					nodeGroupMap["category"] = nodeGroupChildRaw["Category"]
					nodeGroupMap["cpu_core_count"] = nodeGroupChildRaw["CpuCoreCount"]
					nodeGroupMap["enable_attach_local_disk"] = nodeGroupChildRaw["EnableAttachLocalDisk"]
					nodeGroupMap["memory_size_gi_b"] = nodeGroupChildRaw["MemorySizeGiB"]
					nodeGroupMap["node_count"] = nodeGroupChildRaw["Quantity"]
					nodeGroupMap["node_disk_size"] = nodeGroupChildRaw["LocalDiskCapacity"]
					nodeGroupMap["node_disk_type"] = nodeGroupChildRaw["LocalDiskCategory"]
					nodeGroupMap["node_spec"] = nodeGroupChildRaw["NodeSpec"]
					nodeGroupMap["resource_group_name"] = nodeGroupChildRaw["ResourceGroupName"]
					nodeGroupMap["spec_id"] = nodeGroupChildRaw["SpecId"]
					nodeGroupMap["status"] = nodeGroupChildRaw["Status"]

					nodeGroupMaps = append(nodeGroupMaps, nodeGroupMap)
				}
			}
			engineListMap["node_group"] = nodeGroupMaps
			engineListMaps = append(engineListMaps, engineListMap)
		}
	}
	if err := d.Set("engine_list", engineListMaps); err != nil {
		return err
	}
	whiteIpListRaw := objectRaw["WhiteIpList"]
	whiteIpListMaps := make([]map[string]interface{}, 0)
	if whiteIpListRaw != nil {
		for _, whiteIpListChildRaw := range convertToInterfaceArray(whiteIpListRaw) {
			whiteIpListMap := make(map[string]interface{})
			whiteIpListChildRaw := whiteIpListChildRaw.(map[string]interface{})
			whiteIpListMap["group_name"] = whiteIpListChildRaw["GroupName"]
			whiteIpListMap["ip_list"] = whiteIpListChildRaw["IpList"]

			whiteIpListMaps = append(whiteIpListMaps, whiteIpListMap)
		}
	}
	if err := d.Set("white_ip_list", whiteIpListMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudLindormInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateLindormInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("instance_alias") {
		update = true
	}
	request["InstanceAlias"] = d.Get("instance_alias")
	if d.HasChange("deletion_protection") {
		update = true
		request["DeletionProtection"] = d.Get("deletion_protection")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
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
	}
	update = false
	action = "UpdateLindormV2Instance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("engine_list") {
		update = true
	}
	if v, ok := d.GetOk("engine_list"); ok || d.HasChange("engine_list") {
		engineListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["node_group"]
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["NodeDiskSize"] = dataLoop1Tmp["node_disk_size"]
				dataLoop1Map["ResourceGroupName"] = dataLoop1Tmp["resource_group_name"]
				dataLoop1Map["GroupId"] = dataLoop1Tmp["spec_id"]
				dataLoop1Map["NodeCount"] = dataLoop1Tmp["node_count"]
				if dataLoop1Tmp["node_disk_type"] != "" {
					dataLoop1Map["NodeDiskType"] = dataLoop1Tmp["node_disk_type"]
				}
				dataLoop1Map["NodeSpec"] = dataLoop1Tmp["node_spec"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["NodeGroupList"] = localMaps
			dataLoopMap["EngineType"] = dataLoopTmp["engine_type"]
			engineListMapsArray = append(engineListMapsArray, dataLoopMap)
		}
		request["EngineList"] = engineListMapsArray
	}

	if !d.IsNewResource() && d.HasChange("cloud_storage_type") {
		update = true
		request["CloudStorageType"] = d.Get("cloud_storage_type")
	}

	if !d.IsNewResource() && d.HasChange("cloud_storage_size") {
		update = true
		request["CloudStorageSize"] = d.Get("cloud_storage_size")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"LindormErrorCode.OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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
		lindormServiceV2 := LindormServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, lindormServiceV2.LindormInstanceV2StateRefreshFunc(d.Id(), "$.InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstancePayType"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["Duration"] = d.Get("duration")
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["PayType"] = d.Get("payment_type")
	request["PricingCycle"] = d.Get("pricing_cycle")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
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
		lindormServiceV2 := LindormServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, lindormServiceV2.LindormInstanceV2StateRefreshFunc(d.Id(), "$.InstanceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLindormV2WhiteIpList"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("white_ip_list") {
		update = true
	}
	if v, ok := d.GetOk("white_ip_list"); ok || d.HasChange("white_ip_list") {
		whiteIpGroupListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["WhiteIpList"] = dataLoopTmp["ip_list"]
			dataLoopMap["GroupName"] = dataLoopTmp["group_name"]
			whiteIpGroupListMapsArray = append(whiteIpGroupListMapsArray, dataLoopMap)
		}
		request["WhiteIpGroupList"] = whiteIpGroupListMapsArray
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
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
	}

	d.Partial(false)
	return resourceAliCloudLindormInstanceV2Read(d, meta)
}

func resourceAliCloudLindormInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseLindormV2Instance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["Immediately"] = "true"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"Instance.IsDeleted"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	lindormServiceV2 := LindormServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 10*time.Minute, lindormServiceV2.LindormInstanceV2StateRefreshFunc(d.Id(), "$.InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
