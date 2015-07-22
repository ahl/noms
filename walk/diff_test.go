package walk

import (
	"testing"

	"github.com/attic-labs/noms/types"
	"github.com/stretchr/testify/assert"
)

func TestGetReachabilitySetDiff(t *testing.T) {
	assert := assert.New(t)

	// {"string": "string",
	//  "map": {"nested": "string"}
	//  "mtlist": []
	// }
	small := types.NewMap(
		types.NewString("string"), types.NewString("string"),
		types.NewString("map"), types.NewMap(types.NewString("nested"), types.NewString("string")),
		types.NewString("mtlist"), types.NewList())

	setKey := types.NewString("set")
	setElem := types.Int32(7)
	setVal := types.NewSet(setElem)
	big := small.Set(setKey, setVal)

	var hashes []string
	for _, r := range GetReachabilitySetDiff(small, big) {
		hashes = append(hashes, r.String())
	}

	assert.Contains(hashes, setKey.Ref().String())
	assert.Contains(hashes, setElem.Ref().String())
	assert.Contains(hashes, setVal.Ref().String())

	assert.Empty(GetReachabilitySetDiff(small, small))
}