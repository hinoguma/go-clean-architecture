package setup

import (
    "app/threelayerdarchitecture/applogic/domain"
    "app/threelayerdarchitecture/applogic/usecase"
    "app/threelayerdarchitecture/appinfraadapter/repository"
    "app/threelayerdarchitecture/appinfra/memorydb"
    "app/threelayerdarchitecture/usercalladapter/transfer"
)

type Container struct {
    accountClient     *memorydb.AccountClient
    transactionClient *memorydb.TransactionClient
}

func NewContainer() *Container {
    accounts := memorydb.NewAccountClient([]memorydb.AccountRecord{
        {
            ID:       "acct-setup-1",
            UserID:   "user-setup-1",
            Balance:  25_000,
            Currency: "USD",
        },
        {
            ID:       "acct-setup-2",
            UserID:   "user-setup-2",
            Balance:  15_000,
            Currency: "USD",
        },
    })

    transactions := memorydb.NewTransactionClient(nil)

    return &Container{
        accountClient:     accounts,
        transactionClient: transactions,
    }
}

func (c *Container) NewTransferAdapter() transfer.Adapter {
	accountRepo := repository.NewAccountRepository(c.accountClient)
	transactionRepo := repository.NewTransactionRecorder(c.transactionClient)
	service := domain.NewTransferService(nil)
	uc := usecase.NewTransferUseCase(accountRepo, transactionRepo, service)
	return transfer.NewTransferAdapter(uc)
}
