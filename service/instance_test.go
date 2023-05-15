package service

import (
    "encoding/json"
    "errors"
    "github.com/golang/mock/gomock"
    "terraform-provider-m3/client"
    cmock "terraform-provider-m3/client/mock"
    "testing"
)

func TestInstancesService_Run(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                instances := InstancesResultData{
                    Instances: []Instance{
                        {
                            InstanceID:        "123456789",
                            Cloud:             "AWS",
                            InstanceName:      "name",
                            TenantName:        "NORTH",
                            Region:            "NORTH",
                            State:             "starting",
                            Created:           "",
                            Architecture:      "",
                            Image:             "123456789",
                            Shape:             "abc",
                            PrivateIP:         "127.0.0.1",
                            LockedTermination: true,
                            ChefEnabled:       true,
                            InstanceChefUUID:  "abc1234",
                            ChefProfile:       "abc",
                            AdditionalData:    nil,
                        },
                    },
                }
                data, _ := json.Marshal(instances)

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   string(data),
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if instance state not equal starting",

            WantErr: true,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                instances := []Instance{
                    {
                        InstanceID:        "123456789",
                        Cloud:             "AWS",
                        InstanceName:      "name",
                        TenantName:        "NORTH",
                        Region:            "NORTH",
                        State:             "running",
                        Created:           "",
                        Architecture:      "",
                        Image:             "123456789",
                        Shape:             "abc",
                        PrivateIP:         "127.0.0.1",
                        LockedTermination: true,
                        ChefEnabled:       true,
                        InstanceChefUUID:  "abc1234",
                        ChefProfile:       "abc",
                        AdditionalData:    nil,
                    },
                }

                data, _ := json.Marshal(instances)

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   string(data),
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &InstanceRunRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodRunInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            _, err := s.InstanceServicer.Run(testCase.Request.(*InstanceRunRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestInstancesService_Terminate(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if error contain 'No unique instance found by instance ID'",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "No unique instance found by instance ID",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodTerminateInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            err := s.InstanceServicer.Terminate(testCase.Request.(*InstanceTerminateRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestInstancesService_Describe(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                instances := InstancesResultData{
                    Instances: []Instance{
                        {
                            InstanceID:        "123456789",
                            Cloud:             "AWS",
                            InstanceName:      "name",
                            TenantName:        "NORTH",
                            Region:            "NORTH",
                            State:             "starting",
                            Created:           "",
                            Architecture:      "",
                            Image:             "123456789",
                            Shape:             "abc",
                            PrivateIP:         "127.0.0.1",
                            LockedTermination: true,
                            ChefEnabled:       true,
                            InstanceChefUUID:  "abc1234",
                            ChefProfile:       "abc",
                            AdditionalData:    nil,
                        },
                    },
                }

                data, _ := json.Marshal(instances)

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   string(data),
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if in response empty data",

            WantErr: true,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &InstanceDescribeRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeInstance).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            _, err := s.InstanceServicer.Describe(testCase.Request.(*InstanceDescribeRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestInstancesService_UnlockTermination(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(struct {
                    InstanceTerminateRequest
                    Action string `json:"action"`
                }{
                    InstanceTerminateRequest: *(request.(*InstanceTerminateRequest)),
                    Action:                   "DISABLE",
                }, MethodTerminationProtection).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(struct {
                    InstanceTerminateRequest
                    Action string `json:"action"`
                }{
                    InstanceTerminateRequest: *(request.(*InstanceTerminateRequest)),
                    Action:                   "DISABLE",
                }, MethodTerminationProtection).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceTerminateRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(struct {
                    InstanceTerminateRequest
                    Action string `json:"action"`
                }{
                    InstanceTerminateRequest: *(request.(*InstanceTerminateRequest)),
                    Action:                   "DISABLE",
                }, MethodTerminationProtection).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            err := s.InstanceServicer.UnlockTermination(testCase.Request.(*InstanceTerminateRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestInstancesService_UpdateTags(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   " ",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if error contain 'No unique instance found by instance ID'",

            WantErr: true,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "No unique instance found by instance ID",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &InstanceUpdateTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodUpdateTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            err := s.InstanceServicer.UpdateTags(testCase.Request.(*InstanceUpdateTagsRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestInstancesService_DeleteTags(t *testing.T) {
    type MockBehavior func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult)
    type TestCase struct {
        Request      interface{}
        Name         string
        WantErr      bool
        DoResponse   func() *client.M3BatchResult
        MockBehavior MockBehavior
    }

    testTable := []TestCase{
        {
            Name: "OK",

            WantErr: false,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "",
                    Data:   " ",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if error contain 'No unique instance found by instance ID'",

            WantErr: true,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {

                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "No unique instance found by instance ID",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "FAIL",
                    Error:  "some error",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &InstanceDeleteTagsRequest{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "",
                    Error:  "",
                    Data:   "",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteTags).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockTransporter := cmock.NewMockTransporter(ctl)
            testCase.MockBehavior(mockTransporter, testCase.Request, testCase.DoResponse())

            c := &client.Client{Transporter: mockTransporter}
            s := NewService(c)

            err := s.InstanceServicer.DeleteTags(testCase.Request.(*InstanceDeleteTagsRequest))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}
