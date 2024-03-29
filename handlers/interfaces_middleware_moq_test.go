// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"sync"
)

// Ensure, that MiddlewareMock does implement Middleware.
// If this is not the case, regenerate this file with moq.
var _ Middleware = &MiddlewareMock{}

// MiddlewareMock is a mock implementation of Middleware.
//
// 	func TestSomethingThatUsesMiddleware(t *testing.T) {
//
// 		// make and configure a mocked Middleware
// 		mockedMiddleware := &MiddlewareMock{
// 			AuthCookieMiddlewareFunc: func() gin.HandlerFunc {
// 				panic("mock out the AuthCookieMiddleware method")
// 			},
// 			AuthJWTMiddlewareFunc: func(cfg *viper.Viper) gin.HandlerFunc {
// 				panic("mock out the AuthJWTMiddleware method")
// 			},
// 		}
//
// 		// use mockedMiddleware in code that requires Middleware
// 		// and then make assertions.
//
// 	}
type MiddlewareMock struct {
	// AuthCookieMiddlewareFunc mocks the AuthCookieMiddleware method.
	AuthCookieMiddlewareFunc func() gin.HandlerFunc

	// AuthJWTMiddlewareFunc mocks the AuthJWTMiddleware method.
	AuthJWTMiddlewareFunc func(cfg *viper.Viper) gin.HandlerFunc

	// calls tracks calls to the methods.
	calls struct {
		// AuthCookieMiddleware holds details about calls to the AuthCookieMiddleware method.
		AuthCookieMiddleware []struct {
		}
		// AuthJWTMiddleware holds details about calls to the AuthJWTMiddleware method.
		AuthJWTMiddleware []struct {
			// Cfg is the cfg argument value.
			Cfg *viper.Viper
		}
	}
	lockAuthCookieMiddleware sync.RWMutex
	lockAuthJWTMiddleware    sync.RWMutex
}

// AuthCookieMiddleware calls AuthCookieMiddlewareFunc.
func (mock *MiddlewareMock) AuthCookieMiddleware() gin.HandlerFunc {
	if mock.AuthCookieMiddlewareFunc == nil {
		panic("MiddlewareMock.AuthCookieMiddlewareFunc: method is nil but Middleware.AuthCookieMiddleware was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAuthCookieMiddleware.Lock()
	mock.calls.AuthCookieMiddleware = append(mock.calls.AuthCookieMiddleware, callInfo)
	mock.lockAuthCookieMiddleware.Unlock()
	return mock.AuthCookieMiddlewareFunc()
}

// AuthCookieMiddlewareCalls gets all the calls that were made to AuthCookieMiddleware.
// Check the length with:
//     len(mockedMiddleware.AuthCookieMiddlewareCalls())
func (mock *MiddlewareMock) AuthCookieMiddlewareCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAuthCookieMiddleware.RLock()
	calls = mock.calls.AuthCookieMiddleware
	mock.lockAuthCookieMiddleware.RUnlock()
	return calls
}

// AuthJWTMiddleware calls AuthJWTMiddlewareFunc.
func (mock *MiddlewareMock) AuthJWTMiddleware(cfg *viper.Viper) gin.HandlerFunc {
	if mock.AuthJWTMiddlewareFunc == nil {
		panic("MiddlewareMock.AuthJWTMiddlewareFunc: method is nil but Middleware.AuthJWTMiddleware was just called")
	}
	callInfo := struct {
		Cfg *viper.Viper
	}{
		Cfg: cfg,
	}
	mock.lockAuthJWTMiddleware.Lock()
	mock.calls.AuthJWTMiddleware = append(mock.calls.AuthJWTMiddleware, callInfo)
	mock.lockAuthJWTMiddleware.Unlock()
	return mock.AuthJWTMiddlewareFunc(cfg)
}

// AuthJWTMiddlewareCalls gets all the calls that were made to AuthJWTMiddleware.
// Check the length with:
//     len(mockedMiddleware.AuthJWTMiddlewareCalls())
func (mock *MiddlewareMock) AuthJWTMiddlewareCalls() []struct {
	Cfg *viper.Viper
} {
	var calls []struct {
		Cfg *viper.Viper
	}
	mock.lockAuthJWTMiddleware.RLock()
	calls = mock.calls.AuthJWTMiddleware
	mock.lockAuthJWTMiddleware.RUnlock()
	return calls
}
