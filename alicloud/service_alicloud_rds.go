package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type RdsService struct {
	client *connectivity.AliyunClient
}

//
//       _______________                      _______________                       _______________
//       |              | ______param______\  |              |  _____request_____\  |              |
//       |   Business   |                     |    Service   |                      |    SDK/API   |
//       |              | __________________  |              |  __________________  |              |
//       |______________| \    (obj, err)     |______________|  \ (status, cont)    |______________|
//                           |                                    |
//                           |A. {instance, nil}                  |a. {200, content}
//                           |B. {nil, error}                     |b. {200, nil}
//                      					  |c. {4xx, nil}
//
// The API return 200 for resource not found.
// When getInstance is empty, then throw InstanceNotfound error.
// That the business layer only need to check error.
var DBInstanceStatusCatcher = Catcher{OperationDeniedDBInstanceStatus, 60, 5}

func (s *RdsService) DescribeDBInstanceById(id string) (instance *rds.DBInstanceAttribute, err error) {

	request := rds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceAttribute(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBInstanceIdNotFound, InvalidDBInstanceNameNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("DB Instance", id))
		}
		return nil, err
	}
	response, _ := raw.(*rds.DescribeDBInstanceAttributeResponse)
	addDebug(request.GetActionName(), response)
	if response == nil || len(response.Items.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("DB Instance", id))
	}

	return &response.Items.DBInstanceAttribute[0], nil
}

func (s *RdsService) DescribeDatabaseAccount(instanceId, accountName string) (ds *rds.DBInstanceAccount, err error) {

	request := rds.CreateDescribeAccountsRequest()
	request.DBInstanceId = instanceId
	request.AccountName = accountName
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *rds.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}
		response, _ = raw.(*rds.DescribeAccountsResponse)
		return nil
	}); err != nil {
		return nil, err
	}
	addDebug(request.GetActionName(), response)
	if response == nil || len(response.Accounts.DBInstanceAccount) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("Data account %s is not found in the instance %s.", accountName, instanceId))
	}
	return &response.Accounts.DBInstanceAccount[0], nil
}

func (s *RdsService) DescribeDatabaseByName(instanceId, dbName string) (ds *rds.Database, err error) {

	request := rds.CreateDescribeDatabasesRequest()
	request.DBInstanceId = instanceId
	request.DBName = dbName

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDatabases(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DBInternalError, OperationDeniedDBInstanceStatus}) {
				return resource.RetryableError(fmt.Errorf("Describe Database %s timeout and got an error %#v.", dbName, err))
			}
			if s.NotFoundDBInstance(err) || IsExceptedErrors(err, []string{InvalidDBNameNotFound}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(fmt.Sprintf("Database %s is not found in the instance %s.", dbName, instanceId)))
			}
			return resource.NonRetryableError(fmt.Errorf("Describe Databases got an error %#v.", err))
		}
		resp, _ := raw.(*rds.DescribeDatabasesResponse)
		if resp == nil || len(resp.Databases.Database) < 1 {
			return resource.NonRetryableError(GetNotFoundErrorFromString(fmt.Sprintf("Database %s is not found in the instance %s.", dbName, instanceId)))
		}
		ds = &resp.Databases.Database[0]
		return nil
	})

	return ds, err
}

func (s *RdsService) DescribeParameters(instanceId string) (ds *rds.DescribeParametersResponse, err error) {
	request := rds.CreateDescribeParametersRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeParameters(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBInstanceIdNotFound, InvalidDBInstanceNameNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("DB Instance", instanceId))
		}
		return nil, err
	}
	resp, _ := raw.(*rds.DescribeParametersResponse)
	if resp == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Rds Instance Parameter", instanceId))
	}
	return resp, err
}

func (s *RdsService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		d.Set(attribute, param)
		return nil
	}
	response, err := s.DescribeParameters(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Describe DB parameters error: %#v", err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range response.RunningParameters.DBInstanceParameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, i := range response.ConfigParameters.DBInstanceParameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	d.Set(attribute, param)
	return nil
}

func (s *RdsService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	request := rds.CreateModifyParameterRequest()
	request.DBInstanceId = d.Id()
	config := make(map[string]string)
	documented := d.Get(attribute).(*schema.Set).List()
	if len(documented) > 0 {
		for _, i := range documented {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
		_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyParameter(request)
		})
		if err != nil {
			return fmt.Errorf("update parameter got an error: %#v", err)
		}
		// wait instance parameter expect after modifying
		if err := s.WaitForDBParameter(d.Id(), DefaultTimeoutMedium, config); err != nil {
			return fmt.Errorf("WaitForDBParameter %s got error: %#v", d.Id(), err)
		}
	}
	d.SetPartial(attribute)
	return nil
}

