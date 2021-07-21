package payment

import (
	"golang-crowdfunding/config"
	"golang-crowdfunding/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)


type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (midtrans.SnapResponse, error)
}

func NewService() *service{
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (midtrans.SnapResponse, error) {
	midclient := midtrans.NewClient()
    midclient.ServerKey = config.SERVER_KEY
    midclient.ClientKey = config.CLIENT_KEY
    midclient.APIEnvType = midtrans.Sandbox

	convertedTransactionId := strconv.Itoa(transaction.Id)

    snapGateway := midtrans.SnapGateway{
        Client: midclient,
    }

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: convertedTransactionId,
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenRes, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return midtrans.SnapResponse{}, err
	}

	return snapTokenRes, nil
}