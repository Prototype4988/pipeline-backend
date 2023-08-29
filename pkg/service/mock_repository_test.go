// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/instill-ai/pipeline-backend/pkg/repository (interfaces: Repository)

// Package service_test is a generated GoMock package.
package service_test

import (
	context "context"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	datamodel "github.com/instill-ai/pipeline-backend/pkg/datamodel"
	filtering "go.einride.tech/aip/filtering"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateUserPipeline mocks base method.
func (m *MockRepository) CreateUserPipeline(arg0 context.Context, arg1, arg2 string, arg3 *datamodel.Pipeline) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserPipeline", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserPipeline indicates an expected call of CreateUserPipeline.
func (mr *MockRepositoryMockRecorder) CreateUserPipeline(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserPipeline", reflect.TypeOf((*MockRepository)(nil).CreateUserPipeline), arg0, arg1, arg2, arg3)
}

// CreateUserPipelineRelease mocks base method.
func (m *MockRepository) CreateUserPipelineRelease(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4 *datamodel.PipelineRelease) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserPipelineRelease", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserPipelineRelease indicates an expected call of CreateUserPipelineRelease.
func (mr *MockRepositoryMockRecorder) CreateUserPipelineRelease(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserPipelineRelease", reflect.TypeOf((*MockRepository)(nil).CreateUserPipelineRelease), arg0, arg1, arg2, arg3, arg4)
}

// DeleteUserPipelineByID mocks base method.
func (m *MockRepository) DeleteUserPipelineByID(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserPipelineByID", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserPipelineByID indicates an expected call of DeleteUserPipelineByID.
func (mr *MockRepositoryMockRecorder) DeleteUserPipelineByID(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserPipelineByID", reflect.TypeOf((*MockRepository)(nil).DeleteUserPipelineByID), arg0, arg1, arg2, arg3)
}

// DeleteUserPipelineReleaseByID mocks base method.
func (m *MockRepository) DeleteUserPipelineReleaseByID(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserPipelineReleaseByID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserPipelineReleaseByID indicates an expected call of DeleteUserPipelineReleaseByID.
func (mr *MockRepositoryMockRecorder) DeleteUserPipelineReleaseByID(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserPipelineReleaseByID", reflect.TypeOf((*MockRepository)(nil).DeleteUserPipelineReleaseByID), arg0, arg1, arg2, arg3, arg4)
}

// GetPipelineByIDAdmin mocks base method.
func (m *MockRepository) GetPipelineByIDAdmin(arg0 context.Context, arg1 string, arg2 bool) (*datamodel.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineByIDAdmin", arg0, arg1, arg2)
	ret0, _ := ret[0].(*datamodel.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineByIDAdmin indicates an expected call of GetPipelineByIDAdmin.
func (mr *MockRepositoryMockRecorder) GetPipelineByIDAdmin(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineByIDAdmin", reflect.TypeOf((*MockRepository)(nil).GetPipelineByIDAdmin), arg0, arg1, arg2)
}

// GetPipelineByUID mocks base method.
func (m *MockRepository) GetPipelineByUID(arg0 context.Context, arg1 string, arg2 uuid.UUID, arg3 bool) (*datamodel.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineByUID", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*datamodel.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineByUID indicates an expected call of GetPipelineByUID.
func (mr *MockRepositoryMockRecorder) GetPipelineByUID(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineByUID", reflect.TypeOf((*MockRepository)(nil).GetPipelineByUID), arg0, arg1, arg2, arg3)
}

// GetPipelineByUIDAdmin mocks base method.
func (m *MockRepository) GetPipelineByUIDAdmin(arg0 context.Context, arg1 uuid.UUID, arg2 bool) (*datamodel.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineByUIDAdmin", arg0, arg1, arg2)
	ret0, _ := ret[0].(*datamodel.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineByUIDAdmin indicates an expected call of GetPipelineByUIDAdmin.
func (mr *MockRepositoryMockRecorder) GetPipelineByUIDAdmin(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineByUIDAdmin", reflect.TypeOf((*MockRepository)(nil).GetPipelineByUIDAdmin), arg0, arg1, arg2)
}

// GetUserPipelineByID mocks base method.
func (m *MockRepository) GetUserPipelineByID(arg0 context.Context, arg1, arg2, arg3 string, arg4 bool) (*datamodel.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPipelineByID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*datamodel.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPipelineByID indicates an expected call of GetUserPipelineByID.
func (mr *MockRepositoryMockRecorder) GetUserPipelineByID(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPipelineByID", reflect.TypeOf((*MockRepository)(nil).GetUserPipelineByID), arg0, arg1, arg2, arg3, arg4)
}

// GetUserPipelineReleaseByID mocks base method.
func (m *MockRepository) GetUserPipelineReleaseByID(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4 string, arg5 bool) (*datamodel.PipelineRelease, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPipelineReleaseByID", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*datamodel.PipelineRelease)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPipelineReleaseByID indicates an expected call of GetUserPipelineReleaseByID.
func (mr *MockRepositoryMockRecorder) GetUserPipelineReleaseByID(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPipelineReleaseByID", reflect.TypeOf((*MockRepository)(nil).GetUserPipelineReleaseByID), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetUserPipelineReleaseByUID mocks base method.
func (m *MockRepository) GetUserPipelineReleaseByUID(arg0 context.Context, arg1, arg2 string, arg3, arg4 uuid.UUID, arg5 bool) (*datamodel.PipelineRelease, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPipelineReleaseByUID", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*datamodel.PipelineRelease)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPipelineReleaseByUID indicates an expected call of GetUserPipelineReleaseByUID.
func (mr *MockRepositoryMockRecorder) GetUserPipelineReleaseByUID(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPipelineReleaseByUID", reflect.TypeOf((*MockRepository)(nil).GetUserPipelineReleaseByUID), arg0, arg1, arg2, arg3, arg4, arg5)
}

// ListPipelineReleasesAdmin mocks base method.
func (m *MockRepository) ListPipelineReleasesAdmin(arg0 context.Context, arg1 int64, arg2 string, arg3 bool, arg4 filtering.Filter) ([]*datamodel.PipelineRelease, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPipelineReleasesAdmin", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*datamodel.PipelineRelease)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListPipelineReleasesAdmin indicates an expected call of ListPipelineReleasesAdmin.
func (mr *MockRepositoryMockRecorder) ListPipelineReleasesAdmin(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPipelineReleasesAdmin", reflect.TypeOf((*MockRepository)(nil).ListPipelineReleasesAdmin), arg0, arg1, arg2, arg3, arg4)
}

// ListPipelines mocks base method.
func (m *MockRepository) ListPipelines(arg0 context.Context, arg1 string, arg2 int64, arg3 string, arg4 bool, arg5 filtering.Filter) ([]*datamodel.Pipeline, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPipelines", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]*datamodel.Pipeline)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListPipelines indicates an expected call of ListPipelines.
func (mr *MockRepositoryMockRecorder) ListPipelines(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPipelines", reflect.TypeOf((*MockRepository)(nil).ListPipelines), arg0, arg1, arg2, arg3, arg4, arg5)
}

// ListPipelinesAdmin mocks base method.
func (m *MockRepository) ListPipelinesAdmin(arg0 context.Context, arg1 int64, arg2 string, arg3 bool, arg4 filtering.Filter) ([]*datamodel.Pipeline, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPipelinesAdmin", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*datamodel.Pipeline)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListPipelinesAdmin indicates an expected call of ListPipelinesAdmin.
func (mr *MockRepositoryMockRecorder) ListPipelinesAdmin(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPipelinesAdmin", reflect.TypeOf((*MockRepository)(nil).ListPipelinesAdmin), arg0, arg1, arg2, arg3, arg4)
}

// ListUserPipelineReleases mocks base method.
func (m *MockRepository) ListUserPipelineReleases(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4 int64, arg5 string, arg6 bool, arg7 filtering.Filter) ([]*datamodel.PipelineRelease, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserPipelineReleases", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].([]*datamodel.PipelineRelease)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListUserPipelineReleases indicates an expected call of ListUserPipelineReleases.
func (mr *MockRepositoryMockRecorder) ListUserPipelineReleases(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserPipelineReleases", reflect.TypeOf((*MockRepository)(nil).ListUserPipelineReleases), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

// ListUserPipelines mocks base method.
func (m *MockRepository) ListUserPipelines(arg0 context.Context, arg1, arg2 string, arg3 int64, arg4 string, arg5 bool, arg6 filtering.Filter) ([]*datamodel.Pipeline, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserPipelines", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].([]*datamodel.Pipeline)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListUserPipelines indicates an expected call of ListUserPipelines.
func (mr *MockRepositoryMockRecorder) ListUserPipelines(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserPipelines", reflect.TypeOf((*MockRepository)(nil).ListUserPipelines), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// UpdateUserPipelineByID mocks base method.
func (m *MockRepository) UpdateUserPipelineByID(arg0 context.Context, arg1, arg2, arg3 string, arg4 *datamodel.Pipeline) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPipelineByID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPipelineByID indicates an expected call of UpdateUserPipelineByID.
func (mr *MockRepositoryMockRecorder) UpdateUserPipelineByID(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPipelineByID", reflect.TypeOf((*MockRepository)(nil).UpdateUserPipelineByID), arg0, arg1, arg2, arg3, arg4)
}

// UpdateUserPipelineIDByID mocks base method.
func (m *MockRepository) UpdateUserPipelineIDByID(arg0 context.Context, arg1, arg2, arg3, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPipelineIDByID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPipelineIDByID indicates an expected call of UpdateUserPipelineIDByID.
func (mr *MockRepositoryMockRecorder) UpdateUserPipelineIDByID(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPipelineIDByID", reflect.TypeOf((*MockRepository)(nil).UpdateUserPipelineIDByID), arg0, arg1, arg2, arg3, arg4)
}

// UpdateUserPipelineReleaseByID mocks base method.
func (m *MockRepository) UpdateUserPipelineReleaseByID(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4 string, arg5 *datamodel.PipelineRelease) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPipelineReleaseByID", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPipelineReleaseByID indicates an expected call of UpdateUserPipelineReleaseByID.
func (mr *MockRepositoryMockRecorder) UpdateUserPipelineReleaseByID(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPipelineReleaseByID", reflect.TypeOf((*MockRepository)(nil).UpdateUserPipelineReleaseByID), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateUserPipelineReleaseIDByID mocks base method.
func (m *MockRepository) UpdateUserPipelineReleaseIDByID(arg0 context.Context, arg1, arg2 string, arg3 uuid.UUID, arg4, arg5 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPipelineReleaseIDByID", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPipelineReleaseIDByID indicates an expected call of UpdateUserPipelineReleaseIDByID.
func (mr *MockRepositoryMockRecorder) UpdateUserPipelineReleaseIDByID(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPipelineReleaseIDByID", reflect.TypeOf((*MockRepository)(nil).UpdateUserPipelineReleaseIDByID), arg0, arg1, arg2, arg3, arg4, arg5)
}
