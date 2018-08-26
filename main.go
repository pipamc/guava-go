package main

import (
	"context"
	"time"
	retry2 "github.com/pipamc/guava-go/retry"
	"fmt"
	"github.com/pipamc/guava-go/retry/strategy"
)

func ping() error {
	fmt.Println("pong")
	return nil
}

func main() {
	ctx := context.Background()
	ctx2, _ := context.WithTimeout(ctx, 20 * time.Second)
	retry := retry2.NewRetry(ping, nil)
	retry.AppendStrategies(
		strategy.BackOffWithJitter(
			strategy.BinaryExponential(5 * time.Second), strategy.EqualJitter(nil)))
	retry.AppendStrategies(strategy.Limit(2))
	retry.Call(ctx2)
}
