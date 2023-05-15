package utils

import (
    "errors"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "regexp"
    "terraform-provider-m3/client"
)

func Contains(element interface{}, arr []interface{}) bool {
    for _, value := range arr {
        if element == value {
            return true
        }
    }
    return false
}

func MatchEmail(email string) error {
    emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    if emailRegex.MatchString(email) {
        return errors.New("invalid Email")
    }
    return nil
}

func GetTenant(d *schema.ResourceData, conf *client.Config) (string, error) {
    return GetParam(d, conf, "tenant", func(c *client.Config) string {
        return conf.TenantName
    })
}

func GetRegion(d *schema.ResourceData, conf *client.Config) (string, error) {
    return GetParam(d, conf, "region", func(conf *client.Config) string {
        return conf.RegionName
    })
}

func GetCloud(d *schema.ResourceData, c *client.Config) (string, error) {
    return GetParam(d, c, "cloud", func(conf *client.Config) string {
        return conf.Cloud
    })
}

func GetOwner(d *schema.ResourceData, c *client.Config) (string, error) {
    return GetParam(d, c, "owner", func(conf *client.Config) string {
        return conf.UserIdentifier
    })
}

func GetParam(d *schema.ResourceData, conf *client.Config, paramName string, paramExtractor func(config *client.Config) string) (string, error) {
    param := d.Get(paramName).(string)
    if param == "" {
        param = paramExtractor(conf)
        if param == "" {
            return "", errors.New(paramName + " is empty")
        }
        return param, nil
    }
    return param, nil
}
