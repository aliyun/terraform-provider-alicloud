package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAdbCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbClusterCreate,
		Read:   resourceAlicloudAdbClusterRead,
		Update: resourceAlicloudAdbClusterUpdate,
		Delete: resourceAlicloudAdbClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(72 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_version": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"3.0"}, false),
				Optional:     true,
				Default:      "3.0",
			},
			"db_cluster_category": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "Cluster"}, false),
				Required:     true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"db_node_storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pay_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNotRenewal,
				ValidateFunc: validation.StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: adbPostPaidAndRenewDiffSuppressFunc,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: adbPostPaidDiffSuppressFunc,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAdbClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	request, err := buildAdbCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.CreateDBCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adb.CreateDBClusterResponse)
	d.SetId(response.DBClusterId)

	// wait cluster status change from Creating to running
	stateConf := BuildStateConf([]string{"Preparing", "Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 15*time.Minute, adbService.AdbClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAdbClusterUpdate(d, meta)
}

func resourceAlicloudAdbClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	d.Partial(true)

	if err := adbService.setClusterTags(d); err != nil {
		return WrapError(err)
	}

	if (d.Get("pay_type").(string) == string(PrePaid)) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		request := adb.CreateModifyAutoRenewAttributeRequest()
		request.DBClusterId = d.Id()
		request.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			request.Duration = strconv.Itoa(period)
			request.PeriodUnit = string(Month)
			if period > 9 {
				request.Duration = strconv.Itoa(period / 12)
				request.PeriodUnit = string(Year)
			}
		}

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ModifyAutoRenewAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("maintain_time") {
		request := adb.CreateModifyDBClusterMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ModifyDBClusterMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_time")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudAdbClusterRead(d, meta)
	}

	if d.HasChange("description") {
		request := adb.CreateModifyDBClusterDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.DBClusterDescription = d.Get("description").(string)

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ModifyDBClusterDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := adbService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ips")
	}

	if d.HasChange("db_node_class") || d.HasChange("db_node_count") || d.HasChange("db_node_storage") {
		request := adb.CreateModifyDBClusterRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.DBNodeClass = d.Get("db_node_class").(string)
		request.DBNodeStorage = strconv.Itoa(d.Get("db_node_storage").(int))
		request.DBNodeGroupCount = strconv.Itoa(d.Get("db_node_count").(int))

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ModifyDBCluster(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		// wait cluster status change from ClassChanging to Running
		stateConf := BuildStateConf([]string{"ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Minute, adbService.AdbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("db_node_class")
		d.SetPartial("db_node_count")
		d.SetPartial("db_node_storage")
	}

	d.Partial(false)
	return resourceAlicloudAdbClusterRead(d, meta)
}

func resourceAlicloudAdbClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	clusterAttribute, err := adbService.DescribeAdbClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	cluster, err := adbService.DescribeAdbCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	ips, err := adbService.DescribeDBSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", ips)

	d.Set("vswitch_id", clusterAttribute.VSwitchId)
	d.Set("pay_type", getChargeType(clusterAttribute.PayType))
	d.Set("id", clusterAttribute.DBClusterId)
	d.Set("description", clusterAttribute.DBClusterDescription)
	d.Set("db_version", clusterAttribute.DBVersion)
	d.Set("maintain_time", clusterAttribute.MaintainTime)
	d.Set("zone_id", clusterAttribute.ZoneId)
	d.Set("db_node_class", cluster.DBNodeClass)
	d.Set("db_node_count", cluster.DBNodeCount)
	d.Set("db_node_storage", cluster.DBNodeStorage)
	d.Set("db_cluster_category", cluster.Category)
	d.Set("db_cluster_version", cluster.DBVersion)
	tags, err := adbService.DescribeTags(d.Id(), "cluster")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", adbService.tagsToMap(tags))

	if clusterAttribute.PayType == string(Prepaid) {
		clusterAutoRenew, err := adbService.DescribeAdbAutoRenewAttribute(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		renewPeriod := 1
		if clusterAutoRenew != nil {
			renewPeriod = clusterAutoRenew.Duration
		}
		if clusterAutoRenew != nil && clusterAutoRenew.PeriodUnit == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
		period, err := computePeriodByUnit(clusterAttribute.CreationTime, clusterAttribute.ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
		d.Set("renewal_status", clusterAutoRenew.RenewalStatus)
	}

	return nil
}

func resourceAlicloudAdbClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	cluster, err := adbService.DescribeAdbClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	// Pre paid cluster can not be release.
	if PayType(cluster.PayType) == Prepaid {
		return nil
	}

	request := adb.CreateDeleteDBClusterRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Id()
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.DeleteDBCluster(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterStatus", "OperationDenied.InvalidAdbClusterStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Creating", "Running", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, adbService.AdbClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildAdbCreateRequest(d *schema.ResourceData, meta interface{}) (*adb.CreateDBClusterRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := adb.CreateCreateDBClusterRequest()
	request.RegionId = string(client.Region)
	request.DBClusterVersion = Trim(d.Get("db_cluster_version").(string))
	request.DBClusterCategory = Trim(d.Get("db_cluster_category").(string))
	request.DBClusterClass = d.Get("db_node_class").(string)
	request.DBNodeGroupCount = strconv.Itoa(d.Get("db_node_count").(int))
	request.DBNodeStorage = strconv.Itoa(d.Get("db_node_storage").(int))
	request.DBClusterDescription = d.Get("description").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.DBClusterNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VPCId = vsw.VpcId
	}

	payType := Trim(d.Get("pay_type").(string))
	request.PayType = string(Postpaid)
	if payType == string(PrePaid) {
		request.PayType = string(Prepaid)
	}
	if PayType(request.PayType) == Prepaid {
		period := d.Get("period").(int)
		request.UsedTime = strconv.Itoa(period)
		request.Period = string(Month)
		if period > 9 {
			request.UsedTime = strconv.Itoa(period / 12)
			request.Period = string(Year)
		}
	}

	return request, nil
}
