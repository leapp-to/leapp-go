package api

import (
	"encoding/json"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type portInspectParams struct {
	TargetHost  string `json:"target_host"`
	PortRange   string `json:"port_range"`
	ShallowScan bool   `json:"shallow_scan"`
}

func portInspectCmd(request *http.Request) (*executor.Command, error) {
	var params portInspectParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return nil, err
	}

	d := map[string]interface{}{
		"host": ObjValue{params.TargetHost},
		"scan_options": map[string]interface{}{
			"shallow_scan": params.ShallowScan,
			"port_range":   params.PortRange,
			"force_nmap":   !params.ShallowScan,
		},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	c := executor.New("portscan", string(actorInput))
	return c, nil
}
