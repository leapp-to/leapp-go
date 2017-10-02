package api

type ObjValue struct {
	Value interface{} `json:"value"`
}

// TODO: TCPPortsUserMapping and ExcludedTCPPorts types are here just to satisfy the actor's input validation until we refactor them (in there respective actors)

type TCPPortsUserMapping struct {
	Ports []string `json:"ports"`
}

type ExcludedTCPPorts struct {
	TCP map[uint16]uint16 `json:"tcp"`
	UDP map[uint16]uint16 `json:"udp,omitempty"`
}
