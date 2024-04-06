package tea

import "slices"

// _extSequences is used by the map-based algorithm below. It contains
// the sequences plus their alternatives with an escape character
// prefixed, plus the control chars, plus the space.
// It does not contain the NUL character, which is handled specially
// by detectOneMsg.
var _extSequences = func() map[string]Key {
	s := map[string]Key{
		" ":        {Type: KeySpace, Runes: spaceRunes},
		"\x1b ":    {Type: KeySpace, Alt: true, Runes: spaceRunes},
		"\x1b\x1b": {Type: KeyEsc, Alt: true},
	}
	for seq, key := range _sequences {
		key := key
		s[seq] = key
		if !key.Alt {
			key.Alt = true
			s["\x1b"+seq] = key
		}
	}
	for i := _keyNUL + 1; i <= _keyDEL; i++ {
		if i == _keyESC {
			continue
		}
		s[string([]byte{byte(i)})] = Key{Type: i}
		s[string([]byte{'\x1b', byte(i)})] = Key{Type: i, Alt: true}
		if i == _keyUS {
			i = _keyDEL - 1
		}
	}
	return s
}()

// _seqLengths is the sizes of valid sequences, starting with the
// largest size.
var _seqLengths = func() []int {
	seen := map[int]struct{}{}
	lsizes := make([]int, 0, len(seen))
	for seq := range _extSequences {
		if _, ok := seen[len(seq)]; !ok {
			seen[len(seq)] = struct{}{}
			lsizes = append(lsizes, len(seq))
		}
	}
	slices.SortFunc(lsizes, func(i, j int) int { return j - i })
	return lsizes
}()

// detectSequence uses a longest prefix match over the input
// sequence and a hash map.
func detectSequence(input []byte) (hasSeq bool, width int, msg Msg) { //nolint:nonamedreturns // too many returns
	seqs := _extSequences
	for _, sz := range _seqLengths {
		if sz > len(input) {
			continue
		}

		if key, ok := seqs[string(input[:sz])]; ok {
			return true, sz, MsgKey(key)
		}
	}

	// Is this an unknown CSI sequence?
	if loc := unknownCSIRe.FindIndex(input); loc != nil {
		return true, loc[1], msgUnknownCSISequence(input[:loc[1]])
	}

	return false, 0, nil
}
