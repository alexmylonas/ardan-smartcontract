package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ardanlabs/ethereum"
	"github.com/ardanlabs/ethereum/currency"
	"github.com/ardanlabs/smartcontract/app/signature/contract/go/verify"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	keyStoreFile     = "zarf/ethereum/keystore/UTC--2022-05-12T14-47-50.112225000Z--6327a38415c53ffb36c11db55ea74cc9cb4976fd"
	passPhrase       = "123"
	coinMarketCapKey = "a8cd12fb-d056-423f-877b-659046af0aa5"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx := context.Background()

	ownerKey, err := ethereum.PrivateKeyByKeyFile(keyStoreFile, passPhrase)
	if err != nil {
		return err
	}

	eth, err := ethereum.New(ctx, ethereum.NetworkLocalhost, keyStoreFile, passPhrase)
	if err != nil {
		return err
	}

	fmt.Println("\nInput Values")
	fmt.Println("----------------------------------------------------")
	fmt.Println("fromAddress:", eth.Address())

	// =========================================================================

	converter, err := currency.NewConverter(coinMarketCapKey)
	if err != nil {
		converter = currency.NewDefaultConverter()
	}
	oneETHToUSD, oneUSDToETH := converter.Values()

	fmt.Println("oneETHToUSD:", oneETHToUSD)
	fmt.Println("oneUSDToETH:", oneUSDToETH)

	// =========================================================================

	callOpts, err := eth.NewCallOpts(ctx)
	if err != nil {
		return err
	}

	// =========================================================================

	contractIDBytes, err := os.ReadFile("zarf/ethereum/verify.cid")
	if err != nil {
		return fmt.Errorf("importing verify.cid file: %w", err)
	}

	contractID := string(contractIDBytes)
	if contractID == "" {
		return errors.New("need to export the verify.cid file")
	}
	fmt.Println("contractID :", contractID)

	verify, err := verify.NewVerify(common.HexToAddress(contractID), eth.RawClient())
	if err != nil {
		return fmt.Errorf("new contract: %w", err)
	}

	// =========================================================================

	// The data to be encoded and signed.
	id := "asdjh1231"
	participant := common.HexToAddress("0x6327A38415C53FFb36c11db55Ea74cc9cB4976Fd")
	nonce := big.NewInt(1)

	// Encode the data the same way the smart contract will.
	bytes, err := encodeForSolidity(id, participant, nonce)
	if err != nil {
		return fmt.Errorf("encoding data: %w", err)
	}

	// Sign the message with the private key.
	signature, err := ethereum.SignBytes(bytes, ownerKey)
	if err != nil {
		return fmt.Errorf("signing message: %w", err)
	}

	// =========================================================================

	// Decode the hexadecimal signature to bytes for the smart contract calls.
	// It is easier and cheaper to do this here and not in the smart contract.
	sig, err := hex.DecodeString(signature[2:])
	if err != nil {
		return fmt.Errorf("decoding signature string: %w", err)
	}

	// Retrieve the address via the smart contract Address call.
	sigAddress, err := verify.Address(callOpts, id, participant, nonce, sig)
	if err != nil {
		return fmt.Errorf("address from message: %w", err)
	}

	// Retrieve the address via the smart contract Address call.
	matched, err := verify.MatchSender(callOpts, id, participant, nonce, sig)
	if err != nil {
		return fmt.Errorf("address from message: %w", err)
	}

	fmt.Println("\nResults")
	fmt.Println("----------------------------------------------------")
	fmt.Println("Keyfile address   :", eth.Address())
	fmt.Println("Signature address :", sigAddress)
	fmt.Println("Match sender      :", matched)

	return nil
}

// encodeForSolidity will take the arguments and pack them into a byte array that
// conforms to the solidity abi.encode function.
func encodeForSolidity(id string, participant common.Address, nonce *big.Int) ([]byte, error) {
	String, _ := abi.NewType("string", "", nil)
	Address, _ := abi.NewType("address", "", nil)
	Uint, _ := abi.NewType("uint", "", nil)

	arguments := abi.Arguments{
		{
			Type: String,
		},
		{
			Type: Address,
		},
		{
			Type: Uint,
		},
	}

	bytes, err := arguments.Pack(id, participant, nonce)
	if err != nil {
		return nil, fmt.Errorf("arguments pack: %w", err)
	}

	return bytes, nil
}
