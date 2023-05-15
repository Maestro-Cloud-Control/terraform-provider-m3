package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type Script struct {
    TenantName string `json:"tenantDisplayName"`
    Region     string `json:"region"`
    Alias      string `json:"alias"`
    FileName   string `json:"fileName"`
    Content    string `json:"content"`
    Owner      string `json:"owner"`
}

type ScriptResultData struct {
    Scripts Script
}

type ScriptCreateRequest struct {
    *DefaultRequestParams
    FileName      string `json:"fileName"`
    ScriptContent string `json:"content"`
    Email         string `json:"email"`
    Cloud         string `json:"cloud"`
}

type ScriptDeleteRequest struct {
    *DefaultRequestParams
    FileName []string `json:"fileName"`
    Email    string   `json:"email"`
    Cloud    string   `json:"cloud"`
}

type ScriptDescribeRequest struct {
    *DefaultRequestParams
    FileName string `json:"fileName"`
    Email    string `json:"email"`
    Cloud    string `json:"cloud"`
}

type ScriptService struct {
    trans client.Transporter
}

func NewScriptService(t client.Transporter) *ScriptService {
    return &ScriptService{trans: t}
}

func (s *ScriptService) Create(request *ScriptCreateRequest) (*Script, error) {
    payload, err := s.trans.MakePayload(request, MethodCreateScript)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    script := new(Script)
    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err := json.Unmarshal([]byte(singleResult.Data), script)
        if err != nil {
            return nil, err
        }

        return script, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

func (s *ScriptService) Delete(request *ScriptDeleteRequest) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteScript)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return errors.New(singleResult.Error)
    }

    if singleResult.Status == "SUCCESS" {
        return nil
    }

    return errors.New("neither 'result' nor 'error' in response")
}

func (s *ScriptService) Describe(request *ScriptDescribeRequest) (*Script, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeScript)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    scripts := make([]Script, 0, 4)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, errors.New(singleResult.Error)
    }

    if singleResult.Data != "" {

        err = json.Unmarshal([]byte(singleResult.Data), &scripts)
        if err != nil {
            return nil, err
        }
        for _, script := range scripts {
            if script.FileName == request.FileName {
                return &script, err
            }
        }
        return nil, errors.New("404")
        // success
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}
