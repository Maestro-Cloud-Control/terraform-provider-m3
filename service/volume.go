package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "terraform-provider-m3/client"
)

// AvailableVolumeState is for volume in available state
var AvailableVolumeState = "available"

// Volume contains information about fields of volume
type Volume struct {
    TenantName string `json:"tenantName"`
    Region     string `json:"regionName"`
    Name       string `json:"volumeName"`
    VolumeID   string `json:"volumeId"`
    State      string `json:"state"`
    System     bool   `json:"system"`
    SizeLabel  int    `json:"sizeInGb"`
}

// VolumeCreateRequest request to create volume
type VolumeCreateRequest struct {
    *DefaultRequestParams
    VolumeName string `json:"volumeName"`
    SizeInGB   int    `json:"sizeInGB"`
}

// VolumeCreateAndAttachRequest request to create and attach volume
type VolumeCreateAndAttachRequest struct {
    *DefaultRequestParams
    VolumeName string `json:"volumeName"`
    SizeInGB   int    `json:"sizeInGB"`
    InstanceId string `json:"instanceId"`
}

// VolumeDeleteRequest request to remove volume
type VolumeDeleteRequest struct {
    *DefaultRequestParams
    VolumeID string `json:"volumeId"`
}

// VolumeDescribeRequest request to describe volume
type VolumeDescribeRequest struct {
    *DefaultRequestParams
    VolumeIds  []string `json:"volumesIds"`
    InstanceId string   `json:"instanceId"`
}

// VolumeService contains fields needed to implement VolumeServicer interface
type VolumeService struct {
    trans client.Transporter
}

func NewVolumeService(t client.Transporter) *VolumeService {
    return &VolumeService{trans: t}
}

// Create is method to create volume
func (s *VolumeService) Create(request *VolumeCreateRequest) (*Volume, error) {
    payload, err := s.trans.MakePayload(request, MethodCreateVolume)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    volumes := make([]Volume, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {

        err = json.Unmarshal([]byte(singleResult.Data), &volumes)
        if err != nil {
            return nil, err
        }
        if len(volumes) != 1 {
            return nil, fmt.Errorf("volume with name '%s' is not created", request.VolumeName)
        }
        volume := volumes[0]
        return &volume, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

// CreateAndAttach is method to create and attach volume
func (s *VolumeService) CreateAndAttach(request *VolumeCreateAndAttachRequest) (*Volume, error) {
    payload, err := s.trans.MakePayload(request, MethodCreateAndAttachVolume)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    volumes := make([]Volume, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &volumes)
        if err != nil {
            return nil, err
        }
        if len(volumes) != 1 {
            return nil, fmt.Errorf("volume with name '%s' is not created", request.VolumeName)
        }
        volume := volumes[0]
        return &volume, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

// Delete method to delete volume
func (s *VolumeService) Delete(request *VolumeDeleteRequest) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteVolume)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        if strings.Contains(singleResult.Error, "No unique volume found by volume ID") {
            return errors.New("404")
        }
        return fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Status == "SUCCESS" {
        return nil
    }
    return errors.New("neither 'result' nor 'error' in response")
}

// Describe describes volume
func (s *VolumeService) Describe(request *VolumeDescribeRequest) (*Volume, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeVolume)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    volumes := make([]Volume, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &volumes)
        if err != nil {
            return nil, err
        }
        for _, volume := range volumes {
            if volume.VolumeID == request.VolumeIds[0] {
                return &volume, err
            }
        }
        return nil, errors.New("404")
        // success
    }

    if singleResult.Error != "" {
        if strings.Contains(singleResult.Error, "No unique volume found by volume ID") {
            return nil, errors.New("404")
        }
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }
    return nil, errors.New("neither 'result' nor 'error' in response")
}
