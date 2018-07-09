package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

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
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateDBConnectionPrefix,
			},
			"port": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				Default:      "3306",
			},
			"connection_string": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	instance_id := d.Get("instance_id").(string)
	prefix, ok := d.GetOk("connection_prefix")
	if !ok || prefix.(string) == "" {
		prefix = fmt.Sprintf("%stf", instance_id)
	}

	if err := client.AllocateDBPublicConnection(instance_id, prefix.(string), d.Get("port").(string)); err != nil {
		return fmt.Errorf("AllocateInstancePublicConnection got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instance_id, COLON_SEPARATED, prefix.(string)))

	return resourceAlicloudDBConnectionUpdate(d, meta)
}

func resourceAlicloudDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	if strings.HasSuffix(d.Id(), DBConnectionSuffix) {
		d.SetId(strings.Replace(d.Id(), DBConnectionSuffix, "", -1))
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	conn, err := meta.(*AliyunClient).DescribeDBInstanceNetInfoByIpType(parts[0], Public)

	if err != nil {
		if NotFoundDBInstance(err) {
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
	client := meta.(*AliyunClient)
	d.Partial(true)

	if strings.HasSuffix(d.Id(), DBConnectionSuffix) {
		d.SetId(strings.Replace(d.Id(), DBConnectionSuffix, "", -1))
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if d.HasChange("port") && !d.IsNewResource() {
		request := rds.CreateModifyDBInstanceConnectionStringRequest()
		request.DBInstanceId = parts[0]
		request.CurrentConnectionString = fmt.Sprintf("%s%s", parts[1], DBConnectionSuffix)
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)

		// wait instance running before modifying
		if err := client.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			if _, err := client.rdsconn.ModifyDBInstanceConnectionString(request); err != nil {
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
		if err := client.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		d.SetPartial("port")

	}

	d.Partial(false)
	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	if strings.HasSuffix(d.Id(), DBConnectionSuffix) {
		d.SetId(strings.Replace(d.Id(), DBConnectionSuffix, "", -1))
	}

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := client.ReleaseDBPublicConnection(parts[0], fmt.Sprintf("%s%s", parts[1], DBConnectionSuffix))

		if err != nil {
			if IsExceptedError(err, InvalidCurrentConnectionStringNotFound) || IsExceptedError(err, AtLeastOneNetTypeExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Release DB connection timeout and got an error: %#v.", err))
		}
		conn, err := meta.(*AliyunClient).DescribeDBInstanceNetInfoByIpType(parts[0], Public)

		if err != nil {
			if NotFoundDBInstance(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
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
