package config

import "time"

type MatchPolicy struct {
	Url       string
	HttpTypes []string
}

type LoadBalanceType int

const (
	LoadBalanceType_ROUNDROBIN LoadBalanceType = iota
	LoadBalanceType_WEIGHTED_ROUNDROBIN
)

type HealthCheckConfig struct {
	Path              string
	FailureThreshHold int
	Interval          time.Duration
}

type HostUnit struct {
	Host   string
	Weight int
}

type CircuitBreakerConfig struct {
	ExamineWindow         time.Duration
	QuarantineDuration    time.Duration
	FierierToleranceCount uint32
}

type HostLoadPolicy struct {
	LoadBalanceType LoadBalanceType
	HostUnits       []HostUnit
}

type ApiConfig struct {
	Match                MatchPolicy
	HostInfo             HostLoadPolicy
	HealthCheckConfig    *HealthCheckConfig
	CircuitBreakerConfig CircuitBreakerConfig
}
