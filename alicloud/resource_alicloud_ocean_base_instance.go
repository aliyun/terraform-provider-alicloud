package alicloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOceanBaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOceanBaseInstanceCreate,
		Read:   resourceAliCloudOceanBaseInstanceRead,
		Update: resourceAliCloudOceanBaseInstanceUpdate,
		Delete: resourceAliCloudOceanBaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"backup_retain_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"receive_all", "delete_all", "receive_last"}, false),
			},
			"commodity_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_arch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ARM", "X86"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(100, 10000),
			},
			"disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"cloud_essd_pl1", "cloud_essd_pl0"}, false),
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new != "" && old != new {
						oldInstanceClass := strings.TrimSuffix(strings.TrimSpace(old), "B")
						newInstanceClass := strings.TrimSuffix(strings.TrimSpace(new), "B")

						return reflect.DeepEqual(oldInstanceClass, newInstanceClass)
					}
					return false
				},
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ob_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc:     IntInSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60}),
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc:     StringInSlice([]string{"Month", "Hour"}, false),
			},
			"primary_instance": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"series": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"normal", "normal_ssd", "history", "normal_kv", "normal_hg"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_spec_native": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zones": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOceanBaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["DiskSize"] = d.Get("disk_size")
	request["Series"] = d.Get("series")
	request["InstanceClass"] = d.Get("instance_class")
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ChargeType"] = convertOceanBaseInstanceChargeTypeRequest(d.Get("payment_type").(string))
	jsonPathResult10, err := jsonpath.Get("$", d.Get("zones"))
	if err == nil {
		request["Zones"] = convertListToCommaSeparate(jsonPathResult10.(*schema.Set).List())
	}

	if v, ok := d.GetOk("ob_version"); ok {
		request["ObVersion"] = v
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request["DiskType"] = v
	}
	if v, ok := d.GetOk("primary_instance"); ok {
		request["PrimaryInstance"] = v
	}
	if v, ok := d.GetOk("primary_region"); ok {
		request["PrimaryRegion"] = v
	}
	if v, ok := d.GetOk("cpu_arch"); ok {
		request["CpuArch"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ocean_base_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	oceanBaseServiceV2 := OceanBaseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOceanBaseInstanceUpdate(d, meta)
}

func resourceAliCloudOceanBaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseServiceV2 := OceanBaseServiceV2{client}

	objectRaw, err := oceanBaseServiceV2.DescribeOceanBaseInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ocean_base_instance DescribeOceanBaseInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CpuArchitecture"] != nil {
		d.Set("cpu_arch", objectRaw["CpuArchitecture"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["DiskType"] != nil {
		d.Set("disk_type", objectRaw["DiskType"])
	}
	if objectRaw["InstanceClass"] != nil {
		d.Set("instance_class", objectRaw["InstanceClass"])
	}
	if objectRaw["InstanceName"] != nil {
		d.Set("instance_name", objectRaw["InstanceName"])
	}
	if objectRaw["NodeNum"] != nil {
		d.Set("node_num", objectRaw["NodeNum"])
	}
	if objectRaw["Version"] != nil {
		d.Set("ob_version", objectRaw["Version"])
	}
	if convertOceanBaseInstanceInstancePayTypeResponse(objectRaw["PayType"]) != nil {
		d.Set("payment_type", convertOceanBaseInstanceInstancePayTypeResponse(objectRaw["PayType"]))
	}
	if objectRaw["PrimaryInstance"] != nil {
		d.Set("primary_instance", objectRaw["PrimaryInstance"])
	}
	if objectRaw["PrimaryRegion"] != nil {
		d.Set("primary_region", objectRaw["PrimaryRegion"])
	}
	if convertOceanBaseInstanceInstanceSeriesResponse(objectRaw["Series"]) != nil {
		d.Set("series", convertOceanBaseInstanceInstanceSeriesResponse(objectRaw["Series"]))
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["AutoRenewal"] != nil {
		d.Set("auto_renew", objectRaw["AutoRenewal"])
	}

	zones1Raw := make([]interface{}, 0)
	if objectRaw["Zones"] != nil {
		zones1Raw = objectRaw["Zones"].([]interface{})
	}

	d.Set("zones", zones1Raw)

	objectRaw, err = oceanBaseServiceV2.DescribeDescribeInstances(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["CommodityCode"] != nil {
		d.Set("commodity_code", objectRaw["CommodityCode"])
	}
	if objectRaw["Cpu"] != nil {
		d.Set("cpu", formatInt(objectRaw["Cpu"]))
	}
	if objectRaw["DiskSize"] != nil {
		d.Set("disk_size", objectRaw["DiskSize"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}

	return nil
}

func resourceAliCloudOceanBaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstanceName"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	request["InstanceName"] = d.Get("instance_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, false)
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
	}
	update = false
	action = "ModifyInstanceNodeNum"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	oceanBaseServiceV2 := OceanBaseServiceV2{client}
	objectRaw, err := oceanBaseServiceV2.DescribeOceanBaseInstance(d.Id())
	if d.HasChange("node_num") && fmt.Sprint(objectRaw["NodeNum"]) != d.Get("node_num") {
		update = true
	}
	request["NodeNum"] = d.Get("node_num")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) || NeedRetry(err) {
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
		oceanBaseServiceV2 := OceanBaseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ModifyInstanceSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
	}
	request["InstanceClass"] = d.Get("instance_class")
	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")
	if v, ok := d.GetOkExists("upgrade_spec_native"); ok {
		request["UpgradeSpecNative"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Instance.Order.CreateFailed"}) || NeedRetry(err) {
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
		oceanBaseServiceV2 := OceanBaseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ONLINE"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudOceanBaseInstanceRead(d, meta)
}

func resourceAliCloudOceanBaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resource alicloud_ocean_base_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceIds"] = fmt.Sprintf("[\"%s\"]", d.Id())

	if v, ok := d.GetOk("backup_retain_mode"); ok {
		request["BackupRetainMode"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("OceanBasePro", "2019-09-01", action, query, request, false)

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
		if IsExpectedErrors(err, []string{"UnknownError", "IllegalOperation.Resource"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	oceanBaseServiceV2 := OceanBaseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 20*time.Second, oceanBaseServiceV2.OceanBaseInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertOceanBaseInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertOceanBaseInstanceInstancePayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	case "PrePay":
		return "Subscription"
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertOceanBaseInstanceInstanceSeriesResponse(source interface{}) interface{} {
	switch source {
	case "NORMAL":
		return "normal"
	case "NORMAL_SSD":
		return "normal_ssd"
	case "HISTORY":
		return "history"
	}
	return source
}

func convertOceanBaseInstanceChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	case "Subscription":
		return "PrePay"
	}
	return source
}

func convertOceanBaseInstanceSeriesResponse(source interface{}) interface{} {
	switch source {
	case "NORMAL":
		return "normal"
	case "NORMAL_SSD":
		return "normal_ssd"
	case "HISTORY":
		return "history"
	}
	return source
}
