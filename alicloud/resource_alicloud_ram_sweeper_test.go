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

// RAM (Resource Access Manager, product "Ram" / "2015-05-01", RPC style) gap
// sweepers — fills the children left un-swept by the 5 pre-existing inline
// sweepers (user, role, group, policy, account_alias in their own _test.go
// files). RAM is RPC-style: RpcPost(product, "2015-05-01", action,
// query=nil, body=params, autoRetry) — params go in the BODY slot, matching
// the pre-existing ram_policy sweeper.
//
// Eight gap sweepers (List -> Delete, id field):
//   access_key:              ListUsers -> ListAccessKeys(UserName) -> DeleteAccessKey; id=AccessKeyId
//   group_membership:        ListGroups -> ListUsersForGroup(GroupName) -> RemoveUserFromGroup; id=GroupName
//   group_policy_attachment: ListGroups -> ListPoliciesForGroup(GroupName) -> DetachPolicyFromGroup; id=group:PolicyName:PolicyType:GroupName
//   login_profile:           ListUsers -> GetLoginProfile(UserName) -> DeleteLoginProfile; id=UserName
//   role_policy_attachment:   ListRoles -> ListPoliciesForRole(RoleName) -> DetachPolicyFromRole; id=role:PolicyName:PolicyType:RoleName
//   saml_provider:           ListSAMLProviders (Ims 2019-08-15) -> DeleteSAMLProvider; id=SAMLProviderName
//   user_group_attachment:   ListUsers -> ListGroupsForUser(UserName) -> RemoveUserFromGroup; id=GroupName:UserName
//   user_policy_attachment:  ListUsers -> ListPoliciesForUser(UserName) -> DetachPolicyFromUser; id=user:PolicyName:PolicyType:UserName
//
// DEFERRED (no List/Delete or account-level singleton — sweeping would mutate
// account-wide security state, so they are intentionally NOT swept here):
//   access_key_policy:        Ims product; NO ListAccessKeyPolicy API (per-key
//                             GetAccessKeyPolicy only); "delete" is a reset
//                             mutation (SetAccessKeyPolicy with an empty doc),
//                             not a clean delete; child of access_key/user which
//                             are swept via the user cascade.
//   account_password_policy:  account-level singleton (SetPasswordPolicy /
//                             GetPasswordPolicy via Ram); no List, no Delete.
//   password_policy:          Ims account-level singleton (SetPasswordPolicy /
//                             GetPasswordPolicy); no List, no Delete.
//   role_attachment:          DEPRECATED since v1.250.0 (use
//                             alicloud_ecs_ram_role_attachment); ECS-side API
//                             (AttachInstanceRamRole/DetachInstanceRamRole are
//                             ECS product, not Ram); instances swept by the
//                             ECS sweeper, role swept by the role sweeper.
//   security_preference:      Ims account-level singleton (SetSecurityPreference /
//                             GetSecurityPreference); no List, no Delete.

var ramSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"scene-testing-",
	"tftestacc",
}

// ramGoneErrorCodes covers the "already gone" codes returned by RAM/Ims when a
// user/group/role/policy/profile/provider/membership no longer exists.
// NotFoundError handles 404 / RamInstanceNotFound; these codes handle the
// EntityNotExist.* family that NotFoundError misses. IsExpectedErrors matches
// on exact ErrorCode OR substring-of-Message, so the generic "EntityNotExist"
// catches the message-form while the dotted forms catch exact codes.
var ramGoneErrorCodes = []string{
	"EntityNotExist",
	"EntityNotExist.User",
	"EntityNotExist.Group",
	"EntityNotExist.Role",
	"EntityNotExist.Policy",
	"EntityNotExist.LoginProfile",
	"EntityNotExist.SAMLProvider",
	"EntityNotExist.UserOrGroup",
	"NoSuchUser",
	"NoSuchGroup",
	"NoSuchRole",
	"NoSuchPolicy",
	"NoSuchAccessKey",
}

// ramGone reports whether err means the target is already gone. Used in the
// best-effort delete path: log + continue, never abort the fleet sweep.
func ramGone(err error) bool {
	return NotFoundError(err) || IsExpectedErrors(err, ramGoneErrorCodes)
}

