package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/iotexproject/ioTube/witness-service/util"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "", // No explicit use string for the root command
		Short: "CLI application for Bitcoin transaction operations",
		Long:  `This CLI application allows you to sign and send Bitcoin transactions and display public keys.`,
		// The Run function can be omitted or used for default behavior
	}

	// Variables for showpubkey command
	var privateKeyForPubKey string

	// Define the showpubkey command
	var showPubKeyCmd = &cobra.Command{
		Use:   "showpubkey",
		Short: "Display the public key in hex from a given private key",
		Long:  `Display the public key in hexadecimal format from the provided Bitcoin private key.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Private Key: %s\n", privateKeyForPubKey)
			b, err := hex.DecodeString(privateKeyForPubKey)
			if err != nil {
				panic(err)
			}
			_, pbk := btcec.PrivKeyFromBytes(b)
			fmt.Printf("Public Key: %s\n", hex.EncodeToString(pbk.SerializeUncompressed()))
		},
	}

	// Adding flags to the showpubkey command
	showPubKeyCmd.Flags().StringVarP(&privateKeyForPubKey, "privatekey", "p", "", "Private key for generating the public key")
	showPubKeyCmd.MarkFlagRequired("privatekey")

	// Variables for taprootaddr command
	var publicKey string

	// Define the taprootaddr command
	var taprootAddrCmd = &cobra.Command{
		Use:   "taprootaddr",
		Short: "Generate a Taproot address from a given public key",
		Long:  `Generate a Taproot address from the provided public key.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Public Key: %s\n", publicKey)
			pbk, err := util.HexToPubkey(publicKey)
			if err != nil {
				panic(err)
			}
			addr, _, err := util.TaprootAddrFromPubkey(pbk, &chaincfg.TestNet3Params)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Taproot Address: %s\n", addr.EncodeAddress())
		},
	}

	// Adding flags to the taprootaddr command
	taprootAddrCmd.Flags().StringVarP(&publicKey, "publickey", "p", "", "Public key for generating the Taproot address")
	taprootAddrCmd.MarkFlagRequired("publickey")

	// Variables for combinedpubkey command
	var publicKeys []string

	// Define the combinedpubkey command
	var combinedPubKeyCmd = &cobra.Command{
		Use:   "combinedpubkey",
		Short: "Generate a combined public key from an array of public keys",
		Long:  `Generate a combined public key from the provided array of public keys.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Public Keys: %s\n", publicKeys)
			pbks := make([]*btcec.PublicKey, 0)
			for _, pk := range publicKeys {
				pbk, err := util.HexToPubkey(pk)
				if err != nil {
					panic(err)
				}
				pbks = append(pbks, pbk)
			}
			combinedPubkey, err := util.Musig2CombinedPubkey(pbks)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Combined Public Key: %s\n", hex.EncodeToString(combinedPubkey.SerializeUncompressed()))
		},
	}

	// Adding flag to the combinedpubkey command
	combinedPubKeyCmd.Flags().StringSliceVarP(&publicKeys, "publickeys", "p", []string{}, "Array of public keys for generating the combined public key")
	combinedPubKeyCmd.MarkFlagRequired("publickeys")

	// Add commands to the root command
	rootCmd.AddCommand(showPubKeyCmd, taprootAddrCmd, combinedPubKeyCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
