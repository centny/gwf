package pool

import (
	"container/list"
	"github.com/Centny/gwf/util"
	"sync"
)

const GC_T = 300000

type bys_i struct {
	E *list.Element
	T int64
}

type ByteSlice struct {
	P     *BytePool
	size_ int
	ls_   *list.List
	ls_m_ map[interface{}]*bys_i
	ls_l  sync.RWMutex
	zero_ *list.Element
}

func NewByteSlice(p *BytePool, size int) *ByteSlice {
	ls_ := list.New()
	zero_ := ls_.PushBack([]byte{})
	return &ByteSlice{
		P:     p,
		size_: size,
		ls_:   ls_,
		zero_: zero_,
		ls_m_: map[interface{}]*bys_i{},
	}
}
func (b *ByteSlice) Alloc() []byte {
	b.ls_l.Lock()
	defer b.ls_l.Unlock()
	var bys []byte
	tv := b.ls_.Front()
	if tv == b.zero_ {
		bys = make([]byte, b.size_)
		tv = b.ls_.PushBack(bys)
		b.ls_m_[&bys[0]] = &bys_i{
			E: tv,
			T: util.Now(),
		}
	} else {
		b.ls_.MoveToBack(tv)
		bys = tv.Value.([]byte)
		b.ls_m_[&bys[0]].T = util.Now()
	}
	return bys
}
func (b *ByteSlice) Free(bys []byte) {
	b.ls_l.Lock()
	defer b.ls_l.Unlock()
	if tv, ok := b.ls_m_[&bys[0]]; ok {
		b.ls_.MoveToFront(tv.E)
	}
}
func (b *ByteSlice) Size() int64 {
	// b.ls_l.Lock()
	// defer b.ls_l.Unlock()
	return int64(b.ls_.Len()-1) * int64(b.size_)
}
func (b *ByteSlice) GC() (int, int64) {
	b.ls_l.Lock()
	defer b.ls_l.Unlock()
	tn := util.Now()
	rval := []interface{}{}
	for rv, vv := range b.ls_m_ {
		if tn-vv.T > b.P.T {
			rval = append(rval, rv)
		}
	}
	for _, rv := range rval {
		vv := b.ls_m_[rv]
		delete(b.ls_m_, rv)
		b.ls_.Remove(vv.E)
	}
	return len(rval), b.Size()
}

type BytePool struct {
	// Max  int64
	T    int64 //timeout when gc
	Beg  int
	End  int
	ms_  map[int]*ByteSlice
	ms_l sync.RWMutex
}

func NewBytePool(beg, end int) *BytePool {
	return (&BytePool{
		T:   GC_T,
		Beg: beg,
		End: end,
		ms_: map[int]*ByteSlice{},
	}).init(beg, end)
}
func (b *BytePool) init(beg, end int) *BytePool {
	if beg < 1 || end < 1 || (beg%8) != 0 || (end%8) != 0 {
		panic("beg/end must be a multiple of 8")
	}
	for i := (beg / 8); i <= (end / 8); i++ {
		size_ := i * 8
		b.ms_[size_] = NewByteSlice(b, size_)
	}
	return b
}
func (b *BytePool) Alloc(l int) []byte {
	tl := (l / 8) * 8
	if tl < l {
		tl += 8
	}
	if tl < 1 || tl > b.End {
		return nil
	}
	// b.ms_l.Lock()
	// defer b.ms_l.Unlock()
	tv := b.ms_[tl].Alloc()
	if l < tl {
		return tv[:l]
	} else {
		return tv
	}
}

func (b *BytePool) Free(bys []byte) {
	if bys == nil {
		return
	}
	l := len(bys)
	tl := (l / 8) * 8
	if tl < l {
		tl += 8
	}
	if tl < 1 || tl > b.End {
		return
	}
	// b.ms_l.Lock()
	// defer b.ms_l.Unlock()
	b.ms_[tl].Free(bys)
}

func (b *BytePool) Size() int64 {
	var tsize int64 = 0
	for _, bs_ := range b.ms_ {
		tsize += bs_.Size()
	}
	return tsize
}

func (b *BytePool) GC() (int, int64) {
	total := 0
	var tsize int64 = 0
	for _, bs := range b.ms_ {
		t_, ts_ := bs.GC()
		total += t_
		tsize += ts_
	}
	return total, tsize
}
