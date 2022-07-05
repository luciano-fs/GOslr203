/*  Made channels generic
    Created message for exchanging strings */

package main

import(
    "time"
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
    if links == nil {
        links = make(map[int](chan<- interface{}))
    }
    p = &Node{id, in, links}
    return p
}

func (n *Node) Start() {
    for true {
        if m, ok := <-n.In; ok == true { 
            switch m := m.(type) {
            case Link:
                n.Links[m.GetPeer()] = m.GetCh()
                fmt.Println(n.Id, "received a link for", m.GetPeer())
                fmt.Printf("Known adresses:")
                for peer, _ := range n.Links {
                    fmt.Printf(" %d", peer)
                }
                fmt.Printf("\n")
            case string:
                fmt.Println(n.Id, "received the following string", m)
            default:
                fmt.Printf("%d received unknown message of type %T\n",n.Id, m)
            }
		} else {
			break
		}
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

type Broadcast struct {
    Node
}

func (b *Broadcast) Repeat() {
    for true {
        if m, ok := <-b.In; ok == true {
            for _,out := range b.Links {
                out <- m
            }
		} else {
			break
		}
    }
    wg.Done()
}

func NewBroadcast(id int, in <-chan interface{}, links map[int](chan<- interface{})) *Broadcast {
    p := &Broadcast{Node: *NewNode(id, in, links)}
    return p
}

func main() {
    n := 10
    wg.Add(n+1)
    nodes := make([](*Node), n)
    channels := make([](chan interface{}), n)
    outChannels := make(map[int](chan<- interface{}))

    for cnt := 0; cnt < n; cnt++ {
        channels[cnt] = make(chan interface{})
        outChannels[cnt] = channels[cnt]
        nodes[cnt] = NewNode(cnt, channels[cnt], nil)
    }

    inTransmit := make(chan interface{})
    t := NewBroadcast(n, inTransmit, outChannels)

    for cnt := 0; cnt < n; cnt++ {
        go nodes[cnt].Start()
    }
    go t.Repeat()

    for cnt := 0; cnt < n; cnt++ {
        channels[cnt] <- Link{n,inTransmit}
    }

    inTransmit <- "Hello!"
    time.Sleep(time.Second)
    close(inTransmit)

    for cnt := 0; cnt < n; cnt++ {
        close(channels[cnt])
    }

    wg.Wait()
}
