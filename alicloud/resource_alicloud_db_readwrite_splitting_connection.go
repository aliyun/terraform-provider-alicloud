package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDBConnectionPrefix,
			},
			"distribution_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Standard", "Custom"}),
			},
			"weight": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"max_delay_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": &schema.Schema{
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
			return err
		} else {
			request.Weight = string(serial)
		}
	}

	if err := resource.Retry(60*time.Minute, func() *resource.RetryError {
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.AllocateReadWriteSplittingConnection(request)
		})
		if err != nil {
			if IsExceptedErrors(err, DBReadInstanceNotReadyStatus) {
				return resource.RetryableError(fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err)
	}

	d.SetId(request.DBInstanceId)

	// wait read write splitting connection ready after creation
	// for it may take up to 10 hours to create a readonly instance
	if err := rdsService.WaitForDBReadWriteSplitting(request.DBInstanceId, 60*60*10); err != nil {
		return fmt.Errorf("WaitForDBReadWriteSplitting got error: %#v", err)
	}

	return resourceAlicloudDBReadWriteSplittingConnectionUpdate(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	err := rdsService.WaitForDBReadWriteSplitting(d.Id(), 30*60)
	if err != nil {
		return fmt.Errorf("WaitForDBReadWriteSplitting got an error: %#v", err)
	}

	conn, err := rdsService.DescribeReadWriteSplittingConnection(d.Id())
	if err != nil {
		return fmt.Errorf("DescribeReadWriteSplittingConnection got an error: %#v", err)
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_string", conn.ConnectionString)
	d.Set("distribution_type", conn.DistributionType)
	if port, err := strconv.Atoi(conn.Port); err == nil {
		d.Set("port", port)
	}
	if mdt, err := strconv.Atoi(conn.MaxDelayTime); err == nil {
		d.Set("max_delay_time", mdt)
	}
	if w, ok := d.GetOk("weight"); ok {
		documented := w.(map[string]interface{})
		for _, config := range conn.DBInstanceWeights.DBInstanceWeight {
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
	submatch := dbConnectionPrefixWithSuffixRegexp.FindStringSubmatch(conn.ConnectionString)
	if len(submatch) > 1 {
		d.Set("connection_prefix", submatch[1])
	}

	return nil
}

func resourceAlicloudDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	request := rds.CreateModifyReadWriteSplittingConnectionRequest()
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
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
			_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyReadWriteSplittingConnection(request)
			})
			if err != nil {
				if IsExceptedErrors(err, OperationDeniedDBStatus) || IsExceptedErrors(err, DBReadInstanceNotReadyStatus) {
					return resource.RetryableError(fmt.Errorf("Modify DBInstance Connection got an error: %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modify DBInstance Connection got an error: %#v", err))
			}
			return nil
		}); err != nil {
			return err
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	return resource.Retry(30*time.Minute, func() *resource.RetryError {
		request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
		request.DBInstanceId = d.Id()

		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseReadWriteSplittingConnection(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
			}
			if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, InvalidRwSplitNetTypeNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		_, err = rdsService.DescribeReadWriteSplittingConnection(d.Id())
		if err != nil {
			if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		return resource.RetryableError(fmt.Errorf("Release DB connection timeout."))
	})
}
