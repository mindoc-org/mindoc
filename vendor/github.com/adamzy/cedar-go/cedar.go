package cedar

const ValueLimit = int(^uint(0) >> 1)

type node struct {
	Value int
	Check int
}

func (n *node) base() int { return -(n.Value + 1) }

type ninfo struct {
	Sibling, Child byte
}

type block struct {
	Prev, Next, Num, Reject, Trial, Ehead int
}

func (b *block) init() {
	b.Num = 256
	b.Reject = 257
}

type Cedar struct {
	*cedar
}

type cedar struct {
	Array    []node
	Ninfos   []ninfo
	Blocks   []block
	Reject   [257]int
	BheadF   int
	BheadC   int
	BheadO   int
	Capacity int
	Size     int
	Ordered  bool
	MaxTrial int
}

func New() *Cedar {
	da := cedar{
		Array:    make([]node, 256),
		Ninfos:   make([]ninfo, 256),
		Blocks:   make([]block, 1),
		Capacity: 256,
		Size:     256,
		Ordered:  true,
		MaxTrial: 1,
	}

	da.Array[0] = node{-2, 0}
	for i := 1; i < 256; i++ {
		da.Array[i] = node{-(i - 1), -(i + 1)}
	}
	da.Array[1].Value = -255
	da.Array[255].Check = -1

	da.Blocks[0].Ehead = 1
	da.Blocks[0].init()

	for i := 0; i <= 256; i++ {
		da.Reject[i] = i + 1
	}

	return &Cedar{&da}
}

// Get value by key, insert the key if not exist
func (da *cedar) get(key []byte, from, pos int) *int {
	for ; pos < len(key); pos++ {
		if value := da.Array[from].Value; value >= 0 && value != ValueLimit {
			to := da.follow(from, 0)
			da.Array[to].Value = value
		}
		from = da.follow(from, key[pos])
	}
	to := from
	if da.Array[from].Value < 0 {
		to = da.follow(from, 0)
	}
	return &da.Array[to].Value
}

func (da *cedar) follow(from int, label byte) int {
	base := da.Array[from].base()
	to := base ^ int(label)
	if base < 0 || da.Array[to].Check < 0 {
		hasChild := false
		if base >= 0 {
			hasChild = (da.Array[base^int(da.Ninfos[from].Child)].Check == from)
		}
		to = da.popEnode(base, label, from)
		da.pushSibling(from, to^int(label), label, hasChild)
	} else if da.Array[to].Check != from {
		to = da.resolve(from, base, label)
	} else if da.Array[to].Check == from {
	} else {
		panic("cedar: internal error, should not be here")
	}
	return to
}

func (da *cedar) popBlock(bi int, head_in *int, last bool) {
	if last {
		*head_in = 0
	} else {
		b := &da.Blocks[bi]
		da.Blocks[b.Prev].Next = b.Next
		da.Blocks[b.Next].Prev = b.Prev
		if bi == *head_in {
			*head_in = b.Next
		}
	}
}

func (da *cedar) pushBlock(bi int, head_out *int, empty bool) {
	b := &da.Blocks[bi]
	if empty {
		*head_out, b.Prev, b.Next = bi, bi, bi
	} else {
		tail_out := &da.Blocks[*head_out].Prev
		b.Prev = *tail_out
		b.Next = *head_out
		*head_out, *tail_out, da.Blocks[*tail_out].Next = bi, bi, bi
	}
}

func (da *cedar) addBlock() int {
	if da.Size == da.Capacity {
		da.Capacity *= 2

		oldArray := da.Array
		da.Array = make([]node, da.Capacity)
		copy(da.Array, oldArray)

		oldNinfo := da.Ninfos
		da.Ninfos = make([]ninfo, da.Capacity)
		copy(da.Ninfos, oldNinfo)

		oldBlock := da.Blocks
		da.Blocks = make([]block, da.Capacity>>8)
		copy(da.Blocks, oldBlock)
	}

	da.Blocks[da.Size>>8].init()
	da.Blocks[da.Size>>8].Ehead = da.Size

	da.Array[da.Size] = node{-(da.Size + 255), -(da.Size + 1)}
	for i := da.Size + 1; i < da.Size+255; i++ {
		da.Array[i] = node{-(i - 1), -(i + 1)}
	}
	da.Array[da.Size+255] = node{-(da.Size + 254), -da.Size}

	da.pushBlock(da.Size>>8, &da.BheadO, da.BheadO == 0)
	da.Size += 256
	return da.Size>>8 - 1
}

