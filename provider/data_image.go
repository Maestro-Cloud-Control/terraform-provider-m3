package provider

import (
    "encoding/json"
    "errors"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    "regexp"
    "strings"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func dataImage() *schema.Resource {

    return &schema.Resource{
        Read:        DataImageCreate,
        Description: "The Data Image resource is used for specifying images used for creating new images.",
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
            "name": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The image name, can be regular expression.",
            },
            "os_type": {
                Type:         schema.TypeString,
                Optional:     true,
                ForceNew:     true,
                Default:      "",
                Description:  "OS type. \nAllowed values [ L, W ], Linux or Windows.",
                ValidateFunc: validation.StringInSlice([]string{"l", "w"}, true),
            },
            "owner": {
                Type:         schema.TypeString,
                Optional:     true,
                ForceNew:     true,
                Default:      "",
                Description:  "Owner identifier.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"), "invalid Email"),
            },
            "alias": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The image name alias.",
            },
            "only_system_images": {
                Type:        schema.TypeBool,
                Optional:    true,
                ForceNew:    true,
                Default:     false,
                Description: "To disable searching for custom image.",
            },
        },
    }
}

func DataImageCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer DataImageError.WrapP(&err)

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

    images, err := m.Service.DataImageGetList(opts)
    if err != nil {
        return err
    }

    if images == nil {
        return errors.New("empty result")
    }

    b, _ := json.MarshalIndent(images, "", "\t")

    filter := struct {
        Name       string
        Alias      string
        OsType     string
        Owner      string
        OnlySystem bool
    }{
        Name:       d.Get("name").(string),
        Alias:      d.Get("alias").(string),
        OsType:     strings.ToLower(d.Get("os_type").(string)),
        Owner:      d.Get("owner").(string),
        OnlySystem: d.Get("only_system_images").(bool),
    }
    selectedImages := make([]service.Image, 0, 4)
    for _, value := range *images {
        if filter.Name != "" {
            if !regexp.MustCompile(filter.Name).MatchString(value.Name) {
                continue
            }
        }
        if filter.Alias != "" {
            if filter.Alias != value.Alias {
                continue
            }
        }
        if filter.OsType != "" {
            if filter.OsType != value.OsType {
                continue
            }
        }
        if filter.OnlySystem {
            if value.Owner != "" {
                continue
            }
        } else {
            if filter.Owner != "" {
                if filter.Owner != value.Owner {
                    continue
                }
            }
        }

        selectedImages = append(selectedImages, value)
    }

    if len(selectedImages) == 0 {
        return errors.New("no matches, allowed values:\n" + string(b))
    }

    image := selectedImages[0]
    newest := image.CreatedDate
    for _, value := range selectedImages {
        if value.CreatedDate > newest {
            newest = value.CreatedDate
            image = value
        }
    }
    b, _ = json.MarshalIndent(image, "", "\t")
    m.Log.Info("Data image:\n", string(b))
    d.SetId(image.Name)
    return nil
}
