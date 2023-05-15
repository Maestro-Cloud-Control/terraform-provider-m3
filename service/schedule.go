package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type Schedule struct {
    *DefaultRequestParams
    Name          string                     `json:"displayName"`
    ScheduleName  string                     `json:"scheduleName"`
    ScheduleOwner string                     `json:"scheduleOwner"`
    Cron          string                     `json:"cron"`
    Description   string                     `json:"description"`
    Action        string                     `json:"action"`
    Type          string                     `json:"scheduleType"`
    Cloud         string                     `json:"cloud"`
    Instances     []*RequestScheduleInstance `json:"instances"`
    Tag           *RequestScheduleTag        `json:"tag"`
}

type RequestSchedule struct {
    *DefaultRequestParams
    Name         string                     `json:"displayName"`
    ScheduleName string                     `json:"scheduleName"`
    Cron         string                     `json:"cron"`
    Description  string                     `json:"description"`
    Action       string                     `json:"action"`
    Type         string                     `json:"scheduleType"`
    Cloud        string                     `json:"cloud"`
    Instances    []*RequestScheduleInstance `json:"instances"`
    Tag          *RequestScheduleTag        `json:"tag"`
}

type RequestScheduleTag struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type RequestScheduleInstanceLocationInfo struct {
    Region string `json:"region"`
}

type RequestScheduleInstance struct {
    InstanceId           string                               `json:"instanceId"`
    InstanceLocationInfo *RequestScheduleInstanceLocationInfo `json:"instanceLocationInfo"`
}

type ScheduleService struct {
    trans client.Transporter
}

func NewScheduleService(t client.Transporter) *ScheduleService {
    return &ScheduleService{trans: t}
}

func (s *ScheduleService) Create(request *RequestSchedule) (*Schedule, error) {
    req := struct {
        Schedule *RequestSchedule `json:"schedule"`
    }{request}
    payload, err := s.trans.MakePayload(req, MethodCreateSchedule)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    singleResult := r.Results[0]
    schedule := Schedule{}

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {
        err := json.Unmarshal([]byte(singleResult.Data), &schedule)
        if err != nil {
            return nil, err
        }

        schedule.Tag = request.Tag
        return &schedule, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

func (s *ScheduleService) Delete(request *RequestSchedule) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteSchedule)
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

func (s *ScheduleService) Describe(request *RequestSchedule) (*Schedule, error) {
    req := struct {
        *DefaultRequestParams
        Cloud string `json:"cloud"`
    }{
        DefaultRequestParams: request.DefaultRequestParams,
        Cloud:                request.Cloud,
    }
    payload, err := s.trans.MakePayload(req, MethodDescribeSchedule)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    singleResult := r.Results[0]
    schedules := make([]Schedule, 0, 2)

    if singleResult.Error != "" {
        return nil, errors.New(singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &schedules)
        if err != nil {
            return nil, err
        }
        for _, schedule := range schedules {
            if schedule.Name == request.Name {
                return &schedule, nil
            }
        }
        return nil, errors.New("404")
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}
