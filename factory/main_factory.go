package factory

import "bitbucket.org/albert_andrejev/free_info/utils"

//MainFactory dependency injection factory
type MainFactory struct {
	simpleHash *utils.SimpleHash
	x11Hash    *utils.X12Hash
}

//NewMainFactory constructor
func NewMainFactory() *MainFactory {
	return new(MainFactory)
}

//GetSimpleHash simple hashing (sha3+blake2s) factory
func (t MainFactory) GetSimpleHash() *utils.SimpleHash {
	if t.simpleHash == nil {
		t.simpleHash = utils.NewSimpleHash(new(utils.SimpleHashWrap))
	}
	return t.simpleHash
}

//GetX11Hash simple hashing (x11+scrypt) factory
func (t MainFactory) GetX11Hash() *utils.X12Hash {
	if t.x11Hash == nil {
		t.x11Hash = utils.NewX12Hash(new(utils.X12HashWrapper))
	}

	return t.x11Hash
}
