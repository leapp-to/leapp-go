package api

import (
    "encoding/json"
    "net/http"
)

type CheckTargetParams struct {
    TargetHost  string  `json:"target_host,omitempty"`
    Status      bool  `json:"check_target_service_status,omitempty"`
    TargetUser  string  `json:"target_user_name,omitempty"`
}

func buildActorInput(p *CheckTargetParams) (string, error) {
    data := map[string]interface{}{
        "target_host":                  ObjValue{p.TargetHost},
        "check_target_service_status":  ObjValue{p.Status},
        "target_user_name":             ObjValue{p.TargetUser},
    }

    j, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    return string(j), nil
}

func CheckTarget(request *http.Request) (interface{}, error) {
    var params CheckTargetParams

    if err := json.NewDecoder(request.body).Decode(&params); err != nil {
        return nil, err
    }

    actorInput, err := buildActorInput(&params)
    if err != nil {
        return nil, err
    }

    c := executor.New("remote-target-check-group", actorInput)
    r := c.Execute()

    log.Println(r.Stderr)

    var out interface{}
    err = json.Unmarshal([]byte(r.Stdout), &out)
    return out, err
}

