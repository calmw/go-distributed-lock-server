package service

import (
	"fmt"
	"sync"
)

type DisLock struct {
	mu           *sync.Mutex
	LockName     string
	LockClientId string // 加锁的clientID，也需要改clientID释放
	Status       bool
}

var CheckExistLock *sync.Mutex

var Locks map[string]*DisLock // LockName=>*DisLock

func createLockIfNotExist(lockName string) *DisLock {
	CheckExistLock.Lock()
	defer CheckExistLock.Unlock()
	lock, ok := Locks[lockName]
	if !ok {
		Locks[lockName] = &DisLock{
			LockName: lockName,
			mu:       new(sync.Mutex),
		}
	}
	return lock
}

func Lock(lockName, clientId string) (bool, string) {
	lock := createLockIfNotExist(lockName)
	lock.mu.Lock()
	defer lock.mu.Unlock()
	if lock.Status {
		return false, fmt.Sprintf("lockName(%s), waiting client(%s) to be released", lock.LockName, lock.LockClientId)
	}
	lock.Status = true
	lock.LockClientId = clientId

	return true, "ok"
}

func UnLock(lockName, clientId string) (bool, string) {
	lock := createLockIfNotExist(lockName)
	lock.mu.Lock()
	defer lock.mu.Unlock()
	if lock.LockClientId != clientId {
		return false, fmt.Sprintf("lockName(%s), waiting for other client(%s) to be released", lock.LockName, lock.LockClientId)
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
		return false, fmt.Sprintf("lockName(%s), not exist", lock.LockName)
	}
	lock.mu.Lock()
	defer lock.mu.Unlock()

	lock.Status = false

	return true, "ok"
}