func (da *cedar) transferBlock(bi int, head_in, head_out *int) {
	da.popBlock(bi, head_in, bi == da.Blocks[bi].Next)
	da.pushBlock(bi, head_out, *head_out == 0 && da.Blocks[bi].Num != 0)
}

func (da *cedar) popEnode(base int, label byte, from int) int {
	e := base ^ int(label)
	if base < 0 {
		e = da.findPlace()
	}
	bi := e >> 8
	n := &da.Array[e]
	b := &da.Blocks[bi]
	b.Num--
	if b.Num == 0 {
		if bi != 0 {
			da.transferBlock(bi, &da.BheadC, &da.BheadF)
		}
	} else {
		da.Array[-n.Value].Check = n.Check
		da.Array[-n.Check].Value = n.Value
		if e == b.Ehead {
			b.Ehead = -n.Check
		}
		if bi != 0 && b.Num == 1 && b.Trial != da.MaxTrial {
			da.transferBlock(bi, &da.BheadO, &da.BheadC)
		}
	}
	n.Value = ValueLimit
	n.Check = from
	if base < 0 {
		da.Array[from].Value = -(e ^ int(label)) - 1
	}
	return e
}

func (da *cedar) pushEnode(e int) {
	bi := e >> 8
	b := &da.Blocks[bi]
	b.Num++
	if b.Num == 1 {
		b.Ehead = e
		da.Array[e] = node{-e, -e}
		if bi != 0 {
			da.transferBlock(bi, &da.BheadF, &da.BheadC)
		}
	} else {
		prev := b.Ehead
		next := -da.Array[prev].Check
		da.Array[e] = node{-prev, -next}
		da.Array[prev].Check = -e
		da.Array[next].Value = -e
		if b.Num == 2 || b.Trial == da.MaxTrial {
			if bi != 0 {
				da.transferBlock(bi, &da.BheadC, &da.BheadO)
			}
		}
		b.Trial = 0
	}
	if b.Reject < da.Reject[b.Num] {
		b.Reject = da.Reject[b.Num]
	}
	da.Ninfos[e] = ninfo{}
}

// hasChild: wherether the `from` node has children
func (da *cedar) pushSibling(from, base int, label byte, hasChild bool) {
	c := &da.Ninfos[from].Child
	keepOrder := *c == 0
	if da.Ordered {
		keepOrder = label > *c
	}
	if hasChild && keepOrder {
		c = &da.Ninfos[base^int(*c)].Sibling
		for da.Ordered && *c != 0 && *c < label {
			c = &da.Ninfos[base^int(*c)].Sibling
		}
	}
	da.Ninfos[base^int(label)].Sibling = *c
	*c = label
}

func (da *cedar) popSibling(from, base int, label byte) {
	c := &da.Ninfos[from].Child
	for *c != label {
		c = &da.Ninfos[base^int(*c)].Sibling
	}
	*c = da.Ninfos[base^int(*c)].Sibling
}

func (da *cedar) consult(base_n, base_p int, c_n, c_p byte) bool {
	c_n = da.Ninfos[base_n^int(c_n)].Sibling
	c_p = da.Ninfos[base_p^int(c_p)].Sibling
	for c_n != 0 && c_p != 0 {
		c_n = da.Ninfos[base_n^int(c_n)].Sibling
		c_p = da.Ninfos[base_p^int(c_p)].Sibling
	}
	return c_p != 0
}

