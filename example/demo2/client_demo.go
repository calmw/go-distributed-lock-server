package main

import (
	"log"
	"time"
)

const LockClientId = "lock-client-a"

func main() {
	InitLock()
	TaskLock("some_lock")
	defer TaskUnLock("some_lock")
	time.Sleep(time.Second * 10) // 模拟操作
}

func TaskLock(lockName string) {
	t := time.Now()
	log.Printf("%s %s,加锁...", lockName, LockClientId)
	for i := 1; i <= 60; i++ { // 10分钟
		if i == 60 {
			_, _ = ForceUnLock()
			log.Printf("%s,强制释放锁,%d", LockClientId, i)
			break
		}
		res := DisLock(i, lockName, LockClientId)
		if res == 1 {
			break
		} else if res == 2 {
			continue
		}
	}

	log.Printf("%s %s %d, 执行加锁后的业务...", lockName, LockClientId, time.Since(t))
}

func TaskUnLock(lockName string) {
	t := time.Now()
	log.Printf("%s %s,释放锁...", lockName, LockClientId)
	for i := 1; i <= 60; i++ {
		if i == 60 {
			_, _ = ForceUnLock()
			log.Printf("%s,强制释放锁,%d", LockClientId, i)
			break
		}
		res := DisUnLock(i, lockName, LockClientId)
		if res == 1 {
			break
		} else if res == 2 {
			continue
		}
	}

	log.Printf("%s %s %d,执行释放锁后的业务...", lockName, LockClientId, time.Since(t))
}

func DisLock(times int, lockName, LockClientId string) int {
	var result int
	var err error
	for i := 1; i <= 50; i++ {
		if i == 50 {
			result = 2 // 加锁失败
			break
		}
		ok, e := Lock(LockClientId)
		if ok {
			result = 1 // 加锁成功
			err = nil
			break
		}
		err = e
		time.Sleep(time.Millisecond * 200)
	}
	log.Printf("%s-%s 第%d轮加锁 result %d %v", lockName, LockClientId, times, result, err)

	return result
}

func DisUnLock(times int, lockName, LockClientId string) int {
	var result int
	var err error
	for i := 1; i <= 50; i++ {
		if i == 50 {
			result = 2 // 释放锁失败
			break
		}
		ok, e := UnLock(LockClientId)
		if ok {
			err = nil
			result = 1 // 释放锁成功
			break
		}
		err = e
		time.Sleep(time.Millisecond * 200)
	}
	log.Printf("%s-%s 第%d轮释放锁 result %d %v", lockName, LockClientId, times, result, err)

	return result
}
