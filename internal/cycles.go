package internal

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

func RunCycles(cmd *cobra.Command, args []string) error {
	dataType, globals, err := getCycleOptions(cmd, args)
	if err != nil {
		return err
	}

	fmt.Println(dataType)
	fmt.Println(globals)

	return nil
}

// getCycleOptions processes command line options for the Rounds command
func getCycleOptions(cmd *cobra.Command, args []string) (dataType string, globals Globals, err error) {
	globals, err = GetGlobals("csv", cmd, args)
	if err != nil {
		return
	}

	dataType, err = cmd.Flags().GetString("data")
	if err != nil {
		return
	}

	if globals.Script {
		if globals.Sleep == 0 {
			// If it's unset and we're going to be calling the API, set it
			globals.Sleep = 3
		}
	}

	// if len(dataType) > 0 {
	// 	// we got one, check if it's valid. If yes, we're done
	// 	if isValidType(dataType) {
	// 		err = validate(dataType, globals)
	// 		return
	// 	}
	// }
	// for _, arg := range args {
	// 	if isValidType(arg) {
	// 		dataType = arg
	// 		err = validate(dataType, globals)
	// 		return
	// 	}
	// }
	// if len(args) > 0 {
	// 	err = fmt.Errorf("invalid option '%s'", args[0])
	// } else {
	// 	err = fmt.Errorf("flag needs an argument: --data")
	// }

	globals.Rounds = []data.Round{}

	return
}
