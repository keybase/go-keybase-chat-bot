package kbchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"
)

type walletMethod string

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

type sendOptions struct {
	Recipient     string  `json:"recipient"`
	Amount        string  `json:"amount"`
	Currency      *string `json:"currency"`
	Message       *string `json:"message"`
	FromAccountID *string `json:"from-account-id"`
	MemoText      *string `json:"memo-text"`
}

type getRes struct {
	Result stellar1.PaymentCLILocal `json:"result"`
	Error  Error                    `json:"error,omitempty"`
}

type sendRes struct {
	Result stellar1.SendResultCLILocal `json:"result"`
	Error  Error                       `json:"error,omitempty"`
}

type cancelRes struct {
	Result stellar1.RelayClaimResult `json:"result"`
	Error  Error                     `json:"error,omitempty"`
}

func (a *API) GetWalletTxDetails(txID string) (result stellar1.PaymentCLILocal, err error) {
	a.Lock()
	defer a.Unlock()

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

	response := getRes{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}
	if response.Error.Message != "" {
		return result, response.Error
	}
	return response.Result, nil
}

func (a *API) SendWalletTx(recipient string, amount string, currency *string, message *string, fromAccountID *string, memoText *string) (result stellar1.SendResultCLILocal, err error) {
	a.Lock()
	defer a.Unlock()

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

	response := sendRes{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}
	if response.Error.Message != "" {
		return result, response.Error
	}
	return response.Result, nil
}

func (a *API) CancelWalletTx(txID string) (result stellar1.RelayClaimResult, err error) {
	a.Lock()
	defer a.Unlock()

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

	response := cancelRes{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return result, fmt.Errorf("unable to decode wallet output: %s", err.Error())
	}
	if response.Error.Message != "" {
		return result, response.Error
	}
	return response.Result, nil
}
