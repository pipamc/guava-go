package strategy

import (
	"math/rand"
	"time"
	"math"
)

// First, if you don't understand jitter, you can see this article
// https://aws.amazon.com/cn/blogs/architecture/exponential-backoff-and-jitter/
// So why we need jitter? Cause if we just use backoff to control retry time, there are some problems"
// 1. all concurrent request will request in the next round, the downstream will suffer a huge traffic
// 2. if deal lock has accured, next round is still deal lock

type Jitter func(duration time.Duration) time.Duration

// 0 - n
func FullJitter(generator *rand.Rand) Jitter {
	random := createOrGetRandom(generator)
	return func(duration time.Duration) time.Duration {
		return time.Duration(random.Int63n(int64(duration)))
	}
}

// n/2 - n
func EqualJitter(generator *rand.Rand) Jitter {
	random := createOrGetRandom(generator)
	return func(duration time.Duration) time.Duration {
		return (duration / 2) + time.Duration(random.Int63n(int64(duration) / 2))
	}
}

 // min - max
func DeviationJitter(generator *rand.Rand, factor float64) Jitter {
	if factor <= 0 || factor >= 1 {
		factor = 1
	}
	random := createOrGetRandom(generator)
	return func (duration time.Duration) time.Duration {
		min_ := int64(math.Floor(float64(duration) * (1 - factor)))
		max_ := int64(math.Ceil(float64(duration) * (1 + factor)))
		return time.Duration(random.Int63n(max_ - min_) + min_)
	}
}

func createOrGetRandom(random *rand.Rand) *rand.Rand {
	if random != nil {
		return random
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}