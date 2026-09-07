// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package timer provides a generic, per-key, cancelable delayed-callback
// manager. It is the shared extraction of the per-alarm/per-entity timer
// pattern (map of key -> running timer/cancel-channel, guarded by a mutex)
// that recurred independently in more than one consuming project's own
// escalation scheduler.
package timer

import (
	"sync"
	"time"
)

// TimerManager tracks at most one running timer per key.
type TimerManager struct {
	active map[string]*timerState
	mtx    sync.Mutex
}

type timerState struct {
	timer  *time.Timer
	cancel chan struct{}
	onFire func()
}

// NewTimerManager creates an empty TimerManager.
func NewTimerManager() *TimerManager {
	return &TimerManager{active: make(map[string]*timerState)}
}

// Start begins a timer for key that calls onFire after duration elapses.
// Any existing timer for key is cancelled first (Start replaces, not stacks).
func (m *TimerManager) Start(key string, duration time.Duration, onFire func()) {
	cancel := make(chan struct{})
	t := time.NewTimer(duration)

	m.mtx.Lock()
	if existing, ok := m.active[key]; ok {
		close(existing.cancel)
		existing.timer.Stop()
	}
	m.active[key] = &timerState{timer: t, cancel: cancel, onFire: onFire}
	m.mtx.Unlock()

	go m.wait(key, t, cancel, onFire)
}

func (m *TimerManager) wait(key string, t *time.Timer, cancel chan struct{}, onFire func()) {
	select {
	case <-t.C:
		m.mtx.Lock()
		delete(m.active, key)
		m.mtx.Unlock()
		onFire()
	case <-cancel:
		t.Stop()
	}
}

// Cancel stops the timer for key, if one is running. No-op otherwise.
func (m *TimerManager) Cancel(key string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if state, ok := m.active[key]; ok {
		close(state.cancel)
		state.timer.Stop()
		delete(m.active, key)
	}
}

// Reset restarts the timer for key with a new duration, reusing the onFire
// callback it was started with. No-op if key has no running timer.
func (m *TimerManager) Reset(key string, duration time.Duration) {
	m.mtx.Lock()
	existing, ok := m.active[key]
	if !ok {
		m.mtx.Unlock()
		return
	}
	onFire := existing.onFire
	close(existing.cancel)
	existing.timer.Stop()
	delete(m.active, key)
	m.mtx.Unlock()

	m.Start(key, duration, onFire)
}

// Active returns the number of currently running timers.
func (m *TimerManager) Active() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return len(m.active)
}
