package alicloud

import (
	"fmt"
	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const polardbConnectionSuffixRegex = "\\.[a-zA-Z0-9\\-\\.]+\\.rds\\.aliyuncs\\.com"
const polardbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + polardbConnectionSuffixRegex + "$"

var polardbConnectionIdWithSuffixRegexp = regexp.MustCompile(polardbConnectionIdWithSuffixRegex)

func resourceAlicloudPolarDBConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBConnectionCreate,
		Read:   resourceAlicloudPolarDBConnectionRead,
		Update: resourceAlicloudPolarDBConnectionUpdate,
		Delete: resourceAlicloudPolarDBConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
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
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
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

func resourceAlicloudPolarDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	clusterId := d.Get("cluster_id").(string)
	dbEndpointId := d.Get("db_endpoint_id").(string)
	netType := d.Get("net_type").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", dbEndpointId)
	}

	request := polardb.CreateCreateDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.DBEndpointId = dbEndpointId
	request.NetType = netType
	request.ConnectionStringPrefix = prefix
	var raw interface{}
	var err error
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDBEndpointAddress(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId, COLON_SEPARATED, netType))

	if err := polarDBService.WaitForPolarDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := polarDBService.WaitForPolarDBInstance(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBConnectionRead(d, meta)
}

func resourceAlicloudPolarDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := polardbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBConnection(d.Id())

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	d.Set("db_endpoint_id", parts[1])
	d.Set("port", object.Port)
	d.Set("net_type", object.NetType)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceAlicloudPolarDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	submatch := polardbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("connection_prefix") {
		request := polardb.CreateModifyDBEndpointAddressRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBEndpointId = parts[1]
		object, err := polarDBService.DescribePolarDBConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.ConnectionStringPrefix = object.ConnectionString
		request.NetType = object.NetType
		request.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		request.NetType = d.Get("net_type").(string)
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBEndpointAddress(request)
			})
			if err != nil {
				if IsExceptedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		if err := polarDBService.WaitForPolarDBInstance(request.DBClusterId, Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudPolarDBConnectionRead(d, meta)
}

func resourceAlicloudPolarDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	submatch := polardbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := polardb.CreateDeleteDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = split[0]
	request.DBEndpointId = split[1]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := polarDBService.DescribePolarDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.NetType = object.NetType
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBEndpointAddress(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidDBClusterStatus, EndpointStatusNotSupport}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound, InvalidCurrentConnectionStringNotFound, AtLeastOneNetTypeExists}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return polarDBService.WaitForPolarDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
