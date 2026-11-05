package polardb

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEndpointAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointAddressCreate,
		Read:   resourceEndpointAddressRead,
		Update: resourceEndpointAddressUpdate,
		Delete: resourceEndpointAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_endpoint_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Public"}, false),
				Default:      "Public",
				ForceNew:     true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceEndpointAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	clusterId := d.Get("db_cluster_id").(string)
	dbEndpointId := d.Get("db_endpoint_id").(string)
	netType := d.Get("net_type").(string)

	request := polardb.CreateCreateDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.DBEndpointId = dbEndpointId
	request.NetType = netType
	if v, ok := d.GetOk("connection_prefix"); ok && v.(string) != "" {
		request.ConnectionStringPrefix = v.(string)
	}
	var raw interface{}
	var err error
	err = retry.Retry(8*time.Minute, func() *retry.RetryError {
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDBEndpointAddress(request)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_endpoint_address", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	if err := polarDBService.WaitForPolarDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	// wait instance running after allocating
	if err := polarDBService.WaitForPolarDBInstance(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}

	return resourceEndpointAddressUpdate(d, meta)
}

func resourceEndpointAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	object, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), "Public")

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
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
	d.Set("db_endpoint_id", parts[1])
	d.Set("port", object.Port)
	d.Set("net_type", object.NetType)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)
	prefix := strings.Split(object.ConnectionString, ".")
	d.Set("connection_prefix", prefix[0])

	return nil
}

func resourceEndpointAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	if d.IsNewResource() && d.HasChange("port") {
		request := polardb.CreateModifyDBEndpointAddressRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBEndpointId = parts[1]
		object, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), "Public")
		if err != nil {
			return errs.WrapError(err)
		}

		prefix := strings.Split(object.ConnectionString, ".")
		request.ConnectionStringPrefix = prefix[0]
		request.Port = d.Get("port").(string)
		request.NetType = d.Get("net_type").(string)
		if request.Port != object.Port {
			if err := retry.Retry(8*time.Minute, func() *retry.RetryError {
				raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBEndpointAddress(request)
				})
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}

			// wait instance connection_prefix modify success
			if err := polarDBService.WaitForPolarDBConnectionPrefix(d.Id(), request.ConnectionStringPrefix, request.Port, "Public", DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
		}
	}

	if d.IsNewResource() {
		return resourceEndpointAddressRead(d, meta)
	}

	if d.HasChanges("connection_prefix", "port") {
		request := polardb.CreateModifyDBEndpointAddressRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBEndpointId = parts[1]
		object, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), d.Get("net_type").(string))
		if err != nil {
			return errs.WrapError(err)
		}

		request.NetType = object.NetType
		request.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		request.Port = d.Get("port").(string)
		request.NetType = d.Get("net_type").(string)
		if err := retry.Retry(8*time.Minute, func() *retry.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBEndpointAddress(request)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		// wait instance connection_prefix modify success
		if err := polarDBService.WaitForPolarDBConnectionPrefix(d.Id(), request.ConnectionStringPrefix, request.Port, request.NetType, DefaultTimeoutMedium); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceEndpointAddressRead(d, meta)
}

func resourceEndpointAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}

	parts, errParse := helper.ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return errs.WrapError(errParse)
	}
	request := polardb.CreateDeleteDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBEndpointId = parts[1]

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		object, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), d.Get("net_type").(string))
		if err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}
		request.NetType = object.NetType
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBEndpointAddress(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "EndpointStatus.NotSupport"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return polarDBService.WaitForPolarDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
