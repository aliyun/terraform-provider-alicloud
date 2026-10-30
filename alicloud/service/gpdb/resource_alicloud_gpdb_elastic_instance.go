package gpdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceElasticInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticInstanceCreate,
		Read:   resourceElasticInstanceRead,
		Update: resourceElasticInstanceUpdate,
		Delete: resourceElasticInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		DeprecationMessage: "This resource has been deprecated since v1.147.0 and will be removed in the future. Please use 'alicloud_gpdb_instance' instead.",
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"gpdb"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"6.0"}, false),
			},
			"seg_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_essd", "cloud_efficiency"}, false),
			},
			"seg_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 4 || v > 256 || v%4 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 0 and 256 inclusive, and multiple of 4, got: %d", key, v))
					}
					return
				},
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 50 || v > 4000 || v%50 != 0 {
						errs = append(errs, fmt.Errorf("%q must be between 50 and 4000 inclusive, and multiple of 50, got: %d", key, v))
					}
					return
				},
			},
			"instance_spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2C16G", "4C32G", "16C128G", "2C8G", "4C16G", "8C32G", "16C64G"}, false),
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "HighAvailability"}, false),
			},
			"instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
				Default:      "VPC",
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Default:      "PayAsYouGo",
			},
			"payment_duration_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"payment_duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 12),
				DiffSuppressFunc: postPaidDiffSuppressFunc,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceElasticInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	var response map[string]interface{}
	action := "CreateECSDBInstance"
	request := make(map[string]interface{})
	var err error

	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["SegStorageType"] = d.Get("seg_storage_type")
	request["SegNodeNum"] = d.Get("seg_node_num")
	request["StorageSize"] = d.Get("storage_size")
	request["InstanceSpec"] = d.Get("instance_spec")
	request["PayType"] = convertInstancePaymentTypeRequest(d.Get("payment_type").(string))
	if request["PayType"].(string) == "Prepaid" {
		request["Period"] = d.Get("payment_duration_unit")
		paymentDuration := d.Get("payment_duration").(int)
		request["UsedTime"] = strconv.Itoa(paymentDuration)
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("db_instance_category"); ok {
		request["DBInstanceCategory"] = v
	}
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}

	vswitchId := helper.Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vsw, err := describeVSwitchWithTeadsl(client, vswitchId)
		if err != nil {
			return errs.WrapError(err)
		}
		request["VPCId"] = vsw["VpcId"]
		request["VSwitchId"] = vswitchId
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	request["ClientToken"] = helper.BuildClientToken("CreateECSDBInstance")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_elastic_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBInstanceId"]))
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, gpdbService.ElasticInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceElasticInstanceUpdate(d, meta)
}

func resourceElasticInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	instance, err := gpdbService.DescribeElasticInstance(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("engine", instance["Engine"])
	d.Set("engine_version", instance["EngineVersion"])
	d.Set("seg_storage_type", instance["StorageType"])
	d.Set("seg_node_num", instance["SegNodeNum"])
	d.Set("storage_size", instance["StorageSize"])
	d.Set("status", instance["DBInstanceStatus"])
	d.Set("db_instance_description", instance["DBInstanceDescription"])
	d.Set("instance_network_type", instance["InstanceNetworkType"])
	d.Set("vswitch_id", instance["VSwitchId"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("connection_string", instance["ConnectionString"])
	d.Set("port", instance["Port"])
	d.Set("db_instance_category", instance["DBInstanceCategory"])
	d.Set("encryption_key", instance["EncryptionKey"])
	d.Set("encryption_type", instance["EncryptionType"])
	d.Set("tags", tags.TagsToMap(instance["Tags"]))
	cpuCores, err := convert.FormatInt(instance["CpuCores"])
	if err != nil {
		return err
	}
	memorySize, err := convert.FormatInt(instance["MemorySize"])
	if err != nil {
		return err
	}
	d.Set("instance_spec", fmt.Sprintf("%dC%dG", cpuCores, memorySize))

	if v, exist := instance["PayType"]; exist {
		d.Set("payment_type", convertInstancePaymentTypeResponse(v.(string)))
	}

	securityIps, err := gpdbService.DescribeSecurityIps(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("security_ip_list", securityIps)
	return nil
}

func resourceElasticInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	d.Partial(true)
	var err error
	if d.HasChange("tags") {
		if err := gpdbService.SetResourceTags(d, "ALIYUN::GPDB::INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}
	if d.HasChange("db_instance_description") {
		action := "ModifyDBInstanceDescription"
		request := map[string]interface{}{
			"RegionId":              client.RegionId,
			"DBInstanceId":          d.Id(),
			"DBInstanceDescription": d.Get("db_instance_description"),
			"SourceIp":              client.SourceIp,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("security_ip_list") {
		ipList := convert.ExpandStringList(d.Get("security_ip_list").([]interface{}))
		ipStr := strings.Join(ipList[:], helper.CommaSeparated)
		if ipStr == "" {
			ipStr = LOCAL_HOST_IP
		}
		if err := gpdbService.ModifySecurityIps(d.Id(), ipStr); err != nil {
			return errs.WrapError(err)
		}
	}

	d.Partial(false)
	return resourceElasticInstanceRead(d, meta)
}

func resourceElasticInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBInstance"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DBInstanceId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectDBState"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstancePayType"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertInstancePaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertInstancePaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertDBInstanceClassToInstanceSpec(instanceClass string) string {
	splitClass := strings.Split(instanceClass, ".")
	return strings.ToUpper(splitClass[len(splitClass)-1])
}

// postPaidDiffSuppressFunc is a 1:1 local copy of the root alicloud package's
// postPaidDiffSuppressFunc, migrated with the elastic-instance resource per the
// subpackage convention: resource-specific diff-suppress lives package-private.
func postPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}
