package cli

import (
	"cli/blockchain" // Update with the correct import path
	"cli/pow"        // Update with the correct import path
	"cli/transaction"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI handles command-line arguments
type CLI struct{}

// PrintUsage prints the usage instructions
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	getbalance --address ADDRESS : Get balance of ADDRESS")
	fmt.Println("	createblockchain -address ADDRESS : Create a blockhchain and send genesis block reward to ADDRESS")
	fmt.Println(" 	printchain : print all the blocks of the blockchain")
	fmt.Println("	send -from FROM -to TO -amount AMOUNT : send AMOUNT of coins from FROM address to TO")
}

// ValidateArgs validates command-line arguments
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.DB.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain(address)
	defer bc.DB.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance = out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// PrintChain prints all blocks in the blockchain
func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain("")
	defer bc.DB.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := pow.NewPOW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain(from)
	defer bc.DB.Close()

	tx := transaction.NewUTXOTx(from, to, amount, bc)
	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCommand := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCommand := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCommand := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddr := getBalanceCommand.String("address", "", "The address to get balance for")
	createBlockchainAddr := createBlockchainCommand.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCommand.String("from", "", "Source wallet address")
	sendTo := sendCommand.String("to", "", "Destination wallet address")
	sendAmount := sendCommand.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCommand.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCommand.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCommand.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCommand.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCommand.Parsed() {
		if *getBalanceAddr == "" {
			getBalanceCommand.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddr)
	}

	if createBlockchainCommand.Parsed() {
		if *createBlockchainAddr == "" {
			createBlockchainCommand.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddr)
	}

	if printChainCommand.Parsed() {
		cli.printChain()
	}

	if sendCommand.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCommand.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
