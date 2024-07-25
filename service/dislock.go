package service

import "sync"

type DisLock struct {
	mu           *sync.Mutex
	LockName     string
	LockClientId string // 加锁的clientID，也需要改clientID释放
	Status       bool
}

var Locks = map[string]*DisLock{} // LockName=>*DisLock

func Lock(lockName string, clientId string) (bool, string) {
	lock, ok := Locks[lockName]
	if !ok {
		Locks[lockName] = &DisLock{
			LockName:     lockName,
			LockClientId: clientId,
			Status:       true,
			mu:           new(sync.Mutex),
		}
		return true, "ok"
	}
	lock.mu.Lock()
	defer lock.mu.Unlock()
	if lock.Status {
		return false, "waiting for release"
	}
	lock.Status = true
	lock.LockClientId = clientId

	return true, "ok"
}

func UnLock(lockName string, clientId string) (bool, string) {
	lock, ok := Locks[lockName]
	if !ok {
		Locks[lockName] = &DisLock{
			LockName:     lockName,
			LockClientId: clientId,
			Status:       false,
			mu:           new(sync.Mutex),
		}
		return true, "ok"
	}
	lock.mu.Lock()
	defer lock.mu.Unlock()
	if lock.LockClientId != clientId {
		return false, "waiting for other client to be released"
	}

	return true, "ok"
}
