


// A generic red-black tree implementation from the
// functional implementation by Matt Might[1] and Okasaki.
//
// [1] http://matt.might.net/articles/red-black-delete/
//
// Author: Pratik Deoghare


package redblack

import (
	"fmt"
)

type Map[K, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
	Delete(key K)
}

func New[K, V any](less func(K, K) bool) Map[K, V] {
	leaf := &node[K, V]{
		color: B,
	}
	leaf.a = leaf
	leaf.b = leaf
	bbleaf := &node[K, V]{
		color: BB,
	}
	bbleaf.a = leaf
	bbleaf.b = leaf
	return &rbmap[K, V]{
		less:   less,
		leaf:   leaf,
		bbleaf: bbleaf,
		root:   leaf,
	}
}

type color uint8

const (
	R  color = 0
	B  color = 1
	BB color = 2
	NB color = 3
)

type node[K, V any] struct {
	color color
	key   K
	value V
	a     *node[K, V]
	b     *node[K, V]
}
type rbmap[K, V any] struct {
	root   *node[K, V]
	leaf   *node[K, V] // the leaf always Black. We don't touch it. Its a sacred leaf.
	bbleaf *node[K, V] // this is used for deletion
	less   func(K, K) bool
}

func (r rbmap[K, V]) Preorder() {
	r.preorder(r.root, "")
}
func (r rbmap[K, V]) preorder(n *node[K, V], tab string) {
	if n == r.leaf {
		return
	}
	fmt.Println(tab, n.key, "=>", n.value, n.color)
	r.preorder(n.a, ":"+tab)
	r.preorder(n.b, ":"+tab)
}
func (r rbmap[K, V]) Inorder() {
	panic("implement me")
}
func (r rbmap[K, V]) Get(key K) (value V, ok bool) {
	n := r.root
	for n != r.leaf {
		if r.less(key, n.key) {
			n = n.a
		} else if r.less(n.key, key) {
			n = n.b
		} else {
			return n.value, true
		}
	}
	return Nil[V](), false
}
func Nil[T any]() T {
	var zero T
	return zero
}
func (r *rbmap[K, V]) Set(key K, value V) {
	r.root = blacken(r.insert(r.root, key, value))
}
func blacken[K, V any](n *node[K, V]) *node[K, V] {
	n.color = B
	return n
}
func redden[K, V any](n *node[K, V]) *node[K, V] {
	n.color = R
	return n
}
func (r *rbmap[K, V]) insert(n *node[K, V], key K, value V) *node[K, V] {
	if n == r.leaf {
		return &node[K, V]{
			color: R,
			key:   key,
			value: value,
			a:     r.leaf,
			b:     r.leaf,
		}
	}
	if r.less(key, n.key) {
		n.a = r.insert(n.a, key, value)
		n = balance(n)
	} else if r.less(n.key, key) {
		n.b = r.insert(n.b, key, value)
		n = balance(n)
	} else {
		n.value = value
	}
	return n
}
func colors[K, V any](n1, n2, n3 *node[K, V], c1, c2, c3 color) bool {
	return n1.color == c1 && n2.color == c2 && n3.color == c3
}
func balance[K, V any](n *node[K, V]) *node[K, V] {
	var x, y, z *node[K, V]
	var a, b, c, d *node[K, V]
	okasakiCase := false
	switch {
	case colors(n, n.a, n.a.a, B, R, R):
		x, y, z = n.a.a, n.a, n
		a, b, c, d = x.a, x.b, y.b, z.b
		okasakiCase = true
	case colors(n, n.a, n.a.b, B, R, R):
		x, y, z = n.a, n.a.b, n
		a, b, c, d = x.a, y.a, y.b, z.b
		okasakiCase = true
	case colors(n, n.b, n.b.a, B, R, R):
		x, y, z = n, n.b.a, n.b
		a, b, c, d = x.a, y.a, y.b, z.b
		okasakiCase = true
	case colors(n, n.b, n.b.b, B, R, R):
		x, y, z = n, n.b, n.b.b
		a, b, c, d = x.a, y.a, z.a, z.b
		okasakiCase = true
	}
	if okasakiCase {
		x.a, x.b, z.a, z.b = a, b, c, d
		y.a, y.b = x, z
		x.color, y.color, z.color = B, R, B
		return y
	}
	mightCase := false
	switch {
	case colors(n, n.a, n.a.a, BB, R, R):
		x, y, z = n.a.a, n.a, n
		a, b, c, d = x.a, x.b, y.b, z.b
		mightCase = true
	case colors(n, n.a, n.a.b, BB, R, R):
		x, y, z = n.a, n.a.b, n
		a, b, c, d = x.a, y.a, y.b, z.b
		mightCase = true
	case colors(n, n.b, n.b.a, BB, R, R):
		x, y, z = n, n.b.a, n.b
		a, b, c, d = x.a, y.a, y.b, z.b
		mightCase = true
	case colors(n, n.b, n.b.b, BB, R, R):
		x, y, z = n, n.b, n.b.b
		a, b, c, d = x.a, y.a, z.a, z.b
		mightCase = true
	default:
		c1, ok := deleteCaseI(n)
		if ok {
			return c1
		}
		c2, ok := deleteCaseII(n)
		if ok {
			return c2
		}
	}
	if mightCase {
		x.a, x.b, z.a, z.b = a, b, c, d
		y.a, y.b = x, z
		x.color, y.color, z.color = B, B, B
		return y
	}
	return n
}
func deleteCaseI[K, V any](n *node[K, V]) (*node[K, V], bool) {
	cond := n.color == BB && n.b.color == NB && n.b.a.color == B && n.b.b.color == B
	if !cond {
		return n, false
	}
	x, y, z := n, n.b.a, n.b
	a, b, c, d := x.a, y.a, y.b, z.b
	x.a, x.b = a, b
	z.a, z.b = c, redden(d)
	z.color = B
	y.a, y.b = x, balance(z)
	x.color, y.color, z.color = B, B, B
	return y, true
}
func deleteCaseII[K, V any](n *node[K, V]) (*node[K, V], bool) {
	cond := n.color == BB && n.a.color == NB && n.a.a.color == B && n.a.b.color == B
	if !cond {
		return n, false
	}
	x, y, z := n.a, n.a.b, n
	a, b, c, d := x.a, y.a, y.b, z.b
	x.a, x.b = redden(a), b
	z.a, z.b = c, d
	x.color = B
	y.a, y.b = balance(x), z
	x.color, y.color, z.color = B, B, B
	return y, true
}
func (r *rbmap[K, V]) Delete(key K) {
	r.root = blacken(r.del(r.root, key))
}
func (r *rbmap[K, V]) del(n *node[K, V], key K) *node[K, V] {
	if n == r.leaf {
		return r.leaf
	}
	if r.less(key, n.key) {
		n.a = r.del(n.a, key)
		n = r.bubble(n)
	} else if r.less(n.key, key) {
		n.b = r.del(n.b, key)
		n = r.bubble(n)
	} else {
		return r.remove(n)
	}
	return n
}
func (r *rbmap[K, V]) remove(n *node[K, V]) *node[K, V] {
	//fmt.Println("remove: ")
	//r.Preorder()
	//fmt.Println()
	if n == r.leaf {
		return r.leaf
	}
	if n.color == R && n.a == r.leaf && n.b == r.leaf {
		return r.leaf
	}
	if n.color == B && n.a == r.leaf && n.b == r.leaf {
		return r.bbleaf
	}
	if n.color == B && n.a == r.leaf && n.b != r.leaf && n.b.color == R {
		n.b.color = B
		return n.b
	}
	if n.color == B && n.b == r.leaf && n.a != r.leaf && n.a.color == R {
		n.a.color = B
		return n.a
	}
	//chasing same pointers twice. can optimize by
	// making max return a *node and passing that in to removeMax.
	n.key, n.value = r.max(n.a)
	n.a = r.removeMax(n.a)
	n = r.bubble(n)
	return n
}
func (r *rbmap[K, V]) max(n *node[K, V]) (K, V) {
	for n.b != r.leaf {
		n = n.b
	}
	return n.key, n.value
}
func (r *rbmap[K, V]) removeMax(n *node[K, V]) *node[K, V] {
	if n.b == r.leaf {
		return r.remove(n)
	}
	n.b = r.removeMax(n.b)
	return r.bubble(n)
}
func (r *rbmap[K, V]) bubble(n *node[K, V]) *node[K, V] {
	//fmt.Println("remove: ")
	//r.Preorder()
	//fmt.Println()
	if n.a.color == BB || n.b.color == BB {
		n.color = blacker(n.color)
		n.a = r.redder(n.a)
		n.b = r.redder(n.b)
		return balance(n)
	}
	return balance(n)
}
func (r *rbmap[K, V]) redder(n *node[K, V]) *node[K, V] {
	if n == r.bbleaf {
		return r.leaf
	}
	n.color = redder(n.color)
	return n
}
func redder(c color) color {
	switch c {
	case R:
		return NB
	case B:
		return R
	case BB:
		return B
	case NB:
		// can't happen
		panic("impossible")
	}
	panic("why come here")
}
func blacker(c color) color {
	switch c {
	case NB:
		return R
	case R:
		return B
	case B:
		return BB
	default:
		// BB cannot be blackened further
		panic("unmöglish")
	}
}
func (r rbmap[K, V]) CheckInvariants() {
	if r.root.color != B {
		panic("root must be black")
	}
	ys := make([]int, 0)
	xs := &ys
	r.check(r.root, 0, xs)
	i := 1
	for i < len(*xs) {
		if (*xs)[i-1] != (*xs)[i] {
			fmt.Println(xs)
			panic("black height not same for all the leaves")
		}
		i++
	}
}
func (r rbmap[K, V]) check(n *node[K, V], bh int, xs *[]int) {
	if n == r.leaf {
		*xs = append(*xs, bh)
		return
	}
	if n.color == R {
		if !colors(n, n.a, n.b, R, B, B) {
			r.Preorder()
			fmt.Println(n, n.a, n.b)
			panic("red node without both children black")
		}
	}
	if n.color == B {
		bh += 1
	}
	r.check(n.a, bh, xs)
	r.check(n.b, bh, xs)
}