func (da *cedar) setChild(base int, c byte, label byte, flag bool) []byte {
	child := make([]byte, 0, 257)
	if c == 0 {
		child = append(child, c)
		c = da.Ninfos[base^int(c)].Sibling
	}
	if da.Ordered {
		for c != 0 && c <= label {
			child = append(child, c)
			c = da.Ninfos[base^int(c)].Sibling
		}
	}
	if flag {
		child = append(child, label)
	}
	for c != 0 {
		child = append(child, c)
		c = da.Ninfos[base^int(c)].Sibling
	}
	return child
}

func (da *cedar) findPlace() int {
	if da.BheadC != 0 {
		return da.Blocks[da.BheadC].Ehead
	}
	if da.BheadO != 0 {
		return da.Blocks[da.BheadO].Ehead
	}
	return da.addBlock() << 8
}

func (da *cedar) findPlaces(child []byte) int {
	bi := da.BheadO
	if bi != 0 {
		bz := da.Blocks[da.BheadO].Prev
		nc := len(child)
		for {
			b := &da.Blocks[bi]
			if b.Num >= nc && nc < b.Reject {
				for e := b.Ehead; ; {
					base := e ^ int(child[0])
					for i := 0; da.Array[base^int(child[i])].Check < 0; i++ {
						if i == len(child)-1 {
							b.Ehead = e
							return e
						}
					}
					e = -da.Array[e].Check
					if e == b.Ehead {
						break
					}
				}
			}
			b.Reject = nc
			if b.Reject < da.Reject[b.Num] {
				da.Reject[b.Num] = b.Reject
			}
			bi_ := b.Next
			b.Trial++
			if b.Trial == da.MaxTrial {
				da.transferBlock(bi, &da.BheadO, &da.BheadC)
			}
			if bi == bz {
				break
			}
			bi = bi_
		}
	}
	return da.addBlock() << 8
}

func (da *cedar) resolve(from_n, base_n int, label_n byte) int {
	to_pn := base_n ^ int(label_n)
	from_p := da.Array[to_pn].Check
	base_p := da.Array[from_p].base()

	flag := da.consult(base_n, base_p, da.Ninfos[from_n].Child, da.Ninfos[from_p].Child)
	var children []byte
	if flag {
		children = da.setChild(base_n, da.Ninfos[from_n].Child, label_n, true)
	} else {
		children = da.setChild(base_p, da.Ninfos[from_p].Child, 255, false)
	}
	var base int
	if len(children) == 1 {
		base = da.findPlace()
	} else {
		base = da.findPlaces(children)
	}
	base ^= int(children[0])
	var from int
	var base_ int
	if flag {
		from = from_n
		base_ = base_n
	} else {
		from = from_p
		base_ = base_p
	}
	if flag && children[0] == label_n {
		da.Ninfos[from].Child = label_n
	}
	da.Array[from].Value = -base - 1
	for i := 0; i < len(children); i++ {
		to := da.popEnode(base, children[i], from)
		to_ := base_ ^ int(children[i])
		if i == len(children)-1 {
			da.Ninfos[to].Sibling = 0
		} else {
			da.Ninfos[to].Sibling = children[i+1]
		}
		if flag && to_ == to_pn { // new node has no child
			continue
		}
		n := &da.Array[to]
		n_ := &da.Array[to_]
		n.Value = n_.Value
		if n.Value < 0 && children[i] != 0 {
			// this node has children, fix their check
			c := da.Ninfos[to_].Child
			da.Ninfos[to].Child = c
			da.Array[n.base()^int(c)].Check = to
			c = da.Ninfos[n.base()^int(c)].Sibling
			for c != 0 {
				da.Array[n.base()^int(c)].Check = to
				c = da.Ninfos[n.base()^int(c)].Sibling
			}
		}
		if !flag && to_ == from_n { // parent node moved
			from_n = to
		}
		if !flag && to_ == to_pn {
			da.pushSibling(from_n, to_pn^int(label_n), label_n, true)
			da.Ninfos[to_].Child = 0
			n_.Value = ValueLimit
			n_.Check = from_n
		} else {
			da.pushEnode(to_)
		}
	}
	if flag {
		return base ^ int(label_n)
	}
	return to_pn
}
