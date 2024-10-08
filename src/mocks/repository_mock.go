// // Code generated by MockGen. DO NOT EDIT.
// // Source: interface.go

// // Package mocks is a generated GoMock package.
package mocks

// import (
// 	context "context"
// 	gomock "github.com/golang/mock/gomock"
// 	models "github.com/kalpit-sharma-dev/parkinglot-service/src/models"
// 	reflect "reflect"
// )

// // MockDatabaseRepository is a mock of DatabaseRepository interface
// type MockDatabaseRepository struct {
// 	ctrl     *gomock.Controller
// 	recorder *MockDatabaseRepositoryMockRecorder
// }

// // MockDatabaseRepositoryMockRecorder is the mock recorder for MockDatabaseRepository
// type MockDatabaseRepositoryMockRecorder struct {
// 	mock *MockDatabaseRepository
// }

// // NewMockDatabaseRepository creates a new mock instance
// func NewMockDatabaseRepository(ctrl *gomock.Controller) *MockDatabaseRepository {
// 	mock := &MockDatabaseRepository{ctrl: ctrl}
// 	mock.recorder = &MockDatabaseRepositoryMockRecorder{mock}
// 	return mock
// }

// // EXPECT returns an object that allows the caller to indicate expected use
// func (m *MockDatabaseRepository) EXPECT() *MockDatabaseRepositoryMockRecorder {
// 	return m.recorder
// }

// // CreateSlotEvent mocks base method
// func (m *MockDatabaseRepository) CreateSlotEvent(ctx context.Context, req models.Slot) (models.Slot, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "CreateSlotEvent", ctx, req)
// 	ret0, _ := ret[0].(models.Slot)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // CreateSlotEvent indicates an expected call of CreateSlotEvent
// func (mr *MockDatabaseRepositoryMockRecorder) CreateSlotEvent(ctx, req interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSlotEvent", reflect.TypeOf((*MockDatabaseRepository)(nil).CreateSlotEvent), ctx, req)
// }

// // GetAllCarsWithColor mocks base method
// func (m *MockDatabaseRepository) GetAllCarsWithColor(ctx context.Context, reqColor string) ([]models.Vehicle, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "GetAllCarsWithColor", ctx, reqColor)
// 	ret0, _ := ret[0].([]models.Vehicle)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // GetAllCarsWithColor indicates an expected call of GetAllCarsWithColor
// func (mr *MockDatabaseRepositoryMockRecorder) GetAllCarsWithColor(ctx, reqColor interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCarsWithColor", reflect.TypeOf((*MockDatabaseRepository)(nil).GetAllCarsWithColor), ctx, reqColor)
// }

// // GetSlotNumberWithCarID mocks base method
// func (m *MockDatabaseRepository) GetSlotNumberWithCarID(ctx context.Context, reqNumber string) (models.Vehicle, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "GetSlotNumberWithCarID", ctx, reqNumber)
// 	ret0, _ := ret[0].(models.Vehicle)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // GetSlotNumberWithCarID indicates an expected call of GetSlotNumberWithCarID
// func (mr *MockDatabaseRepositoryMockRecorder) GetSlotNumberWithCarID(ctx, reqNumber interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlotNumberWithCarID", reflect.TypeOf((*MockDatabaseRepository)(nil).GetSlotNumberWithCarID), ctx, reqNumber)
// }

// // GetAllSlotNumberWithColor mocks base method
// func (m *MockDatabaseRepository) GetAllSlotNumberWithColor(ctx context.Context, reqColor string) ([]models.Vehicle, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "GetAllSlotNumberWithColor", ctx, reqColor)
// 	ret0, _ := ret[0].([]models.Vehicle)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // GetAllSlotNumberWithColor indicates an expected call of GetAllSlotNumberWithColor
// func (mr *MockDatabaseRepositoryMockRecorder) GetAllSlotNumberWithColor(ctx, reqColor interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSlotNumberWithColor", reflect.TypeOf((*MockDatabaseRepository)(nil).GetAllSlotNumberWithColor), ctx, reqColor)
// }

// // CreateParkEvent mocks base method
// func (m *MockDatabaseRepository) CreateParkEvent(ctx context.Context, req models.Vehicle) (models.Vehicle, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "CreateParkEvent", ctx, req)
// 	ret0, _ := ret[0].(models.Vehicle)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // CreateParkEvent indicates an expected call of CreateParkEvent
// func (mr *MockDatabaseRepositoryMockRecorder) CreateParkEvent(ctx, req interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateParkEvent", reflect.TypeOf((*MockDatabaseRepository)(nil).CreateParkEvent), ctx, req)
// }

// // ExitParkEvent mocks base method
// func (m *MockDatabaseRepository) ExitParkEvent(ctx context.Context, req models.Vehicle) (models.Vehicle, error) {
// 	m.ctrl.T.Helper()
// 	ret := m.ctrl.Call(m, "ExitParkEvent", ctx, req)
// 	ret0, _ := ret[0].(models.Vehicle)
// 	ret1, _ := ret[1].(error)
// 	return ret0, ret1
// }

// // ExitParkEvent indicates an expected call of ExitParkEvent
// func (mr *MockDatabaseRepositoryMockRecorder) ExitParkEvent(ctx, req interface{}) *gomock.Call {
// 	mr.mock.ctrl.T.Helper()
// 	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExitParkEvent", reflect.TypeOf((*MockDatabaseRepository)(nil).ExitParkEvent), ctx, req)
// }
