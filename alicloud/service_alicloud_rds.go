package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/resource"
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
	resp, _ := raw.(*rds.DescribeDBInstanceAttributeResponse)

	if resp == nil || len(resp.Items.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("DB Instance", id))
	}

	return &resp.Items.DBInstanceAttribute[0], nil
}

func (s *RdsService) DescribeDatabaseAccount(instanceId, accountName string) (ds *rds.DBInstanceAccount, err error) {

	request := rds.CreateDescribeAccountsRequest()
	request.DBInstanceId = instanceId
	request.AccountName = accountName
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var resp *rds.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}
		resp, _ = raw.(*rds.DescribeAccountsResponse)
		return nil
	}); err != nil {
		return nil, err
	}

	if resp == nil || len(resp.Accounts.DBInstanceAccount) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("Data account %s is not found in the instance %s.", accountName, instanceId))
	}
	return &resp.Accounts.DBInstanceAccount[0], nil
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

func (s *RdsService) AllocateDBPublicConnection(instanceId, prefix, port string) error {
	request := rds.CreateAllocateInstancePublicConnectionRequest()
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = port

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if IsExceptedError(err, ConnectionOperationDenied) && IsExceptedError(err, ConnectionConflictMessage) {
				return resource.NonRetryableError(fmt.Errorf("Specified connection prefix %s has already been occupied. Please modify it and try again.", prefix))
			}
			if IsExceptedError(err, NetTypeExists) {
				connection, err := s.DescribeDBInstanceNetInfoByIpType(instanceId, Public)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("The connection string with specified prefix %s has already existed. "+
					"Please import it using ID '%s:%s' or specify a new 'connection_prefix' and try again.", prefix, instanceId, connection.ConnectionString))
			} else if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Allocate db connection got an error: %#v.", err))
			}

			return resource.NonRetryableError(fmt.Errorf("Allocate db connection got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := s.WaitForDBConnection(instanceId, Public, 300); err != nil {
		return fmt.Errorf("WaitForDBConnection got error: %#v", err)
	}
	// wait instance running after allocating
	if err := s.WaitForDBInstance(instanceId, Running, 300); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}
	return nil
}

func (s *RdsService) DescribeDBInstanceNetInfos(instanceId string) ([]rds.DBInstanceNetInfo, error) {

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

func (s *RdsService) DescribeDBInstanceNetInfoByIpType(instanceId string, ipType IPType) (*rds.DBInstanceNetInfo, error) {

	resps, err := s.DescribeDBInstanceNetInfos(instanceId)

	if err != nil {
		return nil, err
	}

	if resps == nil {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", instanceId))
	}

	for _, conn := range resps {
		if conn.IPType == string(ipType) {
			return &conn, nil
		}
	}

	return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have specified type %s connection.", instanceId, ipType))
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

func (s *RdsService) WaitForDBConnection(instanceId string, netType IPType, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		resp, err := s.DescribeDBInstanceNetInfoByIpType(instanceId, netType)
		if err != nil && !NotFoundError(err) {
			return err
		}

		if resp != nil && resp.IPType == string(netType) {
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
