package main

import (
	"log"
	"time"
)

const LockClientId = "lock-client-a"

func main() {
	InitLock()
	TaskLock()
	defer TaskUnLock()
	time.Sleep(time.Second * 10) // 模拟操作
}

func TaskLock() {
	log.Printf("%s,加锁... \n", LockClientId)
	for i := 1; i <= 3; i++ {
		res := DisLock(i, LockClientId)
		if res == 1 {
			break
		} else if res == 2 {
			_, _ = ForceUnLock()
			log.Printf("%s,强制释放锁,%d \n", LockClientId, i)
			continue
		}
	}

	log.Printf("%s,执行加锁后的业务... \n", LockClientId)
}

func TaskUnLock() {
	log.Printf("%s,释放锁... \n", LockClientId)
	for i := 1; i <= 3; i++ {
		res := DisUnLock(i, LockClientId)
		if res == 1 {
			break
		} else if res == 2 {
			_, _ = ForceUnLock()
			log.Printf("%s,强制释放锁,%d \n", LockClientId, i)
			continue
		}
	}

	log.Printf("%s,执行释放锁后的业务... \n", LockClientId)
}

func DisLock(times int, clientId string) int {
	var result int
	log.Printf("%s,第%d轮加锁 \n", clientId, times)
	for i := 1; i <= 50; i++ {
		if i == 50 {
			log.Printf("%s,本轮加锁失败 \n", clientId)
			result = 2 // 加锁失败
			break
		}
		ok, err := Lock(clientId)
		if ok {
			log.Printf("%s,第%d轮第%d次加锁成功 \n", clientId, times, i)
			result = 1 // 加锁成功
			break
		} else {
			log.Printf("%s,第%d轮第%d次加锁失败 %v \n", clientId, times, i, err)
		}
		time.Sleep(time.Microsecond * 300)
	}
	log.Printf("%s 第%d轮加锁 result %d \n", clientId, times, result)

	return result
}

func DisUnLock(times int, clientId string) int {
	var result int
	log.Printf("%s,第%d轮释放锁 \n", clientId, times)
	for i := 1; i <= 50; i++ {
		if i == 50 {
			log.Printf("%s,本轮释放锁失败 \n", clientId)
			result = 2 // 释放锁失败
			break
		}
		ok, err := UnLock(clientId)
		if ok {
			log.Printf("%s,第%d轮第%d次释放锁成功 \n", clientId, times, i)
			result = 1 // 释放锁成功
			break
		} else {
			log.Printf("%s,第%d轮第%d次释放锁失败 %v \n", clientId, times, i, err)
		}
		time.Sleep(time.Microsecond * 300)
	}
	log.Printf("%s 第%d轮释放锁 result %d \n", clientId, times, result)

	return result
}
