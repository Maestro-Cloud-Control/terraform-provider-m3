package provider

import (
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    "regexp"
    "strings"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func resourceScript() *schema.Resource {
    return &schema.Resource{
        Create:      resourceScriptCreate,
        Read:        resourceScriptRead,
        Update:      resourceScriptUpdate,
        Delete:      resourceScriptDelete,
        Description: "Upload script to the tenant's library in Maestro.",
        Schema: map[string]*schema.Schema{
            "name": {
                Type:         schema.TypeString,
                Required:     true,
                Description:  "The name of the script.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\\.\\-_A-Za-z0-9]{6,32}$"), "Invalid name"),
            },

            "content": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The content of the new script.",
            },

            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The name of the tenant to which the script will be assigned.",
            },

            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The name of the region to which the script will be uploaded.",
            },

            "extension": {
                Type:         schema.TypeString,
                Required:     true,
                Description:  "The script extension.\nAvailable values [ .sh, .bat, .cmd, .ps1 ].",
                ValidateFunc: validation.StringInSlice([]string{".sh", ".bat", ".cmd", ".ps1"}, false),
            },

            "cloud": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The cloud to which the script will be uploaded.\nAllowed values [ AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK, VSPHERE, VMWARE, YANDEX ].",
                ValidateFunc: validation.StringInSlice([]string{
                    "AWS", "AZURE", "GOOGLE", "OPEN_STACK", "YANDEX",
                    "VMWARE", "VSPHERE", "NUTANIX", "HYPERV"}, true),
            },
        },
    }
}

func resourceScriptCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceScriptError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }
    cloud, err := utils.GetCloud(d, m.Config)
    if err != nil {
        return err
    }
    if err := utils.MatchEmail(m.Config.UserIdentifier); err != nil {
        return err
    }

    opts := &service.ScriptCreateRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        Cloud:         cloud,
        FileName:      d.Get("name").(string) + d.Get("extension").(string),
        ScriptContent: d.Get("content").(string),
        Email:         m.Config.UserIdentifier,
    }

    script, err := m.Service.ScriptServicer.Create(opts)
    if err != nil {
        return err
    }
    if script == nil {
        m.Log.Info(fmt.Sprintf("some troubles with script: %s", d.Get("script_name").(string)))
    }
    return resourceScriptRead(d, meta)
}

func resourceScriptRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceScriptError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }
    cloud, err := utils.GetCloud(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.ScriptDescribeRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        FileName: d.Get("name").(string) + d.Get("extension").(string),
        Email:    m.Config.UserIdentifier,
        Cloud:    strings.ToUpper(cloud),
    }

    script, err := m.Service.ScriptServicer.Describe(opts)
    if err != nil || script == nil {
        return err
    }

    d.SetId(script.FileName)
    return nil
}

func resourceScriptDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceScriptError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }
    cloud, err := utils.GetCloud(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.ScriptDeleteRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        FileName: []string{d.Get("name").(string) + d.Get("extension").(string)},
        Email:    m.Config.UserIdentifier,
        Cloud:    cloud,
    }

    err = m.Service.ScriptServicer.Delete(opts)
    if err != nil {
        if err.Error() == "404" {
            return fmt.Errorf("script: %s not found", d.Id())
        }

        return err
    }
    d.SetId("")
    return nil
}

func resourceScriptUpdate(d *schema.ResourceData, meta interface{}) (err error) {
    defer UpdatingError.WrapP(&err)
    defer ResourceScriptError.WrapP(&err)
    if err := resourceScriptDelete(d, meta); err != nil {
        return err
    }

    if err := resourceScriptCreate(d, meta); err != nil {
        return err
    }
    return nil
}
