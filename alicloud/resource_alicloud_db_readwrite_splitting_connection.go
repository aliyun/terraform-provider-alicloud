package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

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
			"port": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Default: "3306",
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(0, 7200),
			},
			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
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
	if !ok || prefix.(string) == "" {
		prefix = fmt.Sprintf("%srw", request.DBInstanceId)
	}
	request.ConnectionStringPrefix = prefix.(string)

	port, ok := d.GetOk("port")
	if ok && port.(string) != "" {
		request.Port = port.(string)
	}

	distributionType, ok := d.GetOk("distribution_type")
	if ok && distributionType.(string) != "" {
		request.DistributionType = distributionType.(string)
	}

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
			if IsExceptedErrors(err, []string{"OperationDenied.ReadDBInstanceStatus", "ReadDBInstance.Mismatch"}) {
				return resource.RetryableError(fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("AllocateReadWriteSplittingConnection got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, prefix.(string)))

	// wait instance running after modifying
	// for it may take up to 10 hours to create a readonly instance
	if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 60*60*10); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAlicloudDBReadWriteSplittingConnectionUpdate(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	if err := rdsService.WaitForDBInstance(parts[0], Running, 60*60*1); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	if err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		resp, err := rdsService.DescribeDBInstanceNetInfos(parts[0])
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if resp == nil {
			return resource.NonRetryableError(GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", parts[0])))
		}

		found := false
		for _, conn := range resp {
			if conn.ConnectionStringType == "ReadWriteSplitting" {
				d.Set("instance_id", parts[0])
				d.Set("connection_prefix", parts[1])
				d.Set("port", conn.Port)
				d.Set("connection_string", conn.ConnectionString)
				d.Set("distribution_type", conn.DistributionType)
				if mdt, err := strconv.Atoi(conn.MaxDelayTime); err != nil {
					return resource.RetryableError(err)
				} else {
					d.Set("max_delay_time", mdt)
				}
				weights := make(map[string]interface{})
				for _, config := range conn.DBInstanceWeights.DBInstanceWeight {
					if config.Availability != "Available" {
						continue
					}
					weights[config.DBInstanceId] = config.Weight
				}
				d.Set("weight", weights)
				found = true
				break
			}
		}
		if !found {
			return resource.RetryableError(GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any read write splitting connection.", parts[0])))
		}

		return nil
	}); err != nil {
		return fmt.Errorf("read ReadWriteSplittingConnection got an error: %#v", err)
	}

	return nil
}

func resourceAlicloudDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	request := rds.CreateModifyReadWriteSplittingConnectionRequest()
	request.DBInstanceId = parts[0]

	update := false

	if d.HasChange("max_delay_time") {
		request.MaxDelayTime = strconv.Itoa(d.Get("max_delay_time").(int))
		d.SetPartial("max_delay_time")
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
		d.SetPartial("weight")
		update = true
	}

	if d.HasChange("distribution_type") {
		request.DistributionType = d.Get("distribution_type").(string)
		d.SetPartial("distribution_type")
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
				if IsExceptedErrors(err, OperationDeniedDBStatus) || IsExceptedError(err, "ReadDBInstance.Mismatch") {
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

	d.Partial(false)
	return resourceAlicloudDBReadWriteSplittingConnectionRead(d, meta)
}

func resourceAlicloudDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(30*time.Minute, func() *resource.RetryError {
		request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
		request.DBInstanceId = parts[0]

		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseReadWriteSplittingConnection(request)
		})
		if err != nil {
			if IsExceptedError(err, "OperationDenied.DBInstanceStatus") {
				return resource.RetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
			}
			if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, "InvalidRwSplitNetType.NotFound") {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		resp, err := rdsService.DescribeDBInstanceNetInfos(parts[0])

		if err != nil {
			if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		if resp == nil {
			return nil
		}

		found := false
		for _, conn := range resp {
			if conn.ConnectionStringType == "ReadWriteSplitting" {
				found = true
				break
			}
		}

		if !found {
			d.SetId("")
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Release DB connection timeout."))
	})
}
