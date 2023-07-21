package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	storagecontract "dApp/rpc_wrapper"

	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/block"
	"github.com/nspcc-dev/neo-go/pkg/core/native/nativenames"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/neorpc"
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/actor"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/gas"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/management"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/nep17"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/policy"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/manifest"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/nef"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/trigger"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
	"github.com/nspcc-dev/neo-go/pkg/vm/vmstate"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/pkg/errors"
)

const (
	rpcEndpoint = "http://localhost:30333"
	wsEndpoint  = "ws://localhost:30333/ws"
	walletPath  = "../my_wallet.json"
	accPass     = "qwerty"

	nefPath      = "../../neo-go/examples/storage/storage.nef"
	manifestPath = "../../neo-go/examples/storage/storage.manifest.json"
)

var (
	isNeoGoServer  bool
	transferTxH, _ = util.Uint256DecodeStringLE("7f1a2c41f0c03107f7a44ac510fa95fe11dde4c4994d30d61439f73f27e70f0d")
)

func main() {
	// Preparation: follow the guide provided in the README.md file, transfer
	// some GAS from multisignature to a simple account and compile the example
	// storage contract.

	// Simple RPC client: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient
	c, err := rpcclient.New(context.Background(), rpcEndpoint, rpcclient.Options{}) // Default options will be used.
	check(err, "create RPC client")
	defer c.Close()

	err = c.Init()
	check(err, "init RPC client")
	fmt.Printf("Work with simple RPC client:\n\n")

	// Simple RPC methods: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient#hdr-Client
	v, err := c.GetVersion()
	check(err, "retrieve version")
	fmt.Printf("RPC node version: %s\n", v.UserAgent)
	isNeoGoServer = strings.Contains(v.UserAgent, config.UserAgentPrefix) // NeoGo node offers several useful extensions.

	p, err := v.Protocol.MarshalJSON()
	check(err, "marshal protocol config")
	fmt.Printf("RPC node protocol config: %s\n", p)

	bCount, err := c.GetBlockCount()
	check(err, "retrieve block height")
	fmt.Printf("Block count: %d\n", bCount)

	currentB, err := c.GetBlockByIndex(bCount - 1) // Be careful with block count/block index conversions.
	check(err, "retrieve currentB block")

	transferTx, err := c.GetRawTransactionVerbose(transferTxH)
	check(err, "retrieve transfer tx")
	fmt.Printf("Transfer tx script: %s\n", base64.StdEncoding.EncodeToString(transferTx.Script))

	transferB, err := c.GetBlockByHash(transferTx.TransactionMetadata.Blockhash)
	check(err, "retrieve transferB block")

	tr := trigger.Application
	applog, err := c.GetApplicationLog(transferTxH, &tr) // Various triggers are supported.
	check(err, "retrieve applog")
	if len(applog.Executions) != 1 {
		panic("unexpected executions number")
	}
	stack, err := json.MarshalIndent(applog.Executions[0].Stack, "", "\t")
	check(err, "marshal stack")
	fmt.Printf("Transfer applog invocation state: %s\nTransfer applog stack: %s\n", applog.Executions[0].VMState, stack)

	// Almost all JSON RPC API methods are supported by the NeoGo RPC client. For details
	// on unsupported methods read the https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient#hdr-Client.

	// Work with wallets: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/wallet
	w, err := wallet.NewWalletFromFile(walletPath)
	check(err, "open wallet")

	acc := w.Accounts[0]
	err = acc.Decrypt(accPass, w.Scrypt)
	check(err, "decrypt account")

	// Some of the extensions offered by NeoGo RPC server: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#extensions.
	if isNeoGoServer {
		// NEP17/NEP11 transfers paging: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#limits-and-paging-for-getnep11transfers-and-getnep17transfers.
		var (
			start = uint64(0)
			stop  = currentB.Timestamp
			limit = 10
			page  = 0
		)
		nep17T, err := c.GetNEP17Transfers(acc.ScriptHash(), &start, &stop, &limit, &page)
		check(err, "retrieve NEP17 transfers")
		tBytes, err := json.MarshalIndent(nep17T, "", "\t")
		check(err, "marshal NEP17 transfers")
		fmt.Printf("NEP17 transfers of %s:\n%s\n", acc.Address, tBytes)

		// `getblocksysfee` RPC call: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#getblocksysfee-call
		sysfee, err := c.GetBlockSysFee(nep17T.Received[0].Index)
		check(err, "getblocksysfee")
		fmt.Printf("Block #%d system fee: %d\n", nep17T.Received[0].Index, sysfee)

		// For more extensions read the NeoGo RPC server documentation: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#extensions.
	}

	// Web-socket client and server: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#websocket-server
	// https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient#WSClient
	ctx, cancel := context.WithCancel(context.Background())
	wsC, err := rpcclient.NewWS(ctx, wsEndpoint, rpcclient.WSOptions{
		Options:                        rpcclient.Options{},
		CloseNotificationChannelIfFull: false,
	})
	err = wsC.Init()
	check(err, "init WS Client")
	// Do not use `defer wsC.Close()` here, we'll close the client manually later.

	err = wsC.Init()
	check(err, "init WS RPC client")
	fmt.Printf("\nWork with WS RPC client:\n\n")

	// Same simple RPC methods over web-socket:
	gasState, err := wsC.GetContractStateByAddressOrName(nativenames.Gas)
	check(err, "getcontractstate")
	fmt.Printf("GAS contract ID: %d\n", gasState.ID)

	// Unwrap helper package: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap
	dec, err := unwrap.BigInt(wsC.InvokeFunction(gasState.Hash, "decimals", []smartcontract.Parameter{}, nil))
	check(err, "invokefunction")
	fmt.Printf("GAS contract decimals: %s\n", dec)

	// All set of methods supported by standard NeoGo JSON RPC client are supported, but there's also
	// a special extension: notifications subsystem: https://github.com/nspcc-dev/neo-go/blob/master/docs/notifications.md
	bCh := make(chan *block.Block, 5) // Add some buffer to prevent WSC from blocking even regular requests.
	ntfCh := make(chan *state.ContainedNotificationEvent, 5)
	aerCh := make(chan *state.AppExecResult, 5)
	dispatcherToMainCh := make(chan struct{})
	go func(ctx context.Context, bCh chan *block.Block, ntfCh chan *state.ContainedNotificationEvent, exCh chan *state.AppExecResult, exitCh chan struct{}) { // TODO: move to a separate function
	dispatcherLoop:
		for {
			select {
			case b := <-bCh:
				fmt.Printf("Block was received:\n\tIndex: %d\n\tPrimary: %d\n\tTransactions: %d\n",
					b.Index, b.PrimaryIndex, len(b.Transactions))
			case ntf := <-ntfCh:
				params, _ := ntf.Item.MarshalJSON()
				fmt.Printf("Notification from execution was received:\n\tContainer: %s\n\tContract: %s\n\tName: %s\n\tParams: %s\n",
					ntf.Container.StringLE(), ntf.ScriptHash.StringLE(), ntf.Name, params)
			case aer := <-aerCh:
				st, _ := json.Marshal(aer.Stack)
				fmt.Printf("Application execution result was received:\n\tContainer: %s\n\tTrigger: %s\n\tVM state: %s\n\tStack:%s\n\tFault exception: %s\n", aer.Container.StringLE(), aer.Trigger, aer.VMState, st, aer.FaultException)
			case <-ctx.Done():
				break dispatcherLoop
			}
		}
	drainLoop:
		for {
			select {
			case <-bCh:
			case <-ntfCh:
			case <-aerCh:
			default:
				break drainLoop
			}
		}
		close(bCh)
		close(ntfCh)
		close(aerCh)

		// Send notification to the main routine.
		close(dispatcherToMainCh)
	}(ctx, bCh, ntfCh, aerCh, dispatcherToMainCh)

	primary := 0
	till := bCount + 5
	fltB := &neorpc.BlockFilter{
		Primary: &primary,
		Since:   nil,
		Till:    &till,
	}
	bSubID, err := wsC.ReceiveBlocks(fltB, bCh)
	check(err, "subscribe for blocks notifications")

	ctr := gasState.Hash
	name := "Transfer"
	fltNtf := &neorpc.NotificationFilter{
		Contract: &ctr,
		Name:     &name,
	}
	ntfSubID, err := wsC.ReceiveExecutionNotifications(fltNtf, ntfCh)
	check(err, "subscribe for execution notifications")

	st := vmstate.Halt.String()
	fltEx := &neorpc.ExecutionFilter{
		State: &st,
	}
	aerSubId, err := wsC.ReceiveExecutions(fltEx, aerCh)
	check(err, "subscribe for AppExecResults")

	// Pretend our app is running and performing some useful work.
	t := time.NewTimer(5 * time.Second)
	var exitErr error
mainLoop:
	for {
		select {
		case <-t.C:
			break mainLoop
		case <-ctx.Done():
			exitErr = wsC.GetError()
			break mainLoop
		}
	}

	// End of the app work.
	if exitErr != nil {
		fmt.Printf("WS RPC client closing error: %s\n", exitErr)
	} else {
		// We can use wsC.UnsubscribeAll(), but let us show the ID-based unsubscriptions:
		err = wsC.Unsubscribe(bSubID)
		check(err, "unsubscribe from blocks")
		err = wsC.Unsubscribe(ntfSubID)
		check(err, "unsubscribe from execution notifications")
		wsC.Unsubscribe(aerSubId)
		check(err, "unsubscribe from application execution results")
	}
	cancel()
	<-dispatcherToMainCh // Wait for the dispatcher routine to properly finish.

	wsC.Close()
	err = wsC.GetError() // Any error if the closing process wasn't initiated by user via wsC.Close().
	if err != nil {
		fmt.Printf("WS RPC client closing error: %s\n", err)
	}

	// Invoker functionality: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker
	ctx, cancel = context.WithCancel(context.Background())
	wsC, err = rpcclient.NewWS(ctx, wsEndpoint, rpcclient.WSOptions{})
	check(err, "create WS client")

	signers := []transaction.Signer{
		{
			Account: acc.ScriptHash(),
			Scopes:  transaction.CalledByEntry, // An informative article about witnesses and witness scopes: https://neospcc.medium.com/thou-shalt-check-their-witnesses-485d2bf8375d
		},
	}
	inv := invoker.New(wsC, signers)

	// Test invocation of method "transfer" of the GAS contract.
	var (
		from   = acc.ScriptHash()
		to     = util.Uint160{1, 2, 3} // some random account
		amount = 5
		data   = smartcontract.Parameter{Type: smartcontract.AnyType}
	)
	res, err := inv.Call(gasState.Hash, "transfer", from, to, amount, data)
	ok, err := unwrap.Bool(res, err)
	check(err, "perform `transfer` testinvoke")
	fmt.Printf("`transfer` result: %t\n", ok)

	// Historic invocation functionality: https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#invokecontractverifyhistoric-invokefunctionhistoric-and-invokescripthistoric-calls
	invH := invoker.NewHistoricAtHeight(transferB.Index-1, wsC, signers)
	resH, err := invH.Call(gasState.Hash, "transfer", from, to, amount, data)
	ok, err = unwrap.Bool(resH, err)
	check(err, "invoke historic `transfer`")
	fmt.Printf("`transfer` historic result: %t\n", ok)

	// Checking historic balance:
	resH, err = invH.Call(gasState.Hash, "balanceOf", from)
	balance, err := unwrap.BigInt(resH, err)
	check(err, "invoke historic `balanceOf`")
	fmt.Printf("`balanceOf` historic result: %d\n", balance)

	// Painful manual events unwrapping. Do it once and never do it again, because NeoGo has
	// a set of nice helper functions that perform automatic events unwrapping, see the
	// example below.
	var tEvent nep17.TransferEvent
	arr := res.Notifications[0].Item.Value().([]stackitem.Item)
	fromB, _ := arr[0].TryBytes()
	tEvent.From, _ = util.Uint160DecodeBytesBE(fromB)
	toB, _ := arr[1].TryBytes()
	tEvent.To, _ = util.Uint160DecodeBytesBE(toB)
	tEvent.Amount, _ = arr[2].TryInteger()

	// A set of other invoker methods: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker

	// Actor functionality: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/actor
	actSigners := []actor.SignerAccount{
		{
			Signer:  signers[0],
			Account: acc, // Provide the decrypted account if possible. It can be signture/multisignature/contract/dummy.
		},
	}
	act, err := actor.New(wsC, actSigners)
	check(err, "create actor")

	aer, err := act.Wait(act.SendCall(gasState.Hash, "transfer", from, to, amount, data))
	check(err, "send `transfer` call via actor")
	if aer.VMState != vmstate.Halt {
		panic("unexpected `transfer` result")
	}

	aer, err = act.Wait(act.SendTunedCall(gasState.Hash, "transfer", []transaction.Attribute{}, func(r *result.Invoke, t *transaction.Transaction) error {
		err := actor.DefaultCheckerModifier(r, t)
		if err != nil {
			return err
		}

		// Perform some additional checks...
		if len(r.Stack) != 1 {
			return fmt.Errorf("unexpected result stack len: %d", len(r.Stack))
		}
		ok, err := r.Stack[0].TryBool()
		if err != nil {
			return fmt.Errorf("unexpected result stack content: %s", r.Stack[0].Type())
		}
		if !ok {
			return errors.New("false `transfer` result")
		}

		// Change the transaction...
		t.NetworkFee = r.GasConsumed + 1

		return nil
	}, from, to, amount, data))
	check(err, "send `transfer` tuned call via actor")
	if aer.VMState != vmstate.Halt {
		panic("unexpected `transfer` tuned result")
	}

	// A set of other Actor methods: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/actor

	// In general, Actor API is very flexible, but if it's not enough, then you can always
	// use customizable actor:
	// _, _ = actor.NewTuned(wsC, actSigners, actor.Options{})

	// NEP17 package: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/nep17
	nep17Act := nep17.New(act, gasState.Hash)

	d, err := nep17Act.Decimals()
	check(err, "retrieve GAS decimals via NEP17 actor")
	fmt.Printf("GAS decimals: %d\n", d)

	aer, err = act.Wait(nep17Act.MultiTransfer([]nep17.TransferParameters{
		{
			From:   from,
			To:     from,
			Amount: big.NewInt(5),
			Data:   nil,
		},
		{
			From:   from,
			To:     to,
			Amount: big.NewInt(5),
			Data:   nil,
		},
	}))
	check(err, "GAS multitransfer")
	applogTr, err := wsC.GetApplicationLog(aer.Container, nil)
	check(err, "retrieve applog for multitransfer")
	transferEvents, err := nep17.TransferEventsFromApplicationLog(applogTr)
	check(err, "retrieve events from multitransfer applog")
	transferEventsBytes, _ := json.MarshalIndent(transferEvents, "", "\t")
	fmt.Printf("GAS multitransfer events:\n%s\n", transferEventsBytes)

	// Native-specific actor packages (there are a lot of them, but here's the example
	// for native GasToken contract): https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/gas
	gasAct := gas.New(act)
	fromBalance, err := gasAct.BalanceOf(from)
	check(err, "retrieve `acc` balance")
	fmt.Printf("`acc` balance: %d\n", fromBalance)

	policyAct := policy.New(act)
	feePerByte, err := policyAct.GetFeePerByte()
	check(err, "retrieve FeePerByte")
	fmt.Printf("FeePerByte value: %d\n", feePerByte)

	// A set of other native-specific actors are supported:
	// NeoToken: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/neo
	// ContractManagement: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management
	// PolicyContract: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/policy
	// OracleContract: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/oracle
	// RoleManagement: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/rolemgmt
	// Notary: https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/notary

	// Contract deployment example using native-specific management actor:
	// https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management

	// First of all, compile the storage contract example via NeoGo CLI.
	// The example is located at the NeoGo repo: https://github.com/nspcc-dev/neo-go/tree/master/examples/storage
	// To compile storage contract use the following NeoGo CLI command.
	// Modify the paths correspondingly if you have storage contract located at a
	// different place, and adjust the `manifestPath` and `nefPath` constants.
	// ./bin/neo-go contract compile -i ./examples/storage/storage.go -c ./examples/storage/storage.yml -o examples/storage/storage.nef -m ./examples/storage/storage.manifest.json

	// Deploy the compiled contract using native-specific Actor:
	mgmtAct := management.New(act)

	f, err := os.ReadFile(nefPath)
	check(err, "read NEF")
	nefFile, err := nef.FileFromBytes(f)
	check(err, "decode NEF")

	manifestBytes, err := os.ReadFile(manifestPath)
	check(err, "read manifest")
	m := new(manifest.Manifest)
	err = json.Unmarshal(manifestBytes, m)
	check(err, "decode manifest")

	// Check whether the contract was deployed earlier:
	sH := state.CreateContractHash(acc.ScriptHash(), nefFile.Checksum, m.Name)
	sState, err := mgmtAct.GetContract(sH)
	check(err, "failed to retrieve storage contract")
	fmt.Printf("storage contract hash: %s\n", sH)

	// If not, then deploy it:
	if sState == nil {
		aer, err = act.Wait(mgmtAct.Deploy(&nefFile, m, nil))
		check(err, "deploy storage contract")
		if aer.VMState != vmstate.Halt {
			panic(fmt.Sprintf("deploy storage contract: %s", aer.FaultException))
		}

		// Another way to retrieve contract hash is via `deploy` resulting stack:
		sHBytes, err := aer.Stack[0].Value().([]stackitem.Item)[2].TryBytes()
		check(err, "retrieve storage contract hash from stackitem")
		sH, err = util.Uint160DecodeBytesBE(sHBytes)
		check(err, "decode storage contract hash")

		// Put several values into storage of the deployed contract:
		_, err = act.Wait(act.SendCall(sH, "put", "key1", "value1"))
		check(err, "put `key1` into storage contract")
		_, err = act.Wait(act.SendCall(sH, "put", "key2", "value2"))
		check(err, "put `key2` into storage contract")
		_, err = act.Wait(act.SendCall(sH, "put", "key3", "value3"))
		check(err, "put `key3` into storage contract")
	}

	// Iterate over storage values on the RPC-server side:
	res, err = act.CallAndExpandIterator(sH, "findRetIter", 5, "key")
	sItems, err := unwrap.Array(res, err)
	check(err, "expand iterator values")
	fmt.Println("`Storage` contract iterator values with `key` prefix:")
	for _, itm := range sItems {
		kvPair := itm.Value().([]stackitem.Item)
		fmt.Printf("key: `%s`\tvalue: `%s`\n", kvPair[0].Value().([]byte), kvPair[1].Value().([]byte))
	}

	// Work with iterator sessions:
	res, err = act.Call(sH, "findRetIter", "key")
	sessID, iter, err := unwrap.SessionIterator(res, err)
	check(err, "unwrap iterator session")

	iterVals, err := act.TraverseIterator(sessID, &iter, 2)
	check(err, "traverse iterator")
	fmt.Println("`Storage` contract iterator values with `key` prefix retrieved from iterator:")
	for _, itm := range iterVals {
		kvPair := itm.Value().([]stackitem.Item)
		fmt.Printf("key: `%s`\tvalue: `%s`\n", kvPair[0].Value().([]byte), kvPair[1].Value().([]byte))
	}

	err = act.TerminateSession(sessID)
	check(err, "terminate iterator session")

	// NeoGo provides the ability to generate RPC bindings for Neo smart contract
	// written in Go or in any other programming language. However, using NeoGo
	// compiler with Go smart contracts allows to get more precise resulting RPC
	// binding. See the https://github.com/nspcc-dev/neo-go/blob/master/docs/compiler.md#generating-rpc-contract-bindings
	// for details on bindings autogeneration.
	//
	// To compile example Storage contract and generate RPC binding for it use
	// the following NeoGo CLI commands:
	//   ./bin/neo-go contract compile -i ./examples/storage/storage.go -c ./examples/storage/storage.yml -o examples/storage/storage.nef -m ./examples/storage/storage.manifest.json --bindings examples/storage/storage_bindings.config.json
	//   ./bin/neo-go contract generate-rpcwrapper -c examples/storage/storage_bindings.config.json -m examples/storage/storage.manifest.json -o examples/storage/rpc/storage.go
	//
	// After that you can add the resulting auto-generated RPC smart contract
	// binding to your dApp project as a separate package and call the deployed
	// contract directly from your dAPP code which is very convenient:
	storageAct := storagecontract.New(act, sH)
	_, err = act.Wait(storageAct.Put([]byte("key4"), []byte("value4")))
	check(err, "put key4 to the storage contract")
}

func check(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("failed to %s: %w", msg, err))
	}
}
