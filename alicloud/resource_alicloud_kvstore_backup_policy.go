package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudKVStoreBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreBackupPolicyCreate,
		Read:   resourceAlicloudKVStoreBackupPolicyRead,
		Update: resourceAlicloudKVStoreBackupPolicyUpdate,
		Delete: resourceAlicloudKVStoreBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_time": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(BACKUP_TIME),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},
			"backup_period": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				// terraform does not support ValidateFunc of TypeList attr
				// ValidateFunc: validateAllowedStringValue([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}),
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKVStoreBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request := r_kvstore.CreateModifyBackupPolicyRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.PreferredBackupTime = d.Get("backup_time").(string)
	periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	request.PreferredBackupPeriod = backupPeriod

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ModifyBackupPolicy(request); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Create backup policy got an error: %#v", err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(request.InstanceId)

	return resourceAlicloudKVStoreBackupPolicyRead(d, meta)
}

func resourceAlicloudKVStoreBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	policy, err := meta.(*AliyunClient).DescribeRKVInstancebackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", d.Id())
	d.Set("backup_time", policy.PreferredBackupTime)
	d.Set("backup_period", strings.Split(policy.PreferredBackupPeriod, ","))

	return nil
}

func resourceAlicloudKVStoreBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	update := false
	request := r_kvstore.CreateModifyBackupPolicyRequest()
	request.InstanceId = d.Id()

	if d.HasChange("backup_time") {
		request.PreferredBackupTime = d.Get("backup_time").(string)
		update = true
	}

	if d.HasChange("backup_period") {
		periodList := expandStringList(d.Get("backup_period").(*schema.Set).List())
		backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		request.PreferredBackupPeriod = backupPeriod
		update = true
	}

	if update {
		if _, err := conn.ModifyBackupPolicy(request); err != nil {
			return err
		}
	}

	return resourceAlicloudKVStoreBackupPolicyRead(d, meta)
}

func resourceAlicloudKVStoreBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// In case of a delete we are resetting to default values which is Monday - Sunday each 3am-4am
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	request := r_kvstore.CreateModifyBackupPolicyRequest()
	request.InstanceId = d.Id()

	request.PreferredBackupTime = "01:00Z-02:00Z"
	request.PreferredBackupPeriod = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"

	if _, err := conn.ModifyBackupPolicy(request); err != nil {
		return err
	}

	return nil
}
