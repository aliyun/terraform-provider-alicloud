package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDrdsPolardbxInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDrdsPolardbxInstanceCreate,
		Read:   resourceAliCloudDrdsPolardbxInstanceRead,
		Update: resourceAliCloudDrdsPolardbxInstanceUpdate,
		Delete: resourceAliCloudDrdsPolardbxInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(61 * time.Minute),
			Update: schema.DefaultTimeout(61 * time.Minute),
			Delete: schema.DefaultTimeout(61 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cn_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cn_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dn_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dn_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"dn_storage_space": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"5.7", "8.0"}, false),
			},
			"is_read_db_instance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"primary_db_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"secondary_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"specified_dn_scale": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"specified_dn_spec_map_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"custom_local_ssd", "cloud_auto"}, false),
			},
			"switch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"switch_time_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tertiary_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"topology_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"1azone", "3azones"}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudDrdsPolardbxInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tertiary_zone"); ok {
		request["TertiaryZone"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["PayType"] = "POSTPAY"
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["VSwitchId"] = d.Get("vswitch_id")
	request["EngineVersion"] = d.Get("engine_version")
	if request["EngineVersion"] == "" {
		request["EngineVersion"] = "5.7"
	}
	request["NetworkType"] = "vpc"
	request["ZoneId"] = "null"
	request["PrimaryZone"] = d.Get("primary_zone")
	request["CnClass"] = d.Get("cn_class")
	request["TopologyType"] = d.Get("topology_type")
	request["DnClass"] = d.Get("dn_class")
	if v, ok := d.GetOk("secondary_zone"); ok {
		request["SecondaryZone"] = v
	}
	request["CNNodeCount"] = d.Get("cn_node_count")
	request["VPCId"] = d.Get("vpc_id")
	request["DNNodeCount"] = d.Get("dn_node_count")
	if v, ok := d.GetOkExists("is_read_db_instance"); ok {
		request["IsReadDBInstance"] = v
	}
	if v, ok := d.GetOk("primary_db_instance_name"); ok {
		request["PrimaryDBInstanceName"] = v
	}
	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v
	}
	if v, ok := d.GetOk("dn_storage_space"); ok {
		request["DnStorageSpace"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_drds_polardbx_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceName"]))

	drdsServiceV2 := DrdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, drdsServiceV2.DrdsPolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDrdsPolardbxInstanceUpdate(d, meta)
}

func resourceAliCloudDrdsPolardbxInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	drdsServiceV2 := DrdsServiceV2{client}

	objectRaw, err := drdsServiceV2.DescribeDrdsPolardbxInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_drds_polardbx_instance DescribeDrdsPolardbxInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cn_class", objectRaw["CnNodeClassCode"])
	d.Set("cn_node_count", objectRaw["CnNodeCount"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dn_class", objectRaw["DnNodeClassCode"])
	d.Set("dn_node_count", objectRaw["DnNodeCount"])
	d.Set("dn_storage_space", objectRaw["DnStorageSpace"])
	d.Set("engine_version", objectRaw["EngineVersion"])
	d.Set("primary_zone", objectRaw["PrimaryZone"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("secondary_zone", objectRaw["SecondaryZone"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("tertiary_zone", objectRaw["TertiaryZone"])
	d.Set("topology_type", objectRaw["TopologyType"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VPCId"])

	return nil
}

func resourceAliCloudDrdsPolardbxInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdatePolarDBXInstanceNode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("cn_node_count") {
		update = true
	}
	request["CNNodeCount"] = d.Get("cn_node_count")
	if !d.IsNewResource() && d.HasChange("dn_node_count") {
		update = true
	}
	request["DNNodeCount"] = d.Get("dn_node_count")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		drdsServiceV2 := DrdsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, drdsServiceV2.DrdsPolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyDBInstanceClass"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("cn_class") {
		update = true
	}
	request["CnClass"] = d.Get("cn_class")
	if !d.IsNewResource() && d.HasChange("dn_class") {
		update = true
	}
	request["DnClass"] = d.Get("dn_class")
	if !d.IsNewResource() && d.HasChange("dn_storage_space") {
		update = true
		request["DnStorageSpace"] = d.Get("dn_storage_space")
	}
	if d.HasChange("specified_dn_scale") {
		if v, ok := d.GetOkExists("specified_dn_scale"); ok {
			request["SpecifiedDNScale"] = v
		}
	}
	if d.HasChange("specified_dn_spec_map_json") {
		if v, ok := d.GetOk("specified_dn_spec_map_json"); ok {
			request["SpecifiedDNSpecMapJson"] = v
		}
	}
	if d.HasChange("switch_time_mode") {
		if v, ok := d.GetOk("switch_time_mode"); ok {
			request["SwitchTimeMode"] = v
		}
	}
	if d.HasChange("switch_time") {
		if v, ok := d.GetOk("switch_time"); ok {
			request["SwitchTime"] = v
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
		drdsServiceV2 := DrdsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, drdsServiceV2.DrdsPolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "PolarDBXInstance"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
	action = "ModifyDBInstanceDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["DBInstanceDescription"] = d.Get("description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
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
	return resourceAliCloudDrdsPolardbxInstanceRead(d, meta)
}

func resourceAliCloudDrdsPolardbxInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"DBInstance.InOrder"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"DBInstance.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	drdsServiceV2 := DrdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 60*time.Second, drdsServiceV2.DescribeAsyncDrdsPolardbxInstanceStateRefreshFunc(d, response, "$.Items", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

func convertDrdsPolardbxInstanceDBInstancePayTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
