// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	redisStore "github.com/gin-contrib/sessions/redis"
	"sync"
)

// Ensure, that CookieStoreMock does implement CookieStore.
// If this is not the case, regenerate this file with moq.
var _ CookieStore = &CookieStoreMock{}

// CookieStoreMock is a mock implementation of CookieStore.
//
// 	func TestSomethingThatUsesCookieStore(t *testing.T) {
//
// 		// make and configure a mocked CookieStore
// 		mockedCookieStore := &CookieStoreMock{
// 			NewCookieStoreFunc: func() (redisStore.Store, error) {
// 				panic("mock out the NewCookieStore method")
// 			},
// 		}
//
// 		// use mockedCookieStore in code that requires CookieStore
// 		// and then make assertions.
//
// 	}
type CookieStoreMock struct {
	// NewCookieStoreFunc mocks the NewCookieStore method.
	NewCookieStoreFunc func() (redisStore.Store, error)

	// calls tracks calls to the methods.
	calls struct {
		// NewCookieStore holds details about calls to the NewCookieStore method.
		NewCookieStore []struct {
		}
	}
	lockNewCookieStore sync.RWMutex
}

// NewCookieStore calls NewCookieStoreFunc.
func (mock *CookieStoreMock) NewCookieStore() (redisStore.Store, error) {
	if mock.NewCookieStoreFunc == nil {
		panic("CookieStoreMock.NewCookieStoreFunc: method is nil but CookieStore.NewCookieStore was just called")
	}
	callInfo := struct {
	}{}
	mock.lockNewCookieStore.Lock()
	mock.calls.NewCookieStore = append(mock.calls.NewCookieStore, callInfo)
	mock.lockNewCookieStore.Unlock()
	return mock.NewCookieStoreFunc()
}

// NewCookieStoreCalls gets all the calls that were made to NewCookieStore.
// Check the length with:
//     len(mockedCookieStore.NewCookieStoreCalls())
func (mock *CookieStoreMock) NewCookieStoreCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockNewCookieStore.RLock()
	calls = mock.calls.NewCookieStore
	mock.lockNewCookieStore.RUnlock()
	return calls
}