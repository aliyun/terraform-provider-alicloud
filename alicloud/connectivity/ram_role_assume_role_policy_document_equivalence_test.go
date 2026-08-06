package connectivity

import (
	"testing"
)

func TestPolicyEquivalence(t *testing.T) {
	cases := []struct {
		name       string
		policy1    string
		policy2    string
		equivalent bool
		err        bool
	}{
		{
			name:       "Invalid policy JSON",
			policy1:    policyTest0,
			policy2:    policyTest0,
			equivalent: false,
			err:        true,
		},
		{
			name:       "Identical policy text",
			policy1:    policyTest1,
			policy2:    policyTest1,
			equivalent: true,
		},
		{
			name:       "Action block as single item array versus string",
			policy1:    policyTest2a,
			policy2:    policyTest2b,
			equivalent: false,
		},
		{
			name:       "NotAction block",
			policy1:    policyTest3a,
			policy2:    policyTest3b,
			equivalent: true,
		},
		{
			name:       "Principal block as single item array versus string",
			policy1:    policyTest4a,
			policy2:    policyTest4b,
			equivalent: false,
		},
		{
			name:       "Principal in single item array versus string",
			policy1:    policyTest5a,
			policy2:    policyTest5b,
			equivalent: true,
		},
		{
			name:       "Different principal in single item array versus string",
			policy1:    policyTest6a,
			policy2:    policyTest6b,
			equivalent: false,
		},
		{
			name:       "Different Effect",
			policy1:    policyTest7a,
			policy2:    policyTest7b,
			equivalent: false,
		},
		{
			name:       "Different Version",
			policy1:    policyTest8a,
			policy2:    policyTest8b,
			equivalent: false,
		},
		{
			name:       "Multiple principal",
			policy1:    policyTest9a,
			policy2:    policyTest9b,
			equivalent: true,
		},
		{
			name:       "Principal in array, different order",
			policy1:    policyTest10a,
			policy2:    policyTest10b,
			equivalent: false,
		},
		{
			name:       "Condition in string and array",
			policy1:    policyTest11a,
			policy2:    policyTest11b,
			equivalent: true,
		},
		{
			name:       "Multiple statement",
			policy1:    policyTest12a,
			policy2:    policyTest12b,
			equivalent: true,
		},
		{
			name:       "Multiple statement, different order",
			policy1:    policyTest13a,
			policy2:    policyTest13b,
			equivalent: false,
		},
		{
			name:       "invalid stmt key",
			policy1:    policyTest14a,
			policy2:    policyTest14b,
			equivalent: false,
			err:        true,
		},
		{
			name:       "only version",
			policy1:    policyTest15a,
			policy2:    policyTest15b,
			equivalent: false,
		},
		{
			name:       "Policyb has more keys than Policya",
			policy1:    policyTest16a,
			policy2:    policyTest16b,
			equivalent: false,
		},
		{
			name:       "Principal's key ignore case",
			policy1:    policyTest17a,
			policy2:    policyTest17b,
			equivalent: true,
		},
		{
			name:       "Principal's key invalid",
			policy1:    policyTest18a,
			policy2:    policyTest18b,
			equivalent: false,
			err:        true,
		},
		{
			name:       "Statement not in array",
			policy1:    policyTest19a,
			policy2:    policyTest19b,
			equivalent: false,
			err:        true,
		},
		{
			name:       "Principal value invalid format",
			policy1:    policyTest20a,
			policy2:    policyTest20b,
			equivalent: false,
			err:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			equal, err := AssumeRolePolicyDocumentAreEquivalentV2(tc.policy1, tc.policy2)
			if !tc.err && err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if tc.err && err == nil {
				t.Fatal("Expected error, none produced")
			}

			if equal != tc.equivalent {
				t.Fatalf("Bad: %s\n  Expected: %t\n       Got: %t\n", tc.name, tc.equivalent, equal)
			}
		})
	}
}

const policyTest0 = `{
  "Version": "2",
  "Statement": [
    {
  ]
}`

const policyTest1 = `{
 "Version": "1",
 "Statement": [
   {
     "Effect": "Allow", "Principal": {
       "Service": "actiontrail.aliyuncs.com"
     },
     "Action": "sts:AssumeRole"
   }
 ]
}`

const policyTest2a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

