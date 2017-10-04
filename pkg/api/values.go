package api

import "encoding/json"

//ObjValue wraps any value into a JSON object with a key "value".
type ObjValue struct {
	Value interface{} `json:"value"`
}

// PortForwarding is a mapping between a port in the container and a port exposed by the host.
type PortForwarding struct {
	Protocol    string `json:"protocol"`
	ExposedPort uint16 `json:"exposed_port"`
	Port        uint16 `json:"port"`
}

// PortForwardingSlice contains all ports that should be exposed by the host.
type PortForwardingSlice []PortForwarding

// MarshalJSON makes sure that a nil PortForwardingSlice becomes an empty JSON list.
func (f PortForwardingSlice) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]PortForwarding(f))
}

// TCPPortsUserMapping represents the high level JSON object containing all ports that should be exposed by the host.
type TCPPortsUserMapping struct {
	Ports PortForwardingSlice `json:"ports"`
}

// PortMapping represents a map of ports and the names of their respective services.
type PortMapping map[uint16]map[string]string

// MarshalJSON makes sure that a nil PortMapping becomes an empty JSON object.
func (m PortMapping) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[uint16]map[string]string(m))
}

// ExcludedTCPPorts represents the high level JSON object containing all ports that should not be exposed by the target host.
type ExcludedTCPPorts struct {
	TCP PortMapping `json:"tcp"`
	UDP PortMapping `json:"udp"`
}
