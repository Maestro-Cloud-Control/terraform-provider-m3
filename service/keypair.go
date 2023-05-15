package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "terraform-provider-m3/client"
)

type Keypair struct {
    TenantName  string `json:"tenant"`
    Region      string `json:"region"`
    Name        string `json:"name"`
    PublicPart  string `json:"publicPart"`
    PrivatePart string `json:"privatePart"`
    Email       string `json:"email"`
    Cloud       string `json:"cloud"`
    Fingerprint string `json:"fingerprint"`
    AllTenants  bool   `json:"allTenants"`
}

//CREATE:
//  For all tenants and clouds:
//      Name
//      Email
//      KeypairAllTenants
//      KeypairContent
//
//  For specific tenant:
//      Name
//      Email
//      KeypairTenantName
//      KeypairContent
//
//  For all tenants and specific cloud:
//      Name
//      Email
//      KeypairAllTenants
//      KeypairContent
//      KeypairCloud
//
//  For specific tenant and specific cloud:
//      Name
//      Email
//      KeypairTenantName
//      KeypairContent
//      KeypairCloud

//DELETE:
//  Name
//  Email

//READ:
//  Name
//  Email

type KeypairRequest struct {
    *KeypairTenantName
    *KeypairCloud
    *KeypairAllTenants
    *KeypairContent
    Name  string `json:"name"`
    Email string `json:"email"`
}

type KeypairTenantName struct {
    TenantName string `json:"tenantName"`
}
type KeypairCloud struct {
    Cloud string `json:"cloud"`
}
type KeypairAllTenants struct {
    AllTenants bool `json:"allTenants"`
}
type KeypairContent struct {
    Content string `json:"publicKey"`
}

type KeypairService struct {
    trans client.Transporter
}

func NewKeypairService(t client.Transporter) *KeypairService {
    return &KeypairService{trans: t}
}

func (s *KeypairService) Create(request *KeypairRequest) (*Keypair, error) {
    payload, err := s.trans.MakePayload(request, MethodCreateKeypair)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    singleResult := r.Results[0]
    keypair := &Keypair{}

    if singleResult.Error != "" {
        return nil, fmt.Errorf("%+v", singleResult.Error)
    }

    if singleResult.Data != "" {

        err := json.Unmarshal([]byte(singleResult.Data), keypair)
        if err != nil {
            return nil, err
        }

        return keypair, err
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

func (s *KeypairService) Describe(request *KeypairRequest) (*Keypair, error) {
    payload, err := s.trans.MakePayload(request, MethodDescribeKeypair)
    if err != nil {
        return nil, err
    }

    r, err := s.trans.Do(payload)
    if err != nil {
        return nil, err
    }

    keypairs := make([]Keypair, 0, 2)

    singleResult := r.Results[0]

    if singleResult.Error != "" {
        return nil, errors.New(singleResult.Error)
    }

    if singleResult.Data != "" {
        err = json.Unmarshal([]byte(singleResult.Data), &keypairs)
        if err != nil {
            return nil, err
        }

        for _, keypair := range keypairs {
            if keypair.Name == request.Name {
                return &keypair, nil
            }
        }

        return nil, errors.New("404")
    }

    return nil, errors.New("neither 'result' nor 'error' in response")
}

func (s *KeypairService) Delete(request *KeypairRequest) error {
    payload, err := s.trans.MakePayload(request, MethodDeleteKeypair)
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

    if singleResult.Status != "" {
        return nil
    }

    return errors.New("neither 'result' nor 'error' in response")
}
