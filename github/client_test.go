package github

import "testing"

func TestUpdateCommitStatus(t *testing.T) {
	integrationID := 2335
	installationID := 24709
	privKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAozR8ksUtwP+HHvQnbk5ugpW65L1+qc1ClfPMZZUZUA0M1NeH
eFEHJ3i57ArltqcqpWGGqx6Bj8vBLKT23lFGb0wGkJ9xKdQ0Y1TKuufduStyHuKQ
E/ziBte0jRG5hZCIqPTd2rgh/SxussH7GobR/1sAv8xuFa51/TfXdx5Kr+UFYbZ0
oblYDZkJ9wQkJ83Nkhs88zLTJJXIDISt2/RlD27dUGIl12ZvUnk22jE/RQavi7ke
NRZaiwF+XdCgIffp4S2lbFU9i7I14U167AnHBb3lwP3UCjuQWNwCex5RW/tOx5n7
BXvJBLqJoBnAoc6xc1WZdns4GumamiLV0lJu/wIDAQABAoIBADBQVGR6nL3ap1yB
eL7N1ycvPiGC+2E3E3QitSqJRfINuvOzTjZO/mgv9BItc61rtVM6GMzFfiNcdTZL
K43Kz6gZAISFvtFcMZoKOu2zwE6VzvXXqkFQhnTyHE/6pEom7F3gB2E9S0qQgcDH
bhd/O+F1RjaBRUGD4QfyJQXtYLzK/ZO53lNgCU/lm1ZmFmij7fGvctDaasEC84g6
5bczuTuw/krFUm+Sh0Zeouv8vIoqZM8xU3p8TPaGWYUrg1P0qtlOFde+x1sIwqrZ
pN7CNaWcKsj6/6NlXsH8HgiYIb+dNPGYApttZTjlR2OadoIqsF3b7gHaNgOyJbIM
yIymp9ECgYEAzdODhWfbcW2NiRKcFr30F1PDGPBtMr2rSdXdrg9TDomAz4bKnnBz
Ued7zXC7il7qdvnLTqsqTw6MetyzW13IXz+JksSxXTyPZfy12PI+aHRnRptCu+9F
AYaIg8opnAKLUghX/7yLK2C2WVWOnhGsDjB2wI6fQqrD9PCd2BzbAHkCgYEAyv03
q74asf2MEwrJIcpbCyIgttfOK0aa1vLjq+PIMSzfXima8Zk17JjYKGrZjDgoqhSg
JerAViI5Za7V0sDwk54hwVJnd57zww8kiOLpsya8oOaAl3ktPJxfZY7lm8sq3Q4L
/1R/tYiJj9ozmLscBRirGPVDPULA4ijYySxMvTcCgYEAs/Ar1N7822bZC3J3fvJF
iHcz4oOeE3P5YS1VBaxkAht3vvWqAWVxi7MBapMZgViFRcoPUREWhdLEQUzciA2u
9IYJcYP/QvGEs7aAC8+Le6n396QYbVA6VaEVi5GbWsZmoiqlM+/TAvMjt2myqsHs
VuNLjf+hf5jmgyYv+BUR6JECgYA3C6vJGuhKVCNkFoysaR9/SWXtr1/tRFxA5eTv
e/mRvEVmV4n48j85Rcl4TGFqMOB2Htm+7oXx1Z4TAPJjEIcswLkOn7YHLkeUIcsa
g840EtEcIOXGLcoioZUNCU8ijFm3UFPYjaWEKN6E7/sF89eJWkMrpXbyaeO8cK84
/pZgyQKBgD/buwpbtKEDN9vPr0gnLtG+V5Kk2ZTIaedhSjVdUOH7YsuQR6WVi1aV
VJX6sUs8qTQXuG7cUj50s+1uXI0KSppcYV3gY8GU3s8Nw8kVIHxg7fOYPTSJvoLs
20F3vYCnctTnl1kT11lTFjeL3WDzpObiUf1LYp7KqHqwmEsToHL4
-----END RSA PRIVATE KEY-----
`)

	client, err := NewClient(nil, integrationID, installationID, privKey)

	build := BuildCallback{
		Username:    "vsaveliev",
		Repository:  "dig",
		CommitHash:  "68d958691319d732fa4887917b6442e6ef967e32",
		State:       "failure",
		BuildURL:    "https://example.com/build/status",
		Description: "The tested code works!",
		Context:     "continuous-integration/jenkins",
	}

	err = client.UpdateCommitStatus(build)
	if err != nil {
		t.Fatalf("cannot update commit status: %s", err)
	}
}
