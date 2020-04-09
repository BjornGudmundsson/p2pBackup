package peers

import "sync"

type Container interface {
	GetPeerList() []Peer
	Update(p *Peer)
}


type PeerContainer struct {
	mutex sync.Mutex
	container map[string]Peer
	peerFile string
}

func (pc *PeerContainer) GetPeerList() []Peer {
	peers := make([]Peer, 0)
	for _, v := range pc.container {
		peers = append(peers, v)
	}
	return peers
}

func (pc *PeerContainer) Update(p *Peer) {
	//TODO: Do something
}