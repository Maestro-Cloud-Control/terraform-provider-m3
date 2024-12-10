// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package smock is a generated GoMock package.
package smock

import (
    reflect "reflect"
    service "terraform-provider-m3/service"

    gomock "github.com/golang/mock/gomock"
)

// MockVolumeServicer is a mock of VolumeServicer interface.
type MockVolumeServicer struct {
    ctrl     *gomock.Controller
    recorder *MockVolumeServicerMockRecorder
}

// MockVolumeServicerMockRecorder is the mock recorder for MockVolumeServicer.
type MockVolumeServicerMockRecorder struct {
    mock *MockVolumeServicer
}

// NewMockVolumeServicer creates a new mock instance.
func NewMockVolumeServicer(ctrl *gomock.Controller) *MockVolumeServicer {
    mock := &MockVolumeServicer{ctrl: ctrl}
    mock.recorder = &MockVolumeServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVolumeServicer) EXPECT() *MockVolumeServicerMockRecorder {
    return m.recorder
}

// Create mocks base method.
func (m *MockVolumeServicer) Create(arg0 *service.VolumeCreateRequest) (*service.Volume, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", arg0)
    ret0, _ := ret[0].(*service.Volume)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockVolumeServicerMockRecorder) Create(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVolumeServicer)(nil).Create), arg0)
}

// CreateAndAttach mocks base method.
func (m *MockVolumeServicer) CreateAndAttach(arg0 *service.VolumeCreateAndAttachRequest) (*service.Volume, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "CreateAndAttach", arg0)
    ret0, _ := ret[0].(*service.Volume)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// CreateAndAttach indicates an expected call of CreateAndAttach.
func (mr *MockVolumeServicerMockRecorder) CreateAndAttach(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAndAttach", reflect.TypeOf((*MockVolumeServicer)(nil).CreateAndAttach), arg0)
}

// Delete mocks base method.
func (m *MockVolumeServicer) Delete(arg0 *service.VolumeDeleteRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Delete", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockVolumeServicerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVolumeServicer)(nil).Delete), arg0)
}

// Describe mocks base method.
func (m *MockVolumeServicer) Describe(arg0 *service.VolumeDescribeRequest) (*service.Volume, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Volume)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockVolumeServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockVolumeServicer)(nil).Describe), arg0)
}

// MockScriptServicer is a mock of ScriptServicer interface.
type MockScriptServicer struct {
    ctrl     *gomock.Controller
    recorder *MockScriptServicerMockRecorder
}

// MockScriptServicerMockRecorder is the mock recorder for MockScriptServicer.
type MockScriptServicerMockRecorder struct {
    mock *MockScriptServicer
}

// NewMockScriptServicer creates a new mock instance.
func NewMockScriptServicer(ctrl *gomock.Controller) *MockScriptServicer {
    mock := &MockScriptServicer{ctrl: ctrl}
    mock.recorder = &MockScriptServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScriptServicer) EXPECT() *MockScriptServicerMockRecorder {
    return m.recorder
}

// Create mocks base method.
func (m *MockScriptServicer) Create(arg0 *service.ScriptCreateRequest) (*service.Script, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", arg0)
    ret0, _ := ret[0].(*service.Script)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockScriptServicerMockRecorder) Create(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScriptServicer)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockScriptServicer) Delete(arg0 *service.ScriptDeleteRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Delete", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockScriptServicerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScriptServicer)(nil).Delete), arg0)
}

// Describe mocks base method.
func (m *MockScriptServicer) Describe(arg0 *service.ScriptDescribeRequest) (*service.Script, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Script)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockScriptServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockScriptServicer)(nil).Describe), arg0)
}

// MockScheduleServicer is a mock of ScheduleServicer interface.
type MockScheduleServicer struct {
    ctrl     *gomock.Controller
    recorder *MockScheduleServicerMockRecorder
}

// MockScheduleServicerMockRecorder is the mock recorder for MockScheduleServicer.
type MockScheduleServicerMockRecorder struct {
    mock *MockScheduleServicer
}

// NewMockScheduleServicer creates a new mock instance.
func NewMockScheduleServicer(ctrl *gomock.Controller) *MockScheduleServicer {
    mock := &MockScheduleServicer{ctrl: ctrl}
    mock.recorder = &MockScheduleServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduleServicer) EXPECT() *MockScheduleServicerMockRecorder {
    return m.recorder
}

// Create mocks base method.
func (m *MockScheduleServicer) Create(arg0 *service.RequestSchedule) (*service.Schedule, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", arg0)
    ret0, _ := ret[0].(*service.Schedule)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockScheduleServicerMockRecorder) Create(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScheduleServicer)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockScheduleServicer) Delete(arg0 *service.RequestSchedule) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Delete", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockScheduleServicerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScheduleServicer)(nil).Delete), arg0)
}

