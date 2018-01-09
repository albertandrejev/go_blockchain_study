package factory

import "bitbucket.org/albert_andrejev/free_info/utils"

//IMainFactory inteface for creating different factory instances
type IMainFactory interface {
	GetSimpleHash() *utils.SimpleHash
	GetX11Hash() *utils.X12Hash
}
