package swarm

import (
	"github.com/cryptopunkscc/lore/comm/client"
)

type Swarm struct {
	clients map[string]*client.Client
}

func NewSwarm() *Swarm {
	swarm := &Swarm{}
	swarm.clients = make(map[string]*client.Client)
	return swarm
}

func (swarm *Swarm) Add(address string) error {
	swarm.clients[address] = client.NewClient(address)
	return nil
}

func (swarm *Swarm) Remove(address string) {
	delete(swarm.clients, address)
}

func (swarm *Swarm) List() []string {
	res := make([]string, 0)
	for i := range swarm.clients {
		res = append(res, i)
	}
	return res
}

func (swarm *Swarm) FindSource(id string) *client.Client {
	for _, client := range swarm.clients {
		has, _ := client.Item().Info(id)
		if has {
			return client
		}
	}
	return nil
}
