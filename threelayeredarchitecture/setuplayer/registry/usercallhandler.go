package registry

import "app/threelayeredarchitecture/usercallhandlerlayer"

/***********************************
 * usercall handler registry
 ***********************************/
type UserCallHandlerRegistryIF interface {
	GetTransferHandler() usercallhandlerlayer.TransferHandlerIF
}

func NewUserCallHandlerRegistry(
	userCallAdapterRegistry UserCallAdapterRegistryIF,
) UserCallHandlerRegistryIF {
	return UserCallHandlerRegistry{
		userCallAdapterRegistry: userCallAdapterRegistry,
	}
}

type UserCallHandlerRegistry struct {
	userCallAdapterRegistry UserCallAdapterRegistryIF
}

func (registry UserCallHandlerRegistry) GetTransferHandler() usercallhandlerlayer.TransferHandlerIF {
	return usercallhandlerlayer.NewTransferHandler(
		registry.userCallAdapterRegistry.GetTransferAdapter(),
	)
}
