package utils

type BitMap int

func (bm *BitMap) IsSet(val int) bool {
	return int(*bm)&val == val
}

func (bm *BitMap) Set(val int) {
	newval := int(*bm) | val
	*bm = BitMap(newval)
}

func (bm *BitMap) Clear(val int) {
	newval := int(*bm) & ^val
	*bm = BitMap(newval)
}
