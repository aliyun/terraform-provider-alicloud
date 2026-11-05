package cassandra

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCassandraDataCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCassandraDataCenterCreate,
		Read:   resourceAlicloudCassandraDataCenterRead,
		Update: resourceAlicloudCassandraDataCenterUpdate,
		Delete: resourceAlicloudCassandraDataCenterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		DeprecationMessage: "This resource has been deprecated since v1.220.0 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"auto_renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 12}),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_center_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_center_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
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
			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudCassandraDataCenterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := Service{client}

	request := cassandra.CreateCreateDataCenterRequest()
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request.AutoRenewPeriod = requests.NewInteger(v.(int))
	}
	request.ClusterId = d.Get("cluster_id").(string)
	if v, ok := d.GetOk("data_center_name"); ok {
		request.DataCenterName = v.(string)
	}
	if v, ok := d.GetOk("disk_size"); ok {
		request.DiskSize = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request.DiskType = v.(string)
	}
	request.InstanceType = d.Get("instance_type").(string)
	request.NodeCount = requests.NewInteger(d.Get("node_count").(int))
	request.PayType = d.Get("pay_type").(string)
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = v.(string)
	}
	request.RegionId = client.RegionId
	request.VswitchId = d.Get("vswitch_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
	vswitchId := helper.Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		request.VswitchId = vswitchId
		// check vswitchId in zone
		vsw, err := cassandraService.describeVSwitch(vswitchId)
		if err != nil {
			return errs.WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		}
		request.VpcId = vsw.VpcId
	}
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.CreateDataCenter(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cassandra_data_center", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw)
	response, _ := raw.(*cassandra.CreateDataCenterResponse)
	d.SetId(fmt.Sprintf("%v:%v", response.DataCenterId, request.ClusterId))

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 8*time.Minute, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudCassandraDataCenterUpdate(d, meta)
}
func resourceAlicloudCassandraDataCenterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := Service{client}
	object, err := cassandraService.DescribeCassandraDataCenter(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cassandra_data_center cassandraService.DescribeCassandraDataCenter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("cluster_id", object.ClusterId)
	d.Set("data_center_id", object.DataCenterId)
	d.Set("data_center_name", object.DataCenterName)
	d.Set("disk_size", object.DiskSize)
	d.Set("disk_type", object.DiskType)
	d.Set("instance_type", object.InstanceType)
	d.Set("node_count", object.NodeCount)
	d.Set("pay_type", object.PayType)
	d.Set("status", object.Status)
	d.Set("vswitch_id", object.VswitchId)
	d.Set("zone_id", object.ZoneId)
	d.Set("auto_renew", object.AutoRenewal)
	if d.Get("auto_renew").(bool) {
		d.Set("auto_renew_period", object.AutoRenewPeriod)
	} else {
		if v, ok := d.GetOkExists("auto_renew_period"); ok {
			d.Set("auto_renew_period", v.(int))
		}
	}

	d.Set("enable_public", false)
	d.Set("public_points", nil)
	endPoints, err := cassandraService.DescribeCassandraEndpoints(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	for _, contactPoint := range endPoints.ContactPoints.ContactPoint {
		if object.ZoneId == contactPoint.DataCenterId && len(contactPoint.PublicAddresses.PublicAddress) > 0 {
			d.Set("enable_public", true)
			d.Set("public_points", contactPoint.PublicAddresses.PublicAddress)
			break
		}
	}
	return nil
}
func resourceAlicloudCassandraDataCenterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cassandraService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)
	if (d.IsNewResource() && d.Get("enable_public") != nil && d.Get("enable_public").(bool)) || (!d.IsNewResource() && d.HasChange("enable_public")) {
		if d.Get("enable_public").(bool) {
			request := cassandra.CreateAllocatePublicContactPointsRequest()
			request.ClusterId = parts[1]
			request.DataCenterId = parts[0]
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.AllocatePublicContactPoints(request)
			})
			logs.AddDebug(request.GetActionName(), raw)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"AllocatePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		} else {
			request := cassandra.CreateReleasePublicContactPointsRequest()
			request.ClusterId = parts[1]
			request.DataCenterId = parts[0]
			raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
				return cassandraClient.ReleasePublicContactPoints(request)
			})
			logs.AddDebug(request.GetActionName(), raw)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ReleasePubConnFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudCassandraDataCenterRead(d, meta)
	}

	if d.HasChange("data_center_name") {
		request := cassandra.CreateModifyDataCenterRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.DataCenterName = d.Get("data_center_name").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyDataCenter(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("instance_type") {
		request := cassandra.CreateModifyInstanceTypeRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.InstanceType = d.Get("instance_type").(string)
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ModifyInstanceType(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ModifyInstanceTypeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	if d.HasChange("disk_size") {
		request := cassandra.CreateResizeDiskSizeRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.DiskSize = requests.NewInteger(d.Get("disk_size").(int))
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeDiskSize(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ResizeDiskSizeFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	if d.HasChange("node_count") {
		request := cassandra.CreateResizeNodeCountRequest()
		request.ClusterId = parts[1]
		request.DataCenterId = parts[0]
		request.NodeCount = requests.NewInteger(d.Get("node_count").(int))
		raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.ResizeNodeCount(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraDataCenterStateRefreshFunc(d.Id(), []string{"ResizeNodeCountFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAlicloudCassandraDataCenterRead(d, meta)
}
func resourceAlicloudCassandraDataCenterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := cassandra.CreateDeleteDataCenterRequest()
	request.ClusterId = parts[1]
	request.DataCenterId = parts[0]
	raw, err := client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DeleteDataCenter(request)
	})
	logs.AddDebug(request.GetActionName(), raw)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DataCenter.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	cassandraService := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cassandraService.CassandraClusterStateRefreshFunc(parts[1], []string{"DeleteOtherDataCenterFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, parts[1])
	}
	return nil
}
