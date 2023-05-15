package service

import (
    "encoding/json"
    "errors"
    "github.com/golang/mock/gomock"
    "terraform-provider-m3/client"
    "terraform-provider-m3/client/mock"
    "testing"
)

func TestDataImageService_DataImageGetList(t *testing.T) {
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

            Request: &DefaultRequestParams{},

            DoResponse: func() *client.M3BatchResult {
                images := []Image{
                    {
                        TenantName: "North",
                        Region:     "North",
                        Name:       "name",
                        ImageID:    "123456789",
                    },
                }

                data, _ := json.Marshal(images)

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
                m.EXPECT().MakePayload(request, MethodDescribeImage).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error on during make payload",

            WantErr: true,

            Request: &DefaultRequestParams{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeImage).Return(nil, errors.New("some error"))
            },
        },

        {
            Name: "Got error on during DO",

            WantErr: true,

            Request: &DefaultRequestParams{},

            DoResponse: func() *client.M3BatchResult {
                return nil
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeImage).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, errors.New("some error"))
            },
        },

        {
            Name: "Got error if error in response not empty",

            WantErr: true,

            Request: &DefaultRequestParams{},

            DoResponse: func() *client.M3BatchResult {
                raw := &client.M3RawResult{
                    ID:     "123456789",
                    Status: "SUCCESS",
                    Error:  "some error",
                    Data:   "hz",
                }

                result := make([]*client.M3RawResult, 0, 2)

                result = append(result, raw)

                return &client.M3BatchResult{
                    Results: result,
                }
            },

            MockBehavior: func(m *cmock.MockTransporter, request interface{}, DoResponse *client.M3BatchResult) {
                m.EXPECT().MakePayload(request, MethodDescribeImage).Return(nil, nil)
                m.EXPECT().Do(nil).Return(DoResponse, nil)
            },
        },

        {
            Name: "Got error neither result nor error in response",

            WantErr: true,

            Request: &DefaultRequestParams{},

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
                m.EXPECT().MakePayload(request, MethodDescribeImage).Return(nil, nil)
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

            _, err := s.DataImageServicer.DataImageGetList(testCase.Request.(*DefaultRequestParams))

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }

}
