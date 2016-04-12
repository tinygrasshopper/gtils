package httptest

import (
	"io"

	"github.com/pivotalservices/gtils/http"
)

//EntityCaptcherFunc -
type EntityCaptcherFunc func(http.HttpRequestEntity)

//MockGateway -
type MockGateway struct {
	FakeGetAdaptor  http.RequestAdaptor
	FakePutAdaptor  http.RequestAdaptor
	FakePostAdaptor http.RequestAdaptor
	Capture         EntityCaptcherFunc
}

//Get -
func (gateway *MockGateway) Get(entity http.HttpRequestEntity) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakeGetAdaptor
}

//Put -
func (gateway *MockGateway) Put(entity http.HttpRequestEntity, body io.Reader) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakePutAdaptor
}

//Post -
func (gateway *MockGateway) Post(entity http.HttpRequestEntity, body io.Reader) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakePostAdaptor
}
