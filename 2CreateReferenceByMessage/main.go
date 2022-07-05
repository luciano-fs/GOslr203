/* Made channels directional
   Created message that allows user to receive a channel 
   Node treats LinkMsg
*/

package main

import(
    "fmt"
    "sync"
)

var wg sync.WaitGroup

type Node struct {
    Id    int
    In    <-chan Msg
    Links map[int](chan<- Msg)
}

func NewNode(id int, in <-chan Msg, links map[int](chan<- Msg)) (p *Node) {
    if links != nil {
        p = &Node{id, in, links}
    } else {
        p = &Node{id, in, make(map[int](chan<- Msg))}
    }
    return p
}

func (n Node) Start() {
    m := <- n.In
    switch m := m.(type) {
    case LinkMsg:
        n.Links[m.GetPeer()] = m.GetCh()
        fmt.Println("Received a link for", m.GetPeer())
    default:
        fmt.Printf("Received unknown message of type %T\n", m)
    }
    wg.Done()
}

type Msg interface {
    GetSender() int
}

type LinkMsg struct {
    Sender int
    Peer   int
    Ch     chan<- Msg
}

func (lm LinkMsg) GetSender() int {
    return lm.Sender
}

func (lm LinkMsg) GetPeer() int {
    return lm.Peer
}

func (lm LinkMsg) GetCh() (chan<- Msg) {
    return lm.Ch
}

func main() {
    wg.Add(2)

    in1 := make(chan Msg)
    in2 := make(chan Msg)

    a1 := NewNode(1, in1, nil)
    a2 := NewNode(2, in2, nil)

    go a1.Start()
    go a2.Start()

    in1 <- LinkMsg{0, 2, in2}
    in2 <- LinkMsg{0, 1, in1}

    wg.Wait()
}
