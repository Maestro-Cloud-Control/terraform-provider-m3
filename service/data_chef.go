package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type DataChef struct {
    ServerId string `json:"serverId"`
    Roles    []struct {
        RoleName           string   `json:"roleName"`
        MinCpu             int      `json:"minCpu"`
        MinMemoryMb        int      `json:"minMemoryMb"`
        RequiredParameters []string `json:"requiredParameters,omitempty"`
    } `json:"roles"`
    Zones            []string `json:"zones"`
    ChefOrganization string   `json:"chefOrganization"`
}

type DataChefService struct {
    trans client.Transporter
}

func NewDataChefService(t client.Transporter) *DataChefService {
    return &DataChefService{trans: t}
}

func (s *DataChefService) DataChefGetList(request *DefaultRequestParams) (*DataChef, error) {
    payload, err := s.trans.MakePayload(request, MethodGetChefProfiles)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    var dataChef DataChef
    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &dataChef)
        if err != nil {
            return nil, err
        }
        return &dataChef, nil
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}
