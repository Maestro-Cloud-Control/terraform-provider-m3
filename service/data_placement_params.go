package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type DataPlacementParamsService struct {
    trans client.Transporter
}

func NewDataPlacementParamsService(t client.Transporter) *DataPlacementParamsService {
    return &DataPlacementParamsService{trans: t}
}

// PlacementParamsRequest request to receive additional placement params
type PlacementParamsRequest struct {
    *DefaultRequestParams
    Simplify bool `json:"simplify"`
}

type DataOption struct {
    Value string     `json:"value"`
    Title string     `json:"title"`
    Items []DataItem `json:"items"`
}

type DataItem struct {
    Name    string       `json:"name"`
    Options []DataOption `json:"options"`
}

func (s *DataPlacementParamsService) DataPlacementGetList(request *PlacementParamsRequest) (*[]DataItem, error) {
    payload, err := s.trans.MakePayload(request, MethodGetPlacementParameters)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    data := make([]DataItem, 0, 32)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &data)
        if err != nil {
            return nil, err
        }
        return &data, nil
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}
