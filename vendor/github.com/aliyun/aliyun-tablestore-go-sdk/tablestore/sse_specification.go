package tablestore

import "errors"

// 表示服务器端加密的秘钥类型
type SSEKeyType int

func (t *SSEKeyType) String() string {
	switch *t {
	case SSE_KMS_SERVICE:
		return "SSE_KMS_SERVICE"
	case SSE_BYOK:
		return "SSE_BYOK"
	default:
		return ""
	}
}

const (
	// 使用KMS的服务主密钥
	SSE_KMS_SERVICE SSEKeyType = iota

	// 使用KMS的用户主密钥，支持用户自定义秘钥上传
	SSE_BYOK
)

type SSESpecification struct {
	// 是否开启服务器端加密
	Enable bool

	// 当开启服务器端加密时，该参数用于设置秘钥类型
	KeyType *SSEKeyType

	// 当开启服务器端加密且秘钥类型为BYOK时，该参数用于指定KMS用户主密钥的id
	KeyId *string

	// 当开启服务器端加密且秘钥类型为BYOK时，需要通过STS服务授权表格存储获取临时访问令牌访问传入的KMS用户主密钥，
	// 该参数用于指定为此创建的RAM角色的全局资源描述符
	RoleArn *string
}

func (sse *SSESpecification) CheckArguments() error {
	if sse == nil {
		return errors.New("SSESpecification is nil")
	}
	if sse.Enable {
		if sse.KeyType == nil {
			return errors.New("key type is required when enable is true")
		} else {
			if *sse.KeyType != SSE_BYOK {
				if sse.KeyId != nil || sse.RoleArn != nil {
					return errors.New("key id and role arn cannot be set when key type is not SSE_BYOK")
				}
			}

			if *sse.KeyType != SSE_KMS_SERVICE {
				if sse.KeyId == nil || sse.RoleArn == nil {
					return errors.New("key id and role arn are required when key type is not SSE_KMS_SERVICE")
				}
			}
		}
	} else {
		if sse.KeyType != nil {
			return errors.New("key type cannot be set when enable is false")
		}
	}

	return nil
}

func (sse *SSESpecification) SetEnable(enable bool) {
	sse.Enable = enable
}

func (sse *SSESpecification) SetKeyType(keyType SSEKeyType) {
	sse.KeyType = &keyType
}

func (sse *SSESpecification) SetKeyId(keyId string) {
	sse.KeyId = &keyId
}

func (sse *SSESpecification) SetRoleArn(roleArn string) {
	sse.RoleArn = &roleArn
}

type SSEDetails struct {
	// 是否开启服务器端加密
	Enable bool

	// 秘钥类型, 开启服务器端加密时有效
	KeyType SSEKeyType

	// 主密钥在KMS中的id, 可以根据keyId在KMS系统中对秘钥的使用情况进行审计
	// 开启服务器端加密时有效
	KeyId string

	// 授权表格存储临时访问KMS用户主密钥的全局资源描述符
	// 开启服务器端加密且秘钥类型为SSE_BYOK时有效
	RoleArn string
}
