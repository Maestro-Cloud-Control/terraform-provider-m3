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
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    "regexp"
    "strings"
    "terraform-provider-m3/service"
    "terraform-provider-m3/utils"
)

func resourceSchedule() *schema.Resource {
    return &schema.Resource{
        Create:      resourceScheduleCreate,
        Read:        resourceScheduleRead,
        Update:      resourceScheduleUpdate,
        Delete:      resourceScheduleDelete,
        Description: "Resource for configuring a schedule for instances start or stop.",
        Schema: map[string]*schema.Schema{
            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the tenant where the schedule should be applied.",
            },
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The name of the region where the schedule should be applied.",
            },
            "name": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                Description:  "The name of schedule.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\\.\\-_A-Za-z0-9]{0,50}$"), "Invalid `schedule`"),
            },
            "description": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "The schedule description.",
            },
            "action": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                Description:  "The schedule action.\nAllowed values: START, STOP.",
                ValidateFunc: validation.StringInSlice([]string{"START", "STOP"}, true),
            },
            "cloud": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Description: "The cloud in which the schedule should be applied.\nAllowed values [ AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK, VSPHERE, VMWARE, YANDEX ].",
                ValidateFunc: validation.StringInSlice([]string{
                    "AWS", "AZURE", "GOOGLE", "OPEN_STACK", "YANDEX",
                    "VMWARE", "VSPHERE", "NUTANIX", "HYPERV"}, true),
            },
            "cron": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                Description:  "A cron expression identifying the schedule by which the action should be taken.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^([^ ]+) ([^ ]+ ){3,5}([^ ]+)$"), "Cron expression must contain these parameters separated by spaces: minutes, hours, days of month, months, and days of weeks."),
            },
            "instances_id": {
                Type:        schema.TypeList,
                Elem:        &schema.Schema{Type: schema.TypeString},
                Optional:    true,
                Default:     nil,
                Description: "The ID of the instance to be affected by the schedule.\nTo set the schedule for several instances, specify the IDs one by one.",
            },
            "tag_key": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "The Key of tag by which the instances for scheduling should be selected.",
            },
            "tag_value": {
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
                Default:     "",
                Description: "Value of tag by which the instances for scheduling should be selected. Optional.",
            },
        },
    }
}

func resourceScheduleCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceScheduleError.WrapP(&err)

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

    opts := &service.RequestSchedule{
        DefaultRequestParams: &service.DefaultRequestParams{
            Region:     region,
            TenantName: tenant,
        },
        Name:         d.Get("name").(string),
        ScheduleName: strings.ToLower(d.Get("name").(string) + "::" + region),
        Description:  d.Get("description").(string),
        Action:       strings.ToUpper(d.Get("action").(string)),
        Cron:         d.Get("cron").(string),
        Cloud:        strings.ToUpper(cloud),
        Instances:    make([]*service.RequestScheduleInstance, 0, 2),
    }

    if d.Get("tag_key").(string) != "" {
        opts.Type = "My instances with tag"
        opts.Tag = new(service.RequestScheduleTag)
        opts.Tag.Key = d.Get("tag_key").(string)
        opts.Tag.Value = d.Get("tag_value").(string)

    } else if len(d.Get("instances_id").([]interface{})) > 0 {
        opts.Type = "Specified instances"
        for _, val := range d.Get("instances_id").([]interface{}) {
            opts.Instances = append(opts.Instances, &service.RequestScheduleInstance{
                InstanceId: val.(string),
                InstanceLocationInfo: &service.RequestScheduleInstanceLocationInfo{
                    Region: region,
                },
            })
        }

    } else {
        opts.Type = "All my instances in region"
    }

    schedule, err := m.Service.ScheduleServicer.Create(opts)
    if err != nil {
        return err
    }
    if schedule == nil {
        m.Log.Info("some troubles with schedule: ", d.Get("schedule_name"))
    }
    d.SetId(schedule.Name)
    return resourceScheduleRead(d, meta)
}

func resourceScheduleDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceScheduleError.WrapP(&err)

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

    opts := &service.RequestSchedule{
        DefaultRequestParams: &service.DefaultRequestParams{
            Region:     region,
            TenantName: tenant,
        },

        ScheduleName: strings.ToLower(d.Get("name").(string)),
        Instances:    make([]*service.RequestScheduleInstance, 0, 2),
        Cloud:        cloud,
    }

    err = m.Service.ScheduleServicer.Delete(opts)
    if err != nil {
        if err.Error() == "404" {
            m.Log.Info(fmt.Sprintf("schedule %s not found", d.Get("name")))
            return nil
        }

        return err
    }
    d.SetId("")
    return nil
}

func resourceScheduleRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceScheduleError.WrapP(&err)

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

    opts := &service.RequestSchedule{
        DefaultRequestParams: &service.DefaultRequestParams{
            Region:     region,
            TenantName: tenant,
        },
        Cloud: cloud,
        Name:  d.Get("name").(string),
    }

    _, err = m.Service.ScheduleServicer.Describe(opts)
    if err != nil {
        if err.Error() == "404" {
            m.Log.Info(fmt.Sprintf("Schedule %s not found", d.Get("name")))
            d.SetId("")
            return nil
        }
        return err
    }

    return nil
}

func resourceScheduleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
    defer UpdatingError.WrapP(&err)
    defer ResourceKeypairError.WrapP(&err)

    if err := resourceScheduleDelete(d, meta); err != nil {
        return err
    }

    if err := resourceScheduleCreate(d, meta); err != nil {
        return err
    }
    return nil
}
