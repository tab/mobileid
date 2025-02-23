// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -destination=client_mock.go -package=mobileid
//

// Package mobileid is a generated GoMock package.
package mobileid

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockClient) CreateSession(ctx context.Context, phoneNumber, nationalIdentityNumber string) (*Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, phoneNumber, nationalIdentityNumber)
	ret0, _ := ret[0].(*Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockClientMockRecorder) CreateSession(ctx, phoneNumber, nationalIdentityNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockClient)(nil).CreateSession), ctx, phoneNumber, nationalIdentityNumber)
}

// FetchSession mocks base method.
func (m *MockClient) FetchSession(ctx context.Context, sessionId string) (*Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchSession", ctx, sessionId)
	ret0, _ := ret[0].(*Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchSession indicates an expected call of FetchSession.
func (mr *MockClientMockRecorder) FetchSession(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchSession", reflect.TypeOf((*MockClient)(nil).FetchSession), ctx, sessionId)
}

// Validate mocks base method.
func (m *MockClient) Validate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockClientMockRecorder) Validate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockClient)(nil).Validate))
}

// WithHashType mocks base method.
func (m *MockClient) WithHashType(hashType string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithHashType", hashType)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithHashType indicates an expected call of WithHashType.
func (mr *MockClientMockRecorder) WithHashType(hashType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithHashType", reflect.TypeOf((*MockClient)(nil).WithHashType), hashType)
}

// WithLanguage mocks base method.
func (m *MockClient) WithLanguage(language string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithLanguage", language)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithLanguage indicates an expected call of WithLanguage.
func (mr *MockClientMockRecorder) WithLanguage(language any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLanguage", reflect.TypeOf((*MockClient)(nil).WithLanguage), language)
}

// WithRelyingPartyName mocks base method.
func (m *MockClient) WithRelyingPartyName(name string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithRelyingPartyName", name)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithRelyingPartyName indicates an expected call of WithRelyingPartyName.
func (mr *MockClientMockRecorder) WithRelyingPartyName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRelyingPartyName", reflect.TypeOf((*MockClient)(nil).WithRelyingPartyName), name)
}

// WithRelyingPartyUUID mocks base method.
func (m *MockClient) WithRelyingPartyUUID(id string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithRelyingPartyUUID", id)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithRelyingPartyUUID indicates an expected call of WithRelyingPartyUUID.
func (mr *MockClientMockRecorder) WithRelyingPartyUUID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRelyingPartyUUID", reflect.TypeOf((*MockClient)(nil).WithRelyingPartyUUID), id)
}

// WithText mocks base method.
func (m *MockClient) WithText(text string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithText", text)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithText indicates an expected call of WithText.
func (mr *MockClientMockRecorder) WithText(text any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithText", reflect.TypeOf((*MockClient)(nil).WithText), text)
}

// WithTextFormat mocks base method.
func (m *MockClient) WithTextFormat(format string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTextFormat", format)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithTextFormat indicates an expected call of WithTextFormat.
func (mr *MockClientMockRecorder) WithTextFormat(format any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTextFormat", reflect.TypeOf((*MockClient)(nil).WithTextFormat), format)
}

// WithTimeout mocks base method.
func (m *MockClient) WithTimeout(timeout time.Duration) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTimeout", timeout)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithTimeout indicates an expected call of WithTimeout.
func (mr *MockClientMockRecorder) WithTimeout(timeout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTimeout", reflect.TypeOf((*MockClient)(nil).WithTimeout), timeout)
}

// WithURL mocks base method.
func (m *MockClient) WithURL(url string) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithURL", url)
	ret0, _ := ret[0].(Client)
	return ret0
}

// WithURL indicates an expected call of WithURL.
func (mr *MockClientMockRecorder) WithURL(url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithURL", reflect.TypeOf((*MockClient)(nil).WithURL), url)
}