// Describe mocks base method.
func (m *MockScheduleServicer) Describe(arg0 *service.RequestSchedule) (*service.Schedule, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Schedule)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockScheduleServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockScheduleServicer)(nil).Describe), arg0)
}

// MockKeypairServicer is a mock of KeypairServicer interface.
type MockKeypairServicer struct {
    ctrl     *gomock.Controller
    recorder *MockKeypairServicerMockRecorder
}

// MockKeypairServicerMockRecorder is the mock recorder for MockKeypairServicer.
type MockKeypairServicerMockRecorder struct {
    mock *MockKeypairServicer
}

// NewMockKeypairServicer creates a new mock instance.
func NewMockKeypairServicer(ctrl *gomock.Controller) *MockKeypairServicer {
    mock := &MockKeypairServicer{ctrl: ctrl}
    mock.recorder = &MockKeypairServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeypairServicer) EXPECT() *MockKeypairServicerMockRecorder {
    return m.recorder
}

// Create mocks base method.
func (m *MockKeypairServicer) Create(arg0 *service.KeypairRequest) (*service.Keypair, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", arg0)
    ret0, _ := ret[0].(*service.Keypair)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockKeypairServicerMockRecorder) Create(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockKeypairServicer)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockKeypairServicer) Delete(arg0 *service.KeypairRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Delete", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockKeypairServicerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeypairServicer)(nil).Delete), arg0)
}

// Describe mocks base method.
func (m *MockKeypairServicer) Describe(arg0 *service.KeypairRequest) (*service.Keypair, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Keypair)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockKeypairServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockKeypairServicer)(nil).Describe), arg0)
}

// MockInstanceServicer is a mock of InstanceServicer interface.
type MockInstanceServicer struct {
    ctrl     *gomock.Controller
    recorder *MockInstanceServicerMockRecorder
}

// MockInstanceServicerMockRecorder is the mock recorder for MockInstanceServicer.
type MockInstanceServicerMockRecorder struct {
    mock *MockInstanceServicer
}

// NewMockInstanceServicer creates a new mock instance.
func NewMockInstanceServicer(ctrl *gomock.Controller) *MockInstanceServicer {
    mock := &MockInstanceServicer{ctrl: ctrl}
    mock.recorder = &MockInstanceServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstanceServicer) EXPECT() *MockInstanceServicerMockRecorder {
    return m.recorder
}

// DeleteTags mocks base method.
func (m *MockInstanceServicer) DeleteTags(arg0 *service.InstanceDeleteTagsRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "DeleteTags", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// DeleteTags indicates an expected call of DeleteTags.
func (mr *MockInstanceServicerMockRecorder) DeleteTags(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTags", reflect.TypeOf((*MockInstanceServicer)(nil).DeleteTags), arg0)
}

// Describe mocks base method.
func (m *MockInstanceServicer) Describe(arg0 *service.InstanceDescribeRequest) (*service.Instance, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Instance)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockInstanceServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockInstanceServicer)(nil).Describe), arg0)
}

// Run mocks base method.
func (m *MockInstanceServicer) Run(arg0 *service.InstanceRunRequest) (*service.Instance, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Run", arg0)
    ret0, _ := ret[0].(*service.Instance)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockInstanceServicerMockRecorder) Run(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockInstanceServicer)(nil).Run), arg0)
}

// Terminate mocks base method.
func (m *MockInstanceServicer) Terminate(arg0 *service.InstanceTerminateRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Terminate", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Terminate indicates an expected call of Terminate.
func (mr *MockInstanceServicerMockRecorder) Terminate(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Terminate", reflect.TypeOf((*MockInstanceServicer)(nil).Terminate), arg0)
}

// UnlockTermination mocks base method.
func (m *MockInstanceServicer) UnlockTermination(arg0 *service.InstanceTerminateRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "UnlockTermination", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// UnlockTermination indicates an expected call of UnlockTermination.
func (mr *MockInstanceServicerMockRecorder) UnlockTermination(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockTermination", reflect.TypeOf((*MockInstanceServicer)(nil).UnlockTermination), arg0)
}

// UpdateTags mocks base method.
func (m *MockInstanceServicer) UpdateTags(arg0 *service.InstanceUpdateTagsRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "UpdateTags", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// UpdateTags indicates an expected call of UpdateTags.
func (mr *MockInstanceServicerMockRecorder) UpdateTags(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTags", reflect.TypeOf((*MockInstanceServicer)(nil).UpdateTags), arg0)
}

// MockImageServicer is a mock of ImageServicer interface.
type MockImageServicer struct {
    ctrl     *gomock.Controller
    recorder *MockImageServicerMockRecorder
}

// MockImageServicerMockRecorder is the mock recorder for MockImageServicer.
type MockImageServicerMockRecorder struct {
    mock *MockImageServicer
}

// NewMockImageServicer creates a new mock instance.
func NewMockImageServicer(ctrl *gomock.Controller) *MockImageServicer {
    mock := &MockImageServicer{ctrl: ctrl}
    mock.recorder = &MockImageServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageServicer) EXPECT() *MockImageServicerMockRecorder {
    return m.recorder
}

// Create mocks base method.
func (m *MockImageServicer) Create(arg0 *service.ImageCreateRequest) (*service.Image, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", arg0)
    ret0, _ := ret[0].(*service.Image)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockImageServicerMockRecorder) Create(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockImageServicer)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockImageServicer) Delete(arg0 *service.DeleteImageRequest) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Delete", arg0)
    ret0, _ := ret[0].(error)
    return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockImageServicerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockImageServicer)(nil).Delete), arg0)
}

