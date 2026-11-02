package adb

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAdbDbCluster() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAdbDbClusterCreate,
		Read:   ResourceAdbDbClusterRead,
		Update: ResourceAdbDbClusterUpdate,
		Delete: ResourceAdbDbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validators.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: adbPostPaidAndRenewDiffSuppressFunc,
			},
			"compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("mode"); ok && v.(string) == "reserver" {
						return true
					}
					return false
				},
			},
			"db_cluster_category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"Basic", "Cluster", "MixedStorage"}, false),
			},
			"db_cluster_class": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "It duplicates with attribute db_node_class and is deprecated from 1.121.2.",
			},
			"db_cluster_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"3.0"}, false),
				Default:      "3.0",
			},
			"kernel_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"switch_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 1}),
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"db_node_storage": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringLenBetween(2, 256),
			},
			"elastic_io_resource": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("mode"); ok && v.(string) == "reserver" {
						return true
					}
					return false
				},
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"reserver", "flexible"}, false),
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				ConflictsWith: []string{"pay_type"},
			},
			"pay_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"payment_type"},
				ValidateFunc:  validators.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Deprecated:    "Attribute 'pay_type' has been deprecated from the provider version 1.166.0 and it will be remove in the future version. Please use the new attribute 'payment_type' instead.",
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validators.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
				Optional:         true,
			},
			"renewal_status": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validators.StringInSlice([]string{"AutoRenewal", "Normal", "NotRenewal"}, false),
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"elastic_io_resource_size": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"8Core64GB", "12Core96GB", "16Core128GB"}, false),
			},
			"disk_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"disk_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAdbDbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	action := "CreateDBCluster"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("compute_resource"); ok {
		request["ComputeResource"] = v
	}

	request["DBClusterCategory"] = d.Get("db_cluster_category")
	if v, ok := d.GetOk("db_node_class"); ok {
		request["DBClusterClass"] = v
	} else if v, ok := d.GetOk("db_cluster_class"); ok {
		request["DBClusterClass"] = v
	}

	request["DBClusterVersion"] = d.Get("db_cluster_version")
	if v, ok := d.GetOk("db_node_count"); ok {
		request["DBNodeGroupCount"] = v
	}

	if v, ok := d.GetOk("db_node_storage"); ok {
		request["DBNodeStorage"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["DBClusterDescription"] = v
	}

	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertAdbDBClusterPaymentTypeRequest(v.(string))
		if request["PayType"] != string(Postpaid) {
			request["PayType"] = string(Prepaid)
			period := d.Get("period").(int)
			request["UsedTime"] = strconv.Itoa(period)
			request["Period"] = string(Month)
			if period > 9 {
				request["UsedTime"] = strconv.Itoa(period / 12)
				request["Period"] = string(Year)
			}
		}
	} else if v, ok := d.GetOk("pay_type"); ok {
		request["PayType"] = convertAdbDbClusterDBClusterPayTypeRequest(v.(string))
		if request["PayType"] != string(Postpaid) {
			request["PayType"] = string(Prepaid)
			period := d.Get("period").(int)
			request["UsedTime"] = strconv.Itoa(period)
			request["Period"] = string(Month)
			if period > 9 {
				request["UsedTime"] = strconv.Itoa(period / 12)
				request["Period"] = string(Year)
			}
		}
	} else {
		request["PayType"] = "Postpaid"
	}

	if v, ok := d.GetOkExists("elastic_io_resource"); ok {
		request["ElasticIOResource"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["DBClusterNetworkType"] = "VPC"
		request["VSwitchId"] = v
	}

	if (request["ZoneId"] == nil || request["VPCId"] == nil) && request["VSwitchId"] != nil {
		vsw, err := describeVSwitch(client, request["VSwitchId"].(string))
		if err != nil {
			return errs.WrapError(err)
		}
		if request["VPCId"] == nil {
			request["VPCId"] = vsw["VpcId"]
		}
		if request["ZoneId"] == nil {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}

	if v, ok := d.GetOkExists("disk_encryption"); ok {
		request["DiskEncryption"] = v
	}

	if v, ok := d.GetOk("kms_id"); ok {
		request["KmsId"] = v
	}

	if v, ok := d.GetOkExists("enable_ssl"); ok {
		request["EnableSSL"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("CreateDBCluster")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_adb_db_cluster", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBClusterId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 300*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAdbDbClusterUpdate(d, meta)
}

func ResourceAdbDbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbDbCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_analyticdb_for_mysql3.0_db_cluster adbService.DescribeAdbDbCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("compute_resource", object["ComputeResource"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("port", object["Port"])
	d.Set("db_cluster_category", convertAdbDBClusterCategoryResponse(object["Category"].(string)))
	d.Set("db_node_class", object["DBNodeClass"])
	d.Set("db_node_count", object["DBNodeCount"])
	d.Set("db_node_storage", object["DBNodeStorage"])
	d.Set("description", object["DBClusterDescription"])
	elastic_io_resourceVal, err := convert.FormatInt(object["ElasticIOResource"])
	if err != nil {
		return err
	}
	d.Set("elastic_io_resource", elastic_io_resourceVal)
	d.Set("maintain_time", object["MaintainTime"])
	d.Set("mode", object["Mode"])
	d.Set("payment_type", convertAdbDBClusterPaymentTypeResponse(object["PayType"].(string)))
	d.Set("pay_type", convertAdbDbClusterDBClusterPayTypeResponse(object["PayType"].(string)))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["DBClusterStatus"])
	d.Set("tags", tags.TagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("elastic_io_resource_size", object["ElasticIOResourceSize"])
	d.Set("disk_performance_level", object["DiskPerformanceLevel"])
	d.Set("db_cluster_version", object["DBVersion"])
	d.Set("disk_encryption", object["DiskEncryption"])
	d.Set("kms_id", object["KmsId"])

	if object["PayType"].(string) == string(Prepaid) {
		describeAutoRenewAttributeObject, err := adbService.DescribeAutoRenewAttribute(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		renewPeriod := 1
		if describeAutoRenewAttributeObject != nil {
			renewPeriodVal, err := convert.FormatInt(describeAutoRenewAttributeObject["Duration"])
			if err != nil {
				return err
			}
			renewPeriod = renewPeriodVal
		}
		if describeAutoRenewAttributeObject != nil && describeAutoRenewAttributeObject["PeriodUnit"] == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
		//period, err := computePeriodByUnit(object["CreationTime"], object["ExpireTime"], d.Get("period").(int), "Month")
		//if err != nil {
		//	return errs.WrapError(err)
		//}
		//d.Set("period", period)
		d.Set("renewal_status", describeAutoRenewAttributeObject["RenewalStatus"])
	}

	describeDBClusterAccessWhiteListObject, err := adbService.DescribeDBClusterAccessWhiteList(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("security_ips", strings.Split(describeDBClusterAccessWhiteListObject["SecurityIPList"].(string), ","))

	sslObject, err := adbService.DescribeAdbDbClusterSSL(d.Id())
	if err != nil {
		log.Printf("[WARN] Resource alicloud_adb_db_cluster DescribeAdbDbClusterSSL Failed!!! %s", err)
	} else {
		d.Set("enable_ssl", sslObject["SSLEnabled"])
	}
	kernelObject, err := adbService.DescribeAdbDbClusterKernelVersion(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("kernel_version", kernelObject["KernelVersion"])

	return nil
}

func ResourceAdbDbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := adbService.SetResourceTags(d, "ALIYUN::ADB::CLUSTER"); err != nil {
			return errs.WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["DBClusterDescription"] = d.Get("description")
		action := "ModifyDBClusterDescription"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	if d.HasChange("maintain_time") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["MaintainTime"] = d.Get("maintain_time")
		action := "ModifyDBClusterMaintainTime"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		request["NewResourceGroupId"] = d.Get("resource_group_id")
		action := "ModifyDBClusterResourceGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	update := false
	request := map[string]interface{}{
		"DbClusterId": d.Id(),
	}
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && (d.HasChange("pay_type") || d.HasChange("payment_type")) {
		update = true
	}

	if paymentType, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertAdbDBClusterPaymentTypeRequest(paymentType.(string))
	} else if payType, ok := d.GetOk("pay_type"); ok {
		request["PayType"] = convertAdbDbClusterDBClusterPayTypeRequest(payType.(string))
	}
	if request["PayType"] == "Prepaid" {
		request["Period"] = d.Get("period")
		period := d.Get("period").(int)
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(Year)
		}
	}

	if update {
		action := "ModifyDBClusterPayType"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"Preparing", "ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	request = map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if (d.Get("pay_type").(string) == string(PrePaid) || d.Get("payment_type").(string) == "Subscription") && d.HasChange("auto_renew_period") {
		update = true
		if d.Get("renewal_status").(string) == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			request["Duration"] = strconv.Itoa(period)
			request["PeriodUnit"] = string(Month)
			if period > 9 {
				request["Duration"] = strconv.Itoa(period / 12)
				request["PeriodUnit"] = string(Year)
			}
		}
	}
	if (d.Get("pay_type").(string) == string(PrePaid) || d.Get("payment_type").(string) == "Subscription") && d.HasChange("renewal_status") {
		update = true
		request["RenewalStatus"] = d.Get("renewal_status")
	}
	if update {
		action := "ModifyAutoRenewAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"Preparing", "ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterAccessWhiteListReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if d.HasChange("security_ips") {
		update = true
	}
	modifyDBClusterAccessWhiteListReq["SecurityIps"] = convert.ConvertListToCommaSeparate(d.Get("security_ips").(*schema.Set).List())
	if update {
		action := "ModifyDBClusterAccessWhiteList"
		if modifyDBClusterAccessWhiteListReq["SecurityIps"].(string) == "" {
			modifyDBClusterAccessWhiteListReq["SecurityIps"] = LOCAL_HOST_IP
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterAccessWhiteListReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterAccessWhiteListReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	update = false
	modifyDBClusterReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("compute_resource") {
		update = true
		modifyDBClusterReq["ComputeResource"] = d.Get("compute_resource")
	}
	if !d.IsNewResource() && d.HasChange("db_cluster_category") {
		update = true
		modifyDBClusterReq["DBClusterCategory"] = d.Get("db_cluster_category")
	}
	if !d.IsNewResource() && d.HasChange("db_node_class") {
		update = true
		modifyDBClusterReq["DBNodeClass"] = d.Get("db_node_class")
	}
	if !d.IsNewResource() && d.HasChange("db_node_count") {
		update = true
		modifyDBClusterReq["DBNodeGroupCount"] = d.Get("db_node_count")
	}
	if !d.IsNewResource() && d.HasChange("elastic_io_resource") {
		update = true
		modifyDBClusterReq["ElasticIOResource"] = d.Get("elastic_io_resource")
	}

	object, err := adbService.DescribeAdbDbCluster(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	if d.HasChange("elastic_io_resource_size") {
		if v, ok := d.GetOk("elastic_io_resource_size"); ok && v.(string) != object["ElasticIOResourceSize"] {
			update = true

			modifyDBClusterReq["ElasticIOResourceSize"] = v
		}
	}

	modifyDBClusterReq["RegionId"] = client.RegionId
	if update {
		if _, ok := d.GetOk("mode"); ok {
			modifyDBClusterReq["Mode"] = d.Get("mode")
		}
		if _, ok := d.GetOk("modify_type"); ok {
			modifyDBClusterReq["ModifyType"] = d.Get("modify_type")
		}
		action := "ModifyDBCluster"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterReq, false)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OrderProcessing"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"Preparing", "ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

		if d.HasChange("compute_resource") {
			stateConf = helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("compute_resource"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "ComputeResource", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

	}

	update = false
	modifyDBClusterReq = map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("db_node_storage") {
		update = true
		modifyDBClusterReq["DBNodeStorage"] = d.Get("db_node_storage")
	}

	modifyDBClusterReq["RegionId"] = client.RegionId
	if update {
		if _, ok := d.GetOk("mode"); ok {
			modifyDBClusterReq["Mode"] = d.Get("mode")
		}
		if _, ok := d.GetOk("modify_type"); ok {
			modifyDBClusterReq["ModifyType"] = d.Get("modify_type")
		}
		action := "ModifyDBCluster"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterReq, false)
			if err != nil {
				// There is service bug and needs checking IncorrectDBInstanceState.
				// If the bug is fixed, the IncorrectDBInstanceState checking can be removed.
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OrderProcessing"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"Preparing", "ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterReq = map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if d.HasChange("disk_performance_level") {
		if v, ok := d.GetOk("disk_performance_level"); ok && v.(string) != object["DiskPerformanceLevel"] {
			update = true

			modifyDBClusterReq["DiskPerformanceLevel"] = v
		}
	}

	modifyDBClusterReq["RegionId"] = client.RegionId
	if update {
		if _, ok := d.GetOk("mode"); ok {
			modifyDBClusterReq["Mode"] = d.Get("mode")
		}
		if _, ok := d.GetOk("modify_type"); ok {
			modifyDBClusterReq["ModifyType"] = d.Get("modify_type")
		}
		action := "ModifyDBCluster"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterReq, false)
			if err != nil {
				// There is service bug and needs checking IncorrectDBInstanceState.
				// If the bug is fixed, the IncorrectDBInstanceState checking can be removed.
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OrderProcessing"}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{"Preparing", "ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	upgradeKernelVersionReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("kernel_version") {
		update = true
		upgradeKernelVersionReq["DBVersion"] = d.Get("kernel_version")
		if v, ok := d.GetOkExists("switch_mode"); ok {
			upgradeKernelVersionReq["SwitchMode"] = v
		}
	}

	if update {
		action := "UpgradeKernelVersion"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, upgradeKernelVersionReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, upgradeKernelVersionReq)
			return nil
		})

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	update = false
	modifyDBClusterSSLReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("enable_ssl") {
		update = true

		if v, ok := d.GetOkExists("enable_ssl"); ok {
			modifyDBClusterSSLReq["EnableSSL"] = v

			if fmt.Sprint(v) == "true" {
				modifyDBClusterSSLReq["ConnectionString"] = d.Get("connection_string")
			}
		}
	}

	if update {
		action := "ModifyDBClusterSSL"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterSSLReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterSSLReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	d.Partial(false)

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, adbService.AdbDbClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAdbDbClusterRead(d, meta)
}

func ResourceAdbDbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	action := "DeleteDBCluster"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}
	var taskId string
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OrderProcessing"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		taskId = response["TaskId"].(json.Number).String()
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{"Waiting", "Running", "Failed", "Retry", "Pause", "Stop"}, []string{"Finished", "Closed", "Cancel"}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, adbService.AdbTaskStateRefreshFunc(d.Id(), taskId))
	if _, err = stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertAdbDbClusterDBClusterPayTypeRequest(source string) string {
	switch source {
	case "PostPaid":
		return "Postpaid"
	case "PrePaid":
		return "Prepaid"
	}
	return source
}

func convertAdbDbClusterDBClusterPayTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PostPaid"
	case "Prepaid":
		return "PrePaid"
	}
	return source
}

func convertAdbDBClusterPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}

func convertAdbDBClusterPaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertAdbDBClusterCategoryResponse(source string) string {
	switch source {
	case "MIXED_STORAGE":
		return "MixedStorage"
	}
	return source
}
