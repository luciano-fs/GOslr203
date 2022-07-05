/*  Made channels generic
    Created message for exchanging strings */

package main

import(
    "fmt"
    "sync"
)

var wg sync.WaitGroup

type Node struct {
    Id    int
    In    <-chan interface{}
    Links map[int](chan<- interface{})
}

func NewNode(id int, in <-chan interface{}, links map[int](chan<- interface{})) (p *Node) {
    if links != nil {
        p = &Node{id, in, links}
    } else {
        p = &Node{id, in, make(map[int](chan<- interface{}))}
    }
    return p
}

func (n Node) Start() {
    m := <- n.In
    switch m := m.(type) {
    case Link:
        n.Links[m.GetPeer()] = m.GetCh()
        fmt.Println("Received a link for", m.GetPeer())
        fmt.Printf("Known adresses:")
        for peer, _ := range n.Links {
            fmt.Printf(" %d", peer)
        }
        fmt.Printf("\n")
    case string:
        fmt.Println("Received the following string", m)
    default:
        fmt.Printf("Received unknown message of type %T\n", m)
    }
    wg.Done()
}

type Link struct {
    Peer   int
    Ch     chan<- interface{}
}

func (l Link) GetPeer() int {
    return l.Peer
}

func (l Link) GetCh() (chan<- interface{}) {
    return l.Ch
}

func main() {
    wg.Add(2)

    in1 := make(chan interface{})
    in2 := make(chan interface{})

    a1 := NewNode(1, in1, nil)
    a2 := NewNode(2, in2, nil)

    go a1.Start()
    go a2.Start()

    in1 <- "Hello!"
    in2 <- Link{1, in1}

    wg.Wait()
}

