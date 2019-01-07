package alicloud

import (
	"fmt"
	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const dbConnectionSuffixRegex = "\\.mysql\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const dbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionIdWithSuffixRegexp = regexp.MustCompile(dbConnectionIdWithSuffixRegex)

func resourceAlicloudDBConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBConnectionCreate,
		Read:   resourceAlicloudDBConnectionRead,
		Update: resourceAlicloudDBConnectionUpdate,
		Delete: resourceAlicloudDBConnectionDelete,
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

func resourceAlicloudDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	instance_id := d.Get("instance_id").(string)
	prefix, ok := d.GetOk("connection_prefix")
	if !ok || prefix.(string) == "" {
		prefix = fmt.Sprintf("%stf", instance_id)
	}

	if err := rdsService.AllocateDBPublicConnection(instance_id, prefix.(string), d.Get("port").(string)); err != nil {
		return fmt.Errorf("AllocateInstancePublicConnection got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instance_id, COLON_SEPARATED, prefix.(string)))

	return resourceAlicloudDBConnectionUpdate(d, meta)
}

func resourceAlicloudDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	conn, err := rdsService.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", conn.Port)
	d.Set("connection_string", conn.ConnectionString)
	d.Set("ip_address", conn.IPAddress)

	return nil
}

func resourceAlicloudDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if d.HasChange("port") && !d.IsNewResource() {
		request := rds.CreateModifyDBInstanceConnectionStringRequest()
		request.DBInstanceId = parts[0]
		connectionString, err := getCurrentConnectionString(parts[0], meta)
		if err != nil {
			return fmt.Errorf("getCurrentConnectionString got error: %#v", err)
		}
		request.CurrentConnectionString = connectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		// wait instance running before modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if IsExceptedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(fmt.Errorf("Modify DBInstance Connection Port got an error: %#v.", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modify DBInstance Connection Port got an error: %#v.", err))
			}
			return nil
		}); err != nil {
			return err
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		d.SetPartial("port")

	}

	d.Partial(false)
	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		connectionString, err := getCurrentConnectionString(parts[0], meta)
		if err != nil {
			if rdsService.NotFoundDBInstance(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("getCurrentConnectionString got error: %#v", err))
		}
		err = rdsService.ReleaseDBPublicConnection(parts[0], connectionString)

		if err != nil {
			if IsExceptedErrors(err, []string{InvalidCurrentConnectionStringNotFound, AtLeastOneNetTypeExists}) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Release DB connection timeout and got an error: %#v.", err))
		}
		conn, err := rdsService.DescribeDBInstanceNetInfoByIpType(parts[0], Public)

		if err != nil {
			if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		if conn == nil {
			d.SetId("")
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Release DB connection timeout."))
	})
}

func getCurrentConnectionString(dbInstanceId string, meta interface{}) (string, error) {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	resp, err := rdsService.DescribeDBInstanceNetInfoByIpType(dbInstanceId, Public)
	if err != nil {
		return "", err
	}
	return resp.ConnectionString, nil
}
