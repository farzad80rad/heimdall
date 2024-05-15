package loadBalancer

import "time"

type LoadBalancer interface {
	Next() string
	SetHostStatus(host string, isActive bool)
	DisableHostForDuration(host string, duration time.Duration)
}
