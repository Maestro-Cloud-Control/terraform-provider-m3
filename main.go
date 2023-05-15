package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "terraform-provider-m3/logger"
    "terraform-provider-m3/provider"
)

func main() {
    logger.NewTFLog()
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: provider.Provider,
        Logger:       logger.NewTFLog(),
    })
}
