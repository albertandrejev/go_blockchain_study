package factory

import "bitbucket.org/albert_andrejev/free_info/utils"

//IMainFactory inteface for creating different factory instances
type IMainFactory interface {
	GetSimpleHash() utils.ISimpleHash
	GetX11Hash() utils.IX12Hash
}