func (s *RdsService) DescribeDBInstanceNetInfo(instanceId string) ([]rds.DBInstanceNetInfo, error) {

	request := rds.CreateDescribeDBInstanceNetInfoRequest()
	request.DBInstanceId = instanceId
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceNetInfo(request)
	})

	if err != nil {
		return nil, err
	}
	resp, _ := raw.(*rds.DescribeDBInstanceNetInfoResponse)
	if len(resp.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", instanceId))
	}

	return resp.DBInstanceNetInfos.DBInstanceNetInfo, nil
}

func (s *RdsService) DescribeDBConnection(id string) (*rds.DBInstanceNetInfo, error) {
	split := strings.Split(id, COLON_SEPARATED)
	object, err := s.DescribeDBInstanceNetInfo(split[0])

	if err != nil {
		if IsExceptedError(err, InvalidCurrentConnectionStringNotFound) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapError(err)
	}

	if object != nil {
		for _, o := range object {
			if strings.HasPrefix(o.ConnectionString, split[1]) {
				return &o, nil
			}
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("DBConnection", id)), NotFoundMsg, ProviderERROR)
}
func (s *RdsService) DescribeReadWriteSplittingConnection(instanceId string) (*rds.DBInstanceNetInfo, error) {
	resp, err := s.DescribeDBInstanceNetInfo(instanceId)
	if err != nil && !NotFoundError(err) {
		return nil, err
	}

	if resp != nil {
		for _, conn := range resp {
			if conn.ConnectionStringType != "ReadWriteSplitting" {
				continue
			}
			if conn.MaxDelayTime == "" {
				continue
			}
			if _, err := strconv.Atoi(conn.MaxDelayTime); err != nil {
				return nil, err
			}
			return &conn, nil
		}
	}

	return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have read write splitting connection.", instanceId))
}

func (s *RdsService) GrantAccountPrivilege(instanceId, account, dbName, privilege string) error {
	request := rds.CreateGrantAccountPrivilegeRequest()
	request.DBInstanceId = instanceId
	request.AccountName = account
	request.DBName = dbName
	request.AccountPrivilege = privilege

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		rq := request
		_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.GrantAccountPrivilege(rq)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Grant DB %s account %s privilege got an error: %#v.", dbName, account, err))
			}
			return resource.NonRetryableError(fmt.Errorf("Grant DB %s account %s privilege got an error: %#v.", dbName, account, err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := s.WaitForAccountPrivilege(instanceId, account, dbName, privilege, 300); err != nil {
		return fmt.Errorf("Wait for grantting DB %s account %s privilege got an error: %#v.", dbName, account, err)
	}

	return nil
}

func (s *RdsService) RevokeAccountPrivilege(instanceId, account, dbName string) error {

	request := rds.CreateRevokeAccountPrivilegeRequest()
	request.DBInstanceId = instanceId
	request.AccountName = account
	request.DBName = dbName

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ag := request
		_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.RevokeAccountPrivilege(ag)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", dbName, account, err))
			}
			return resource.NonRetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", dbName, account, err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := s.WaitForAccountPrivilegeRevoked(instanceId, account, dbName, 300); err != nil {
		return fmt.Errorf("Wait for revoking DB %s account %s privilege got an error: %#v.", dbName, account, err)
	}

	return nil
}

func (s *RdsService) ReleaseDBPublicConnection(instanceId, connection string) error {

	request := rds.CreateReleaseInstancePublicConnectionRequest()
	request.DBInstanceId = instanceId
	request.CurrentConnectionString = connection

	_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ReleaseInstancePublicConnection(request)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *RdsService) ModifyDBBackupPolicy(instanceId, backupTime, backupPeriod, retentionPeriod, backupLog, LogBackupRetentionPeriod string) error {

	request := rds.CreateModifyBackupPolicyRequest()
	request.DBInstanceId = instanceId
	request.PreferredBackupPeriod = backupPeriod
	request.BackupRetentionPeriod = retentionPeriod
	request.PreferredBackupTime = backupTime
	request.BackupLog = backupLog
	request.LogBackupRetentionPeriod = LogBackupRetentionPeriod

	_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		return err
	}

	if err := s.WaitForDBInstance(instanceId, Running, 600); err != nil {
		return err
	}
	return nil
}

func (s *RdsService) ModifyDBSecurityIps(instanceId, ips string) error {

	request := rds.CreateModifySecurityIpsRequest()
	request.DBInstanceId = instanceId
	request.SecurityIps = ips

	_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifySecurityIps(request)
	})
	if err != nil {
		return err
	}

	if err := s.WaitForDBInstance(instanceId, Running, 600); err != nil {
		return err
	}
	return nil
}

