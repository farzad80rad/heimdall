package loadBalancer

import (
	"heimdall/internal/config"
	"sync"
	"time"
)

type weightedHost struct {
	url          string
	enabled      bool
	currentScore int
	MaxScore     int
}

type weightedRoundRobin struct {
	index             int
	nextAvailableHost chan string
	hosts             []*weightedHost
	sync.Mutex
}

func NewWeightedRoundRobin(hosts []config.HostUnit) LoadBalancer {
	hostsInfo := make([]*weightedHost, len(hosts))
	calculatedGcd := hosts[0].Weight
	sum := 0
	for _, unit := range hosts {
		sum += unit.Weight
		calculatedGcd = gcd(unit.Weight, calculatedGcd)
	}

	for index, h := range hosts {
		w := 1
		if sum != 0 {
			w = h.Weight / calculatedGcd
		}
		hostsInfo[index] = &weightedHost{
			url:          h.Host,
			enabled:      true,
			currentScore: 0,
			MaxScore:     w,
		}
	}
	r := &weightedRoundRobin{
		nextAvailableHost: make(chan string),
		hosts:             hostsInfo,
	}
	go r.createNexAvailableHost()
	return r
}

func (r *weightedRoundRobin) Next() string {
	return <-r.nextAvailableHost
}

func (r *weightedRoundRobin) createNexAvailableHost() {
	consecutiveFailures := 0
	for {
		r.index = (r.index + 1) % len(r.hosts)
		h := r.hosts[r.index]
		if h.enabled && h.currentScore < h.MaxScore {
			consecutiveFailures = 0
			for h.currentScore < h.MaxScore {
				h.currentScore++
				r.nextAvailableHost <- r.hosts[r.index].url
			}
			h.currentScore = 0
		}

		consecutiveFailures++
		if consecutiveFailures == len(r.hosts) {
			// just to prevent cpu bursting when there is no host available
			time.Sleep(3 * time.Second)
		}
	}

}

func (r *weightedRoundRobin) SetHostStatus(host string, isActive bool) {
	r.Lock()
	defer r.Unlock()
	for _, s := range r.hosts {
		if s.url == host {
			if s.enabled != isActive {
				s.enabled = isActive
			}
			return
		}
	}
}

func (r *weightedRoundRobin) DisableHostForDuration(host string, duration time.Duration) {
	r.Lock()
	defer r.Unlock()
	for _, s := range r.hosts {
		if s.url == host {
			if s.enabled {
				s.enabled = false
				go func() {
					time.Sleep(duration)
					s.enabled = true
				}()
			}
			return
		}
	}
}

func gcd(x, y int) int {
	var t int
	for {
		t = x % y
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
