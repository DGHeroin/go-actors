package main

import (
    "github.com/DGHeroin/go-actors"
    "log"
    "time"
)

func main() {
    t0 := time.Now()
    sys := actors.NewSystem(actors.WithWorkers(10))
    a1 := sys.ManageActor(&PlayerActor{})
    a2 := sys.ManageActor(&PlayerActor{})

    a1.SendMessage("111")
    a2.SendMessage("222")

    log.Println("stopping...")

    a1.SendMessage("333")
    a2.SendMessage("444")

    a1.Release()

    a1.SendMessage("555")
    a2.SendMessage("666")

    sys.Release()
    log.Println("stopped", time.Since(t0))
}

type PlayerActor struct {
    actors.Actor
}

func (p *PlayerActor) HandleMessage(msg interface{}) {
    switch v := msg.(type) {
    case string:
        log.Println("is string:", v)
        time.Sleep(time.Second)
    }
}
