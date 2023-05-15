package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
    "terraform-provider-m3/client"
    "terraform-provider-m3/utils"
)

// InstanceStates contains possible instance states
var InstanceStates = struct {
    Starting    string
    Stopping    string
    Stopped     string
    Running     string
    Terminating string
    Cloning     string
}{
    Starting:    "starting",
    Stopping:    "stopping",
    Stopped:     "stopped",
    Running:     "running",
    Terminating: "terminating",
    Cloning:     "cloning",
}

// It's mainly for isStateAllowed function, because I don't want to edit this something after I'll change InstanceStates
// Just transform struct to array
var reflectedInstanceStates = reflect.ValueOf(InstanceStates)
var allowedInstanceStates = make([]string, reflectedInstanceStates.NumField())

func init() {
    for i := 0; i < reflectedInstanceStates.NumField(); i++ {
        allowedInstanceStates[i] = reflectedInstanceStates.Field(i).String()
    }
}

// Instance is structure that contains all information about instance
type Instance struct {
    InstanceID        string                 `json:"instanceId"`
    Cloud             string                 `json:"cloud"`
    InstanceName      string                 `json:"instanceName"`
    TenantName        string                 `json:"tenant"`
    Region            string                 `json:"region"`
    State             string                 `json:"state"`
    Created           string                 `json:"creationDate"`
    Architecture      string                 `json:"architecture"`
    Image             string                 `json:"imageId"`
    Shape             string                 `json:"instanceType"`
    PrivateIP         string                 `json:"privateIpAddress"`
    LockedTermination bool                   `json:"lockedTermination"`
    ChefEnabled       bool                   `json:"installChefClient"`
    InstanceChefUUID  string                 `json:"insanceChefUuid"`
    ChefProfile       string                 `json:"chefProfile"`
    AvailabilityZone  string                 `json:"availabilityZone"`
    ResourceGroup     string                 `json:"resourceGroup"`
    VolumesIds        []string               `json:"volumesIds"`
    AdditionalData    map[string]interface{} `json:"additionalData"`
    Tags              []Tag                  `json:"tags"`
}

type Tag struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

//InstancesResultData is structure that contais info about instances
type InstancesResultData struct {
    Instances []Instance `json:"instances"`
}

// InstanceRunRequest request to run instance
type InstanceRunRequest struct {
    *DefaultRequestParams
    InstanceName     string `json:"instanceName"`
    KeyName          string `json:"keyName"`
    Image            string `json:"imageId"`
    Shape            string `json:"shape"`
    Owner            string `json:"owner"`
    InstancesCount   int    `json:"count"`
    ChefEnabled      bool   `json:"installChefClient"`
    InstanceChefUUID string `json:"insanceChefUuid"`
    ChefProfile      string `json:"chefProfile"`
    *StopAfter
    *TerminateAfter
    LockedTermination bool                   `json:"lockedTermination"`
    AdditionalData    map[string]interface{} `json:"additionalData"`
    Tags              map[string]interface{} `json:"tags"`
}

// StopAfter optional parameter for delayed stop
type StopAfter struct {
    StopAfter int `json:"stopAfter"`
}

// TerminateAfter optional parameter for delayed terminate
type TerminateAfter struct {
    TerminateAfter int `json:"terminateAfter"`
}

// InstanceTerminateRequest request to terminate instance
type InstanceTerminateRequest struct {
    *DefaultRequestParams
    InstanceID string `json:"instanceId"`
}

// InstanceDescribeRequest request to describe instance
type InstanceDescribeRequest struct {
    *DefaultRequestParams
    InstanceIds []string `json:"instanceIds"`
}

type InstanceUpdateTagsRequest struct {
    *DefaultRequestParams
    InstanceName     string                 `json:"instanceName"`
    Cloud            string                 `json:"cloud"`
    Id               string                 `json:"instanceId"`
    AvailabilityZone string                 `json:"availabilityZone"`
    ResourceGroup    string                 `json:"resourceGroup"`
    Overwrite        bool                   `json:"overwrite"`
    VolumeIds        []string               `json:"volumeIds"`
    Tags             map[string]interface{} `json:"tags"`
}
type InstanceDeleteTagsRequest struct {
    *DefaultRequestParams
    Cloud            string   `json:"cloud"`
    Id               string   `json:"instanceId"`
    AvailabilityZone string   `json:"availabilityZone"`
    ResourceGroup    string   `json:"resourceGroup"`
    Tags             []string `json:"tags"`
}

// InstancesService contains fields needed to implement InstancesServicer interface
type InstancesService struct {
    trans client.Transporter
}

