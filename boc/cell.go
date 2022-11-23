package boc

import "fmt"

const CellBits = 1023
const CellMaxRefs = 4

type Cell struct {
	bits BitString
	refs [4]*Cell
}

func NewCell() *Cell {
	return &Cell{
		bits: NewBitString(CellBits), // NewBitString()是由自定义的BitString提供的
		refs: [4]*Cell{},
	}
}
