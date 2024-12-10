package service

import "terraform-provider-m3/client"

//go:generate go install github.com/golang/mock/mockgen@v1.6.0
//go:generate mockgen -source ./service/service.go -destination ./service/mock/service_mock.go -package smock

// DefaultRequestParams contains fields that every request needs to be executed correctly
type DefaultRequestParams struct {
    Region     string `json:"region"`
    TenantName string `json:"tenantName"`
}

// VolumeServicer interface that provides methods to work with volumes
type VolumeServicer interface {
    Create(*VolumeCreateRequest) (*Volume, error)
    CreateAndAttach(*VolumeCreateAndAttachRequest) (*Volume, error)
    Delete(*VolumeDeleteRequest) error
    Describe(*VolumeDescribeRequest) (*Volume, error)
}

// ScriptServicer interface that provides methods to work with scripts
type ScriptServicer interface {
    Create(*ScriptCreateRequest) (*Script, error)
    Delete(*ScriptDeleteRequest) error
    Describe(*ScriptDescribeRequest) (*Script, error)
}

// ScheduleServicer interface that provides methods to work with schedules
type ScheduleServicer interface {
    Create(*RequestSchedule) (*Schedule, error)
    Delete(*RequestSchedule) error
    Describe(*RequestSchedule) (*Schedule, error)
}

// KeypairServicer interface that provides methods to work with keypairs
type KeypairServicer interface {
    Create(*KeypairRequest) (*Keypair, error)
    Delete(*KeypairRequest) error
    Describe(*KeypairRequest) (*Keypair, error)
}

// InstanceServicer interface that provides methods to work with instances
type InstanceServicer interface {
    Run(*InstanceRunRequest) (*Instance, error)
    Terminate(*InstanceTerminateRequest) error
    Describe(*InstanceDescribeRequest) (*Instance, error)
    UnlockTermination(*InstanceTerminateRequest) error
    UpdateTags(*InstanceUpdateTagsRequest) error
    DeleteTags(*InstanceDeleteTagsRequest) error
}

// ImageServicer interface that provides methods to work with images
type ImageServicer interface {
    Create(*ImageCreateRequest) (*Image, error)
    Delete(*DeleteImageRequest) error
    Describe(*ImageDescribeRequest) (*Image, error)
}

// DataImageServicer interface that provides methods to work with DataImages
type DataImageServicer interface {
    DataImageGetList(*DefaultRequestParams) (*[]Image, error)
}

// DataPlacementServicer interface that provides methods to work with DataPlacementParams
type DataPlacementServicer interface {
    DataPlacementGetList(request *PlacementParamsRequest) (*[]DataItem, error)
}

// DataChefServicer interface that provides methods to work with ChefProfiles
type DataChefServicer interface {
    DataChefGetList(request *DefaultRequestParams) (*DataChef, error)
}

type Service struct {
    VolumeServicer
    ScriptServicer
    ScheduleServicer
    KeypairServicer
    InstanceServicer
    ImageServicer
    DataImageServicer
    DataPlacementServicer
    DataChefServicer
}

func NewService(c *client.Client) *Service {
    return &Service{
        VolumeServicer:        NewVolumeService(c.Transporter),
        ScriptServicer:        NewScriptService(c.Transporter),
        ScheduleServicer:      NewScheduleService(c.Transporter),
        KeypairServicer:       NewKeypairService(c.Transporter),
        InstanceServicer:      NewInstancesService(c.Transporter),
        ImageServicer:         NewImageService(c.Transporter),
        DataImageServicer:     NewDataImageService(c.Transporter),
        DataPlacementServicer: NewDataPlacementParamsService(c.Transporter),
        DataChefServicer:      NewDataChefService(c.Transporter),
    }
}
