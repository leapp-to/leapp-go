package api

import (
	"encoding/json"
	"github.com/leapp-to/leapp-go/pkg/executor"
	"log"
	"net/http"
)

type PortInspectParams struct {
	TargetHost  string `json:"target_host,omitempty"`
	PortRange   string `json:"port_range,omitempty"`
	ShallowScan bool   `json:"shallow_scan,omitempty"`
}

func PortInspect(request *http.Request) (interface{}, error) {
	var portInspectParams PortInspectParams

	if err := json.NewDecoder(request.Body).Decode(&portInspectParams); err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"host": ObjValue{portInspectParams.TargetHost},
		"scan_options": map[string]interface{}{
			"shallow_scan": portInspectParams.ShallowScan,
			"port_range":   portInspectParams.PortRange}}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	exec := executor.New("port-inspect", string(json_data))
	result := exec.Execute()

	log.Println(result.Stderr)

	var out interface{}
	json.Unmarshal([]byte(result.Stdout), &out)

	return out, err
}
