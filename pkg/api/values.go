package api

import "encoding/json"

type ObjValue struct {
	Value interface{} `json:"value"`
}

type PortForwarding struct {
	Protocol    string `json:"protocol"`
	ExposedPort uint16 `json:"exposed_port"`
	Port        uint16 `json:"port"`
}

type PortForwardingSlice []PortForwarding

func (forwarding PortForwardingSlice) MarshalJSON() ([]byte, error) {
	if forwarding == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]PortForwarding(forwarding))
}

type TCPPortsUserMapping struct {
	Ports PortForwardingSlice `json:"ports"`
}

type PortMapping map[uint16]map[string]string

func (mapping PortMapping) MarshalJSON() ([]byte, error) {
	if mapping == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[uint16]map[string]string(mapping))
}

type ExcludedTCPPorts struct {
	TCP PortMapping `json:"tcp"`
	UDP PortMapping `json:"udp"`
}
