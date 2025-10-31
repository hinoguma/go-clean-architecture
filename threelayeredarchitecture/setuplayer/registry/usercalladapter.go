package registry

import "app/threelayeredarchitecture/usercalladapterlayer"

/***********************************
 * usercall adapter registry
 ***********************************/
type UserCallAdapterRegistryIF interface {
	GetTransferAdapter() usercalladapterlayer.TransferAdapterIF
}

func NewUserCallAdapterRegistry(
	usecaseRegistry UseCaseRegistryIF,
) UserCallAdapterRegistryIF {
	return UserCallAdapterRegistry{
		usecaseRegistry: usecaseRegistry,
	}
}

type UserCallAdapterRegistry struct {
	usecaseRegistry UseCaseRegistryIF
}

func (registry UserCallAdapterRegistry) GetTransferAdapter() usercalladapterlayer.TransferAdapterIF {
	return usercalladapterlayer.NewTransferAdapter(
		registry.usecaseRegistry.GetTransferUseCase(),
	)
}
