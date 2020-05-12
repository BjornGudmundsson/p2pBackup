package peers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func TestGetCommitMessages(t *testing.T) {
	user := os.Getenv("user")
	pw := os.Getenv("pw")
	repo := os.Getenv("repo")
	email := os.Getenv("email")
	if pw == "" || user == "" {
		t.Log("Have to give environment variables to access the git service of choice")
		return
	}
	r, e := cloneRepo(repo, user, pw)
	assert.Nil(t, e, "Should be able to clone the repo")
	tree, e := r.Worktree()
	assert.Nil(t, e, "Should be able to get worktree")
	msg := "Testing Golang commit"
	fmt.Println(user, pw, repo)
	e = commitMessage(msg, user, email, pw, tree)
	assert.Nil(t, e, "Should be able to commit the message")
}

func TestPushMessageParallel(t *testing.T) {
	user := os.Getenv("user")
	pw := os.Getenv("pw")
	repo := os.Getenv("repo")
	email := os.Getenv("email")
	if pw == "" || user == "" {
		t.Log("Have to give environment variables to access the git service of choice")
		return
	}
	msg := "testing push function"
	e := PushMessageParallel(repo, user, pw, email, msg)
	assert.Nil(t, e, "Should be able to push the message")
	msgs, e := GetCommitMessages(repo, user, pw, 2 * time.Minute)
	assert.Nil(t, e, "Should be able to get commit messages")
	passed := false
	for _, m := range msgs {
		if msg == m {
			passed = true
			break
		}
	}
	assert.True(t, passed, "Commit message should be on the repo in the last 2 minutes")
}

func TestUpdateContainer(t *testing.T) {
	user := os.Getenv("user")
	pw := os.Getenv("pw")
	repo := os.Getenv("repo")
	email := os.Getenv("email")
	if pw == "" || user == ""{
		t.Log("Have to give environment variables to access the git service of choice")
		return
	}
	time.Sleep(2 * time.Minute)
	ip1, port1 := "127.0.0.1", "3000"
	ip2, port2 := "127.0.0.1", "3001"
	m1 := ip1 + seperator + port1
	m2 := ip2 + seperator + port2
	p1, e := NewTCPPeer(m1)
	assert.Nil(t, e, "can't make peer")
	p2, e := NewTCPPeer(m2)
	assert.Nil(t, e, "can't make peer")
	container := &PeerContainer{
		mutex:     sync.Mutex{},
		container: make(map[string]Peer),
		peerFile:  "peers.txt",
	}
	container.New([]Peer{p1, p2})
	oldContainer := container.GetPeerList()
	container.Storage()
	peersFromFile, e := GetPeerList("peers.txt")
	assert.Nil(t, e, "Should be able to get the peer list")
	assert.True(t, arraysAreEqual(peersFromFile, oldContainer), "")
	e = PushMessageParallel(repo, user, pw, email, m1)
	assert.Nil(t, e, "Should be able to push")
	e = PushMessageParallel(repo, user, pw, email, m2)
	assert.Nil(t, e, "should be able to push")
	msgs, e := GetCommitMessages(repo, user, pw, 2 * time.Minute)
	assert.Nil(t, e, "Should be able to get the logs")
	peers, e := PeersFromStrings(msgs)
	assert.True(t, arraysAreEqual(peers, oldContainer), "Should equal the old container")
	assert.Nil(t, e, "Should be able to get peer list")
	assert.True(t, isInArray(oldContainer, p1), "peer 1 should be in the old container")
	assert.True(t, isInArray(oldContainer, p2), "Peer 1 should be in the old container")
	time.Sleep(2 * time.Minute)//Sleep to make an "epoch" pass
	e = PushMessageParallel(repo, user, pw, email, m1)
	assert.Nil(t, e, "Should be able to push")
	m3 := ip1 + seperator + "3003"
	e = PushMessageParallel(repo, user, pw, email, m3)
	assert.Nil(t, e, "Should be able to push")
	nextMessages, e := GetCommitMessages(repo, user, pw, 2 * time.Minute)
	assert.Nil(t, e, "should be able to get the logs")
	newPeers, e := PeersFromStrings(nextMessages)
	assert.Nil(t, e, "Should be able to get peer list")
	container.New(newPeers)
	newContainer := container.GetPeerList()
	assert.False(t, arraysAreEqual(oldContainer, newContainer), "New and old container should not be equal")
	assert.True(t, len(oldContainer) == len(newContainer), "Containers should be of the same length")
	container.Storage()
	peersFromFile, e = GetPeerList("peers.txt")
	assert.Nil(t, e, "Should be able to get the peer list")
	assert.True(t, arraysAreEqual(newContainer, peersFromFile), "Peers should equal after being written to file")
}


func PeersFromStrings(msgs []string) ([]Peer, error) {
	peers := make([]Peer, 0)
	for _, msg := range msgs {
		p, e := NewTCPPeer(msg)
		fmt.Println(msg)
		if e != nil {
			fmt.Println(e)
		} else {
			peers = append(peers, p)
		}
	}
	return peers, nil
}

func arraysAreEqual(p1, p2 []Peer) bool {
	if len(p1) != len(p2) {
		return false
	}
	for _, p := range p1 {
		if !isInArray(p2, p) {
			return false
		}
	}
	return true
}

func isInArray(peers []Peer, p Peer) bool {
	if p == nil || peers == nil {
		return false
	}
	for _, peer := range peers {
		if peer.String() == p.String() {
			return true
		}
	}
	return false
}
