package convoy

type SDSResp struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Tags      Tags   `json:"tags"`
}

type Tags struct {
	AZ                  string `json:"az,omitempty"`
	Canary              string `json:"canary,omitempty"`
	LoadBalancingWeight int    `json:"load_balancing_weight,omitempty"`
}