// ramGetField returns the first non-empty string value among names on item, as
// a plain string. Returns "" if none found or the value is "<nil>".
func ramGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// ramBool returns the bool value at key; false if missing or non-bool.
func ramBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// ramSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (1511928242963727, cn-beijing) is
// fleet-wiped, sparing only NODELETING-marked resources. protectedSweepVpcName
// is always-active in BOTH modes (it normalizes away non-alphanumeric, so e.g.
// "no-deleting-cen" -> "nodeletingcen" CONTAINS "nodeleting" -> protected).
func ramSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping ram resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range ramSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ram resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// ramRpcPost issues an RpcPost to the RAM product (2015-05-01) with retry +
// incrementalWait. Params go in the BODY slot (query=nil), matching the
// pre-existing ram_policy sweeper (RpcPost signature: product, version, action,
// query, body, autoRetry).
func ramRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// ramDeletePost issues a delete RpcPost (RAM) with retry. Best-effort callers
// log errors and continue; this never aborts the fleet sweep.
func ramDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	_, err := ramRpcPost(client, action, params)
	return err
}

// ramImsPost issues an RpcPost to the Ims product (2019-08-15) with retry. The
// saml_provider resource is implemented against Ims (not Ram), so its list +
// delete go through here.
func ramImsPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// ramImsDeletePost issues a delete RpcPost (Ims) with retry.
func ramImsDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	_, err := ramImsPost(client, action, params)
	return err
}

// ramListEntities paginates a RAM top-level list API (ListUsers / ListGroups /
// ListRoles) that returns an item array at arrayPath with IsTruncated/Marker
// pagination. Returns ALL items (callers gate each item with ramSweepNameOk).
func ramListEntities(client *connectivity.AliyunClient, action, arrayPath string) ([]interface{}, error) {
	var out []interface{}
	params := map[string]interface{}{"MaxItems": PageSizeLarge}
	for {
		response, err := ramRpcPost(client, action, params)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", action, err)
		}
		if raw, e := jsonpath.Get(arrayPath, response); e == nil {
			if items, ok := raw.([]interface{}); ok {
				out = append(out, items...)
			}
		}
		if !ramBool(response, "IsTruncated") {
			break
		}
		if marker := ramGetField(response, "Marker"); marker != "" {
			params["Marker"] = marker
			continue
		}
		break
	}
	return out, nil
}

// ramListUsersForGroup paginates ListUsersForGroup (GroupName) —
// $.Users.User[].UserName with IsTruncated/Marker pagination.
func ramListUsersForGroup(client *connectivity.AliyunClient, groupName string) ([]interface{}, error) {
	var out []interface{}
	params := map[string]interface{}{"GroupName": groupName, "MaxItems": PageSizeLarge}
	for {
		response, err := ramRpcPost(client, "ListUsersForGroup", params)
		if err != nil {
			return nil, err
		}
		if raw, e := jsonpath.Get("$.Users.User", response); e == nil {
			if items, ok := raw.([]interface{}); ok {
				out = append(out, items...)
			}
		}
		if !ramBool(response, "IsTruncated") {
			break
		}
		if marker := ramGetField(response, "Marker"); marker != "" {
			params["Marker"] = marker
			continue
		}
		break
	}
	return out, nil
}

