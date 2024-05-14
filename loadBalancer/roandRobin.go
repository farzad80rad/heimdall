package loadBalancer

import (
	"fmt"
	"time"
)

type host struct {
	url         string
	enabled     bool
	failedCount int
}

type roundRobin struct {
	index             int
	nextAvailableHost chan string
	hosts             []*host
	hostsMap          map[string]int
}

func NewRoundRobin(hosts []string) LoadBalancer {
	hostsInfo := make([]*host, len(hosts))
	indexingMap := make(map[string]int, 5*len(hosts))
	for index, h := range hosts {
		hostsInfo[index] = &host{
			url:     h,
			enabled: true,
		}
		indexingMap[h] = index
	}
	r := &roundRobin{
		nextAvailableHost: make(chan string),
		hosts:             hostsInfo,
		hostsMap:          indexingMap,
	}
	go r.createNexAvailableHost()
	return r
}

func (r *roundRobin) DisableHost(host string, duration time.Duration) {
	index := r.hostsMap[host]
	r.hosts[index].enabled = false
	go func() {
		time.Sleep(duration)
		r.hosts[index].enabled = true
	}()
}

func (r *roundRobin) Next() string {
	return <-r.nextAvailableHost
}

func (r *roundRobin) createNexAvailableHost() {
	consecutiveFailures := 0
	for {
		for {
			r.index = (r.index + 1) % len(r.hosts)
			if r.hosts[r.index].enabled {
				consecutiveFailures = 0
				break
			}
			consecutiveFailures++
			if consecutiveFailures == len(r.hosts) {
				fmt.Println("sleep")
				// just to prevent cpu bursting when there is no available host
				time.Sleep(time.Second)
			}
		}
		fmt.Println("host is active", *r.hosts[r.index])
		r.nextAvailableHost <- r.hosts[r.index].url
	}
}
