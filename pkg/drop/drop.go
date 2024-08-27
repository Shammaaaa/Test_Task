package drop

type Drop interface {
	Drop() error
}

type Debug interface {
	DropMsg() string
}

type Droppable struct {
	droppers []Drop
}

func (n *Droppable) AddDroppers(droppers ...Drop) {
	for _, dropper := range droppers {
		n.AddDropper(dropper)
	}
}

func (n *Droppable) AddDropper(dropper Drop) {
	n.droppers = append(n.droppers, dropper)
}

func (n *Droppable) EachDroppers(callback func(Drop)) {
	for _, dropper := range n.droppers {
		callback(dropper)
	}
}
