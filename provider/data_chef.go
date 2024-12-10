package provider

import (
    "errors"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func dataChef() *schema.Resource {

    return &schema.Resource{
        Read:        DataChefCreate,
        Description: "The Data Chef resource is used for specifying chef profiler for instances.",
        Schema: map[string]*schema.Schema{
            "tenant": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The tenant name.",
            },
            "region": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The region name.",
            },
            "profile": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The profile name.",
            },
            "parameters": {
                Type:        schema.TypeMap,
                Required:    true,
                ForceNew:    true,
                Description: "Parameters for setup chef.",
            },
        },
    }
}

func DataChefCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer DataChefError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }

    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.DefaultRequestParams{
        TenantName: tenant,
        Region:     region,
    }

    chefs, err := m.Service.DataChefGetList(opts)
    if err != nil {
        return err
    }

    if chefs == nil {
        return errors.New("empty result")
    }

    filter := struct {
        Chef       string
        Parameters map[string]interface{}
    }{
        Chef:       d.Get("profile").(string),
        Parameters: d.Get("parameters").(map[string]interface{}),
    }

    for _, role := range chefs.Roles {
        if role.RoleName == filter.Chef {
            for _, reqParameter := range role.RequiredParameters {
                if _, consist := filter.Parameters[reqParameter]; !consist {
                    return fmt.Errorf("not enough required parameters, required parameters: %s", reqParameter)
                }
            }

            d.SetId(role.RoleName)
            return nil
        }
    }

    return fmt.Errorf("chef profile did not find, aloved values: %v", chefs.Roles)
}
