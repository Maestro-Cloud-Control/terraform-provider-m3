package logger

import "github.com/hashicorp/go-hclog"

func NewTFLog() hclog.Logger {
    opt := hclog.DefaultOptions
    log := hclog.New(opt)
    return log
}
