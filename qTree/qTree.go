package qTree

import (
	"fmt"
	"log"
)

type Particle struct {
	Mass, X, Y, Vx, Vy float64
}

type QTNode struct {
	P                                        *Particle
	Which_child, Size                        int
	X_com, Y_com, Total_mass, Lb, Rb, Ub, Db float64
	Parent                                   *QTNode
	Child                                    []*QTNode
}

func Min_float(a float64, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max_float(a float64, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min_int(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max_int(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func root_boundary(pArray []Particle, nParticles int) []float64 {
	ans := make([]float64, 4)
	ans[0] = pArray[0].X
	ans[1] = pArray[0].X
	ans[2] = pArray[0].Y
	ans[3] = pArray[0].Y
	for i := 1; i < nParticles; i++ {
		ans[0] = Min_float(ans[0], pArray[i].X)
		ans[1] = Max_float(ans[1], pArray[i].X)
		ans[2] = Min_float(ans[2], pArray[i].Y)
		ans[3] = Max_float(ans[3], pArray[i].Y)
	}
	xdiff := ans[1] - ans[0]
	ydiff := ans[3] - ans[2]
	if xdiff > ydiff {
		ans[3] = ans[2] + xdiff
	} else {
		ans[1] = ans[0] + ydiff
	}
	return ans
}

func Create_node(parent *QTNode, child_index int, pArray []Particle, nParticles int) *QTNode {
	new_node := &QTNode{Parent: parent, Which_child: child_index, Child: make([]*QTNode, 4), P: nil, Size: 0, X_com: 0.0, Y_com: 0.0, Total_mass: 0.0, Lb: 0.0, Rb: 0.0, Db: 0.0, Ub: 0.0}
	for i := 0; i < 4; i++ {
		new_node.Child[i] = nil
	}
	if parent == nil {
		boundary := root_boundary(pArray, nParticles)
		new_node.Lb = boundary[0]
		new_node.Rb = boundary[1]
		new_node.Db = boundary[2]
		new_node.Ub = boundary[3]
	} else {
		vb := (parent.Lb + parent.Rb) / 2.0
		hb := (parent.Db + parent.Ub) / 2.0
		switch child_index {
		case 0:
			{ // north west region
				new_node.Lb = parent.Lb
				new_node.Rb = vb
				new_node.Db = hb
				new_node.Ub = parent.Ub
			}
		case 1:
			{ // north east region
				new_node.Lb = vb
				new_node.Rb = parent.Rb
				new_node.Db = hb
				new_node.Ub = parent.Ub
			}
		case 2:
			{ // south west region
				new_node.Lb = parent.Lb
				new_node.Rb = vb
				new_node.Db = parent.Db
				new_node.Ub = hb
			}
		case 3:
			{ // south east region
				new_node.Lb = vb
				new_node.Rb = parent.Rb
				new_node.Db = parent.Db
				new_node.Ub = hb
			}
		}
	}
	return new_node
}

func Which_child_contains(n *QTNode, p *Particle) int {
	for i := 0; i < 4; i++ {
		if p.X >= n.Child[i].Lb && p.X <= n.Child[i].Rb && p.Y >= n.Child[i].Db && p.Y <= n.Child[i].Ub {
			return i
		}
	}
	fmt.Printf("%+v\n\n%+v\n", n, p)
	log.Panic("which child failed")
	return -1
}

func QTree_insert(p *Particle, root *QTNode) {
	if root.Size == 0 { // empty valid node found, insert here
		root.P = p
	} else {
		if root.Size == 1 { // only one particle in node
			for i := 0; i < 4; i++ {
				root.Child[i] = Create_node(root, i, nil, -1)
			}
			QTree_insert(root.P, root.Child[Which_child_contains(root, root.P)])
			QTree_insert(p, root.Child[Which_child_contains(root, p)])
		} else {
			QTree_insert(p, root.Child[Which_child_contains(root, p)])
		}
	}
	root.Size += 1
	root.X_com = ((root.X_com * root.Total_mass) + p.Mass*p.X) / (root.Total_mass + p.Mass)
	root.Y_com = ((root.Y_com * root.Total_mass) + p.Mass*p.Y) / (root.Total_mass + p.Mass)
	root.Total_mass += p.Mass
}

func Remove_empty_nodes(root *QTNode) *QTNode {
	if root == nil {
		return nil
	}
	if root.Size == 0 {
		return nil
	}
	for i := 0; i < 4; i++ {
		root.Child[i] = Remove_empty_nodes(root.Child[i])
	}
	return root
}

func QTree_print(n *QTNode, lvl int) {
	if n == nil {
		fmt.Println("sx")
		return
	}
	fmt.Printf("lvl = %d\n", lvl)
	if n.P != nil {
		fmt.Printf("p.x = %f, p.y = %f, p.mass = %f\n", n.P.X, n.P.Y, n.P.Mass)
	}
	fmt.Printf("which_child = %d, size = %d\nxcom = %f, ycom = %f, tmass = %f\nlb = %f, rb = %f, db = %f, ub = %f\n\n",
		n.Which_child, n.Size,
		n.X_com, n.Y_com, n.Total_mass,
		n.Lb, n.Rb, n.Db, n.Ub)
}

func Preorder_traversal(root *QTNode, lvl int) {
	if root == nil {
		return
	}
	QTree_print(root, lvl)
	for i := 0; i < 4; i++ {
		Preorder_traversal(root.Child[i], lvl+1)
	}
}
