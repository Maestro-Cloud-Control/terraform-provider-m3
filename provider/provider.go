package provider

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/zeebo/errs"
    "terraform-provider-m3/client"
    "terraform-provider-m3/logger"
    "terraform-provider-m3/service"
)

var (
    CreatingError            = errs.Class("Creating")
    ReadingError             = errs.Class("Reading")
    UpdatingError            = errs.Class("Updating")
    DeletingError            = errs.Class("Deleting")
    DataChefError         = errs.Class("data_chef")
    DataImageError           = errs.Class("data_image")
    DataPlacementParamsError = errs.Class("data_placement_params")
    ResourceImageError       = errs.Class("resource_image")
    ResourceInstanceError    = errs.Class("resource_instance")
    ResourceKeypairError     = errs.Class("resource_keypair")
    ResourceScheduleError    = errs.Class("resource_schedule")
    ResourceScriptError      = errs.Class("resource_script")
    ResourceVolumeError      = errs.Class("resource_volume")
)

type Meta struct {
    Service *service.Service
    Config  *client.Config
    Log     logger.Log
}

func newMeta(s *service.Service, conf *client.Config, log logger.Log) *Meta {
    return &Meta{Service: s, Config: conf, Log: log}
}

func Provider() *schema.Provider {

    return &schema.Provider{

        Schema: map[string]*schema.Schema{
            "url": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "URL to Maestro3 API.",
            },
            "access_key": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The access key.",
            },
            "secret_key": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The secret key.",
            },
            "user_identifier": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The user identifier, preferably email.",
            },
            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The tenant name.",
            },
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The region name.",
            },
            "cloud": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The cloud. \nAllowed values: [AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK, VSPHERE, VMWARE, YANDEX].",
            },
        },
        ResourcesMap: map[string]*schema.Resource{
            "m3_instance": resourceInstance(),
            "m3_image":    resourceImage(),
            "m3_volume":   resourceVolume(),
            "m3_script":   resourceScript(),
            "m3_schedule": resourceSchedule(),
            "m3_keypair":  resourceKeypair(),
        },
        DataSourcesMap: map[string]*schema.Resource{
            "m3_data_image": dataImage(),
            "m3_data_chef":  dataChef(),
            "m3_data_placement_params": dataPlacementParams(),
        },
        ConfigureFunc: providerConfigure,
    }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
    conf := client.NewConfig(
        d.Get("url").(string),
        d.Get("user_identifier").(string),
        d.Get("access_key").(string),
        d.Get("secret_key").(string),
        d.Get("tenant").(string),
        d.Get("region").(string),
        d.Get("cloud").(string),
    )
    c := client.NewClient(conf)
    s := service.NewService(c)
    m := newMeta(s, conf, logger.NewTFLog())
    return m, nil
}
