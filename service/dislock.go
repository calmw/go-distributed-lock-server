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

var existLock *sync.Mutex

var locks map[string]*DisLock // LockName=>*DisLock

func InitLock() {
	existLock = new(sync.Mutex)
	locks = map[string]*DisLock{}
}

func createLockIfNotExist(lockName string) *DisLock {
	existLock.Lock()
	defer existLock.Unlock()
	lock, ok := locks[lockName]
	if !ok {
		locks[lockName] = &DisLock{
			LockName: lockName,
			mu:       new(sync.Mutex),
		}
		lock = locks[lockName]
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

// ForceLock 强制加锁，避免比如某个客户端退出未释放锁，导致其他客户端拿不到锁
func ForceLock(lockName, clientId string) (bool, string) {
	lock := createLockIfNotExist(lockName)
	lock.mu.Lock()
	defer lock.mu.Unlock()
	lock.Status = true
	lock.LockClientId = clientId

	return true, "ok"
}

// ForceUnLock 强制释放锁，避免比如某个客户端退出未释放锁，导致其他客户端拿不到锁
func ForceUnLock(lockName string) (bool, string) {
	lock, ok := locks[lockName]
	if !ok {
		locks[lockName] = &DisLock{
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
