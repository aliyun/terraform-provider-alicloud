package adb

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	adbSdk "github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAdbConnection() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAdbConnectionCreate,
		Read:   ResourceAdbConnectionRead,
		Delete: ResourceAdbConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
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

func ResourceAdbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	dbClusterId := d.Get("db_cluster_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", dbClusterId)
	}

	request := adbSdk.CreateAllocateClusterPublicConnectionRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = dbClusterId
	request.ConnectionStringPrefix = prefix
	var raw interface{}
	var err error
	err = retry.Retry(8*time.Minute, func() *retry.RetryError {
		raw, err = client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
			return adbClient.AllocateClusterPublicConnection(request)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_adb_connection", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", dbClusterId, COLON_SEPARATED, request.ConnectionStringPrefix))

	if err := adbService.WaitForAdbConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	// wait instance running after allocating
	if err := adbService.WaitForCluster(dbClusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}

	return ResourceAdbConnectionRead(d, meta)
}

func ResourceAdbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbConnection(d.Id())

	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("db_cluster_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func ResourceAdbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := adbSdk.CreateReleaseClusterPublicConnectionRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = split[0]

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var raw interface{}
		raw, err := client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
			return adbClient.ReleaseClusterPublicConnection(request)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return adbService.WaitForAdbConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
