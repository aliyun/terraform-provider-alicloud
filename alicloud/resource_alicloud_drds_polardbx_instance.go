// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
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
			Create: schema.DefaultTimeout(65 * time.Minute),
			Update: schema.DefaultTimeout(65 * time.Minute),
			Delete: schema.DefaultTimeout(65 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cn_class": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"polarx.x4.medium.2e", "polarx.x4.large.2e", "polarx.x8.large.2e", "polarx.x4.xlarge.2e", "polarx.x8.xlarge.2e", "polarx.x4.2xlarge.2e", "polarx.x8.2xlarge.2e", "polarx.x4.4xlarge.2e", "polarx.x8.4xlarge.2e", "polarx.st.8xlarge.2e", "polarx.st.12xlarge.2e"}, false),
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
			"dn_class": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"mysql.n4.medium.25", "mysql.n4.large.25", "mysql.x8.large.25", "mysql.n4.xlarge.25", "mysql.x8.xlarge.25", "mysql.n4.2xlarge.25", "mysql.x8.2xlarge.25", "mysql.x4.4xlarge.25", "mysql.x8.4xlarge.25", "mysql.st.8xlarge.25", "mysql.st.12xlarge.25"}, false),
			},
			"dn_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntAtLeast(2),
			},
			"primary_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["NetworkType"] = "vpc"
	request["VPCId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["EngineVersion"] = "5.7"
	request["TopologyType"] = d.Get("topology_type")
	request["PrimaryZone"] = d.Get("primary_zone")
	if v, ok := d.GetOk("secondary_zone"); ok {
		request["SecondaryZone"] = v
	}
	if v, ok := d.GetOk("tertiary_zone"); ok {
		request["TertiaryZone"] = v
	}
	request["CnClass"] = d.Get("cn_class")
	request["DnClass"] = d.Get("dn_class")
	request["CNNodeCount"] = d.Get("cn_node_count")
	request["DNNodeCount"] = d.Get("dn_node_count")
	request["PayType"] = "POSTPAY"
	request["ZoneId"] = "null"
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

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
	d.Set("dn_class", objectRaw["DnNodeClassCode"])
	d.Set("dn_node_count", objectRaw["DnNodeCount"])
	d.Set("primary_zone", objectRaw["PrimaryZone"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("secondary_zone", objectRaw["SecondaryZone"])
	d.Set("status", objectRaw["Status"])
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
	update := false
	d.Partial(true)
	action := "UpdatePolarDBXInstanceNode"
	var err error
	request = make(map[string]interface{})
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		drdsServiceV2 := DrdsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, drdsServiceV2.DrdsPolardbxInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cn_node_count")
		d.SetPartial("dn_node_count")
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "PolarDBXInstance"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardbx", "2020-02-02", action, nil, request, false)

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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	d.Partial(false)
	return resourceAliCloudDrdsPolardbxInstanceRead(d, meta)
}

func resourceAliCloudDrdsPolardbxInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["DBInstanceName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardbx", "2020-02-02", action, nil, request, false)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"DBInstance.InOrder"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"DBInstance.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	drdsServiceV2 := DrdsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Minute, drdsServiceV2.DrdsPolardbxInstanceDeleteJobStateRefreshFunc(d, response, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
