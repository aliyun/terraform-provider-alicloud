package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"regexp"
	"strings"
	"time"
)

func resourceAlicloudPolarDBPrimaryEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBPrimaryEndpointCreate,
		Read:   resourceAlicloudPolarDBPrimaryEndpointRead,
		Update: resourceAlicloudPolarDBPrimaryEndpointUpdate,
		Delete: resourceAlicloudPolarDBPrimaryEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable", "Update"}, false),
				Optional:     true,
			},
			"ssl_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Public", "Private", "Inner"}, false),
			},
			"ssl_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_auto_rotate": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
			},
			"ssl_certificate_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_endpoint_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-),  must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBPrimaryEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clusterId := d.Get("db_cluster_id").(string)
	request := polardb.CreateDescribeDBClusterEndpointsRequest()

	request.RegionId = client.RegionId
	request.DBClusterId = clusterId

	raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.DescribeDBClusterEndpoints(request)
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_primary_endpoint", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)

	var endpoints []polardb.DBEndpoint
	for _, item := range response.Items {
		if item.EndpointType == "Primary" {
			endpoints = append(endpoints, item)
		}
	}

	if len(endpoints) == 0 {
		return WrapError(fmt.Errorf("PolarDB Cluster %s does not have primary endpoint", clusterId))
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, endpoints[0].DBEndpointId))

	return resourceAlicloudPolarDBPrimaryEndpointUpdate(d, meta)
}

func resourceAlicloudPolarDBPrimaryEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]
	object, err := polarDBService.DescribePolarDBClusterEndpoint(d.Id())

	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_cluster_id", dbClusterId)
	d.Set("db_endpoint_id", dbEndpointId)
	d.Set("endpoint_type", object.EndpointType)
	d.Set("db_endpoint_description", object.DBEndpointDescription)

	dbClusterSSL, err := polarDBService.DescribePolarDBClusterSSL(d)

	var sslConnectionString string
	var sslExpireTime string
	var sslEnabled string
	if len(dbClusterSSL.Items) < 1 {
		sslConnectionString = ""
		sslExpireTime = ""
		sslEnabled = ""
	} else if len(dbClusterSSL.Items) == 1 && dbClusterSSL.Items[0].DBEndpointId == "" {
		sslConnectionString = dbClusterSSL.Items[0].SSLConnectionString
		sslExpireTime = dbClusterSSL.Items[0].SSLExpireTime
		sslEnabled = convertPolarDBSSLEnableResponse(dbClusterSSL.Items[0].SSLEnabled)
	} else {
		for _, item := range dbClusterSSL.Items {
			if item.DBEndpointId == dbEndpointId {
				sslConnectionString = item.SSLConnectionString
				sslExpireTime = item.SSLExpireTime
				sslEnabled = convertPolarDBSSLEnableResponse(item.SSLEnabled)
			}
		}
	}
	sslAutoRotate := dbClusterSSL.SSLAutoRotate

	if err := d.Set("ssl_connection_string", sslConnectionString); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ssl_expire_time", sslExpireTime); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ssl_auto_rotate", sslAutoRotate); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ssl_enabled", sslEnabled); err != nil {
		return WrapError(err)
	}
	polarDBService.fillingPolarDBEndpointSslCertificateUrl(sslEnabled, d)

	privateAdress, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), "Private")

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("port", privateAdress.Port)
	prefix := strings.Split(privateAdress.ConnectionString, ".")
	d.Set("connection_prefix", prefix[0])

	return nil
}

func resourceAlicloudPolarDBPrimaryEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]
	if d.HasChange("db_endpoint_description") {
		modifyEndpointRequest := polardb.CreateModifyDBClusterEndpointRequest()
		modifyEndpointRequest.RegionId = client.RegionId
		modifyEndpointRequest.DBClusterId = dbClusterId
		modifyEndpointRequest.DBEndpointId = dbEndpointId

		configItem := make(map[string]string)
		if d.HasChange("db_endpoint_description") {
			modifyEndpointRequest.DBEndpointDescription = d.Get("db_endpoint_description").(string)
			configItem["DBEndpointDescription"] = d.Get("db_endpoint_description").(string)
		}
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBClusterEndpoint(modifyEndpointRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyEndpointRequest.GetActionName(), raw, modifyEndpointRequest.RpcRequest, modifyEndpointRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyEndpointRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait cluster endpoint config modified
		if err := polarDBService.WaitPolardbEndpointConfigEffect(
			d.Id(), configItem, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChanges("ssl_enabled", "ssl_auto_rotate") {
		if d.Get("ssl_enabled") == "" && d.Get("net_type") != "" {
			return WrapErrorf(Error("Need to specify ssl_enabled as Enable or Disable, if you want to modify the net_type."), DefaultErrorMsg, d.Id(), "ModifyDBClusterSSL", ProviderERROR)
		}
		modifySSLRequest := polardb.CreateModifyDBClusterSSLRequest()
		modifySSLRequest.SSLEnabled = d.Get("ssl_enabled").(string)
		modifySSLRequest.NetType = d.Get("net_type").(string)
		modifySSLRequest.DBClusterId = dbClusterId
		modifySSLRequest.DBEndpointId = dbEndpointId
		modifySSLRequest.SSLAutoRotate = d.Get("ssl_auto_rotate").(string)
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBClusterSSL(modifySSLRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifySSLRequest.GetActionName(), raw, modifySSLRequest.RpcRequest, modifySSLRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifySSLRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from SSL_MODIFYING to Running
		stateConf := BuildStateConf([]string{"SSLModifying"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(dbClusterId, []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, dbClusterId)
		}
	}

	if d.HasChanges("connection_prefix", "port") {
		request := polardb.CreateModifyDBEndpointAddressRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBEndpointId = parts[1]
		object, err := polarDBService.DescribePolarDBConnectionV2(d.Id(), "Private")
		if err != nil {
			return WrapError(err)
		}

		request.NetType = "Private"
		request.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		request.Port = d.Get("port").(string)
		prefix := strings.Split(object.ConnectionString, ".")
		if (request.Port != "" && request.Port != object.Port) || (request.ConnectionStringPrefix != "" && request.ConnectionStringPrefix != prefix[0]) {
			if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
				raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBEndpointAddress(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			// wait instance connection_prefix modify success
			if err := polarDBService.WaitForPolarDBConnectionPrefix(d.Id(), request.ConnectionStringPrefix, request.Port, "Private", DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}

		}
	}
	return resourceAlicloudPolarDBPrimaryEndpointRead(d, meta)
}

func resourceAlicloudPolarDBPrimaryEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	//  Terraform can not destroy it..
	return nil
}
