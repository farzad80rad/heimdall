package loadBalancer

import "time"

type LoadBalancer interface {
	Next() string
	DisableHost(host string, duration time.Duration)
}
