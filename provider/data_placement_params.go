package provider

import (
    "encoding/json"
    "errors"
    "github.com/google/uuid"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func dataPlacementParams() *schema.Resource {

    return &schema.Resource{
        Read:        DataParamsCreate,
        Description: "The Data placement params resource is used for describing native ids for vm placement.",
        Schema: map[string]*schema.Schema{
            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The tenant name.",
            },
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The region name.",
            },
            "cluster_name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The cluster name.",
            },
            "datastore_name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The datastore name.",
            },
            "folder_name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The folder name.",
            },
            "resource_pool_name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The resource pool name.",
            },
            "host_name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The host name.",
            },
            "placement_data": {
                Type:        schema.TypeMap,
                Optional:    true,
                ForceNew:    true,
                Computed:    true,
                Description: "Placement data with native resource ids",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

var (
    PlacementMapping = map[string]string{
        "clusterId":   "cluster_name",
        "datastoreId": "datastore_name",
        "folderId":    "folder_name",
        "resPoolId":   "resource_pool_name",
        "hostId":      "host_name",
    }
)

func processDataItem(item *service.DataItem, d *schema.ResourceData, dataMap map[string]string) error {
    if item == nil {
        return nil
    }
    paramTitle := d.Get(PlacementMapping[item.Name]).(string)
    if len(paramTitle) == 0 {
        return errors.New("Required param " + PlacementMapping[item.Name] + " not specified")
    }
    var processed = false
    for _, option := range item.Options {
        if option.Title == paramTitle {
            processed = true
            dataMap[item.Name] = option.Value
            if option.Items != nil && len(option.Items) > 0 {
                for _, nestedItem := range option.Items {
                    err := processDataItem(&nestedItem, d, dataMap)
                    if err != nil {
                        return err
                    }
                }
            }
        }
    }
    if !processed {
        return errors.New("Failed to find parameter where " + PlacementMapping[item.Name] + " is " + paramTitle)
    }
    return nil
}

func DataParamsCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer DataPlacementParamsError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.PlacementParamsRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        Simplify: true,
    }

    placementData, err := m.Service.DataPlacementGetList(opts)
    if err != nil {
        return err
    }

    if placementData == nil {
        return errors.New("empty result")
    }
    dataMap := make(map[string]string)
    for _, item := range *placementData {
        err = processDataItem(&item, d, dataMap)
        if err != nil {
            return err
        }
    }
    b, _ := json.MarshalIndent(dataMap, "", "\t")

    m.Log.Info("Placement params:\n", string(b))
    err = d.Set("placement_data", dataMap)
    if err != nil {
        return err
    }
    id := uuid.New()
    d.SetId(id.String())
    return nil
}
