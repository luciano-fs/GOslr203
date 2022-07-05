package main

type Node struct {
    id int
    in chan Msg
    links map[int](chan<- Msg)
}

func NewNode(id int, in chan Msg, links map[int](chan Msg)) *Node{
    p := &Node{id, in, links}
    return p
}

type Msg interface {
    GetSender() int
}

func main() {
    in1 := make(chan Msg)
    in2 := make(chan Msg)

    links2 := make(map[int](chan Msg))
    links2[1] = in1

    a1 := NewNode(1, in1, nil)
    a2 := NewNode(2, in2, links2)
}