// Describe mocks base method.
func (m *MockImageServicer) Describe(arg0 *service.ImageDescribeRequest) (*service.Image, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Describe", arg0)
    ret0, _ := ret[0].(*service.Image)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockImageServicerMockRecorder) Describe(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockImageServicer)(nil).Describe), arg0)
}

// MockDataImageServicer is a mock of DataImageServicer interface.
type MockDataImageServicer struct {
    ctrl     *gomock.Controller
    recorder *MockDataImageServicerMockRecorder
}

// MockDataImageServicerMockRecorder is the mock recorder for MockDataImageServicer.
type MockDataImageServicerMockRecorder struct {
    mock *MockDataImageServicer
}

// NewMockDataImageServicer creates a new mock instance.
func NewMockDataImageServicer(ctrl *gomock.Controller) *MockDataImageServicer {
    mock := &MockDataImageServicer{ctrl: ctrl}
    mock.recorder = &MockDataImageServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataImageServicer) EXPECT() *MockDataImageServicerMockRecorder {
    return m.recorder
}

// DataImageGetList mocks base method.
func (m *MockDataImageServicer) DataImageGetList(arg0 *service.DefaultRequestParams) (*[]service.Image, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "DataImageGetList", arg0)
    ret0, _ := ret[0].(*[]service.Image)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// DataImageGetList indicates an expected call of DataImageGetList.
func (mr *MockDataImageServicerMockRecorder) DataImageGetList(arg0 interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataImageGetList", reflect.TypeOf((*MockDataImageServicer)(nil).DataImageGetList), arg0)
}

// MockDataPlacementServicer is a mock of DataPlacementServicer interface.
type MockDataPlacementServicer struct {
    ctrl     *gomock.Controller
    recorder *MockDataPlacementServicerMockRecorder
}

// MockDataPlacementServicerMockRecorder is the mock recorder for MockDataPlacementServicer.
type MockDataPlacementServicerMockRecorder struct {
    mock *MockDataPlacementServicer
}

// NewMockDataPlacementServicer creates a new mock instance.
func NewMockDataPlacementServicer(ctrl *gomock.Controller) *MockDataPlacementServicer {
    mock := &MockDataPlacementServicer{ctrl: ctrl}
    mock.recorder = &MockDataPlacementServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataPlacementServicer) EXPECT() *MockDataPlacementServicerMockRecorder {
    return m.recorder
}

// DataPlacementGetList mocks base method.
func (m *MockDataPlacementServicer) DataPlacementGetList(request *service.PlacementParamsRequest) (*[]service.DataItem, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "DataPlacementGetList", request)
    ret0, _ := ret[0].(*[]service.DataItem)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// DataPlacementGetList indicates an expected call of DataPlacementGetList.
func (mr *MockDataPlacementServicerMockRecorder) DataPlacementGetList(request interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataPlacementGetList", reflect.TypeOf((*MockDataPlacementServicer)(nil).DataPlacementGetList), request)
}

// MockDataChefServicer is a mock of DataChefServicer interface.
type MockDataChefServicer struct {
    ctrl     *gomock.Controller
    recorder *MockDataChefServicerMockRecorder
}

// MockDataChefServicerMockRecorder is the mock recorder for MockDataChefServicer.
type MockDataChefServicerMockRecorder struct {
    mock *MockDataChefServicer
}

// NewMockDataChefServicer creates a new mock instance.
func NewMockDataChefServicer(ctrl *gomock.Controller) *MockDataChefServicer {
    mock := &MockDataChefServicer{ctrl: ctrl}
    mock.recorder = &MockDataChefServicerMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataChefServicer) EXPECT() *MockDataChefServicerMockRecorder {
    return m.recorder
}

// DataChefGetList mocks base method.
func (m *MockDataChefServicer) DataChefGetList(request *service.DefaultRequestParams) (*service.DataChef, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "DataChefGetList", request)
    ret0, _ := ret[0].(*service.DataChef)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// DataChefGetList indicates an expected call of DataChefGetList.
func (mr *MockDataChefServicerMockRecorder) DataChefGetList(request interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataChefGetList", reflect.TypeOf((*MockDataChefServicer)(nil).DataChefGetList), request)
}