func (s *RdsService) DescribeDBSecurityIps(instanceId string) (ips []rds.DBInstanceIPArray, err error) {

	request := rds.CreateDescribeDBInstanceIPArrayListRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceIPArrayList(request)
	})
	if err != nil {
		return nil, err
	}
	resp, _ := raw.(*rds.DescribeDBInstanceIPArrayListResponse)
	return resp.Items.DBInstanceIPArray, nil
}

func (s *RdsService) GetSecurityIps(instanceId string) ([]string, error) {
	arr, err := s.DescribeDBSecurityIps(instanceId)
	if err != nil {
		return nil, err
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range arr {
		ips += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *RdsService) DescribeSecurityGroupConfiguration(instanceid string) (string, error) {
	request := rds.CreateDescribeSecurityGroupConfigurationRequest()
	request.DBInstanceId = instanceid
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeSecurityGroupConfiguration(request)
	})

	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, instanceid, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*rds.DescribeSecurityGroupConfigurationResponse)
	addDebug(request.GetActionName(), response)
	if response != nil && len(response.Items.EcsSecurityGroupRelation) > 0 {
		return response.Items.EcsSecurityGroupRelation[0].SecurityGroupId, nil
	}
	return "", nil
}

func (s *RdsService) ModifySecurityGroupConfiguration(instanceid string, groupid string) error {
	request := rds.CreateModifySecurityGroupConfigurationRequest()
	request.DBInstanceId = instanceid
	//openapi required that input "Empty" if groupid is ""
	if len(groupid) == 0 {
		groupid = "Empty"
	}
	request.SecurityGroupId = groupid
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifySecurityGroupConfiguration(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceid, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return nil
}

// return multiIZ list of current region
func (s *RdsService) DescribeMultiIZByRegion() (izs []string, err error) {
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeRegions(rds.CreateDescribeRegionsRequest())
	})
	if err != nil {
		return nil, fmt.Errorf("error to list regions not found")
	}
	resp, _ := raw.(*rds.DescribeRegionsResponse)
	regions := resp.Regions.RDSRegion

	zoneIds := []string{}
	for _, r := range regions {
		if r.RegionId == string(s.client.Region) && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, r.ZoneId)
		}
	}

	return zoneIds, nil
}

func (s *RdsService) DescribeBackupPolicy(instanceId string) (policy *rds.DescribeBackupPolicyResponse, err error) {

	request := rds.CreateDescribeBackupPolicyRequest()
	request.DBInstanceId = instanceId

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeBackupPolicy(request)
	})
	return raw.(*rds.DescribeBackupPolicyResponse), err
}

func (s *RdsService) DescribeDBInstanceMonitor(instanceId string) (monitoringPeriod int, err error) {
	raw, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		request := rds.CreateDescribeDBInstanceMonitorRequest()
		request.DBInstanceId = instanceId
		return client.DescribeDBInstanceMonitor(request)
	})
	if err != nil {
		return 0, fmt.Errorf("Error reading monitoring period of db instance: %#v", err)
	}

	resp, _ := raw.(*rds.DescribeDBInstanceMonitorResponse)
	monPeriod, err := strconv.Atoi(resp.Period)
	if err != nil {
		return 0, fmt.Errorf("Error converting (Atoi) monitoring period of db instance: %#v", err)
	}
	return monPeriod, nil
}

// WaitForInstance waits for instance to given status
func (s *RdsService) WaitForDBInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeDBInstanceById(instanceId)
		if err != nil && !NotFoundError(err) && !IsExceptedError(err, InvalidDBInstanceIdNotFound) {
			return err
		}
		if instance != nil && strings.ToLower(instance.DBInstanceStatus) == strings.ToLower(string(status)) {
			break
		}

		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("RDS Instance", instanceId))
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *RdsService) WaitForDBParameter(instanceId string, timeout int, expects map[string]string) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		response, err := s.DescribeParameters(instanceId)
		if err != nil {
			return err
		}

		var actuals = make(map[string]string)
		for _, i := range response.RunningParameters.DBInstanceParameter {
			actuals[i.ParameterName] = i.ParameterValue
		}
		for _, i := range response.ConfigParameters.DBInstanceParameter {
			actuals[i.ParameterName] = i.ParameterValue
		}

		match := true
		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("RDS Instance", instanceId))
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *RdsService) WaitForDBConnection(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBConnection(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
		if object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
	}
}

