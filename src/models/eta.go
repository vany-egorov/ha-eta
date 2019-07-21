package models

import "strconv"

type ETA uint64

func (it ETA) String() string { return strconv.FormatUint(uint64(it), 10) }
func (it ETA) Bytes() []byte  { return []byte(it.String()) }
