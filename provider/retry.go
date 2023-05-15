package provider

import (
    "errors"
    "time"
)

type waitAction func() (interface{}, error)
type waitCompareFunc func(error) bool

func defaultWaitCompareFunc() waitCompareFunc {
    return func(err error) bool {
        if err != nil {
            return false
        }
        return true
    }
}

func defaultInverseWaitCompareFunc() waitCompareFunc {
    return func(err error) bool {
        if err != nil {
            return true
        }
        return false
    }
}

type wait struct {
    CompareFn waitCompareFunc
    Action    waitAction
    Attempts  int
    Delay     int
}

func (w wait) Wait() (interface{}, error) {
    if w.Delay == 0 {
        w.Delay = 60
    }
    if w.Attempts == 0 {
        w.Attempts = 30
    }

    for i := 1; i <= w.Attempts; i++ {
        result, err := w.Action()
        if w.CompareFn(err) {
            return result, nil
        }

        time.Sleep(time.Duration(w.Delay) * time.Second)
    }

    return nil, errors.New("timeout")
}
