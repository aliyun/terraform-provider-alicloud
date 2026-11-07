package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// AMQP (消息队列 RabbitMQ版, product "amqp-open" / "2019-12-12") child-resource
// sweepers. Six child resources are swept: exchange, queue, binding,
// static_account, open_source_account, open_source_permission. The instance
// sweeper is intentionally omitted: instance Delete is BSS RefundInstance
// (退订, requires explicit approval) and skips PayAsYouGo instances, so it
// has no sweep value here; releasing an instance cascades all its children.
//
// Fleet-wipe pattern: SWEEP_ALL (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses
// the tf-prefix filter so all non-NODELETING children are swept; NORMAL mode
// keeps the prefix filter. The protectedSweepVpcName NODELETING guard (21
// markers) is always-active in both modes. All sweepers return nil — errors
// are logged and never propagated, so one failing child does not abort the
// cascade.

func init() {
	resource.AddTestSweepers("alicloud_amqp_exchange", &resource.Sweeper{
		Name: "alicloud_amqp_exchange",
		F:    testSweepAmqpExchanges,
	})
	resource.AddTestSweepers("alicloud_amqp_queue", &resource.Sweeper{
		Name: "alicloud_amqp_queue",
		F:    testSweepAmqpQueues,
	})
	resource.AddTestSweepers("alicloud_amqp_binding", &resource.Sweeper{
		Name: "alicloud_amqp_binding",
		F:    testSweepAmqpBindings,
	})
	resource.AddTestSweepers("alicloud_amqp_static_account", &resource.Sweeper{
		Name: "alicloud_amqp_static_account",
		F:    testSweepAmqpStaticAccounts,
	})
	resource.AddTestSweepers("alicloud_amqp_open_source_account", &resource.Sweeper{
		Name: "alicloud_amqp_open_source_account",
		F:    testSweepAmqpOpenSourceAccounts,
	})
	resource.AddTestSweepers("alicloud_amqp_open_source_permission", &resource.Sweeper{
		Name: "alicloud_amqp_open_source_permission",
		F:    testSweepAmqpOpenSourcePermissions,
	})
}

var amqpSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// amqpGetField returns the first non-empty string value of item[name] across
// the given name variants. amqp-open response field casing is inconsistent
// across endpoints (ListInstances/ListVirtualHosts/ListExchanges use
// PascalCase; ListAccounts uses camelCase), so callers pass both variants.
func amqpGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// amqpSweepNameOk applies the always-active NODELETING guard and the
// sweepAll()-aware tf-prefix filter. Returns true when the resource should
// be swept. SWEEP_ALL bypasses the prefix filter (fleet-wipe: clean all
// non-NODELETING); NORMAL mode keeps the prefix filter.
func amqpSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping amqp resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range amqpSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping amqp resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// listAmqpServingInstanceIDs enumerates SERVING amqp instances — the only
// state in which child resources can be listed/deleted.
func listAmqpServingInstanceIDs(client *connectivity.AliyunClient) ([]string, error) {
	action := "ListInstances"
	request := map[string]interface{}{"MaxResults": PageSizeXLarge}
	instanceIds := make([]string, 0)
	for {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		var err error
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("ListInstances: %s", err)
		}
		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			return nil, fmt.Errorf("ListInstances jsonpath $.Data.Instances: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				if amqpGetField(item, "Status") == "SERVING" {
					if id := amqpGetField(item, "InstanceId"); id != "" {
						instanceIds = append(instanceIds, id)
					}
				}
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return instanceIds, nil
}

