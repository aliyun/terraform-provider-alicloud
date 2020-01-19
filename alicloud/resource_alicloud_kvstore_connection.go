package alicloud

import (
	"fmt"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r_kvstore"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"regexp"
)

const kvstoreConnectionSuffixRegex = "\\.redis\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const kvstoreConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + kvstoreConnectionSuffixRegex + "$"

var kvstoreConnectionIdWithSuffixRegexp = regexp.MustCompile(kvstoreConnectionIdWithSuffixRegex)

func resourceAlicloudKVstoreConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVstoreConnectionCreate,
		Read:   resourceAlicloudKVstoreConnectionRead,
		Update: resourceAlicloudKVstoreConnectionUpdate,
		Delete: resourceAlicloudKVstoreConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_string_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
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
		},
	}
}

func resourceAlicloudKVstoreConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_string_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}

	request := r_kvstore.CreateAllocateInstancePublicConnectionRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)
	var raw interface{}
	var err error
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))

	if err := kvstoreService.WaitForKVstorePublicConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance Normal after allocating
	if err := kvstoreService.WaitForKVstoreInstance(instanceId, Normal, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudKVstoreConnectionRead(d, meta)
}

func resourceAlicloudKVstoreConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := kvstoreConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstorePublicConnection(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("connection_string_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)

	return nil
}

func resourceAlicloudKVstoreConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	submatch := kvstoreConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") {
		request := r_kvstore.CreateModifyDBInstanceConnectionStringRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = parts[0]
		object, err := kvstoreService.DescribeKVstorePublicConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.NewConnectionString = parts[1]
		request.Port = d.Get("port").(string)
		request.IPType = "Public"

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{"IncorrectDBInstanceState"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance Normal after modifying
		if err := kvstoreService.WaitForKVstoreInstance(request.DBInstanceId, Normal, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudKVstoreConnectionRead(d, meta)
}

func resourceAlicloudKVstoreConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}

	submatch := kvstoreConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := r_kvstore.CreateReleaseInstancePublicConnectionRequest()
	request.RegionId = client.RegionId
	request.InstanceId = split[0]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := kvstoreService.DescribeKVstorePublicConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString
		var raw interface{}
		raw, err = client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ReleaseInstancePublicConnection(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{"IncorrectDBInstanceState"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidCurrentConnectionStringNotFound, AtLeastOneNetTypeExists}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return kvstoreService.WaitForKVstorePublicConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
