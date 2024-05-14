package config

import "time"

type MatchPolicy struct {
	Url       string
	HttpTypes []string
}

type CircuitBreakerConfig struct {
	ExamineWindow         time.Duration
	QuarantineDuration    time.Duration
	FierierToleranceCount uint32
}

type HostInfo struct {
	HostAddress []string
}

type ApiConfig struct {
	Match                MatchPolicy
	HostInfo             HostInfo
	CircuitBreakerConfig CircuitBreakerConfig
}