// listAmqpVirtualHostNames enumerates virtual host names in an instance.
func listAmqpVirtualHostNames(client *connectivity.AliyunClient, instanceId string) ([]string, error) {
	action := "ListVirtualHosts"
	request := map[string]interface{}{
		"InstanceId": instanceId,
		"MaxResults": PageSizeLarge,
	}
	vhosts := make([]string, 0)
	for {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		var err error
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("ListVirtualHosts(%s): %s", instanceId, err)
		}
		resp, err := jsonpath.Get("$.Data.VirtualHosts", response)
		if err != nil {
			return nil, fmt.Errorf("ListVirtualHosts jsonpath: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				if name := amqpGetField(item, "Name"); name != "" {
					vhosts = append(vhosts, name)
				}
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return vhosts, nil
}

// listAmqpOpenSourceAccountNames enumerates open-source account names (the
// UserName field) in an instance — the middle level of the
// open_source_permission walk. $.Data is normally an array; a map-shaped
// response (keyed by instanceId) is handled as a fallback.
func listAmqpOpenSourceAccountNames(client *connectivity.AliyunClient, instanceId string) ([]string, error) {
	action := "ListOpenSourceAccounts"
	request := map[string]interface{}{"InstanceId": instanceId}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ListOpenSourceAccounts(%s): %s", instanceId, err)
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, fmt.Errorf("ListOpenSourceAccounts jsonpath: %s", err)
	}
	names := make([]string, 0)
	switch data := resp.(type) {
	case []interface{}:
		for _, v := range data {
			if item, ok := v.(map[string]interface{}); ok {
				if name := amqpGetField(item, "Name", "UserName"); name != "" {
					names = append(names, name)
				}
			}
		}
	case map[string]interface{}:
		for _, v := range data {
			if arr, ok := v.([]interface{}); ok {
				for _, a := range arr {
					if item, ok := a.(map[string]interface{}); ok {
						if name := amqpGetField(item, "Name", "UserName"); name != "" {
							names = append(names, name)
						}
					}
				}
			}
		}
	}
	return names, nil
}

func testSweepAmqpExchanges(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for exchange sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp exchange: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		vhosts, err := listAmqpVirtualHostNames(client, instanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list virtual hosts for instance %s: %s", instanceId, err)
			continue
		}
		if len(vhosts) == 0 {
			log.Printf("[INFO] amqp exchange: 0 virtual hosts in instance %s", instanceId)
			continue
		}
		for _, vhost := range vhosts {
			action := "ListExchanges"
			request := map[string]interface{}{
				"InstanceId":  instanceId,
				"VirtualHost": vhost,
				"MaxResults":  PageSizeLarge,
			}
			for {
				var response map[string]interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Failed to list exchanges (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				resp, err := jsonpath.Get("$.Data.Exchanges", response)
				if err != nil {
					log.Printf("[ERROR] ListExchanges jsonpath (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				result, _ := resp.([]interface{})
				for _, v := range result {
					item, _ := v.(map[string]interface{})
					name := amqpGetField(item, "Name")
					if !amqpSweepNameOk(name) {
						continue
					}
					log.Printf("[INFO] Deleting amqp exchange: %s (instance %s, vhost %s)", name, instanceId, vhost)
					deleteRequest := map[string]interface{}{
						"InstanceId":   instanceId,
						"VirtualHost":  vhost,
						"ExchangeName": name,
					}
					if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteExchange", nil, deleteRequest, false); e != nil {
						log.Printf("[ERROR] Failed to delete amqp exchange (%s (instance %s, vhost %s)): %s", name, instanceId, vhost, e)
						continue
					}
					log.Printf("[INFO] Delete amqp exchange success: %s (instance %s, vhost %s)", name, instanceId, vhost)
				}
				if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
					request["NextToken"] = nextToken
				} else {
					break
				}
			}
		}
	}
	return nil
}

