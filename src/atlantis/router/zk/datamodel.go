package zk

import (
	"atlantis/router/config"
)

type ZkPool struct {
	Name     string
	Internal bool
	Config   config.PoolConfig
}

func ToZkPool(p config.Pool) (ZkPool, map[string]config.Host) {
	zkPool := ZkPool{
		Name:     p.Name,
		Internal: p.Internal,
		Config:   p.Config,
	}

	return zkPool, p.Hosts
}

func (z ZkPool) Pool(hosts map[string]config.Host) config.Pool {
	return config.Pool{
		Name:     z.Name,
		Internal: z.Internal,
		Hosts:    hosts,
		Config:   z.Config,
	}
}
