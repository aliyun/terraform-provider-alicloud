package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/rds"
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	if !ok || prefix == "" {
		prefix = instance_id
	}

	if err := client.AllocateDBPublicConnection(instance_id, prefix.(string), d.Get("port").(string)); err != nil {
		return fmt.Errorf("AllocateInstancePublicConnection got an error: %#v", err)
	}

	connection, err := client.DescribeDBInstanceNetInfoByIpType(instance_id, rds.Public)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", instance_id, connection.ConnectionString))

	return resourceAlicloudDBConnectionUpdate(d, meta)
}

func resourceAlicloudDBConnectionRead(d *schema.ResourceData, meta interface{}) error {

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	conn, err := meta.(*AliyunClient).DescribeDBInstanceNetInfoByIpType(parts[0], rds.Public)

	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", strings.Split(conn.ConnectionString, DOT_SEPARATED)[0])
	d.Set("port", conn.Port)
	d.Set("connection_string", conn.ConnectionString)
	d.Set("ip_address", conn.IPAddress)

	return nil
}

func resourceAlicloudDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	parts := strings.Split(d.Id(), COLON_SEPARATED)

	update := false
	args := &rds.ModifyDBInstanceConnectionStringArgs{
		DBInstanceId:            parts[0],
		CurrentConnectionString: parts[1],
		ConnectionStringPrefix:  strings.Split(parts[1], DOT_SEPARATED)[0],
		Port: d.Get("port").(string),
	}
	if d.HasChange("connection_prefix") && !d.IsNewResource() {
		args.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		update = true
		d.SetPartial("connection_prefix")
	}

	if d.HasChange("port") && !d.IsNewResource() {
		update = true
		d.SetPartial("port")
	}

	if update {
		if err := meta.(*AliyunClient).rdsconn.ModifyDBInstanceConnectionString(args); err != nil {
			return fmt.Errorf("ModifyDBInstanceConnectionString got an error: %#v", err)
		}
		connection, err := meta.(*AliyunClient).DescribeDBInstanceNetInfoByIpType(args.DBInstanceId, rds.Public)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%s:%s", args.DBInstanceId, connection.ConnectionString))
	}

	d.Partial(false)
	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := client.ReleaseDBPublicConnection(parts[0], parts[1])

		if err != nil {
			if IsExceptedError(err, InvalidCurrentConnectionStringNotFound) || IsExceptedError(err, AtLeastOneNetTypeExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Release DB connection timeout and got an error: %#v.", err))
		}
		conn, err := meta.(*AliyunClient).DescribeDBInstanceNetInfoByIpType(parts[0], rds.Public)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Release DB connection got an error: %#v.", err))
		}

		if conn == nil {
			d.SetId("")
			return nil
		}

		return resource.RetryableError(fmt.Errorf("DB in use - trying again while it is deleted."))
	})
}
