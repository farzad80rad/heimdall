package loadBalancer

import "time"

type LoadBalancer interface {
	// Next will return the next host by specified load balancing algorithm
	Next() string
	SetHostStatus(host string, isActive bool)
	DisableHostForDuration(host string, duration time.Duration)
}
