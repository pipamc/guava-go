package strategy

import (
	"time"
	"math"
)

type BackOff func(retryTime uint) time.Duration

func Incremental(init time.Duration, increment time.Duration) BackOff {
	return func(retryTime uint) time.Duration {
		return init + (increment * time.Duration(retryTime))
	}
}

func Linear(factor time.Duration) BackOff {
	return func (retryTime uint) time.Duration {
		return factor * time.Duration(retryTime)
	}
}

func Exponential(factor time.Duration, base float64) BackOff {
	return func(retryTime uint) time.Duration {
		return factor * time.Duration(math.Pow(base, float64(retryTime)))
	}
}

func BinaryExponential(factor time.Duration) BackOff {
	return Exponential(factor, 2)
}

func Fibonacci(factor time.Duration) BackOff {
	return func (retryTime uint) time.Duration {
		return factor * time.Duration(fibonacciNum(retryTime))
	}
}

func fibonacciNum(n uint) uint {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		a := 0
		b := 1
		for i := 2; i <= int(n) ; i++ {
			c := b
			b += a
			a = c
		}
		return uint(b)
	}
}
