package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
)

type Discovery interface {
	// refresh from remote registry
	Refresh() error
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}

// 目前服务发现支持轮询和随机选择
type MultiServersDiscovery struct {
	r       *rand.Rand // 随机选择需要
	mu      sync.RWMutex
	servers []string
	index   int // 轮询需要
}

func (m *MultiServersDiscovery) Refresh() error {
	return nil
}

func (m *MultiServersDiscovery) Update(servers []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.servers = servers
	return nil
}

func (m *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available server")
	}
	switch mode {
	case RandomSelect:
		return m.servers[m.r.Intn(n)], nil
	case RoundRobinSelect:
		s := m.servers[m.index%n]
		m.index = (m.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: no support select mode")
	}
}

func (m *MultiServersDiscovery) GetAll() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	servers := make([]string, len(m.servers), len(m.servers))
	copy(servers, m.servers)
	return servers, nil
}

func NewMultiServersDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		servers: servers,
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

var _ Discovery = (*MultiServersDiscovery)(nil)
