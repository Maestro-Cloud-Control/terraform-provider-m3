package provider

import (
    "errors"
    "github.com/golang/mock/gomock"
    "terraform-provider-m3/service"
    smock "terraform-provider-m3/service/mock"
    "testing"
)

func TestWait_Wait(t *testing.T) {
    type (
        MockBehavior func(m *smock.MockKeypairServicer)
        TestCase     struct {
            Name            string
            WantErr         bool
            MockBehavior    MockBehavior
            WaitAction      func(*service.Service) waitAction
            WaitCompareFunc waitCompareFunc
        }
    )

    testTable := []TestCase{
        {
            Name:            "OK",
            WantErr:         false,
            WaitCompareFunc: defaultWaitCompareFunc(),
            WaitAction: func(s *service.Service) waitAction {
                return func() (interface{}, error) {
                    return s.KeypairServicer.Describe(nil)
                }
            },
            MockBehavior: func(m *smock.MockKeypairServicer) {
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, nil)
            },
        },

        {
            Name:            "Got an error if run out of attempts",
            WantErr:         true,
            WaitCompareFunc: defaultWaitCompareFunc(),
            WaitAction: func(s *service.Service) waitAction {
                return func() (interface{}, error) {
                    return s.KeypairServicer.Describe(nil)
                }
            },
            MockBehavior: func(m *smock.MockKeypairServicer) {
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
                m.EXPECT().Describe(nil).Return(nil, errors.New("some error"))
            },
        },
    }

    for _, testCase := range testTable {
        t.Run(testCase.Name, func(t *testing.T) {
            ctl := gomock.NewController(t)
            defer ctl.Finish()

            mockKeypairServicer := smock.NewMockKeypairServicer(ctl)
            testCase.MockBehavior(mockKeypairServicer)

            s := service.Service{}
            s.KeypairServicer = mockKeypairServicer

            w := wait{
                Delay:     1,
                Attempts:  5,
                Action:    testCase.WaitAction(&s),
                CompareFn: testCase.WaitCompareFunc,
            }
            _, err := w.Wait()

            if (!testCase.WantErr && err != nil) || (testCase.WantErr && err == nil) {
                t.Fatal()
            }

        })
    }
}
