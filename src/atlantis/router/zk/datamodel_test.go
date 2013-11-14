package zk

import (
	"atlantis/router/config"
	"testing"
)

func TestToZkPoolToPool(t *testing.T) {
	conf := config.PoolConfig{
		HealthzEvery:   "1m",
		HealthzTimeout: "9s",
		RequestTimeout: "1s",
		Status:         "OK",
	}

	host1 := config.Host{
		Address: "localhost:8081",
	}
	host2 := config.Host{
		Address: "localhost:8082",
	}
	hosts := map[string]config.Host{
		"host1": host1,
		"host2": host2,
	}

	pool := config.Pool{
		Name:     "test",
		Internal: false,
		Hosts:    hosts,
		Config:   conf,
	}

	zkPool, hosts := ToZkPool(pool)
	recon := zkPool.Pool(hosts)

	if zkPool.Name != "test" || recon.Name != "test" {
		t.Errorf("should preserve name")
	}

	if zkPool.Internal != false || recon.Internal != false {
		t.Errorf("should preserve internal")
	}

	if hosts["host1"].Address != "localhost:8081" || recon.Hosts["host2"].Address != "localhost:8082" {
		t.Errorf("should transform hosts to and fro")
	}

	if zkPool.Config.HealthzEvery != "1m" || recon.Config.RequestTimeout != "1s" {
		t.Errorf("should preserve config")
	}

}
