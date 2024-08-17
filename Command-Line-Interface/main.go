package main

import (
	"cli/blockchain"  // Update with the correct import path
	cli "cli/cmdline" // Update with the correct import path
)

func main() {
	// Initialize the blockchain
	bc := blockchain.NewBlockchain()
	defer bc.DB.Close() // Ensure the database is closed when done

	// Initialize the command-line interface with the blockchain instance
	cli := cli.CLI{BC: bc}

	// Run the CLI to process command-line arguments and commands
	cli.Run()
}
