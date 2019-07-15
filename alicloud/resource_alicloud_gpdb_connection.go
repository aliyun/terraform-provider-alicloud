package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudGpdbConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGpdbConnectionCreate,
		Read:   resourceAlicloudGpdbConnectionRead,
		Update: resourceAlicloudGpdbConnectionUpdate,
		Delete: resourceAlicloudGpdbConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: validateDBConnectionPrefix,
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
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

func resourceAlicloudGpdbConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%s-tf", instanceId)
	}

	request := gpdb.CreateAllocateInstancePublicConnectionRequest()
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)

	err := resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))

	if err := gpdbService.WaitForGpdbConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := gpdbService.WaitForGpdbInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudGpdbConnectionRead(d, meta)
}

func resourceAlicloudGpdbConnectionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}
	object, err := gpdbService.DescribeGpdbConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceAlicloudGpdbConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") {
		client := meta.(*connectivity.AliyunClient)
		gpdbService := GpdbService{client}

		request := gpdb.CreateModifyDBInstanceConnectionStringRequest()
		request.DBInstanceId = parts[0]
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
				return gpdbClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if IsExceptedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		if err := gpdbService.WaitForGpdbInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudGpdbConnectionRead(d, meta)
}

func resourceAlicloudGpdbConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := gpdb.CreateReleaseInstancePublicConnectionRequest()
	request.DBInstanceId = parts[0]

	client := meta.(*connectivity.AliyunClient)
	gpdbService := GpdbService{client}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := gpdbService.DescribeGpdbConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString

		var raw interface{}
		raw, err = client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ReleaseInstancePublicConnection(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidGpdbNameNotFound, InvalidGpdbInstanceIdNotFound, InvalidCurrentConnectionStringNotFound, AtLeastOneNetTypeExists}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return gpdbService.WaitForGpdbConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
