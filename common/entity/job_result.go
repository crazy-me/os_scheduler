package entity

// AgentPush 数据上报实体
type AgentPush struct {
	Ident  string `json:"ident"`
	Alias  string `json:"alias"`
	Metric string `json:"metric"`
	Tags   struct {
		Device string `json:"device"`
	} `json:"tags"`
	Time  int64  `json:"time"`
	Value string `json:"value"`
}

// Network 网络设备指标实体
type Network struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CpuRate string `json:"cpu_rate"`
		MemRate string `json:"mem_rate"`
		Runtime string `json:"runtime"`
	} `json:"data"`
}

// Server 服务器设备指标实体
type Server struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SshStatus        string `json:"ssh_status"`
		SnmpStatus       string `json:"snmp_status"`
		RunStatus        string `json:"run_status"`
		RunTime          string `json:"run_time"`
		PingStatus       string `json:"ping_status"`
		Loss             string `json:"loss"`
		PingResponseTime string `json:"ping_response_time"`
		CpuRate          string `json:"cpu_rate"`
		MemRate          string `json:"mem_rate"`
		DiskFullSpace    string `json:"disk_full_space"`
		DiskUsedSpace    string `json:"disk_used_space"`
		DiskFreeSpace    string `json:"disk_free_space"`
		MemFree          string `json:"mem_free"`
		MemUsed          string `json:"mem_used"`
		MemTotal         string `json:"mem_total"`
		SwapFree         string `json:"swap_free"`
		SwapUsed         string `json:"swap_used"`
		SwapTotal        string `json:"swap_total"`
	} `json:"data"`
}
