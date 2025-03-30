// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"
	"time"

	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"
)

type FakeMaker struct {
	CreateTokenStub        func(string, string, time.Duration) (string, *tokenMaker.Payload, error)
	createTokenMutex       sync.RWMutex
	createTokenArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 time.Duration
	}
	createTokenReturns struct {
		result1 string
		result2 *tokenMaker.Payload
		result3 error
	}
	createTokenReturnsOnCall map[int]struct {
		result1 string
		result2 *tokenMaker.Payload
		result3 error
	}
	VerifyTokenStub        func(string) (*tokenMaker.Payload, error)
	verifyTokenMutex       sync.RWMutex
	verifyTokenArgsForCall []struct {
		arg1 string
	}
	verifyTokenReturns struct {
		result1 *tokenMaker.Payload
		result2 error
	}
	verifyTokenReturnsOnCall map[int]struct {
		result1 *tokenMaker.Payload
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMaker) CreateToken(arg1 string, arg2 string, arg3 time.Duration) (string, *tokenMaker.Payload, error) {
	fake.createTokenMutex.Lock()
	ret, specificReturn := fake.createTokenReturnsOnCall[len(fake.createTokenArgsForCall)]
	fake.createTokenArgsForCall = append(fake.createTokenArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 time.Duration
	}{arg1, arg2, arg3})
	stub := fake.CreateTokenStub
	fakeReturns := fake.createTokenReturns
	fake.recordInvocation("CreateToken", []interface{}{arg1, arg2, arg3})
	fake.createTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeMaker) CreateTokenCallCount() int {
	fake.createTokenMutex.RLock()
	defer fake.createTokenMutex.RUnlock()
	return len(fake.createTokenArgsForCall)
}

func (fake *FakeMaker) CreateTokenCalls(stub func(string, string, time.Duration) (string, *tokenMaker.Payload, error)) {
	fake.createTokenMutex.Lock()
	defer fake.createTokenMutex.Unlock()
	fake.CreateTokenStub = stub
}

func (fake *FakeMaker) CreateTokenArgsForCall(i int) (string, string, time.Duration) {
	fake.createTokenMutex.RLock()
	defer fake.createTokenMutex.RUnlock()
	argsForCall := fake.createTokenArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeMaker) CreateTokenReturns(result1 string, result2 *tokenMaker.Payload, result3 error) {
	fake.createTokenMutex.Lock()
	defer fake.createTokenMutex.Unlock()
	fake.CreateTokenStub = nil
	fake.createTokenReturns = struct {
		result1 string
		result2 *tokenMaker.Payload
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeMaker) CreateTokenReturnsOnCall(i int, result1 string, result2 *tokenMaker.Payload, result3 error) {
	fake.createTokenMutex.Lock()
	defer fake.createTokenMutex.Unlock()
	fake.CreateTokenStub = nil
	if fake.createTokenReturnsOnCall == nil {
		fake.createTokenReturnsOnCall = make(map[int]struct {
			result1 string
			result2 *tokenMaker.Payload
			result3 error
		})
	}
	fake.createTokenReturnsOnCall[i] = struct {
		result1 string
		result2 *tokenMaker.Payload
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeMaker) VerifyToken(arg1 string) (*tokenMaker.Payload, error) {
	fake.verifyTokenMutex.Lock()
	ret, specificReturn := fake.verifyTokenReturnsOnCall[len(fake.verifyTokenArgsForCall)]
	fake.verifyTokenArgsForCall = append(fake.verifyTokenArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.VerifyTokenStub
	fakeReturns := fake.verifyTokenReturns
	fake.recordInvocation("VerifyToken", []interface{}{arg1})
	fake.verifyTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeMaker) VerifyTokenCallCount() int {
	fake.verifyTokenMutex.RLock()
	defer fake.verifyTokenMutex.RUnlock()
	return len(fake.verifyTokenArgsForCall)
}

func (fake *FakeMaker) VerifyTokenCalls(stub func(string) (*tokenMaker.Payload, error)) {
	fake.verifyTokenMutex.Lock()
	defer fake.verifyTokenMutex.Unlock()
	fake.VerifyTokenStub = stub
}

func (fake *FakeMaker) VerifyTokenArgsForCall(i int) string {
	fake.verifyTokenMutex.RLock()
	defer fake.verifyTokenMutex.RUnlock()
	argsForCall := fake.verifyTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeMaker) VerifyTokenReturns(result1 *tokenMaker.Payload, result2 error) {
	fake.verifyTokenMutex.Lock()
	defer fake.verifyTokenMutex.Unlock()
	fake.VerifyTokenStub = nil
	fake.verifyTokenReturns = struct {
		result1 *tokenMaker.Payload
		result2 error
	}{result1, result2}
}

func (fake *FakeMaker) VerifyTokenReturnsOnCall(i int, result1 *tokenMaker.Payload, result2 error) {
	fake.verifyTokenMutex.Lock()
	defer fake.verifyTokenMutex.Unlock()
	fake.VerifyTokenStub = nil
	if fake.verifyTokenReturnsOnCall == nil {
		fake.verifyTokenReturnsOnCall = make(map[int]struct {
			result1 *tokenMaker.Payload
			result2 error
		})
	}
	fake.verifyTokenReturnsOnCall[i] = struct {
		result1 *tokenMaker.Payload
		result2 error
	}{result1, result2}
}

func (fake *FakeMaker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createTokenMutex.RLock()
	defer fake.createTokenMutex.RUnlock()
	fake.verifyTokenMutex.RLock()
	defer fake.verifyTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMaker) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ tokenMaker.Maker = new(FakeMaker)
