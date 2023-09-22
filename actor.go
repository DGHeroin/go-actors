package actors

import (
    "sync"
    "sync/atomic"
)

type ActorHandler interface {
    HandleMessage(msg interface{})
}

type Actor struct {
    Id          int64
    ActorSystem *ActorSystem
    handler     ActorHandler
    mutex       sync.Mutex
    mailbox     []interface{}
    state       int32
    isRelease   bool
}

func (a *Actor) SendMessage(msg interface{}) {
    if a.isRelease {
        return
    }
    a.mutex.Lock()
    a.mailbox = append(a.mailbox, msg)
    a.mutex.Unlock()

    a.ActorSystem.addSchedule(a)
}

func (a *Actor) wakeup() {
    if atomic.CompareAndSwapInt32(&a.state, 0, 1) {
        a.processMessage()
        atomic.StoreInt32(&a.state, 0)
    }
}
func (a *Actor) processMessage() {
    for {
        if a.MailCount() == 0 {
            return
        }
        a.mutex.Lock()
        mails := a.mailbox
        a.mailbox = []interface{}{}
        a.mutex.Unlock()

        a.dispatchMail(mails)
    }
}
func (a *Actor) dispatchMail(mails []interface{}) {
    defer func() {
        recover()
    }()

    for _, mail := range mails {
        a.handler.HandleMessage(mail)
    }
}
func (a *Actor) MailCount() int {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    return len(a.mailbox)
}

func (a *Actor) Release() {
    a.isRelease = true
    a.wakeup()
    a.ActorSystem.removeActor(a)
}
