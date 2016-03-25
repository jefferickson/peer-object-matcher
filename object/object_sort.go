package object

// Implement the sort interface:
// Len()
// Less()
// Swap()

func (pcs peerComps) Len() int {
	return len(pcs)
}

// first sort by distance, then hashed ID
func (pcs peerComps) Less(i, j int) bool {
	if pcs[i].Distance == pcs[j].Distance {
		return pcs[i].PeerObject.HashedID < pcs[j].PeerObject.HashedID
	} else {
		return pcs[i].Distance < pcs[j].Distance
	}
}

func (pcs peerComps) Swap(i, j int) {
	pcs[i], pcs[j] = pcs[j], pcs[i]
}
