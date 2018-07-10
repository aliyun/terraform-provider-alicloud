package alicloud

import (
	"fmt"
	"strings"
	"time"

	"log"

	"strconv"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSlbServerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbServerGroupCreate,
		Read:   resourceAliyunSlbServerGroupRead,
		Update: resourceAliyunSlbServerGroupUpdate,
		Delete: resourceAliyunSlbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-server-group",
			},

			"servers": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ids": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							MaxItems: 20,
							MinItems: 1,
						},
						"port": &schema.Schema{
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(1, 65535),
						},
						"weight": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validateIntegerInRange(0, 100),
						},
					},
				},
				//Set: func(v interface{}) int {
				//	var buf bytes.Buffer
				//	m := v.(map[string]interface{})
				//	buf.WriteString(fmt.Sprintf("%s-", m["server_ids"]))
				//	buf.WriteString(fmt.Sprintf("%d-", m["weight"]))
				//	buf.WriteString(fmt.Sprintf("%d-", m["port"]))
				//	return hashcode.String(buf.String())
				//},
				MaxItems: 20,
				MinItems: 1,
			},
		},
	}
}

func resourceAliyunSlbServerGroupCreate(d *schema.ResourceData, meta interface{}) error {

	var groupId string
	if group, err := meta.(*AliyunClient).slbconn.CreateVServerGroup(&slb.CreateVServerGroupArgs{
		RegionId:         getRegion(d, meta),
		LoadBalancerId:   d.Get("load_balancer_id").(string),
		VServerGroupName: d.Get("name").(string),
		BackendServers:   convertServersToString(d.Get("servers").(*schema.Set).List()),
	}); err != nil {
		return fmt.Errorf("CreateVServerGroup got an error: %#v", err)
	} else {
		groupId = group.VServerGroupId
	}

	d.SetId(groupId)

	return resourceAliyunSlbServerGroupUpdate(d, meta)
}

func resourceAliyunSlbServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	group, err := meta.(*AliyunClient).DescribeSlbVServerGroupAttribute(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", group.VServerGroupName)
	d.Set("load_balancer_id", d.Get("load_balancer_id").(string))

	servers := make([]map[string]interface{}, 0)
	portAndWeight := make(map[string][]string)
	for _, server := range group.BackendServers.BackendServer {
		key := fmt.Sprintf("%d%s%d", server.Port, COLON_SEPARATED, server.Weight)
		if v, ok := portAndWeight[key]; !ok {
			portAndWeight[key] = []string{server.ServerId}
		} else {
			v = append(v, server.ServerId)
			portAndWeight[key] = v
		}
	}
	for key, value := range portAndWeight {
		k := strings.Split(key, COLON_SEPARATED)
		p, e := strconv.Atoi(k[0])
		if e != nil {
			return fmt.Errorf("Convertting port %s to int got an error: %#v.", k[0], e)
		}
		w, e := strconv.Atoi(k[1])
		if e != nil {
			return fmt.Errorf("Convertting weight %s to int got an error: %#v.", k[1], e)
		}
		s := map[string]interface{}{
			"server_ids": value,
			"port":       p,
			"weight":     w,
		}
		servers = append(servers, s)
	}

	if err := d.Set("servers", servers); err != nil {
		return err
	}

	return nil
}

func resourceAliyunSlbServerGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	slbconn := meta.(*AliyunClient).slbconn

	d.Partial(true)

	slb_id := d.Get("load_balancer_id").(string)
	name := d.Get("name").(string)
	update := false

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		update = true
	}

	if d.HasChange("servers") && !d.IsNewResource() {
		o, n := d.GetChange("servers")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			log.Printf("[INFO] Remove old servers: %#v", remove)
			if _, err := slbconn.RemoveVServerGroupBackendServers(&slb.RemoveVServerGroupBackendServersArgs{
				LoadBalancerId: slb_id,
				RegionId:       getRegion(d, meta),
				VServerGroupId: d.Id(),
				BackendServers: convertServersToString(remove),
			}); err != nil {
				return fmt.Errorf("RemoveVServerGroupBackendServers got an error: %#v", err)
			}
		}
		if len(add) > 0 {
			log.Printf("[INFO] Add new servers: %#v", add)
			if _, err := slbconn.AddVServerGroupBackendServers(&slb.AddVServerGroupBackendServersArgs{
				LoadBalancerId: slb_id,
				RegionId:       getRegion(d, meta),
				VServerGroupId: d.Id(),
				BackendServers: convertServersToString(add),
			}); err != nil {
				return fmt.Errorf("AddVServerGroupBackendServers got an error: %#v", err)
			}
		}
		if len(add) < 1 && len(remove) < 1 {
			update = true
		}

		d.SetPartial("servers")
	}

	if update {
		log.Printf("[INFO] Update attribute: name %s and backend servers %#v", name, d.Get("servers").(*schema.Set).List())
		if _, err := slbconn.SetVServerGroupAttribute(&slb.SetVServerGroupAttributeArgs{
			RegionId:         getRegion(d, meta),
			LoadBalancerId:   slb_id,
			VServerGroupId:   d.Id(),
			VServerGroupName: name,
			BackendServers:   convertServersToString(d.Get("servers").(*schema.Set).List()),
		}); err != nil {
			return fmt.Errorf("SetVServerGroupAttribute got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAliyunSlbServerGroupRead(d, meta)
}

func resourceAliyunSlbServerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	slbconn := meta.(*AliyunClient).slbconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := slbconn.DeleteVServerGroup(&slb.DeleteVServerGroupArgs{
			RegionId:       getRegion(d, meta),
			VServerGroupId: d.Id(),
		}); err != nil {
			if IsExceptedError(err, VServerGroupNotFoundMessage) || IsExceptedError(err, InvalidParameter) {
				return nil
			}
			if IsExceptedError(err, RspoolVipExist) {
				return resource.RetryableError(fmt.Errorf("DeleteVServerGroup got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}

		group, err := slbconn.DescribeVServerGroupAttribute(&slb.DescribeVServerGroupAttributeArgs{
			RegionId:       getRegion(d, meta),
			VServerGroupId: d.Id(),
		})
		if err != nil {
			if IsExceptedError(err, VServerGroupNotFoundMessage) || IsExceptedError(err, InvalidParameter) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While deleting VServer Group, DescribeVServerGroupAttribute got an error: %#v", err))
		}
		if group != nil {
			return resource.RetryableError(fmt.Errorf("DeleteVServerGroup got an error: %#v", err))
		}
		return nil
	})
}

func convertServersToString(items []interface{}) string {

	if len(items) < 1 {
		return ""
	}
	var servers []string
	for _, server := range items {
		s := server.(map[string]interface{})

		var server_ids []interface{}
		var port, weight int
		if v, ok := s["server_ids"]; ok {
			server_ids = v.([]interface{})
		}
		if v, ok := s["port"]; ok {
			port = v.(int)
		}
		if v, ok := s["weight"]; ok {
			weight = v.(int)
		}

		for _, id := range server_ids {
			str := fmt.Sprintf("{'ServerId':'%s','Port':'%d','Weight':'%d'}", strings.Trim(id.(string), " "), port, weight)

			servers = append(servers, str)
		}

	}
	return fmt.Sprintf("[%s]", strings.Join(servers, COMMA_SEPARATED))
}
