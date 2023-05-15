package service

import (
    "encoding/json"
    "errors"
    "github.com/golang/mock/gomock"
    "terraform-provider-m3/client"
    cmock "terraform-provider-m3/client/mock"
    "testing"
)

func TestScheduleService_Create(t *testing.T) {
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

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
            },

            DoResponse: func() *client.M3BatchResult {
                def := &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                }

                schedule := Schedule{
                    Name:                 "name",
                    ScheduleName:         "schedule_name",
                    DefaultRequestParams: def,
                }

                data, _ := json.Marshal(schedule)

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
                req := struct {
                    Schedule *RequestSchedule `json:"schedule"`
                }{request.(*RequestSchedule)}
                m.EXPECT().MakePayload(req, MethodCreateSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &RequestSchedule{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                req := struct {
                    Schedule *RequestSchedule `json:"schedule"`
                }{request.(*RequestSchedule)}
                m.EXPECT().MakePayload(req, MethodCreateSchedule).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &RequestSchedule{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                req := struct {
                    Schedule *RequestSchedule `json:"schedule"`
                }{request.(*RequestSchedule)}
                m.EXPECT().MakePayload(req, MethodCreateSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &RequestSchedule{},

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
                req := struct {
                    Schedule *RequestSchedule `json:"schedule"`
                }{request.(*RequestSchedule)}
                m.EXPECT().MakePayload(req, MethodCreateSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &RequestSchedule{},

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
                req := struct {
                    Schedule *RequestSchedule `json:"schedule"`
                }{request.(*RequestSchedule)}
                m.EXPECT().MakePayload(req, MethodCreateSchedule).Return(nil, nil)
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

            _, err := s.ScheduleServicer.Create(testCase.Request.(*RequestSchedule))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestScheduleService_Describe(t *testing.T) {
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

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
                DefaultRequestParams: &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                },
            },

            DoResponse: func() *client.M3BatchResult {
                def := &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                }

                schedule := []Schedule{
                    {
                        Name:                 "name",
                        ScheduleName:         "schedule_name",
                        DefaultRequestParams: def,
                    },
                }

                data, _ := json.Marshal(schedule)

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
                req := struct {
                    *DefaultRequestParams
                    Cloud string `json:"cloud"`
                }{
                    DefaultRequestParams: request.(*RequestSchedule).DefaultRequestParams,
                    Cloud:                request.(*RequestSchedule).Cloud,
                }
                m.EXPECT().MakePayload(req, MethodDescribeSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
                DefaultRequestParams: &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                },
            },

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                req := struct {
                    *DefaultRequestParams
                    Cloud string `json:"cloud"`
                }{
                    DefaultRequestParams: request.(*RequestSchedule).DefaultRequestParams,
                    Cloud:                request.(*RequestSchedule).Cloud,
                }
                m.EXPECT().MakePayload(req, MethodDescribeSchedule).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
                DefaultRequestParams: &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                },
            },

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                req := struct {
                    *DefaultRequestParams
                    Cloud string `json:"cloud"`
                }{
                    DefaultRequestParams: request.(*RequestSchedule).DefaultRequestParams,
                    Cloud:                request.(*RequestSchedule).Cloud,
                }
                m.EXPECT().MakePayload(req, MethodDescribeSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &RequestSchedule{},

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
                req := struct {
                    *DefaultRequestParams
                    Cloud string `json:"cloud"`
                }{
                    DefaultRequestParams: request.(*RequestSchedule).DefaultRequestParams,
                    Cloud:                request.(*RequestSchedule).Cloud,
                }
                m.EXPECT().MakePayload(req, MethodDescribeSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
                DefaultRequestParams: &DefaultRequestParams{
                    TenantName: "NORTH",
                    Region:     "NORTH",
                },
            },

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
                req := struct {
                    *DefaultRequestParams
                    Cloud string `json:"cloud"`
                }{
                    DefaultRequestParams: request.(*RequestSchedule).DefaultRequestParams,
                    Cloud:                request.(*RequestSchedule).Cloud,
                }
                m.EXPECT().MakePayload(req, MethodDescribeSchedule).Return(nil, nil)
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

            _, err := s.ScheduleServicer.Describe(testCase.Request.(*RequestSchedule))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}

func TestScheduleService_Delete(t *testing.T) {
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

            Request: &RequestSchedule{
                Name:         "name",
                ScheduleName: "schedule_name",
            },

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
                m.EXPECT().MakePayload(request, MethodDeleteSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &RequestSchedule{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteSchedule).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &RequestSchedule{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDeleteSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &RequestSchedule{},

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
                m.EXPECT().MakePayload(request, MethodDeleteSchedule).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error if neither result nor error in response",

            WantErr: true,

            Request: &RequestSchedule{},

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
                m.EXPECT().MakePayload(request, MethodDeleteSchedule).Return(nil, nil)
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

            err := s.ScheduleServicer.Delete(testCase.Request.(*RequestSchedule))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}
