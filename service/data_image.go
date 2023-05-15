package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type DataImageService struct {
    trans client.Transporter
}

func NewDataImageService(t client.Transporter) *DataImageService {
    return &DataImageService{trans: t}
}

func (s *DataImageService) DataImageGetList(request *DefaultRequestParams) (*[]Image, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeImage)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    images := make([]Image, 0, 32)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &images)
        if err != nil {
            return nil, err
        }
        return &images, nil
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}
