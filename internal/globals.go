package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/chifra"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/data"
	"github.com/spf13/cobra"
)

type Globals struct {
	Format  string
	Update  bool
	Script  bool
	Verbose bool
	Sleep   time.Duration
	Rounds  []data.Round
}

func GetGlobals(defFmt string, cmd *cobra.Command, args []string) (ret Globals, err error) {
	if !utils.IsTerminal() || os.Getenv("NO_COLORS") == "true" {
		chifra.ColorsOff()
	}

	if !file.FolderExists(data.DataFolder()) {
		err = fmt.Errorf("data folder (%s) not found in current working folder", data.DataFolder())
		return
	}

	envFmt := os.Getenv("TB_DEFAULT_FMT")
	if len(envFmt) > 0 {
		defFmt = envFmt
	}
	if ret.Format, err = cmd.Flags().GetString("fmt"); err != nil {
		return
	}
	if ret.Format == "" {
		ret.Format = defFmt
	}

	if ret.Update, err = cmd.Flags().GetBool("update"); err != nil {
		return
	}

	if ret.Script, err = cmd.Flags().GetBool("script"); err != nil {
		return
	}

	if ret.Verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
		return
	}

	var s uint64
	if s, err = cmd.Flags().GetUint64("sleep"); err != nil {
		return
	}
	ret.Sleep = time.Duration(s)

	var round uint64
	if round, err = cmd.Flags().GetUint64("round"); err != nil {
		return
	}
	ret.Rounds, err = data.GetRounds(int(round), 25)
	if ret.Update {
		last := ret.Rounds[len(ret.Rounds)-1]
		ret.Rounds = []data.Round{}
		ret.Rounds = append(ret.Rounds, last)
	}

	return
}
