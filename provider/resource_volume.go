package provider

import (
    "errors"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func resourceVolume() *schema.Resource {
    return &schema.Resource{
        Create:      resourceVolumeCreate,
        Read:        resourceVolumeRead,
        Delete:      resourceVolumeDelete,
        Description: "Creates a new storage volume and attaches it to the specified instance.",
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
                Required:    true,
                ForceNew:    true,
                Description: "The volume name.",
            },
            "size_in_gb": {
                Type:        schema.TypeInt,
                Required:    true,
                ForceNew:    true,
                Description: "The maximum size of the volume at creation, in GB.",
            },
            "instance_id": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The ID of the instance to which the volume will be added.\nIf not specified, the volume will not be attached to any instance.",
            },
        },
    }
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceVolumeError.WrapP(&err)

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
    instanceId := d.Get("instance_id").(string)
    var volume *service.Volume
    var neededState string

    if len(instanceId) == 0 {
        instanceId = ""
        opts := &service.VolumeCreateRequest{
            DefaultRequestParams: defaultParams,
            VolumeName:           d.Get("name").(string),
            SizeInGB:             d.Get("size_in_gb").(int),
        }
        volume, err = m.Service.VolumeServicer.Create(opts)
        neededState = service.AvailableImageState
    } else {
        opts := &service.VolumeCreateAndAttachRequest{
            DefaultRequestParams: defaultParams,
            VolumeName:           d.Get("name").(string),
            SizeInGB:             d.Get("size_in_gb").(int),
            InstanceId:           instanceId,
        }
        volume, err = m.Service.VolumeServicer.CreateAndAttach(opts)
        neededState = service.InUseState
    }

    m.Log.Info(fmt.Sprintf("Creating volume: %s", d.Id()))

    if err != nil {
        return err
    }

    w := wait{
        Action: func() (interface{}, error) {
            volume, err := m.Service.VolumeServicer.Describe(
                &service.VolumeDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    VolumeIds:            []string{d.Id()},
                    InstanceId:           instanceId,
                })
            if err != nil {
                return nil, err
            }
            if volume.State != neededState {
                return nil, errors.New(fmt.Sprintf("volume state: not %s", neededState))
            }
            return volume, nil
        },
        CompareFn: defaultWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        if err.Error() == "404" {
            return fmt.Errorf("volume %s not found", d.Get("volume_name").(string))
        }
        return fmt.Errorf("error wait for state, volume: %s", d.Get("volume_name").(string))
    }
    d.SetId(volume.VolumeID)

    m.Log.Info(fmt.Sprintf("Volume created ID: %s", d.Id()))
    return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceVolumeError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    instanceId := d.Get("instance_id").(string)
    if len(instanceId) != 0 {
        return nil
    }
    opts := &service.VolumeDescribeRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        VolumeIds: []string{d.Id()},
    }

    _, err = m.Service.VolumeServicer.Describe(opts)
    if err != nil {
        if err.Error() == "404" {
            return err
        }
        d.SetId("")
        return err
    }
    return nil
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceVolumeError.WrapP(&err)

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

    m.Log.Info("Deleting volume: %s", d.Id())

    deleteOpts := &service.VolumeDeleteRequest{
        DefaultRequestParams: defaultParams,
        VolumeID:             d.Id(),
    }
    err = m.Service.VolumeServicer.Delete(deleteOpts)
    // Volume not found
    if err != nil {
        if err.Error() == "404" {
            return fmt.Errorf("volume %s not found", d.Id())
        }
        d.SetId("")
        return err
    }

    w := wait{
        Action: func() (interface{}, error) {
            return m.Service.VolumeServicer.Describe(
                &service.VolumeDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    VolumeIds:            []string{d.Id()},
                })
        },
        CompareFn: defaultInverseWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        return fmt.Errorf("error wait for state image: %s", d.Id())
    }

    m.Log.Info(fmt.Sprintf("volume terminated: %s", d.Id()))
    d.SetId("")
    return nil
}
