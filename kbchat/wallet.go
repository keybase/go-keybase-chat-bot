package kbchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"
)

type walletMethod string

type sendOptions struct {
	Recipient     string  `json:"recipient"`
	Amount        string  `json:"amount"`
	Currency      *string `json:"currency"`
	Message       *string `json:"message"`
	FromAccountID *string `json:"from-account-id"`
	MemoText      *string `json:"memo-text"`
}

type SendOutput struct {
	Result stellar1.SendResultCLILocal `json:"result"`
}

type CancelOutput struct {
	Result stellar1.RelayClaimResult `json:"result"`
}

type walletTxIDAPIReq struct {
	Method walletMethod     `json:"method"`
	Params walletTxIDParams `json:"params"`
}

type walletTxIDParams struct {
	Options txIDOptions `json:"options"`
}

type txIDOptions struct {
	TxID string `json:"txid"`
}

type walletSendAPIReq struct {
	Method walletMethod     `json:"method"`
	Params walletSendParams `json:"params"`
}

type walletSendParams struct {
	Options sendOptions `json:"options"`
}

func (a *API) GetWalletTxDetails(txID string) (result stellar1.PaymentCLILocal, err error) {
	a.Lock()
	defer a.Unlock()

	type res struct {
		Result stellar1.PaymentCLILocal `json:"result"`
	}

	opts := txIDOptions{
		TxID: txID,
	}
	args := walletTxIDAPIReq{Method: "details", Params: walletTxIDParams{Options: opts}}
	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}
	cmd := a.runOpts.Command("wallet", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return result, err
	}

	response := res{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}

	return response.Result, nil
}

func (a *API) SendWalletTx(recipient string, amount string, currency *string, message *string, fromAccountID *string, memoText *string) (result stellar1.SendResultCLILocal, err error) {
	a.Lock()
	defer a.Unlock()

	type res struct {
		Result stellar1.SendResultCLILocal `json:"result"`
	}

	opts := sendOptions{
		Recipient:     recipient,
		Amount:        amount,
		Currency:      currency,
		Message:       message,
		FromAccountID: fromAccountID,
		MemoText:      memoText,
	}
	args := walletSendAPIReq{Method: "send", Params: walletSendParams{Options: opts}}
	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}
	cmd := a.runOpts.Command("wallet", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return result, err
	}

	response := res{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}

	return response.Result, nil
}

func (a *API) CancelWalletTx(txID string) (result stellar1.RelayClaimResult, err error) {
	a.Lock()
	defer a.Unlock()

	type res struct {
		Result stellar1.RelayClaimResult `json:"result"`
	}

	opts := txIDOptions{
		TxID: txID,
	}
	args := walletTxIDAPIReq{Method: "cancel", Params: walletTxIDParams{Options: opts}}
	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}

	cmd := a.runOpts.Command("wallet", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return result, err
	}

	response := res{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}

	return response.Result, nil
}
