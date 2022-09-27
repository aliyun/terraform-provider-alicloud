package common

// CredentialInf is interface for get AccessKeyID,AccessKeySecret,SecurityToken
type Credentials interface {
	GetAccessKeyID() string
	GetAccessKeySecret() string
	GetSecurityToken() string
}

// CredentialInfBuild is interface for get CredentialInf
type CredentialsProvider interface {
	GetCredentials() Credentials
}

type DefaultCredentials struct {
	AccessKeyID string
	AccessKeySecret string
	SecurityToken string
}

func (defCre *DefaultCredentials) GetAccessKeyID() string {
	return defCre.AccessKeyID
}

func (defCre *DefaultCredentials) GetAccessKeySecret() string {
	return defCre.AccessKeySecret
}

func (defCre *DefaultCredentials) GetSecurityToken() string {
	return defCre.SecurityToken
}

type DefaultCredentialsProvider struct {
	AccessKeyID string
	AccessKeySecret string
	SecurityToken string
}

func (defBuild *DefaultCredentialsProvider) GetCredentials() Credentials {
	return &DefaultCredentials{AccessKeyID: defBuild.AccessKeyID, AccessKeySecret: defBuild.AccessKeySecret, SecurityToken: defBuild.SecurityToken}
}
