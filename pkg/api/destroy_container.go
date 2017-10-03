package api

import (
	"encoding/json"
	"github.com/leapp-to/leapp-go/pkg/executor"
	"log"
	"net/http"
)

type DestroyContainerParams struct {
	ContainerName string `json:"container_name"`
	TargetHost    string `json:"target_host"`
	TargetUser    string `json:"target_user"`
}

func DestroyContainer(request *http.Request) (interface{}, error) {
	var destroyParams DestroyContainerParams

	if err := json.NewDecoder(request.Body).Decode(&destroyParams); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"container_name":   ObjValue{destroyParams.ContainerName},
		"target_host":      ObjValue{destroyParams.TargetHost},
		"target_user_name": ObjValue{destroyParams.TargetUser}}

	json_data, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	exec := executor.New("destroy-container", string(json_data))
	result := exec.Execute()

	log.Println(result.Stderr)

	var out interface{}
	err = json.Unmarshal([]byte(result.Stdout), &out)

	return out, err
}
