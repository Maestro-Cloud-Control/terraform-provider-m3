/*
 * Copyright 2022 Softline Group Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package provider

import (
    "errors"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    uuid "github.com/nu7hatch/gouuid"
    "regexp"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func resourceInstance() *schema.Resource {
    return &schema.Resource{
        Update:      resourceInstanceUpdate,
        Create:      resourceInstanceCreate,
        Read:        resourceInstanceRead,
        Delete:      resourceInstanceDelete,
        Description: "Creates instances of the specified configuration",
        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The name of the new instance.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-z0-9]{6,63}$"),
                    "Instance name must start with a letter, can be 6-63 characters long, and can include Latin alphanumeric characters ('a-z', '0-9')."),
            },
            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the tenant where the instance is to be launched.",
            },
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the region where the instance is to be run.",
            },
            "image": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The name of the image that will be used for the instance configuration.",
            },
            "key": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The name of the key pair to be used for the instance. Optional for Azure cloud",
            },
            "shape": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "Required if InstanceType is not specified. Instance shape is a Maestro name for a capacity configuration, mapped to InstanceType and corresponding attributes in other CPs. Some possible values: LARGE, MICRO, SMALL, etc.",
            },
            "enable_chef": {
                Type:        schema.TypeBool,
                Optional:    true,
                ForceNew:    true,
                Default:     false,
                Description: "Enabling chef application.",
            },
            "chef_profile": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The name of the chef application.",
            },
            "additional_data": {
                Type:        schema.TypeMap,
                Optional:    true,
                ForceNew:    true,
                Description: "The size of an additional storage volume in GB.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "owner": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "Owner identifier.",
            },
            "tags": {
                Type:        schema.TypeMap,
                Optional:    true,
                Description: "Key value parameter simplifying instance identification.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "instances_count": {
                Type:        schema.TypeInt,
                Optional:    true,
                ForceNew:    true,
                Default:     1,
                Description: "The number of instances that will be run. The default value is 1 (used if the parameter is not specified).",
            },
            "stop_after": {
                Type:         schema.TypeInt,
                Optional:     true,
                ForceNew:     true,
                Default:      0,
                ValidateFunc: validation.IntBetween(0, 720),
                Description:  "The expiration parameter which specifies when the machine will stop, in hours after creation.",
            },
            "terminate_after": {
                Type:         schema.TypeInt,
                Optional:     true,
                ForceNew:     true,
                Default:      0,
                ValidateFunc: validation.IntBetween(0, 720),
                Description:  "Termination parameter which specifies when the instance will be terminated, in hours after creation.",
            },
            "lock_termination": {
                Type:        schema.TypeBool,
                Optional:    true,
                ForceNew:    true,
                Default:     false,
                Description: "Locking the instance from termination.\nAllowed for clouds: [AWS, AZURE, GOOGLE].",
            },
            "cloud": {
                Type:        schema.TypeString,
                Computed:    true,
                ForceNew:    true,
                Description: "The cloud. \nAllowed values: [AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK, VSPHERE, VMWARE, YANDEX].",
            },
        },
    }
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceInstanceError.WrapP(&err)
    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }
    owner, err := utils.GetOwner(d, m.Config)
    if err != nil {
        return err
    }

    stopAfter := d.Get("stop_after").(int)
    terminateAfter := d.Get("terminate_after").(int)
    if terminateAfter != 0 && stopAfter != 0 && stopAfter >= terminateAfter {
        return fmt.Errorf("impossible stop instance: %s ,when it will be terminate.", d.Get("instance_name").(string))
    }

    defaultParams := &service.DefaultRequestParams{
        TenantName: tenant,
        Region:     region,
    }
    u, err := uuid.NewV4()
    instanceChefUUID := tenant + "." + region + "." + u.String()

    opts := &service.InstanceRunRequest{
        DefaultRequestParams: defaultParams,
        InstanceName:         d.Get("name").(string),
        KeyName:              d.Get("key").(string),
        Image:                d.Get("image").(string),
        Shape:                d.Get("shape").(string),
        Owner:                owner,
        ChefEnabled:          d.Get("enable_chef").(bool),
        InstanceChefUUID:     instanceChefUUID,
        ChefProfile:          d.Get("chef_profile").(string),
        LockedTermination:    d.Get("lock_termination").(bool),
        AdditionalData:       d.Get("additional_data").(map[string]interface{}),
        Tags:                 d.Get("tags").(map[string]interface{}),
        InstancesCount:       d.Get("instances_count").(int),
    }
    if stopAfter != 0 {
        opts.StopAfter = &service.StopAfter{StopAfter: stopAfter}
    }
    if terminateAfter != 0 {
        opts.TerminateAfter = &service.TerminateAfter{TerminateAfter: terminateAfter}
    }
    instance, err := m.Service.InstanceServicer.Run(opts)
    if err != nil {
        return err
    }
    d.SetId(instance.InstanceID)
    err = d.Set("lock_termination", instance.LockedTermination)
    if err != nil {
        return err
    }
    err = d.Set("cloud", instance.Cloud)
    if err != nil {
        return err
    }
    w := wait{
        Action: func() (interface{}, error) {
            instance, err := m.Service.InstanceServicer.Describe(
                &service.InstanceDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    InstanceIds:          []string{d.Id()},
                })
            if err != nil {
                return nil, err
            }
            if instance.State != service.InstanceStates.Running {
                return nil, errors.New("instance state: not running")
            }
            return instance, nil
        },
        CompareFn: defaultWaitCompareFunc(),
    }
    _, err = w.Wait()
    if err != nil {
        return err
    }

    m.Log.Info(fmt.Sprintf("Instance created ID: %s", d.Id()))
    return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceInstanceError.WrapP(&err)

    m := meta.(*Meta)
    tenant, err := utils.GetTenant(d, m.Config)
    if err != nil {
        return err
    }
    region, err := utils.GetRegion(d, m.Config)
    if err != nil {
        return err
    }

    opts := &service.InstanceDescribeRequest{
        DefaultRequestParams: &service.DefaultRequestParams{
            TenantName: tenant,
            Region:     region,
        },
        InstanceIds: []string{d.Id()},
    }

    instance, err := m.Service.InstanceServicer.Describe(opts)

    if err != nil && err.Error() != "404" {
        return err
    }
    if instance == nil {
        d.SetId("")
        return nil
    }
    return nil
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceInstanceError.WrapP(&err)

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

    m.Log.Info(fmt.Sprintf("Deleting instance: %s", d.Id()))

    terminateOpts := &service.InstanceTerminateRequest{
        DefaultRequestParams: defaultParams,
        InstanceID:           d.Id(),
    }
    if d.Get("lock_termination").(bool) {
        err := m.Service.InstanceServicer.UnlockTermination(terminateOpts)
        if err != nil {
            if err.Error() == "404" {
                return errors.New("instance not found")
            }
            return err
        }
    }

    err = m.Service.InstanceServicer.Terminate(terminateOpts)
    // Instance not found
    if err != nil {
        if err.Error() == "404" {
            return errors.New("instance  not found")
        }
        return err
    }

    w := wait{
        Action: func() (interface{}, error) {
            return m.Service.InstanceServicer.Describe(
                &service.InstanceDescribeRequest{
                    DefaultRequestParams: defaultParams,
                    InstanceIds:          []string{d.Id()},
                })
        },
        CompareFn: defaultInverseWaitCompareFunc(),
    }
    _, err = w.Wait()

    if err != nil {
        if err.Error() == "404" {
            return nil
        }
        return fmt.Errorf("error wait for state %s instance: %s", d.Id(), err)
    }

    d.SetId("")
    m.Log.Info(fmt.Sprintf("Instance terminated: %s", d.Id()))
    return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
    defer UpdatingError.WrapP(&err)
    defer ResourceInstanceError.WrapP(&err)

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

    describeOpts := &service.InstanceDescribeRequest{
        DefaultRequestParams: defaultParams,
        InstanceIds:          []string{d.Id()},
    }

    instance, err := m.Service.InstanceServicer.Describe(describeOpts)

    if len(d.Get("tags").(map[string]interface{})) < 1 {
        tags := make([]string, 0, 4)

        for _, tag := range instance.Tags {
            tags = append(tags, fmt.Sprintf("%s=%s", tag.Key, tag.Value))
        }

        deleteOpts := &service.InstanceDeleteTagsRequest{
            DefaultRequestParams: defaultParams,
            Id:                   d.Id(),
            Cloud:                instance.Cloud,
            AvailabilityZone:     instance.AvailabilityZone,
            ResourceGroup:        instance.ResourceGroup,
            Tags:                 tags,
        }

        return m.Service.InstanceServicer.DeleteTags(deleteOpts)
    }

    updateOpts := &service.InstanceUpdateTagsRequest{
        DefaultRequestParams: defaultParams,
        Id:                   d.Id(),
        Cloud:                instance.Cloud,
        AvailabilityZone:     instance.AvailabilityZone,
        ResourceGroup:        instance.ResourceGroup,
        VolumeIds:            instance.VolumesIds,
        InstanceName:         d.Get("name").(string),
        Tags:                 d.Get("tags").(map[string]interface{}),
        Overwrite:            true,
    }

    return m.Service.InstanceServicer.UpdateTags(updateOpts)
}
