package peers

import (
	"sync"
)

type Container interface {
	GetPeerList() []Peer//Get all currently contained peer.
	Update(p Peer)//Update the container by either removing or adding a peer.
	Storage()//Write the container to long term storage.
}


type PeerContainer struct {
	mutex sync.Mutex
	container map[string]Peer
	peerFile string
}

func (pc *PeerContainer) GetPeerList() []Peer {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	peers := make([]Peer, 0)
	for _, v := range pc.container {
		peers = append(peers, v)
	}
	return peers
}

func (pc *PeerContainer) Update(p Peer) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	k := p.String()
	if p.Available() {
		pc.container[k] = p
	} else {
		delete(pc.container, k)
	}
}

func (pc *PeerContainer) Storage() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	//TODO: Add such that the current state of the container is written to long term storage.
}

func NewContainerFromFile(fn string) (Container, error) {
	m := make(map[string]Peer)
	container := &PeerContainer{
		container: m,
		peerFile: fn,
	}
	peerList, e := GetPeerList(fn)
	if e != nil {
		return nil, e
	}
	for _, p := range peerList {
		m[p.String()] = p
	}
	return container, nil
}