func (s *RdsService) WaitForDBReadWriteSplitting(instanceId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		_, err := s.DescribeReadWriteSplittingConnection(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}

		if err == nil {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *RdsService) WaitForAccount(instanceId string, accountName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {

		account, err := s.DescribeDatabaseAccount(instanceId, accountName)
		if err != nil {
			return err
		}

		if account != nil && account.AccountStatus == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)

	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilege(instanceId, accountName, dbName, privilege string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {

		account, err := s.DescribeDatabaseAccount(instanceId, accountName)
		if err != nil {
			return err
		}

		ready := false
		if account != nil {
			for _, dp := range account.DatabasePrivileges.DatabasePrivilege {
				if dp.DBName == dbName && dp.AccountPrivilege == privilege {
					ready = true
					break
				}
			}
		}

		if ready {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)

	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilegeRevoked(instanceId, accountName, dbName string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		account, err := s.DescribeDatabaseAccount(instanceId, accountName)
		if err != nil {
			return err
		}

		exist := false
		if account != nil {
			for _, dp := range account.DatabasePrivileges.DatabasePrivilege {
				if dp.DBName == dbName {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)

	}
	return nil
}

// turn period to TimeType
func (s *RdsService) TransformPeriod2Time(period int, chargeType string) (ut int, tt common.TimeType) {
	if chargeType == string(Postpaid) {
		return 1, common.Day
	}

	if period >= 1 && period <= 9 {
		return period, common.Month
	}

	if period == 12 {
		return 1, common.Year
	}

	if period == 24 {
		return 2, common.Year
	}
	return 0, common.Day

}

// turn TimeType to Period
func (s *RdsService) TransformTime2Period(ut int, tt common.TimeType) (period int) {
	if tt == common.Year {
		return 12 * ut
	}

	return ut

}

func (s *RdsService) flattenDBSecurityIPs(list []rds.DBInstanceIPArray) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIPList,
		}
		result = append(result, l)
	}
	return result
}

func (s *RdsService) NotFoundDBInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidDBInstanceIdNotFound, InvalidDBInstanceNameNotFound}) {
		return true
	}

	return false
}

const tagsMaxNumPerTime = 5

func (s *RdsService) setInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			if err := s.removeTags(remove, d.Id()); err != nil {
				return err
			}
		}

		if len(create) > 0 {
			if err := s.addTags(create, d.Id()); err != nil {
				return err
			}
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *RdsService) addTags(tags []Tag, instanceId string) error {
	return s.doBatchTags(s.addTagsPerTime, tags, instanceId)
}

func (s *RdsService) addTagsPerTime(tags []Tag, instanceId string) error {
	request := rds.CreateAddTagsToResourceRequest()
	request.DBInstanceId = instanceId
	request.Tags = s.tagsToString(tags)

	_, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		return client.AddTagsToResource(request)
	})
	if err != nil {
		return fmt.Errorf("AddTags got an error: %#v", err)
	}

	return nil
}

func (s *RdsService) removeTags(tags []Tag, instanceId string) error {
	return s.doBatchTags(s.removeTagsPerTime, tags, instanceId)
}

func (s *RdsService) removeTagsPerTime(tags []Tag, instanceId string) error {
	request := rds.CreateRemoveTagsFromResourceRequest()
	request.DBInstanceId = instanceId
	request.Tags = s.tagsToString(tags)

	_, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		return client.RemoveTagsFromResource(request)
	})
	if err != nil {
		return fmt.Errorf("RemoveTags got an error: %#v", err)
	}

	return nil
}

func (s *RdsService) doBatchTags(batchFunc func([]Tag, string) error, tags []Tag, instanceId string) error {
	num := len(tags)
	if num <= 0 {
		return nil
	}

	start, end := 0, 0
	for end < num {
		start = end
		end += tagsMaxNumPerTime
		if end > num {
			end = num
		}
		if err := batchFunc(tags[start:end], instanceId); err != nil {
			return err
		}
	}
	return nil
}

func (s *RdsService) describeTags(d *schema.ResourceData) (tags []Tag, err error) {
	request := rds.CreateDescribeTagsRequest()
	request.DBInstanceId = d.Id()

	raw, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		return client.DescribeTags(request)
	})
	if err != nil {
		tmp := make([]Tag, 0)
		return tmp, err
	}

	resp, _ := raw.(*rds.DescribeTagsResponse)
	return s.respToTags(resp.Items.TagInfos), nil
}

func (s *RdsService) respToTags(tagSet []rds.TagInfos) (tags []Tag) {
	result := make([]Tag, 0, len(tagSet))
	for _, t := range tagSet {
		tag := Tag{
			Key:   t.TagKey,
			Value: t.TagValue,
		}
		result = append(result, tag)
	}

	return result
}

func (s *RdsService) tagsToMap(tags []Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func (s *RdsService) ignoreTag(t Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *RdsService) tagsToString(tags []Tag) string {
	v, _ := json.Marshal(s.tagsToMap(tags))

	return string(v)
}
