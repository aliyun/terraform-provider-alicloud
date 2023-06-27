// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsInstanceCreate,
		Read:   resourceAliCloudEnsInstanceRead,
		Update: resourceAliCloudEnsInstanceUpdate,
		Delete: resourceAliCloudEnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"auto_renew": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"True", "False"}, false),
			},
			"auto_renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 12),
			},
			"carrier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_district_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 9, 12}),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip_identification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"schedule_area_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduling_price_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduling_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"unique_suffix": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "RunInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOkExists("unique_suffix"); ok {
		request["UniqueSuffix"] = v
	}
	if v, ok := d.GetOk("user_data"); ok {
		request["UserData"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOkExists("public_ip_identification"); ok {
		request["PublicIpIdentification"] = v
	}
	if v, ok := d.GetOk("net_district_code"); ok {
		request["NetDistrictCode"] = v
	}
	if v, ok := d.GetOk("carrier"); ok {
		request["Carrier"] = v
	}
	request["ScheduleAreaLevel"] = d.Get("schedule_area_level")
	if v, ok := d.GetOk("scheduling_strategy"); ok {
		request["SchedulingStrategy"] = v
	}
	request["InternetMaxBandwidthOut"] = d.Get("internet_max_bandwidth_out")
	request["Amount"] = d.Get("amount")
	if v, ok := d.GetOk("scheduling_price_strategy"); ok {
		request["SchedulingPriceStrategy"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("password_inherit"); ok {
		request["PasswordInherit"] = v
	}
	request["InstanceChargeType"] = d.Get("instance_charge_type")
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDiskMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Size"] = dataLoopTmp["size"]
			dataLoopMap["Category"] = dataLoopTmp["category"]
			dataDiskMaps = append(dataDiskMaps, dataLoopMap)
		}
		dataDiskMapsJson, err := json.Marshal(dataDiskMaps)
		if err != nil {
			return WrapError(err)
		}
		request["DataDisk"] = string(dataDiskMapsJson)
	}

	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("system_disk"); ok {
		nodeNative2, _ := jsonpath.Get("$[0].size", v)
		if nodeNative2 != "" {
			objectDataLocalMap["Size"] = nodeNative2
		}
	}
	objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
	if err != nil {
		return WrapError(err)
	}
	request["SystemDisk"] = string(objectDataLocalMapJson)

	if v, ok := d.GetOk("instance_charge_strategy"); ok {
		request["InstanceChargeStrategy"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_instance", action, AlibabaCloudSdkGoERROR)
	}

	// 先用jsonpath取id
	id, _ := jsonpath.Get("$.InstanceIds[0]", response)
	d.SetId(fmt.Sprint(id))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsInstanceRead(d, meta)
}

func resourceAliCloudEnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_instance DescribeEnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("host_name", objectRaw["HostName"])
	d.Set("image_id", objectRaw["ImageId"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("instance_type", objectRaw["SpecName"])
	d.Set("internet_max_bandwidth_out", objectRaw["InternetMaxBandwidthOut"])
	d.Set("payment_type", objectRaw["InstanceResourceType"])
	d.Set("status", objectRaw["Status"])
	dataDisk3Raw, _ := jsonpath.Get("$.DataDisk.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDisk3Raw != nil {
		for _, dataDiskChild1Raw := range dataDisk3Raw.([]interface{}) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChild1Raw := dataDiskChild1Raw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChild1Raw["Category"]
			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	d.Set("data_disk", dataDiskMaps)

	e := jsonata.MustCompile("ApiOutput.Instances.Instance.DataDisk.DataDisk.Size/1024")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("data_disk", evaluation)
	e = jsonata.MustCompile("ApiOutput.Instances.Instance.SystemDisk.Size/1024")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("system_disk", evaluation)

	return nil
}

func resourceAliCloudEnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Instance.")
	return nil
}

func resourceAliCloudEnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ReleasePrePaidInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 50*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