func NewInstancesService(t client.Transporter) *InstancesService {
    return &InstancesService{trans: t}
}

// Run method is needed to Run instance
func (s *InstancesService) Run(request *InstanceRunRequest) (*Instance, error) {
    body := *request
    payload, err := s.trans.MakePayload(&body, MethodRunInstance)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    instances := InstancesResultData{}

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &instances)
        if err != nil {
            return nil, err
        }
        i := instances.Instances[0]
        if i.State != InstanceStates.Starting {
            return nil, fmt.Errorf("instance must be in '%v' state got '%v' instead",
                InstanceStates.Starting, i.State)
        }
        // success
        instance := &Instance{
            InstanceID:        i.InstanceID,
            Cloud:             i.Cloud,
            InstanceName:      i.InstanceName,
            TenantName:        i.TenantName,
            Region:            i.Region,
            State:             i.State,
            Created:           i.Created,
            Architecture:      i.Architecture,
            Image:             i.Image,
            Shape:             i.Shape,
            PrivateIP:         i.PrivateIP,
            LockedTermination: request.LockedTermination,
            ChefEnabled:       request.ChefEnabled,
            InstanceChefUUID:  request.InstanceChefUUID,
            ChefProfile:       request.ChefProfile,
            AdditionalData:    request.AdditionalData,
            Tags:              i.Tags,
        }
        instance.LockedTermination = instance.LockedTermination && utils.Contains(instance.Cloud, []interface{}{"AWS", "AZURE", "GOOGLE"})

        return instance, err

    }
    return nil, errors.New("neither 'result' nor 'error' in response")
}

// Terminate method is used to terminate instance
func (s *InstancesService) Terminate(request *InstanceTerminateRequest) error {
    payload, err := s.trans.MakePayload(request, MethodTerminateInstance)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        // Somehow there's no instance, probably someone terminated it from web UI
        if singleResult.StatusCode == 404 {
            return errors.New("404")
        }
        return fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Status != "" {
        return nil
    }

    return errors.New("neither 'result' nor 'error' in response")
}

//Describe method is used to describe instance
func (s *InstancesService) Describe(request *InstanceDescribeRequest) (*Instance, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeInstance)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    instances := InstancesResultData{}

    singleResult := r.Results[0]

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &instances)
        if err != nil {
            return nil, err
        }
        if len(instances.Instances) > 0 {
            i := instances.Instances[0]
            // success
            instance := &Instance{
                InstanceID:       i.InstanceID,
                Cloud:            i.Cloud,
                InstanceName:     i.InstanceName,
                TenantName:       i.TenantName,
                Region:           i.Region,
                State:            i.State,
                Created:          i.Created,
                Architecture:     i.Architecture,
                Image:            i.Image,
                Shape:            i.Shape,
                PrivateIP:        i.PrivateIP,
                AvailabilityZone: i.AvailabilityZone,
                ResourceGroup:    i.ResourceGroup,
                VolumesIds:       i.VolumesIds,
                Tags:             i.Tags,
            }

            return instance, err
        }

        // Somehow there's no instance, probably someone terminated it from web UI
        return nil, errors.New("404")

    }

    if singleResult.Error != "" {
        // Somehow there's no instance, probably someone terminated it from web UI
        if singleResult.StatusCode == 404 {
            return nil, errors.New("404")
        }
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

func (s *InstancesService) UnlockTermination(request *InstanceTerminateRequest) error {
    body := struct {
        InstanceTerminateRequest
        Action string `json:"action"`
    }{
        InstanceTerminateRequest: *(request),
        Action:                   "DISABLE",
    }

    payload, err := s.trans.MakePayload(body, MethodTerminationProtection)
    if err != nil {
        return err
    }

    _, err = s.trans.Do(payload)
    return err
}

func (s *InstancesService) UpdateTags(request *InstanceUpdateTagsRequest) error {
    payload, err := s.trans.MakePayload(request, MethodUpdateTags)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    if r.Results[0].Error != "" {
        return errors.New(r.Results[0].Error)
    }

    if r.Results[0].Data == "" {
        return errors.New("neither 'result' nor 'error' in response")
    }

    return nil
}

func (s *InstancesService) DeleteTags(request *InstanceDeleteTagsRequest) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteTags)
    if err != nil {
        return err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return err
    }

    if r.Results[0].Error != "" {
        return errors.New(r.Results[0].Error)
    }

    if r.Results[0].Data == "" {
        return errors.New("neither 'result' nor 'error' in response")
    }

    return nil
}
