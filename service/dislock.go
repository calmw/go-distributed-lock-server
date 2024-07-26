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
	if clientId == "admin" {
		lock.Status = false
		return true, "ok"
	}
	if lock.LockClientId != clientId {
		return false, "waiting for other client to be released"
	}
	lock.Status = false

	return true, "ok"
}

// ForceUnLock 强制释放锁，避免比如某个客户端退出未释放锁，导致其他客户端拿不到锁
func ForceUnLock(lockName string) (bool, string) {
	lock, ok := Locks[lockName]
	if !ok {
		Locks[lockName] = &DisLock{
			LockName: lockName,
			Status:   false,
			mu:       new(sync.Mutex),
		}
		return true, "ok"
	}
	lock.mu.Lock()
	defer lock.mu.Unlock()

	lock.Status = false

	return true, "ok"
}
