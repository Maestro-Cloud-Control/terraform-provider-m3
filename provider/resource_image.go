package provider

import (
    "errors"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func resourceImage() *schema.Resource {
    return &schema.Resource{
        Create:      resourceImageCreate,
        Read:        resourceImageRead,
        Delete:      resourceImageDelete,
        Description: "Creates an image based on an existing instance.",
        Schema: map[string]*schema.Schema{
            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the tenant to which the source instance belongs.",
            },
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the region where the source instance is hosted.",
            },
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The image name. Length: 3-63. \nImage name requirements:\n  - Google regions: lowercase alphanumeric characters including hyphens, if not as a first or last symbol. \n  - OpenStack, Azure and AWS regions: alphanumeric characters including .()[]-@_",
            },
            "source_instance_id": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The ID of the instance used as a source for the image.",
            },
            "description": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The description for the image.",
            },
        },
    }
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceImageError.WrapP(&err)
    m := meta.(*Meta)

    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    defaultParams := &service.DefaultRequestParams{
        TenantName: tenant,
        Region:     region,
    }
    opts := &service.ImageCreateRequest{
        DefaultRequestParams: defaultParams,
        InstanceID:           d.Get("source_instance_id").(string),
        ImageName:            d.Get("name").(string),
        Description:          d.Get("description").(string),
    }

    image, err := m.Service.ImageServicer.Create(opts)
    if err != nil {
        return err
    }

    w := wait{
        Action: func() (interface{}, error) {
            image, err := m.Service.ImageServicer.Describe(
                &service.ImageDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    ImageIds:             []string{d.Id()},
                })
            if err != nil {
                return nil, err
            }
            if image.State != service.AvailableImageState {
                return nil, errors.New("image state: not available")
            }
            return image, nil
        },
        CompareFn: defaultWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        if err.Error() == "404" {
            return errors.New("image not found")
        }

        return fmt.Errorf("error wait for state %s", err)
    }
    d.SetId(image.ImageID)

    m.Log.Info(fmt.Sprintf("Image created ID: %s", d.Id()))
    return resourceImageRead(d, meta)
}

func resourceImageRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceImageError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.ImageDescribeRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        ImageIds: []string{d.Id()},
    }

    _, err = m.Service.ImageServicer.Describe(opts)
    if err != nil {
        if err.Error() == "404" {
            return errors.New("image not found")
        }

        return fmt.Errorf("error wait for state %s", err)
    }

    return nil
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceImageError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    defaultParams := &service.DefaultRequestParams{
        TenantName: tenant,
        Region:     region,
    }

    m.Log.Info(fmt.Sprintf("Deleting Image: %s", d.Id()))

    deleteOpts := &service.DeleteImageRequest{
        DefaultRequestParams: defaultParams,
        ImageID:              d.Id(),
    }
    err = m.Service.ImageServicer.Delete(deleteOpts)
    // Image not found
    if err != nil {
        if err.Error() == "404" {
            return fmt.Errorf("image %s not found", d.Id())
        }

        return err
    }

    w := wait{
        Action: func() (interface{}, error) {
            return m.Service.ImageServicer.Describe(
                &service.ImageDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    ImageIds:             []string{d.Id()},
                })
        },
        CompareFn: defaultInverseWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        return fmt.Errorf("error wait for state %s image: %s", d.Id(), err)
    }

    d.SetId("")
    m.Log.Info("Image terminated: ", d.Id())
    return nil
}
