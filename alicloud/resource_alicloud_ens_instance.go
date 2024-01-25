// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
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
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: IntBetween(0, 100),
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_use_coupon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"carrier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Operator, required for regional scheduling. Optional values:-cmcc (mobile)-unicom-telecom"),
			},
			"data_disk": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency", "cloud_ssd", "local_hdd", "local_ssd"}, true),
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
				Computed: true,
				ForceNew: true,
			},
			"force_stop": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Whether to force the identity when operating the instance. Optional values:-true: Force-false (default): non-mandatory"),
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"include_data_disks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_charge_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The instance billing policy. Optional values:-instance: instance granularity (the subscription method does not support instance)-user: user Dimension (user is not transmitted or supported in the prepaid mode)"),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^.{2,128}$"), "The instance name. Example value: test-InstanceName. It must be 2 to 128 characters in length and must start with an uppercase or lowercase letter or a Chinese character. It cannot start with http:// or https. Can contain Chinese, English, numbers, half-width colons (:), underscores (_), periods (.), or hyphens (-)The default value is the InstanceId of the instance."),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Instance bandwidth billing method. If the billing method can be selected for the first purchase, the subsequent value of this field will be processed by default according to the billing method selected for the first time. Optional values:-BandwidthByDay: Daily peak bandwidth-95bandwidthbymonth: 95 peak bandwidth"),
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "The IP type. Value:-ipv4 (default):IPv4-ipv6:IPv6-ipv4Andipv6:IPv4 and IPv6"),
			},
			"net_district_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_work_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, true),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip_identification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"schedule_area_level": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Scheduling level, through which node-level scheduling or area scheduling is performed. Optional values:-Node-level scheduling: Region-Regional scheduling: Big (region),Middle (province),Small (city)"),
			},
			"scheduling_price_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduling_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Scheduling policy. Optional values:-Concentrate for node-level scheduling-For regional scheduling, Concentrate, Disperse"),
			},
			"security_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
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
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "RunInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

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
	if v := d.Get("system_disk"); !IsNil(v) {
		nodeNative2, _ := jsonpath.Get("$[0].size", d.Get("system_disk"))
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["Size"] = nodeNative2
		}
		request["SystemDisk"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
	}

	if v, ok := d.GetOk("instance_charge_strategy"); ok {
		request["InstanceChargeStrategy"] = v
	}
	request["InstanceChargeType"] = convertEnsInstanceInstanceChargeTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("net_work_id"); ok {
		request["NetWorkId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("security_id"); ok {
		request["SecurityId"] = v
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = v
	}
	if v, ok := d.GetOk("billing_cycle"); ok {
		request["BillingCycle"] = v
	}
	if v, ok := d.GetOk("ip_type"); ok {
		request["IpType"] = v
	}
	if v, ok := d.GetOk("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsInstanceUpdate(d, meta)
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
	d.Set("payment_type", convertEnsInstanceInstancesInstanceInstanceResourceTypeResponse(objectRaw["InstanceResourceType"]))
	d.Set("status", objectRaw["Status"])

	networkAttributes1RawObj, _ := jsonpath.Get("$.NetworkAttributes", objectRaw)
	networkAttributes1Raw := make(map[string]interface{})
	if networkAttributes1RawObj != nil {
		networkAttributes1Raw = networkAttributes1RawObj.(map[string]interface{})
	}
	d.Set("net_work_id", networkAttributes1Raw["NetworkId"])
	d.Set("vswitch_id", networkAttributes1Raw["VSwitchId"])

	privateIpAddress1RawObj, _ := jsonpath.Get("$.PrivateIpAddresses.PrivateIpAddress[*]", objectRaw)
	privateIpAddress1Raw := make([]interface{}, 0)
	if privateIpAddress1RawObj != nil {
		privateIpAddress1Raw = privateIpAddress1RawObj.([]interface{})
	}

	if len(privateIpAddress1Raw) > 0 {
		d.Set("private_ip_address", privateIpAddress1Raw[0].(map[string]interface{})["Ip"])
	}

	securityGroupIds1Raw, _ := jsonpath.Get("$.SecurityGroupIds.SecurityGroupId", objectRaw)
	if len(securityGroupIds1Raw.([]interface{})) > 0 {
		d.Set("security_id", securityGroupIds1Raw.([]interface{})[0])
	}

	dataDisk3Raw, _ := jsonpath.Get("$.DataDisk.DataDisk", objectRaw)
	dataDiskMaps := make([]map[string]interface{}, 0)
	if dataDisk3Raw != nil {
		for _, dataDiskChild1Raw := range dataDisk3Raw.([]interface{}) {
			dataDiskMap := make(map[string]interface{})
			dataDiskChild1Raw := dataDiskChild1Raw.(map[string]interface{})
			dataDiskMap["category"] = dataDiskChild1Raw["Category"]
			dataDiskMap["size"] = dataDiskChild1Raw["DiskSize"]

			size, _ := dataDiskChild1Raw["Size"].(json.Number).Int64()
			dataDiskMap["size"] = size / 1024
			dataDiskMaps = append(dataDiskMaps, dataDiskMap)
		}
	}
	d.Set("data_disk", dataDiskMaps)
	systemDiskMaps := make([]map[string]interface{}, 0)
	systemDiskMap := make(map[string]interface{})
	systemDisk1Raw := make(map[string]interface{})
	if objectRaw["SystemDisk"] != nil {
		systemDisk1Raw = objectRaw["SystemDisk"].(map[string]interface{})
	}
	if len(systemDisk1Raw) > 0 {
		systemDiskMap["category"] = systemDisk1Raw["Category"]
		size, _ := systemDisk1Raw["Size"].(json.Number).Int64()
		systemDiskMap["size"] = size / 1024
		systemDiskMaps = append(systemDiskMaps, systemDiskMap)
	}
	d.Set("system_disk", systemDiskMaps)
	return nil
}

func resourceAliCloudEnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstanceAttribute"
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}

	if !d.IsNewResource() && d.HasChange("host_name") {
		update = true
		request["HostName"] = d.Get("host_name")
	}

	if v, ok := d.GetOk("user_data"); ok {
		request["UserData"] = v
	}
	if !d.IsNewResource() && d.HasChange("password") {
		update = true
		request["Password"] = d.Get("password")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("host_name"))}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "HostName", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_name")
		d.SetPartial("host_name")
		d.SetPartial("password")
	}
	update = false
	action = "ModifyPrepayInstanceSpec"
	conn, err = client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_type") && (d.Get("payment_type").(string) == "Subscription") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_type"))}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "SpecName", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_type")
	}
	update = false
	action = "ModifyPostPaidInstanceSpec"
	conn, err = client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("instance_type") && (d.Get("payment_type").(string) == "PayAsYouGo") {
		update = true
	}
	request["InstanceType"] = d.Get("instance_type")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("instance_type"))}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "SpecName", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("instance_type")
	}
	update = false
	action = "ModifyInstanceChargeType"
	conn, err = client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["InstanceChargeType"] = convertEnsInstanceInstanceChargeTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("include_data_disks"); ok {
		request["IncludeDataDisks"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("payment_type"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "InstanceResourceType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("payment_type")
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		ensServiceV2 := EnsServiceV2{client}
		object, err := ensServiceV2.DescribeEnsInstance(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Stopped" {
				action = "StopInstance"
				conn, err = client.NewEnsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["InstanceId"] = d.Id()
				if v, ok := d.GetOk("force_stop"); ok {
					request["ForceStop"] = v
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Running" {
				action = "StartInstance"
				conn, err = client.NewEnsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["InstanceId"] = d.Id()
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ensServiceV2.EnsInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	d.Partial(false)
	return resourceAliCloudEnsInstanceRead(d, meta)
}

func resourceAliCloudEnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		client := meta.(*connectivity.AliyunClient)
		action := "ReleasePrePaidInstance"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		conn, err := client.NewEnsClient()
		if err != nil {
			return WrapError(err)
		}
		request = make(map[string]interface{})
		query["InstanceId"] = d.Id()

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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
		query := make(map[string]interface{})
		conn, err := client.NewEnsClient()
		if err != nil {
			return WrapError(err)
		}
		request = make(map[string]interface{})
		query["InstanceId"] = d.Id()

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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

func convertEnsInstanceInstancesInstanceInstanceResourceTypeResponse(source interface{}) interface{} {
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
func convertEnsInstanceInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
