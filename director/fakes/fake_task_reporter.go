// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-init/director"
)

type FakeTaskReporter struct {
	TaskStartedStub        func(int)
	taskStartedMutex       sync.RWMutex
	taskStartedArgsForCall []struct {
		arg1 int
	}
	TaskFinishedStub        func(int, string)
	taskFinishedMutex       sync.RWMutex
	taskFinishedArgsForCall []struct {
		arg1 int
		arg2 string
	}
	TaskOutputChunkStub        func(int, []byte)
	taskOutputChunkMutex       sync.RWMutex
	taskOutputChunkArgsForCall []struct {
		arg1 int
		arg2 []byte
	}
}

func (fake *FakeTaskReporter) TaskStarted(arg1 int) {
	fake.taskStartedMutex.Lock()
	fake.taskStartedArgsForCall = append(fake.taskStartedArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.taskStartedMutex.Unlock()
	if fake.TaskStartedStub != nil {
		fake.TaskStartedStub(arg1)
	}
}

func (fake *FakeTaskReporter) TaskStartedCallCount() int {
	fake.taskStartedMutex.RLock()
	defer fake.taskStartedMutex.RUnlock()
	return len(fake.taskStartedArgsForCall)
}

func (fake *FakeTaskReporter) TaskStartedArgsForCall(i int) int {
	fake.taskStartedMutex.RLock()
	defer fake.taskStartedMutex.RUnlock()
	return fake.taskStartedArgsForCall[i].arg1
}

func (fake *FakeTaskReporter) TaskFinished(arg1 int, arg2 string) {
	fake.taskFinishedMutex.Lock()
	fake.taskFinishedArgsForCall = append(fake.taskFinishedArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	fake.taskFinishedMutex.Unlock()
	if fake.TaskFinishedStub != nil {
		fake.TaskFinishedStub(arg1, arg2)
	}
}

func (fake *FakeTaskReporter) TaskFinishedCallCount() int {
	fake.taskFinishedMutex.RLock()
	defer fake.taskFinishedMutex.RUnlock()
	return len(fake.taskFinishedArgsForCall)
}

func (fake *FakeTaskReporter) TaskFinishedArgsForCall(i int) (int, string) {
	fake.taskFinishedMutex.RLock()
	defer fake.taskFinishedMutex.RUnlock()
	return fake.taskFinishedArgsForCall[i].arg1, fake.taskFinishedArgsForCall[i].arg2
}

func (fake *FakeTaskReporter) TaskOutputChunk(arg1 int, arg2 []byte) {
	fake.taskOutputChunkMutex.Lock()
	fake.taskOutputChunkArgsForCall = append(fake.taskOutputChunkArgsForCall, struct {
		arg1 int
		arg2 []byte
	}{arg1, arg2})
	fake.taskOutputChunkMutex.Unlock()
	if fake.TaskOutputChunkStub != nil {
		fake.TaskOutputChunkStub(arg1, arg2)
	}
}

func (fake *FakeTaskReporter) TaskOutputChunkCallCount() int {
	fake.taskOutputChunkMutex.RLock()
	defer fake.taskOutputChunkMutex.RUnlock()
	return len(fake.taskOutputChunkArgsForCall)
}

func (fake *FakeTaskReporter) TaskOutputChunkArgsForCall(i int) (int, []byte) {
	fake.taskOutputChunkMutex.RLock()
	defer fake.taskOutputChunkMutex.RUnlock()
	return fake.taskOutputChunkArgsForCall[i].arg1, fake.taskOutputChunkArgsForCall[i].arg2
}

var _ director.TaskReporter = new(FakeTaskReporter)