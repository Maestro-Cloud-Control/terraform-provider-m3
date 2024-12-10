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

func resourceKeypair() *schema.Resource {
    return &schema.Resource{
        Create:      resourceKeypairCreate,
        Read:        resourceKeypairRead,
        Update:      resourceKeypairUpdate,
        Delete:      resourceKeypairDelete,
        Description: "Registers an SSH key for further usage",
        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The name of the public SSH key.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-_]{6,32}$`),
                    "SSH key name must be 6-32 characters long, and include only these characters: '0-9', 'a-z', 'A-Z', '-_'."),
            },

            "tenant": {
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "",
                Description: "The tenant for which the key is to be registered.",
            },

            "public_key": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The public SSH key content.",
                ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(.*?)(AAAAB3NzaC1yc[0-9A-Za-z+/\\n]{256,}[=]{0,3})(.*)$`),
                    "Failed to parse the SSH public key. Please, check that your value is in the SSH or SSH2 format."),
            },

            "cloud": {
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "",
                Description: "The cloud for which the key is to be registered.\nAllowed values [ AWS, AZURE, GOOGLE, NUTANIX, OPEN_STACK ].",
                ValidateFunc: validation.StringInSlice([]string{
                    "AWS", "AZURE", "GOOGLE", "OPEN_STACK", "NUTANIX"}, true),
            },
        },
    }
}

func resourceKeypairCreate(d *schema.ResourceData, meta interface{}) (err error) {
    defer CreatingError.WrapP(&err)
    defer ResourceKeypairError.WrapP(&err)

    m := meta.(*Meta)
    if err := utils.MatchEmail(m.Config.UserIdentifier); err != nil {
        return err
    }
    opts := service.KeypairRequest{
        Name:           d.Get("name").(string),
        KeypairContent: &service.KeypairContent{Content: d.Get("public_key").(string)},
        Email:          m.Config.UserIdentifier,
    }

    if d.Get("tenant").(string) != "" {
        opts.KeypairTenantName = &service.KeypairTenantName{TenantName: d.Get("tenant").(string)}
    } else {
        opts.KeypairAllTenants = &service.KeypairAllTenants{AllTenants: true}
    }

    if d.Get("cloud").(string) != "" {
        opts.KeypairCloud = &service.KeypairCloud{Cloud: strings.ToUpper(d.Get("cloud").(string))}
    }

    keypair, err := m.Service.KeypairServicer.Create(&opts)
    if err != nil {
        return err
    }
    if keypair == nil {
        m.Log.Info("Some troubles with keypair.")
    }
    return resourceKeypairRead(d, meta)
}

func resourceKeypairRead(d *schema.ResourceData, meta interface{}) (err error) {
    defer ReadingError.WrapP(&err)
    defer ResourceKeypairError.WrapP(&err)

    m := meta.(*Meta)
    opts := &service.KeypairRequest{
        Email: m.Config.UserIdentifier,
        Name:  d.Get("name").(string),
    }
    if d.Get("tenant").(string) != "" {
        opts.KeypairTenantName = &service.KeypairTenantName{TenantName: d.Get("tenant").(string)}
    } else {
        opts.KeypairAllTenants = &service.KeypairAllTenants{AllTenants: true}
    }

    if d.Get("cloud").(string) != "" {
        opts.KeypairCloud = &service.KeypairCloud{Cloud: strings.ToUpper(d.Get("cloud").(string))}
    }
    w := wait{
        Action: func() (interface{}, error) {
            return m.Service.KeypairServicer.Describe(opts)
        },
        CompareFn: defaultWaitCompareFunc(),
    }
    result, err := w.Wait()
    if err != nil || result == nil {
        return err
    }

    d.SetId(result.(*service.Keypair).Name)
    return nil
}

func resourceKeypairDelete(d *schema.ResourceData, meta interface{}) (err error) {
    defer DeletingError.WrapP(&err)
    defer ResourceKeypairError.WrapP(&err)

    m := meta.(*Meta)
    opts := &service.KeypairRequest{
        Name:  d.Get("name").(string),
        Email: m.Config.UserIdentifier,
    }

    err = m.Service.KeypairServicer.Delete(opts)
    if err != nil {
        if err.Error() == "404" {
            return fmt.Errorf("script %s not found", d.Id())
        }

        return err
    }
    d.SetId("")
    return nil
}

func resourceKeypairUpdate(d *schema.ResourceData, meta interface{}) (err error) {
    defer UpdatingError.WrapP(&err)
    defer ResourceKeypairError.WrapP(&err)

    if err := resourceKeypairDelete(d, meta); err != nil {
        return err
    }

    if err := resourceKeypairCreate(d, meta); err != nil {
        return err
    }
    return nil
}