const policyTest2b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
      "Action": ["sts:AssumeRole"]
    }
  ]
}`

const policyTest3a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest3b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest4a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"Action": "sts:AssumeRole",
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest4b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest5a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest5b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest6a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest6b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "test.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest7a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Deny",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest7b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest8a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest8b = `{
  "Version": "2",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest9a = `{
"Statement": [
 {
   "Effect": "Allow",
   "NotAction": "sts:AssumeRole",
   "Principal": {
     "RAM":
       "acs:ram::11111:user/admin_user",
     "Service": [
       "actiontrail.aliyuncs.com"
     ]
   }
 }
],
"Version": "1"
}`

const policyTest9b = `{
"Statement": [
 {
   "Effect": "Allow",
   "NotAction": "sts:AssumeRole",
   "Principal": {
     "RAM": [
       "acs:ram::11111:user/admin_user"
     ],
     "Service": "actiontrail.aliyuncs.com"
   }
 }
],
"Version": "1"
}`

const policyTest10a = `{
"Statement": [
 {
   "Effect": "Allow",
   "NotAction": "sts:AssumeRole",
   "Principal": {
     "RAM":
       "acs:ram::11111:user/admin_user",
     "Service": [
       "actiontrail.aliyuncs.com",
		"fortest.aliyuncs.com"
     ]
   }
 }
],
"Version": "1"
}`

const policyTest10b = `{
"Statement": [
 {
   "Effect": "Allow",
   "NotAction": "sts:AssumeRole",
   "Principal": {
     "RAM": [
       "acs:ram::11111:user/admin_user"
     ],
     "Service": [
		"fortest.aliyuncs.com",
		"actiontrail.aliyuncs.com"
     ]
   }
 }
],
"Version": "1"
}`

const policyTest11a = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM":
      "acs:ram::1758326709434815:user/admin_user",
    "Service": [
      "actiontrail.aliyuncs.com"
    ]
  }
}
],
"Version": "1"
}`

const policyTest11b = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM": [
      "acs:ram::1758326709434815:user/admin_user"
    ],
    "Service": "actiontrail.aliyuncs.com"
  }
}
],
"Version": "1"
}`

const policyTest12a = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM":
      "acs:ram::1758326709434815:user/admin_user",
    "Service": [
      "actiontrail.aliyuncs.com"
    ]
  }
},
{
  "Action": "sts:AssumeRole",
  "Effect": "Allow",
  "Principal": {
    "Service": [
      "actiontrail.aliyuncs.com"
    ]
  }
}
],
"Version": "1"
}`

const policyTest12b = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM": [
      "acs:ram::1758326709434815:user/admin_user"
    ],
    "Service": "actiontrail.aliyuncs.com"
  }
},
{
  "Action": "sts:AssumeRole",
  "Effect": "Allow",
  "Principal": {
    "Service": "actiontrail.aliyuncs.com"
  }
}
],
"Version": "1"
}`

const policyTest13a = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM":
      "acs:ram::1758326709434815:user/admin_user",
    "Service": [
      "actiontrail.aliyuncs.com"
    ]
  }
},
{
  "Action": "sts:AssumeRole",
  "Effect": "Allow",
  "Principal": {
    "Service": [
      "actiontrail.aliyuncs.com"
    ]
  }
}
],
"Version": "1"
}`

const policyTest13b = `{
"Statement": [
{
  "Action": "sts:AssumeRole",
  "Effect": "Allow",
  "Principal": {
    "Service": "actiontrail.aliyuncs.com"
  }
},
{
  "Condition": {
    "StringEquals":              {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", 
		"172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "RAM": [
      "acs:ram::1758326709434815:user/admin_user"
    ],
    "Service": "actiontrail.aliyuncs.com"
  }
}
],
"Version": "1"
}`

const policyTest14a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole",
		"fortest": "test"
    }
  ]
}`

const policyTest14b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest15a = `{
  "Version": "1"
}`

const policyTest15b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest16a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

const policyTest16b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com",
		"RAM": "acs:ram::1758326709434815:user/admin_user"
      },
      "Action": ["sts:AssumeRole"]
    }
  ]
}`

const policyTest17a = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "rAm":
      "acs:ram::1758326709434815:user/admin_user",
    "sErviCe": [
      "actiontrail.aliyuncs.com"
    ]
  }
}
],
"Version": "1"
}`

const policyTest17b = `{
"Statement": [
{
  "Condition": {
    "StringEquals": {
      "saml:recipient": "https://signin.aliyun.com/saml-role/sso"
    },
	"IpAddress": {
		"acs:SourceIp": ["192.168.0.0/16", "172.16.0.0/12"]
	}
  },
  "Effect": "Allow",
  "NotAction": "sts:AssumeRole",
  "Principal": {
    "raM":
      "acs:ram::1758326709434815:user/admin_user",
    "serVICE": [
      "actiontrail.aliyuncs.com"
    ]
  }
}
],
"Version": "1"
}`

const policyTest18a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "actiontrail.aliyuncs.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

const policyTest18b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "ForTest": "actiontrail.aliyuncs.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

const policyTest19a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest19b = `{
  "Version": "1",
  "Statement": 
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
}`

const policyTest20a = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": ["actiontrail.aliyuncs.com"]
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`

const policyTest20b = `{
  "Version": "1",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": {"actiontrail.aliyuncs.com"}
      },
		"NotAction": "ram:GetRole"
    }
  ]
}`
