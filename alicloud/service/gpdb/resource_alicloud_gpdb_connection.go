package gpdb

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.ValidateDBConnectionPort,
				Default:      "3306",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%s-tf", instanceId)
	}

	request := gpdb.CreateAllocateInstancePublicConnectionRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)

	err := retry.Retry(8*time.Minute, func() *retry.RetryError {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_connection", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))
	// wait instance running after allocating
	stateConf := helper.BuildStateConf([]string{"Creating", "NetAddressCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbService.InstanceStateRefreshFunc(instanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceConnectionRead(d, meta)
}

func resourceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	object, err := gpdbService.DescribeConnection(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	if d.HasChange("port") {
		client := meta.(*connectivity.AliyunClient)
		gpdbService := NewService(client)

		request := gpdb.CreateModifyDBInstanceConnectionStringRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = parts[0]
		object, err := gpdbService.DescribeConnection(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		if err := retry.Retry(8*time.Minute, func() *retry.RetryError {
			raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
				return gpdbClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		stateConf := helper.BuildStateConf([]string{"NET_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, gpdbService.InstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return resourceConnectionRead(d, meta)
}

func resourceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := gpdb.CreateReleaseInstancePublicConnectionRequest()
	request.DBInstanceId = parts[0]

	client := meta.(*connectivity.AliyunClient)
	gpdbService := NewService(client)
	request.RegionId = client.RegionId
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		object, err := gpdbService.DescribeConnection(d.Id())
		if err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString

		var raw interface{}
		raw, err = client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ReleaseInstancePublicConnection(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"NetAddressDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbService.InstanceStateRefreshFunc(request.DBInstanceId, []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return errs.WrapError(gpdbService.WaitForConnection(d.Id(), Deleted, DefaultTimeoutMedium))
}
