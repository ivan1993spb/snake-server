package connections

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type ConnectionGroupManager struct {
	groups      map[int]*ConnectionGroup
	groupsMutex *sync.RWMutex
	groupLimit  int
	connsLimit  int
	connsCount  int
	logger      logrus.FieldLogger
}

func NewConnectionGroupManager(logger logrus.FieldLogger, groupLimit, connsLimit int) (*ConnectionGroupManager, error) {
	if groupLimit > 0 {
		return &ConnectionGroupManager{
			groups:      map[int]*ConnectionGroup{},
			groupsMutex: &sync.RWMutex{},
			groupLimit:  groupLimit,
			connsLimit:  connsLimit,
			logger:      logger,
		}, nil
	}

	return nil, errors.New("cannot create connection group manager: invalid group limit")
}

func (m *ConnectionGroupManager) unsafeIsFull() bool {
	return len(m.groups) == m.groupLimit
}

func (m *ConnectionGroupManager) IsFull() bool {
	m.groupsMutex.RLock()
	defer m.groupsMutex.RUnlock()
	return m.unsafeIsFull()
}

type ErrAddGroup string

func (e ErrAddGroup) Error() string {
	return "cannot add group: " + string(e)
}

var (
	ErrGroupLimitReached = ErrAddGroup("limit group count reached")
	ErrCannotGetID       = ErrAddGroup("cannot get id for group")
	ErrConnsLimitReached = ErrAddGroup("cannot reserve connections for group: connections count reached")
)

func (m *ConnectionGroupManager) Add(group *ConnectionGroup) (int, error) {
	m.groupsMutex.Lock()
	defer m.groupsMutex.Unlock()

	if m.unsafeIsFull() {
		return 0, ErrGroupLimitReached
	}

	if group.GetLimit() > m.connsLimit-m.connsCount {
		if m.connsLimit-m.connsCount == 0 {
			return 0, ErrConnsLimitReached
		}
		group.SetLimit(m.connsLimit - m.connsCount)
	}

	m.connsCount += group.GetLimit()

	for id := 0; id <= len(m.groups); id++ {
		if _, occupied := m.groups[id]; !occupied {
			m.groups[id] = group
			return id, nil
		}
	}

	return 0, ErrCannotGetID
}

type ErrDeleteGroup string

func (e ErrDeleteGroup) Error() string {
	return "cannot delete group: " + string(e)
}

var (
	ErrDeleteNotEmptyGroup = ErrDeleteGroup("group is not empty")
	ErrDeleteNotFoundGroup = ErrDeleteGroup("group not found")
)

func (m *ConnectionGroupManager) Delete(group *ConnectionGroup) error {
	if !group.IsEmpty() {
		return ErrDeleteNotEmptyGroup
	}

	m.groupsMutex.Lock()
	defer m.groupsMutex.Unlock()

	m.connsCount -= group.GetLimit()

	for id := range m.groups {
		if m.groups[id] == group {
			delete(m.groups, id)
			return nil
		}
	}

	return ErrDeleteNotFoundGroup
}

var ErrNotFoundGroup = errors.New("not found group")

func (m *ConnectionGroupManager) Get(id int) (*ConnectionGroup, error) {
	m.groupsMutex.RLock()
	defer m.groupsMutex.RUnlock()

	if group, ok := m.groups[id]; ok {
		return group, nil
	}

	return nil, ErrNotFoundGroup
}

func (m *ConnectionGroupManager) Groups() map[int]*ConnectionGroup {
	m.groupsMutex.RLock()
	defer m.groupsMutex.RUnlock()
	groups := map[int]*ConnectionGroup{}
	for id, group := range m.groups {
		groups[id] = group
	}
	return groups
}

func (m *ConnectionGroupManager) GroupLimit() int {
	return m.groupLimit
}

func (m *ConnectionGroupManager) GroupCount() int {
	m.groupsMutex.RLock()
	defer m.groupsMutex.RUnlock()
	return len(m.groups)
}

func (m *ConnectionGroupManager) Capacity() float32 {
	m.groupsMutex.RLock()
	defer m.groupsMutex.RUnlock()

	var count = 0
	for _, group := range m.groups {
		count += group.GetCount()
	}

	return float32(count) / float32(m.connsLimit)
}