func init() {
	resource.AddTestSweepers("alicloud_ram_access_key", &resource.Sweeper{
		Name: "alicloud_ram_access_key", F: testSweepRamAccessKeys,
	})
	resource.AddTestSweepers("alicloud_ram_group_membership", &resource.Sweeper{
		Name: "alicloud_ram_group_membership", F: testSweepRamGroupMemberships,
	})
	resource.AddTestSweepers("alicloud_ram_group_policy_attachment", &resource.Sweeper{
		Name: "alicloud_ram_group_policy_attachment", F: testSweepRamGroupPolicyAttachments,
	})
	resource.AddTestSweepers("alicloud_ram_login_profile", &resource.Sweeper{
		Name: "alicloud_ram_login_profile", F: testSweepRamLoginProfiles,
	})
	resource.AddTestSweepers("alicloud_ram_role_policy_attachment", &resource.Sweeper{
		Name: "alicloud_ram_role_policy_attachment", F: testSweepRamRolePolicyAttachments,
	})
	resource.AddTestSweepers("alicloud_ram_saml_provider", &resource.Sweeper{
		Name: "alicloud_ram_saml_provider", F: testSweepRamSAMLProviders,
	})
	resource.AddTestSweepers("alicloud_ram_user_group_attachment", &resource.Sweeper{
		Name: "alicloud_ram_user_group_attachment", F: testSweepRamUserGroupAttachments,
	})
	resource.AddTestSweepers("alicloud_ram_user_policy_attachment", &resource.Sweeper{
		Name: "alicloud_ram_user_policy_attachment", F: testSweepRamUserPolicyAttachments,
	})
}

