package block

import "github.com/notebox/nb-crdt-go/common"

type ReplicaNonce []common.Nonce // []common.Nonce is a tuple of [block-nonce, text-nonce]
type Version map[common.ReplicaID]ReplicaNonce

func (v Version) IsNewerOrEqualThan(other Version) bool {
	for replicaID, version := range other {
		curr, ok := v[replicaID]
		if !ok || curr[0] < version[0] {
			return false
		}
	}
	return true
}

func (v Version) IsNewerThanExceptFor(other Version, replicaID common.ReplicaID) bool {
	result := false
	for rid, version := range other {
		if rid == replicaID {
			continue
		}
		curr, ok := v[rid]
		if !ok || curr[0] < version[0] {
			return false
		}
		if curr[0] > version[0] {
			result = true
		}
	}
	return result
}

func (v Version) Add(replicaID common.ReplicaID, nonces ReplicaNonce) bool {
	curr, ok := v[replicaID]
	if ok && curr[0] >= nonces[0] {
		return false
	}
	v[replicaID] = nonces
	return true
}

func (v Version) Merge(other Version) bool {
	var updated bool
	for replicaID, nonces := range other {
		ok := v.Add(replicaID, nonces)
		if ok && !updated {
			updated = true
		}
	}
	return updated
}
