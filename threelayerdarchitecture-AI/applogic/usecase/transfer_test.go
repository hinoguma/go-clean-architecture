package usecase_test

import (
    "context"
    "testing"

    "app/threelayerdarchitecture/applogic/domain"
    "app/threelayerdarchitecture/applogic/usecase"
    "app/threelayerdarchitecture/appinfraadapter/repository"
    "app/threelayerdarchitecture/appinfra/memorydb"
)

func TestTransferUseCase_Success(t *testing.T) {
    accountClient := memorydb.NewAccountClient([]memorydb.AccountRecord{
        {ID: "acct-1", UserID: "user-1", Balance: 10_000, Currency: "USD"},
        {ID: "acct-2", UserID: "user-2", Balance: 2_000, Currency: "USD"},
    })
    transactionClient := memorydb.NewTransactionClient(nil)

    accountRepo := repository.NewAccountRepository(accountClient)
    txnRepo := repository.NewTransactionRecorder(transactionClient)
    uc := usecase.NewTransferUseCase(accountRepo, txnRepo, domain.NewTransferService(nil))

    result, err := uc.Execute(context.Background(), usecase.TransferRequest{
        InitiatingUserID: "user-1",
        FromAccountID:    "acct-1",
        ToAccountID:      "acct-2",
        Amount:           1_000,
    })
    if err != nil {
        t.Fatalf("execute transfer: %v", err)
    }
    if result.TransactionID == "" {
        t.Fatalf("expected transaction id to be generated")
    }

    if record, err := transactionClient.Get(result.TransactionID); err != nil {
        t.Fatalf("transaction not saved: %v", err)
    } else if record.Amount != 1_000 {
        t.Fatalf("unexpected transaction amount: %d", record.Amount)
    }
}
