package provider

import (
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
            "additional_data": {
                Type:        schema.TypeMap,
                Optional:    true,
                ForceNew:    true,
                Description: "Additional cloud related parameter.",
                Elem: &schema.Schema{Type: schema.TypeString},
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
            AdditionalData:       d.Get("additional_data").(map[string]interface{}),
        }
        m.Log.Info(fmt.Sprintf("Creating volume without instance: %s", opts.VolumeName))
        volume, err = m.Service.VolumeServicer.Create(opts)
        neededState = service.AvailableVolumeState
    } else {
        opts := &service.VolumeCreateAndAttachRequest{
            DefaultRequestParams: defaultParams,
            VolumeName:           d.Get("name").(string),
            SizeInGB:             d.Get("size_in_gb").(int),
            InstanceId:           instanceId,
            AdditionalData:       d.Get("additional_data").(map[string]interface{}),
        }
        m.Log.Info(fmt.Sprintf("Creating volume with instance attachment: %s, instance: %s", opts.VolumeName, instanceId))
        volume, err = m.Service.VolumeServicer.CreateAndAttach(opts)
        neededState = service.InUseState
    }

    if err != nil {
        m.Log.Info(fmt.Sprintf("Failed to create volume: %s, error: %v", d.Get("name").(string), err))
        return err
    }

    m.Log.Info(fmt.Sprintf("Volume created, ID: %s, waiting for state: %s", volume.VolumeID, neededState))

    w := wait{
        Action: func() (interface{}, error) {
            m.Log.Info(fmt.Sprintf("Describing volume: %s, instance: %s, tenant: %s, region: %s",
                volume.VolumeID, instanceId, tenant, region))
            describedVolume, err := m.Service.VolumeServicer.Describe(
                &service.VolumeDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    VolumeIds:            []string{volume.VolumeID},
                    InstanceId:           instanceId,
                })
            if err != nil {
                m.Log.Info(fmt.Sprintf("Volume describe failed for volume ID: %s, error: %v", volume.VolumeID, err))
                return nil, err
            }
            m.Log.Info(fmt.Sprintf("Volume describe result: ID: %s, State: %s, Expected: %s",
                describedVolume.VolumeID, describedVolume.State, neededState))
            if describedVolume.State != neededState {
                return nil, fmt.Errorf("volume state: not %s", neededState)
            }
            return describedVolume, nil
        },
        CompareFn: defaultWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        if err.Error() == "404" {
            m.Log.Info(fmt.Sprintf("Volume not found during wait: %s", d.Get("volume_name").(string)))
            return fmt.Errorf("volume %s not found", d.Get("volume_name").(string))
        }
        m.Log.Info(fmt.Sprintf("Error waiting for volume state: %s, error: %v", d.Get("volume_name").(string), err))
        return fmt.Errorf("error wait for state, volume: %s", d.Get("volume_name").(string))
    }
    d.SetId(volume.VolumeID)

    m.Log.Info(fmt.Sprintf("Volume created successfully, ID: %s", d.Id()))
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

    m.Log.Info(fmt.Sprintf("Deleting volume: %s", d.Id()))

    instanceId := d.Get("instance_id").(string)
    describeOpts := &service.VolumeDescribeRequest{
        DefaultRequestParams: defaultParams,
        VolumeIds:            []string{d.Id()},
        InstanceId:           instanceId,
    }

    volume, err := m.Service.VolumeServicer.Describe(describeOpts)
    if err != nil {
        if err.Error() == "404" {
            m.Log.Info(fmt.Sprintf("Volume not found, skipping delete: %s", d.Id()))
            d.SetId("")
            return nil
        }
        // Retry without instanceId if describe failed with instanceId
        if len(instanceId) > 0 {
            m.Log.Info(fmt.Sprintf("Describe failed with instanceId, retrying without instanceId: %s", d.Id()))
            describeOptsWithoutInstance := &service.VolumeDescribeRequest{
                DefaultRequestParams: defaultParams,
                VolumeIds:            []string{d.Id()},
            }
            volume, err = m.Service.VolumeServicer.Describe(describeOptsWithoutInstance)
            if err != nil {
                if err.Error() == "404" {
                    m.Log.Info(fmt.Sprintf("Volume not found, skipping delete: %s", d.Id()))
                    d.SetId("")
                    return nil
                }
                return err
            }
        } else {
            return err
        }
    }

    m.Log.Info(fmt.Sprintf("Volume state: %s", volume.State))

    if volume.State == service.InUseState {
        m.Log.Info(fmt.Sprintf("Volume is attached, detaching: %s", d.Id()))
        detachOpts := &service.VolumeDetachRequest{
            DefaultRequestParams: defaultParams,
            VolumeID:             d.Id(),
            InstanceId:           instanceId,
        }
        err = m.Service.VolumeServicer.Detach(detachOpts)
        if err != nil {
            if err.Error() == "404" {
                m.Log.Info(fmt.Sprintf("Volume not found during detach, skipping: %s", d.Id()))
                d.SetId("")
                return nil
            }
            return err
        }

        m.Log.Info(fmt.Sprintf("Waiting for volume to become detached: %s", d.Id()))
        w := wait{
            Action: func() (interface{}, error) {
                describedVolume, describeErr := m.Service.VolumeServicer.Describe(
                    &service.VolumeDescribeRequest{
                        DefaultRequestParams: defaultParams,
                        VolumeIds:            []string{d.Id()},
                    })
                if describeErr != nil {
                    if describeErr.Error() == "404" {
                        m.Log.Info(fmt.Sprintf("Volume not found during detach wait, considering detached: %s", d.Id()))
                        return nil, nil
                    }
                    m.Log.Info(fmt.Sprintf("Volume describe failed for volume ID: %s, error: %v", d.Id(), describeErr))
                    return nil, describeErr
                }
                m.Log.Info(fmt.Sprintf("Volume state after detach: %s", describedVolume.State))
                if describedVolume.State != service.AvailableVolumeState {
                    return nil, fmt.Errorf("volume state: not %s", service.AvailableVolumeState)
                }
                return describedVolume, nil
            },
            CompareFn: defaultWaitCompareFunc(),
        }
        _, err = w.Wait()
        if err != nil {
            return fmt.Errorf("error waiting for volume to detach: %s", d.Id())
        }
    }

    if err := deleteVolumeWithRetry(m, defaultParams, d.Id()); err != nil {
        return err
    }

    m.Log.Info(fmt.Sprintf("Volume deleted successfully: %s", d.Id()))
    d.SetId("")
    return nil
}

func deleteVolumeWithRetry(m *Meta, defaultParams *service.DefaultRequestParams, volumeID string) error {
    m.Log.Info(fmt.Sprintf("Deleting volume with retry: %s", volumeID))

    deleteOpts := &service.VolumeDeleteRequest{
        DefaultRequestParams: defaultParams,
        VolumeID:             volumeID,
    }

    w := wait{
        Action: func() (interface{}, error) {
            err := m.Service.VolumeServicer.Delete(deleteOpts)
            if err != nil {
                if err.Error() == "404" {
                    m.Log.Info(fmt.Sprintf("Volume not found during delete, skipping: %s", volumeID))
                    return nil, nil
                }
                m.Log.Info(fmt.Sprintf("Error deleting volume (will retry if attempts remain): %s, error: %v", volumeID, err))
                return nil, err
            }
            return nil, nil
        },
        CompareFn: defaultWaitCompareFunc(),
        Attempts:  10,
        Delay:     60,
    }

    _, err := w.Wait()
    if err != nil {
        return fmt.Errorf("error waiting for volume deletion: %s", volumeID)
    }

    return nil
}
