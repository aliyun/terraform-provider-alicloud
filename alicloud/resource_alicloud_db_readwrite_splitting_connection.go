package alicloud

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const dbConnectionPrefixWithSuffixRegex = "^([a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionPrefixWithSuffixRegexp = regexp.MustCompile(dbConnectionPrefixWithSuffixRegex)

func resourceAlicloudDBReadWriteSplittingConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBReadWriteSplittingConnectionCreate,
		Read:   resourceAlicloudDBReadWriteSplittingConnectionRead,
		Update: resourceAlicloudDBReadWriteSplittingConnectionUpdate,
		Delete: resourceAlicloudDBReadWriteSplittingConnectionDelete,
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
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"distribution_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Custom"}, false),
			},
			"weight": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	request := rds.CreateAllocateReadWriteSplittingConnectionRequest()
	request.RegionId = string(client.Region)
	request.DBInstanceId = Trim(d.Get("instance_id").(string))
	request.MaxDelayTime = strconv.Itoa(d.Get("max_delay_time").(int))

	prefix, ok := d.GetOk("connection_prefix")
	if ok && prefix.(string) != "" {
		request.ConnectionStringPrefix = prefix.(string)
	}

	port, ok := d.GetOk("port")
	if ok {
		request.Port = strconv.Itoa(port.(int))
	}

	request.DistributionType = d.Get("distribution_type").(string)

	if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
		if serial, err := json.Marshal(weight); err != nil {
			return WrapError(err)
		} else {
			request.Weight = string(serial)
		}
	}

	if err := resource.Retry(60*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.AllocateReadWriteSplittingConnection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, DBReadInstanceNotReadyStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(request.DBInstanceId)

	// wait read write splitting connection ready after creation
	// for it may take up to 10 hours to create a readonly instance
	if err := rdsService.WaitForDBReadWriteSplitting(request.DBInstanceId, "", 60*60*10); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBReadWriteSplittingConnectionUpdate(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	proxy, proxyErr := rdsService.DescribeDBProxy(d.Id())
	if proxyErr != nil {
		return WrapError(proxyErr)
	}
	if proxy.DBProxyInstanceType == "2" {
		return resourceAlicloudDBProxyEndpointRead(d, rdsService, proxy.DBProxyInstanceName)
	}

	err := rdsService.WaitForDBReadWriteSplitting(d.Id(), "", DefaultLongTimeout)
	if err != nil {
		return WrapError(err)
	}

	object, err := rdsService.DescribeDBReadWriteSplittingConnection(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_string", object.ConnectionString)
	d.Set("distribution_type", object.DistributionType)
	if port, err := strconv.Atoi(object.Port); err == nil {
		d.Set("port", port)
	}
	if mdt, err := strconv.Atoi(object.MaxDelayTime); err == nil {
		d.Set("max_delay_time", mdt)
	}
	if w, ok := d.GetOk("weight"); ok {
		documented := w.(map[string]interface{})
		for _, config := range object.DBInstanceWeights.DBInstanceWeight {
			if config.Availability != "Available" {
				delete(documented, config.DBInstanceId)
				continue
			}
			if config.Weight != "0" {
				if _, ok := documented[config.DBInstanceId]; ok {
					documented[config.DBInstanceId] = config.Weight
				}
			}
		}
		d.Set("weight", documented)
	}
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(object.ConnectionString)
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	return nil
}

func resourceAlicloudDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	request := rds.CreateModifyReadWriteSplittingConnectionRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()

	update := false

	if d.HasChange("max_delay_time") {
		request.MaxDelayTime = strconv.Itoa(d.Get("max_delay_time").(int))
		update = true
	}

	if !update && d.IsNewResource() {
		return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
	}

	if d.HasChange("weight") {
		if weight, ok := d.GetOk("weight"); ok && weight != nil && len(weight.(map[string]interface{})) > 0 {
			if serial, err := json.Marshal(weight); err != nil {
				return err
			} else {
				request.Weight = string(serial)
			}
		}
		update = true
	}

	if d.HasChange("distribution_type") {
		request.DistributionType = d.Get("distribution_type").(string)
		update = true
	}

	if update {
		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 60*60); err != nil {
			return WrapError(err)
		}

		if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyReadWriteSplittingConnection(request)
			})
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) || IsExpectedErrors(err, DBReadInstanceNotReadyStatus) {
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
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Id()

	if err := resource.Retry(30*time.Minute, func() *resource.RetryError {

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseReadWriteSplittingConnection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidRwSplitNetType.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(rdsService.WaitForDBReadWriteSplitting(d.Id(), Deleted, DefaultLongTimeout))
}

func resourceAlicloudDBProxyEndpointRead(d *schema.ResourceData, rdsService RdsService, endPointName string) error {
	endpointInfo, endpointError := rdsService.DescribeDBProxyEndpoint(d.Id(), endPointName)
	if endpointError != nil {
		return WrapError(endpointError)
	}
	d.Set("instance_id", d.Id())
	d.Set("connection_string", endpointInfo.DBProxyConnectString)
	d.Set("distribution_type", endpointInfo.ReadOnlyInstanceDistributionType)
	if port, err := strconv.Atoi(endpointInfo.DBProxyConnectStringPort); err == nil {
		d.Set("port", port)
	}

	if mdt, err := strconv.Atoi(endpointInfo.ReadOnlyInstanceMaxDelayTime); err == nil {
		d.Set("max_delay_time", mdt)
	}
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(endpointInfo.DBProxyConnectString)
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	var documented map[string]interface{}
	if w, ok := d.GetOk("weight"); ok {
		documented = w.(map[string]interface{})
	} else {
		documented = make(map[string]interface{})
	}
	var weight []map[string]interface{}
	rawData := []byte(endpointInfo.ReadOnlyInstanceWeight)
	parseErr := json.Unmarshal(rawData, &weight)
	if parseErr != nil {
		return WrapError(parseErr)
	}
	for _, configNode := range weight {
		var dbInstanceId string
		if instanceId, ok := configNode["DBInstanceId"]; ok {
			dbInstanceId = instanceId.(string)
		}
		if _, ok := configNode["Availability"]; ok && configNode["Availability"] != "Available" {
			delete(documented, dbInstanceId)
			continue
		}
		if _, ok := configNode["Weight"]; ok && configNode["Weight"] != "0" {
			if _, ok := documented[dbInstanceId]; ok {
				documented[dbInstanceId] = configNode["Weight"]
			}
		}
	}
	d.Set("weight", documented)
	return nil
}
