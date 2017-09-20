package msg

type MigrateMachine struct {
	SourceHost       string          `json:"source_host,omitempty"`
	TargetHost       string          `json:"target_target,omitempty"`
	ContainerName    string          `json:"container_name,omitempty"`
	SourceUser       string          `json:"source_user,omitempty"`
	TargetUser       string          `json:"target_user,omitempty"`
	ExcludePaths     []string        `json:"excluded_paths,omitempty"`
	TcpPorts         map[int16]int16 `json:"tcp_ports,omitempty"`
	ExcludedTcpPorts []int16         `json:"excluded_tcp_ports,omitempty"`
	ForceCreate      bool            `json:"force_create,omitempty"`
	DisableStart     bool            `json:"disable_start,omitempty"`
	Debug            bool            `json:"debug,omitempty"`
}

type Result struct {
	Success bool   `json:"success,omitempty"`
	Msg     string `json:"msg,omitempty"`
}