// testSweepRamAccessKeys: ListUsers -> ListAccessKeys(UserName) ->
// DeleteAccessKey(UserAccessKeyId, UserName). Access key ids are random
// (LTAI...) and cannot be name-filtered, so the gate is on UserName: every
// access key belonging to a sweepable (tf-prefixed / SWEEP_ALL non-NODELETING)
// user is deleted. The pre-existing user sweeper already deletes keys as part
// of deleting tf users, but this standalone sweeper also clears orphaned
// access_key resources independently.
func testSweepRamAccessKeys(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	users, err := ramListEntities(client, "ListUsers", "$.Users.User")
	if err != nil {
		log.Printf("[ERROR] ram ListUsers: %s", err)
		return nil
	}
	sweeped := false
	for _, u := range users {
		item, _ := u.(map[string]interface{})
		userName := ramGetField(item, "UserName")
		if !ramSweepNameOk(userName) {
			continue
		}
		// ListAccessKeys has no pagination (a user holds at most 2 keys).
		resp, e := ramRpcPost(client, "ListAccessKeys", map[string]interface{}{
			"UserName": userName,
		})
		if e != nil {
			if ramGone(e) {
				log.Printf("[INFO] ram user %s already gone (ListAccessKeys): %s", userName, e)
			} else {
				log.Printf("[ERROR] ram ListAccessKeys (%s): %s", userName, e)
			}
			continue
		}
		raw, e := jsonpath.Get("$.AccessKeys.AccessKey", resp)
		if e != nil {
			continue
		}
		keys, _ := raw.([]interface{})
		for _, k := range keys {
			keyItem, _ := k.(map[string]interface{})
			accessKeyId := ramGetField(keyItem, "AccessKeyId")
			if accessKeyId == "" {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Deleting ram access key %s for user %s", accessKeyId, userName)
			if de := ramDeletePost(client, "DeleteAccessKey", map[string]interface{}{
				"UserName":        userName,
				"UserAccessKeyId": accessKeyId,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram access key %s already gone: %s", accessKeyId, de)
				} else {
					log.Printf("[ERROR] Failed to delete ram access key (%s): %s", accessKeyId, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamGroupMemberships: ListGroups -> ListUsersForGroup(GroupName) ->
// RemoveUserFromGroup(UserName, GroupName). The resource is group-centric (id
// is the GroupName); a tf-prefixed user is removed from a sweepable group. Both
// group and user must pass the gate.
func testSweepRamGroupMemberships(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	groups, err := ramListEntities(client, "ListGroups", "$.Groups.Group")
	if err != nil {
		log.Printf("[ERROR] ram ListGroups: %s", err)
		return nil
	}
	sweeped := false
	for _, g := range groups {
		item, _ := g.(map[string]interface{})
		groupName := ramGetField(item, "GroupName")
		if !ramSweepNameOk(groupName) {
			continue
		}
		users, e := ramListUsersForGroup(client, groupName)
		if e != nil {
			if ramGone(e) {
				log.Printf("[INFO] ram group %s already gone (ListUsersForGroup): %s", groupName, e)
			} else {
				log.Printf("[ERROR] ram ListUsersForGroup (%s): %s", groupName, e)
			}
			continue
		}
		for _, u := range users {
			userItem, _ := u.(map[string]interface{})
			userName := ramGetField(userItem, "UserName")
			if !ramSweepNameOk(userName) {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Removing ram user %s from group %s", userName, groupName)
			if de := ramDeletePost(client, "RemoveUserFromGroup", map[string]interface{}{
				"UserName":  userName,
				"GroupName": groupName,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram membership %s/%s already gone: %s", groupName, userName, de)
				} else {
					log.Printf("[ERROR] Failed to remove ram user %s from group %s: %s", userName, groupName, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamGroupPolicyAttachments: ListGroups -> ListPoliciesForGroup ->
// DetachPolicyFromGroup(PolicyType, PolicyName, GroupName). The gate is on the
// POLICY name (not the group): a tf-prefixed policy attached to a
// non-tf/system-managed group is detached so the policy can be deleted. Only
// NODELETING groups are skipped outright (never detach from protected groups).
func testSweepRamGroupPolicyAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	groups, err := ramListEntities(client, "ListGroups", "$.Groups.Group")
	if err != nil {
		log.Printf("[ERROR] ram ListGroups: %s", err)
		return nil
	}
	sweeped := false
	for _, g := range groups {
		item, _ := g.(map[string]interface{})
		groupName := ramGetField(item, "GroupName")
		if groupName == "" || protectedSweepVpcName(groupName) {
			log.Printf("[INFO] Skipping ram group (NODELETING): %s", groupName)
			continue
		}
		resp, e := ramRpcPost(client, "ListPoliciesForGroup", map[string]interface{}{
			"GroupName": groupName,
		})
		if e != nil {
			if ramGone(e) {
				continue
			}
			log.Printf("[ERROR] ram ListPoliciesForGroup (%s): %s", groupName, e)
			continue
		}
		raw, e := jsonpath.Get("$.Policies.Policy", resp)
		if e != nil {
			continue
		}
		policies, _ := raw.([]interface{})
		for _, p := range policies {
			polItem, _ := p.(map[string]interface{})
			policyName := ramGetField(polItem, "PolicyName")
			policyType := ramGetField(polItem, "PolicyType")
			if !ramSweepNameOk(policyName) {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Detaching ram policy %s (%s) from group %s", policyName, policyType, groupName)
			if de := ramDeletePost(client, "DetachPolicyFromGroup", map[string]interface{}{
				"PolicyType": policyType,
				"PolicyName": policyName,
				"GroupName":  groupName,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram group policy attachment %s/%s already gone: %s", groupName, policyName, de)
				} else {
					log.Printf("[ERROR] Failed to detach ram policy %s from group %s: %s", policyName, groupName, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamLoginProfiles: ListUsers -> GetLoginProfile(UserName) ->
// DeleteLoginProfile(UserName). GetLoginProfile returns
// EntityNotExist.LoginProfile when the user has no profile, so absence is a
// no-op. The pre-existing user sweeper does NOT delete login profiles, so this
// standalone sweeper fills a real gap (a tf user with a login profile may
// otherwise block clean user deletion).
func testSweepRamLoginProfiles(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	users, err := ramListEntities(client, "ListUsers", "$.Users.User")
	if err != nil {
		log.Printf("[ERROR] ram ListUsers: %s", err)
		return nil
	}
	sweeped := false
	for _, u := range users {
		item, _ := u.(map[string]interface{})
		userName := ramGetField(item, "UserName")
		if !ramSweepNameOk(userName) {
			continue
		}
		resp, e := ramRpcPost(client, "GetLoginProfile", map[string]interface{}{
			"UserName": userName,
		})
		if e != nil {
			if ramGone(e) {
				continue
			}
			log.Printf("[ERROR] ram GetLoginProfile (%s): %s", userName, e)
			continue
		}
		if _, ge := jsonpath.Get("$.LoginProfile.UserName", resp); ge != nil {
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting ram login profile for user %s", userName)
		if de := ramDeletePost(client, "DeleteLoginProfile", map[string]interface{}{
			"UserName": userName,
		}); de != nil {
			if ramGone(de) {
				log.Printf("[INFO] ram login profile %s already gone: %s", userName, de)
			} else {
				log.Printf("[ERROR] Failed to delete ram login profile (%s): %s", userName, de)
			}
			continue
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamRolePolicyAttachments: ListRoles -> ListPoliciesForRole ->
// DetachPolicyFromRole(PolicyType, PolicyName, RoleName). Same shape as the
// group attachment sweeper: gate on the POLICY name so a tf policy attached to
// a system-managed role is detached; only NODELETING roles are skipped outright.
func testSweepRamRolePolicyAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	roles, err := ramListEntities(client, "ListRoles", "$.Roles.Role")
	if err != nil {
		log.Printf("[ERROR] ram ListRoles: %s", err)
		return nil
	}
	sweeped := false
	for _, r := range roles {
		item, _ := r.(map[string]interface{})
		roleName := ramGetField(item, "RoleName")
		if roleName == "" || protectedSweepVpcName(roleName) {
			log.Printf("[INFO] Skipping ram role (NODELETING): %s", roleName)
			continue
		}
		resp, e := ramRpcPost(client, "ListPoliciesForRole", map[string]interface{}{
			"RoleName": roleName,
		})
		if e != nil {
			if ramGone(e) {
				continue
			}
			log.Printf("[ERROR] ram ListPoliciesForRole (%s): %s", roleName, e)
			continue
		}
		raw, e := jsonpath.Get("$.Policies.Policy", resp)
		if e != nil {
			continue
		}
		policies, _ := raw.([]interface{})
		for _, p := range policies {
			polItem, _ := p.(map[string]interface{})
			policyName := ramGetField(polItem, "PolicyName")
			policyType := ramGetField(polItem, "PolicyType")
			if !ramSweepNameOk(policyName) {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Detaching ram policy %s (%s) from role %s", policyName, policyType, roleName)
			if de := ramDeletePost(client, "DetachPolicyFromRole", map[string]interface{}{
				"PolicyType": policyType,
				"PolicyName": policyName,
				"RoleName":   roleName,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram role policy attachment %s/%s already gone: %s", roleName, policyName, de)
				} else {
					log.Printf("[ERROR] Failed to detach ram policy %s from role %s: %s", policyName, roleName, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamSAMLProviders: ListSAMLProviders (Ims 2019-08-15) ->
// DeleteSAMLProvider(SAMLProviderName). The saml_provider resource is
// implemented against Ims, so list + delete go through ramImsPost.
// $.SAMLProviders.SAMLProvider[].SAMLProviderName with IsTruncated/Marker.
func testSweepRamSAMLProviders(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	sweeped := false
	params := map[string]interface{}{"MaxItems": PageSizeLarge}
	for {
		resp, e := ramImsPost(client, "ListSAMLProviders", params)
		if e != nil {
			log.Printf("[ERROR] ims ListSAMLProviders: %s", e)
			break
		}
		raw, e := jsonpath.Get("$.SAMLProviders.SAMLProvider", resp)
		if e == nil {
			if providers, ok := raw.([]interface{}); ok {
				for _, p := range providers {
					item, _ := p.(map[string]interface{})
					name := ramGetField(item, "SAMLProviderName")
					if !ramSweepNameOk(name) {
						continue
					}
					sweeped = true
					log.Printf("[INFO] Deleting ram saml provider %s", name)
					if de := ramImsDeletePost(client, "DeleteSAMLProvider", map[string]interface{}{
						"SAMLProviderName": name,
					}); de != nil {
						if ramGone(de) {
							log.Printf("[INFO] ram saml provider %s already gone: %s", name, de)
						} else {
							log.Printf("[ERROR] Failed to delete ram saml provider (%s): %s", name, de)
						}
						continue
					}
				}
			}
		}
		if !ramBool(resp, "IsTruncated") {
			break
		}
		if marker := ramGetField(resp, "Marker"); marker != "" {
			params["Marker"] = marker
			continue
		}
		break
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamUserGroupAttachments: ListUsers -> ListGroupsForUser(UserName) ->
// RemoveUserFromGroup(UserName, GroupName). The resource is user-centric (id
// is GroupName:UserName); a sweepable user is removed from sweepable groups.
// ListGroupsForUser has no pagination ($.Groups.Group[].GroupName). Both user
// and group must pass the gate.
func testSweepRamUserGroupAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	users, err := ramListEntities(client, "ListUsers", "$.Users.User")
	if err != nil {
		log.Printf("[ERROR] ram ListUsers: %s", err)
		return nil
	}
	sweeped := false
	for _, u := range users {
		item, _ := u.(map[string]interface{})
		userName := ramGetField(item, "UserName")
		if !ramSweepNameOk(userName) {
			continue
		}
		resp, e := ramRpcPost(client, "ListGroupsForUser", map[string]interface{}{
			"UserName": userName,
		})
		if e != nil {
			if ramGone(e) {
				log.Printf("[INFO] ram user %s already gone (ListGroupsForUser): %s", userName, e)
			} else {
				log.Printf("[ERROR] ram ListGroupsForUser (%s): %s", userName, e)
			}
			continue
		}
		raw, e := jsonpath.Get("$.Groups.Group", resp)
		if e != nil {
			continue
		}
		groups, _ := raw.([]interface{})
		for _, g := range groups {
			groupItem, _ := g.(map[string]interface{})
			groupName := ramGetField(groupItem, "GroupName")
			if !ramSweepNameOk(groupName) {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Removing ram user %s from group %s", userName, groupName)
			if de := ramDeletePost(client, "RemoveUserFromGroup", map[string]interface{}{
				"UserName":  userName,
				"GroupName": groupName,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram user/group %s/%s already gone: %s", userName, groupName, de)
				} else {
					log.Printf("[ERROR] Failed to remove ram user %s from group %s: %s", userName, groupName, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// testSweepRamUserPolicyAttachments: ListUsers -> ListPoliciesForUser ->
// DetachPolicyFromUser(PolicyType, PolicyName, UserName). Same shape as the
// group/role attachment sweepers: gate on the POLICY name so a tf policy
// attached to a non-tf/system-managed user is detached; only NODELETING users
// are skipped outright.
func testSweepRamUserPolicyAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	users, err := ramListEntities(client, "ListUsers", "$.Users.User")
	if err != nil {
		log.Printf("[ERROR] ram ListUsers: %s", err)
		return nil
	}
	sweeped := false
	for _, u := range users {
		item, _ := u.(map[string]interface{})
		userName := ramGetField(item, "UserName")
		if userName == "" || protectedSweepVpcName(userName) {
			log.Printf("[INFO] Skipping ram user (NODELETING): %s", userName)
			continue
		}
		resp, e := ramRpcPost(client, "ListPoliciesForUser", map[string]interface{}{
			"UserName": userName,
		})
		if e != nil {
			if ramGone(e) {
				continue
			}
			log.Printf("[ERROR] ram ListPoliciesForUser (%s): %s", userName, e)
			continue
		}
		raw, e := jsonpath.Get("$.Policies.Policy", resp)
		if e != nil {
			continue
		}
		policies, _ := raw.([]interface{})
		for _, p := range policies {
			polItem, _ := p.(map[string]interface{})
			policyName := ramGetField(polItem, "PolicyName")
			policyType := ramGetField(polItem, "PolicyType")
			if !ramSweepNameOk(policyName) {
				continue
			}
			sweeped = true
			log.Printf("[INFO] Detaching ram policy %s (%s) from user %s", policyName, policyType, userName)
			if de := ramDeletePost(client, "DetachPolicyFromUser", map[string]interface{}{
				"PolicyType": policyType,
				"PolicyName": policyName,
				"UserName":   userName,
			}); de != nil {
				if ramGone(de) {
					log.Printf("[INFO] ram user policy attachment %s/%s already gone: %s", userName, policyName, de)
				} else {
					log.Printf("[ERROR] Failed to detach ram policy %s from user %s: %s", policyName, userName, de)
				}
				continue
			}
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}
