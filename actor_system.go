package actors

import (
    "sync"
    "sync/atomic"
)

type ActorSystem struct {
    id     int64
    mutex  sync.RWMutex
    actors map[int64]*Actor
    ch     chan *Actor
    wg     sync.WaitGroup
}
type Option func(*option)
type option struct {
    workerNum int
}

func WithWorkers(n int) Option {
    return func(o *option) {
        if n < 0 {
            n = 100
        }
        o.workerNum = n
    }
}
func NewSystem(opts ...Option) *ActorSystem {
    o := &option{}
    for _, f := range opts {
        f(o)
    }

    sys := &ActorSystem{
        id:     1,
        actors: map[int64]*Actor{},
        ch:     make(chan *Actor),
    }

    workerFn := func(workId int) {
        defer sys.wg.Done()
        for actor := range sys.ch {
            actor.wakeup()
        }
    }
    for i := 0; i < o.workerNum; i++ {
        sys.wg.Add(1)
        go workerFn(i)
    }
    return sys
}

func (sys *ActorSystem) ManageActor(ptr interface{}) *Actor {
    handler, ok := ptr.(ActorHandler)
    if !ok {
        return nil
    }
    actor := toActor(ptr)
    if actor == nil {
        return nil
    }

    id := atomic.AddInt64(&sys.id, 1)
    if id == 0 {
        id = atomic.AddInt64(&sys.id, 1)
    }

    actor.Id = id
    actor.ActorSystem = sys
    actor.handler = handler

    sys.mutex.Lock()
    sys.actors[id] = actor
    sys.mutex.Unlock()
    return actor
}
func (sys *ActorSystem) addSchedule(actor *Actor) {
    sys.ch <- actor
}

func (sys *ActorSystem) Release() {
    sys.mutex.Lock()
    var actors []*Actor
    for _, actor := range sys.actors {
        actors = append(actors, actor)
    }
    sys.mutex.Unlock()

    for _, actor := range actors {
        actor.wakeup()
    }

    close(sys.ch)
    sys.wg.Wait()
}

func (sys *ActorSystem) removeActor(a *Actor) {
    sys.mutex.Lock()
    defer sys.mutex.Unlock()
    delete(sys.actors, a.Id)
}
