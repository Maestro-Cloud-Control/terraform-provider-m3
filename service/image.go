package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "terraform-provider-m3/client"
)

//AvailableImageState is for image in Available state
var AvailableImageState = "Available"
var InUseState = "in-use"

// Image struct contains information about image
type Image struct {
    TenantName  string `json:"tenant"`
    Region      string `json:"region"`
    Alias       string `json:"alias"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CreatedDate int    `json:"createdDate"`
    ImageID     string `json:"imageId"`
    OsType      string `json:"osType"`
    ImageType   string `json:"imageType"`
    State       string `json:"imageState"`
    Cloud       string `json:"cloud"`
    Owner       string `json:"owner"`
}

// ImageCreateRequest request to create image
type ImageCreateRequest struct {
    *DefaultRequestParams
    InstanceID  string `json:"instanceId"`
    ImageName   string `json:"name"`
    Description string `json:"description"`
    Owner       string `json:"owner"`
}

// DeleteImageRequest request to remove image
type DeleteImageRequest struct {
    *DefaultRequestParams
    ImageID string `json:"imageId"`
}

// ImageDescribeRequest request to remove image
type ImageDescribeRequest struct {
    *DefaultRequestParams
    ImageIds []string `json:"imageIds"`
}

// ImageService contains fields needed to implement ImageServicer interface
type ImageService struct {
    trans client.Transporter
}

func NewImageService(t client.Transporter) *ImageService {
    return &ImageService{trans: t}
}

//Create is method to create image
func (s *ImageService) Create(request *ImageCreateRequest) (*Image, error) {
    payload, err := s.trans.MakePayload(request, MethodCreateImage)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    images := make([]Image, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &images)
        if err != nil {
            return nil, err
        }

        image := images[0]
        return &image, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

// Delete method to delete image
func (s *ImageService) Delete(request *DeleteImageRequest) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteImage)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        if strings.Contains(singleResult.Error, "no unique image found by image ID") {
            return errors.New("404")
        }
        return fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Status == "SUCCESS" {
        return nil
    }

    return errors.New("neither 'result' nor 'error' in response")
}

// Describe is method for describe image
func (s *ImageService) Describe(request *ImageDescribeRequest) (*Image, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeImage)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    images := make([]Image, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &images)
        if err != nil {
            return nil, err
        }
        for _, image := range images {
            if image.ImageID == request.ImageIds[0] {
                return &image, err
            }
        }
        return nil, errors.New("404")
    }

    if singleResult.Error != "" {
        if strings.Contains(singleResult.Error, "No unique image found by image ID") {
            return nil, errors.New("404")
        }
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }
    return nil, errors.New("neither 'result' nor 'error' in response")
}
