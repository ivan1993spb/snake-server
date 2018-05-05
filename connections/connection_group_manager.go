package connections

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type ConnectionGroupManager struct {
	logger      *logrus.Logger
	groups      map[int]*ConnectionGroup
	groupsMutex *sync.RWMutex
	groupLimit  int
}

// TODO: Is it necessary to pass logger in group manager ?
func NewConnectionGroupManager(logger *logrus.Logger, groupLimit int) (*ConnectionGroupManager, error) {
	if groupLimit > 0 {
		return &ConnectionGroupManager{
			logger:      logger,
			groups:      map[int]*ConnectionGroup{},
			groupsMutex: &sync.RWMutex{},
			groupLimit:  groupLimit,
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
)

func (m *ConnectionGroupManager) Add(group *ConnectionGroup) (int, error) {
	m.groupsMutex.Lock()
	defer m.groupsMutex.Unlock()

	if m.unsafeIsFull() {
		return 0, ErrGroupLimitReached
	}

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

func (m *ConnectionGroupManager) Delete(id int) error {
	m.groupsMutex.Lock()
	defer m.groupsMutex.Unlock()

	if group, ok := m.groups[id]; ok {
		if !group.IsEmpty() {
			return ErrDeleteNotEmptyGroup
		}

		delete(m.groups, id)
		return nil
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
