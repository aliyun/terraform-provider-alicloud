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
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"carrier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"cmcc", "unicom", "telecom"}, false),
			},
			"data_disk": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "local_hdd", "local_ssd"}, false),
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntBetween(20, 32000),
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
				Computed: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"user", "instance"}, false),
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"BandwidthByDay", "95BandwidthByMonth", "PayByBandwidth4thMonth"}, false),
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"net_district_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"100102", "100106"}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"password_inherit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12}),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Month", "Day"}, false),
			},
			"public_ip_identification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"quantity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_area_level": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Big", "Middle", "Small", "Region"}, false),
			},
			"scheduling_price_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PriceHighPriority", "PriceLowPriority"}, false),
			},
			"scheduling_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Concentrate", "Disperse"}, false),
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
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntBetween(20, 32000),
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
	if v, ok := d.GetOkExists("auto_renew"); ok {
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
	if v, ok := d.GetOk("scheduling_price_strategy"); ok {
		request["SchedulingPriceStrategy"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("password_inherit"); ok {
		request["PasswordInherit"] = v
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("image_id"); ok {
		request["ImageId"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDiskMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
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
	request["InstanceChargeType"] = convertEnsInstanceChargeTypeRequest(d.Get("payment_type").(string))
	request["Amount"] = "1"
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

	id, _ := jsonpath.Get("$.InstanceIds[0]", response)
	d.SetId(fmt.Sprint(id))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
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
	d.Set("payment_type", convertEnsInstancesInstanceInstanceResourceTypeResponse(objectRaw["InstanceResourceType"]))
	d.Set("status", objectRaw["Status"])
	dataDisk3Raw, _ := jsonpath.Get("$.DataDisk.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDisk3Raw != nil {
		for _, dataDiskChild1Raw := range dataDisk3Raw.([]interface{}) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChild1Raw := dataDiskChild1Raw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChild1Raw["Category"]
			size, _ := dataDiskChild1Raw["Size"].(json.Number).Int64()
			dataDiskMap["size"] = size / 1024
			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	d.Set("data_disk", dataDiskMaps)

	systemDiskRaw, _ := jsonpath.Get("$.SystemDisk", objectRaw)
	systemDiskMaps := make([]map[string]interface{}, 0)
	if systemDiskRaw != nil {
		systemDiskMap := make(map[string]interface{})
		systemDiskChild1Raw := systemDiskRaw.(map[string]interface{})
		size, _ := systemDiskChild1Raw["Size"].(json.Number).Int64()
		systemDiskMap["size"] = size / 1024
		systemDiskMaps = append(systemDiskMaps, systemDiskMap)
	}
	d.Set("system_disk", systemDiskMaps)

	return nil
}

func resourceAliCloudEnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Instance.")
	return nil
}

func resourceAliCloudEnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
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

	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PayAsYouGo" {
		client := meta.(*connectivity.AliyunClient)
		action := "ReleasePostPaidInstance"
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
	return nil
}

func convertEnsInstancesInstanceInstanceResourceTypeResponse(source interface{}) interface{} {
	switch source {
	case "EnsInstance":
		return "Subscription"
	case "EnsService":
		return "PayAsYouGo"
	case "EnsPostPaidInstance":
		return "PayAsYouGo"
	case "EnsPostInstance":
		return "PayAsYouGo"
	case "BuildMachine":
		return "PayAsYouGo"
	}
	return source
}
func convertEnsInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