func testSweepAmqpQueues(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for queue sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp queue: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		vhosts, err := listAmqpVirtualHostNames(client, instanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list virtual hosts for instance %s: %s", instanceId, err)
			continue
		}
		if len(vhosts) == 0 {
			log.Printf("[INFO] amqp queue: 0 virtual hosts in instance %s", instanceId)
			continue
		}
		for _, vhost := range vhosts {
			action := "ListQueues"
			request := map[string]interface{}{
				"InstanceId":  instanceId,
				"VirtualHost": vhost,
				"MaxResults":  PageSizeLarge,
			}
			for {
				var response map[string]interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Failed to list queues (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				resp, err := jsonpath.Get("$.Data.Queues", response)
				if err != nil {
					log.Printf("[ERROR] ListQueues jsonpath (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				result, _ := resp.([]interface{})
				for _, v := range result {
					item, _ := v.(map[string]interface{})
					name := amqpGetField(item, "Name")
					if !amqpSweepNameOk(name) {
						continue
					}
					log.Printf("[INFO] Deleting amqp queue: %s (instance %s, vhost %s)", name, instanceId, vhost)
					deleteRequest := map[string]interface{}{
						"InstanceId":  instanceId,
						"VirtualHost": vhost,
						"QueueName":   name,
					}
					if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteQueue", nil, deleteRequest, false); e != nil {
						log.Printf("[ERROR] Failed to delete amqp queue (%s (instance %s, vhost %s)): %s", name, instanceId, vhost, e)
						continue
					}
					log.Printf("[INFO] Delete amqp queue success: %s (instance %s, vhost %s)", name, instanceId, vhost)
				}
				if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
					request["NextToken"] = nextToken
				} else {
					break
				}
			}
		}
	}
	return nil
}

func testSweepAmqpBindings(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for binding sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp binding: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		vhosts, err := listAmqpVirtualHostNames(client, instanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list virtual hosts for instance %s: %s", instanceId, err)
			continue
		}
		if len(vhosts) == 0 {
			log.Printf("[INFO] amqp binding: 0 virtual hosts in instance %s", instanceId)
			continue
		}
		for _, vhost := range vhosts {
			action := "ListBindings"
			request := map[string]interface{}{
				"InstanceId":  instanceId,
				"VirtualHost": vhost,
				"MaxResults":  PageSizeLarge,
			}
			for {
				var response map[string]interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(5*time.Minute, func() *retry.RetryError {
					response, err = client.RpcGet("amqp-open", "2019-12-12", action, request, nil)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Failed to list bindings (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				resp, err := jsonpath.Get("$.Data.Bindings", response)
				if err != nil {
					log.Printf("[ERROR] ListBindings jsonpath (instance %s, vhost %s): %s", instanceId, vhost, err)
					break
				}
				result, _ := resp.([]interface{})
				for _, v := range result {
					item, _ := v.(map[string]interface{})
					sourceExchange := amqpGetField(item, "SourceExchange", "sourceExchange")
					destinationName := amqpGetField(item, "DestinationName", "destinationName")
					bindingType := amqpGetField(item, "BindingType", "bindingType")
					bindingKey := amqpGetField(item, "BindingKey", "bindingKey", "RoutingKey", "routingKey")
					// A binding has no single name; prefix-check on the
					// source exchange (test exchanges carry the tf-prefix).
					if !amqpSweepNameOk(sourceExchange) {
						continue
					}
					log.Printf("[INFO] Deleting amqp binding: %s -> %s (instance %s, vhost %s)", sourceExchange, destinationName, instanceId, vhost)
					deleteRequest := map[string]interface{}{
						"InstanceId":      instanceId,
						"VirtualHost":     vhost,
						"SourceExchange":  sourceExchange,
						"DestinationName": destinationName,
						"BindingType":     bindingType,
						"BindingKey":      bindingKey,
					}
					if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteBinding", nil, deleteRequest, false); e != nil {
						log.Printf("[ERROR] Failed to delete amqp binding (%s -> %s, instance %s, vhost %s): %s", sourceExchange, destinationName, instanceId, vhost, e)
						continue
					}
					log.Printf("[INFO] Delete amqp binding success: %s -> %s (instance %s, vhost %s)", sourceExchange, destinationName, instanceId, vhost)
				}
				if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
					request["NextToken"] = nextToken
				} else {
					break
				}
			}
		}
	}
	return nil
}

func testSweepAmqpStaticAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for static_account sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp static_account: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		action := "ListAccounts"
		request := map[string]interface{}{"InstanceId": instanceId}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list accounts (instance %s): %s", instanceId, err)
			continue
		}
		// $.Data is a map keyed by instanceId whose values are []account.
		data, _ := response["Data"].(map[string]interface{})
		for _, accts := range data {
			accounts, _ := accts.([]interface{})
			for _, v := range accounts {
				item, _ := v.(map[string]interface{})
				userName := amqpGetField(item, "UserName", "userName")
				createTimestamp := amqpGetField(item, "CreateTimestamp", "createTimestamp")
				accessKey := amqpGetField(item, "AccessKey", "accessKey")
				if !amqpSweepNameOk(userName) {
					continue
				}
				log.Printf("[INFO] Deleting amqp static_account: %s (instance %s)", accessKey, instanceId)
				deleteRequest := map[string]interface{}{
					"UserName":        userName,
					"CreateTimestamp": createTimestamp,
				}
				if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteAccount", nil, deleteRequest, false); e != nil {
					log.Printf("[ERROR] Failed to delete amqp static_account (%s, instance %s): %s", accessKey, instanceId, e)
					continue
				}
				log.Printf("[INFO] Delete amqp static_account success: %s (instance %s)", accessKey, instanceId)
			}
		}
	}
	return nil
}

func testSweepAmqpOpenSourceAccounts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for open_source_account sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp open_source_account: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		userNames, err := listAmqpOpenSourceAccountNames(client, instanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list open source accounts (instance %s): %s", instanceId, err)
			continue
		}
		for _, userName := range userNames {
			if !amqpSweepNameOk(userName) {
				continue
			}
			log.Printf("[INFO] Deleting amqp open_source_account: %s (instance %s)", userName, instanceId)
			deleteRequest := map[string]interface{}{
				"InstanceId": instanceId,
				"UserName":   userName,
			}
			if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteOpenSourceAccount", nil, deleteRequest, false); e != nil {
				log.Printf("[ERROR] Failed to delete amqp open_source_account (%s, instance %s): %s", userName, instanceId, e)
				continue
			}
			log.Printf("[INFO] Delete amqp open_source_account success: %s (instance %s)", userName, instanceId)
		}
	}
	return nil
}

func testSweepAmqpOpenSourcePermissions(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceIds, err := listAmqpServingInstanceIDs(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list amqp instances for open_source_permission sweep: %s", err)
		return nil
	}
	if len(instanceIds) == 0 {
		log.Printf("[INFO] amqp open_source_permission: 0 SERVING instances in region %s", region)
		return nil
	}

	for _, instanceId := range instanceIds {
		userNames, err := listAmqpOpenSourceAccountNames(client, instanceId)
		if err != nil {
			log.Printf("[ERROR] Failed to list open source accounts (instance %s): %s", instanceId, err)
			continue
		}
		for _, userName := range userNames {
			action := "ListOpenSourcePermissions"
			request := map[string]interface{}{
				"InstanceId": instanceId,
				"UserName":   userName,
			}
			var response map[string]interface{}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Failed to list open source permissions (instance %s, user %s): %s", instanceId, userName, err)
				continue
			}
			resp, err := jsonpath.Get("$.Data", response)
			if err != nil {
				log.Printf("[ERROR] ListOpenSourcePermissions jsonpath (instance %s, user %s): %s", instanceId, userName, err)
				continue
			}
			iter := func(arr []interface{}) {
				for _, v := range arr {
					item, _ := v.(map[string]interface{})
					vhost := amqpGetField(item, "Vhost", "VHost")
					if !amqpSweepNameOk(vhost) {
						continue
					}
					log.Printf("[INFO] Deleting amqp open_source_permission: %s (instance %s, user %s)", vhost, instanceId, userName)
					deleteRequest := map[string]interface{}{
						"InstanceId": instanceId,
						"UserName":   userName,
						"Vhost":      vhost,
					}
					if _, e := client.RpcPost("amqp-open", "2019-12-12", "DeleteOpenSourcePermission", nil, deleteRequest, false); e != nil {
						log.Printf("[ERROR] Failed to delete amqp open_source_permission (%s, instance %s, user %s): %s", vhost, instanceId, userName, e)
						continue
					}
					log.Printf("[INFO] Delete amqp open_source_permission success: %s (instance %s, user %s)", vhost, instanceId, userName)
				}
			}
			switch data := resp.(type) {
			case []interface{}:
				iter(data)
			case map[string]interface{}:
				for _, v := range data {
					if arr, ok := v.([]interface{}); ok {
						iter(arr)
					}
				}
			}
		}
	}
	return nil
}
