package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
)

type portInspectParams struct {
	TargetHost  string `json:"target_host"`
	PortRange   string `json:"port_range"`
	ShallowScan bool   `json:"shallow_scan"`
}

func portInspectHandler(request *http.Request) (interface{}, error) {
	var params portInspectParams

	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		return "", err
	}

	d := map[string]interface{}{
		"host": ObjValue{params.TargetHost},
		"scan_options": map[string]interface{}{
			"shallow_scan": params.ShallowScan,
			"port_range":   params.PortRange,
		},
	}

	actorInput, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	c := executor.New("port-inspect", string(actorInput))
	r := c.Execute()

	log.Println(r.Stderr)

	var out interface{}
	json.Unmarshal([]byte(r.Stdout), &out)
	return out, err
}
