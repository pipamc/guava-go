package strategy

import (
	"time"
)

type Strategy func(retryTime uint) bool

func Infinite() Strategy {
	return func(retryTime uint) bool {
		return true
	}
}

func Delay(duration time.Duration) Strategy {
	return func(retryTime uint) bool {
		if retryTime == 0 {
			time.Sleep(duration)
		}
		return true
	}
}

func Limit(limit uint) Strategy {
	return func(retryTime uint) bool {
		return retryTime < limit
	}
}

func BackOffOnly(backoff BackOff) Strategy {
	return func(retryTime uint) bool {
		time.Sleep(backoff(retryTime))
		return true
	}
}
func BackOffWithJitter(backoff BackOff, jitter Jitter) Strategy {
	return func(retryTime uint) bool {
		time.Sleep(jitter(backoff(retryTime)))
		return true
	}
}

func Wait(durations ...time.Duration) Strategy {
	return func(retryTime uint) bool {
		if retryTime > 0 && len(durations) > 0 {
			durationIndex := int(retryTime - 1)

			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}

			time.Sleep(durations[durationIndex])
		}
		return true
	}
}
