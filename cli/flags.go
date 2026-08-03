package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagData struct {
	FmtCompat  bool `mapstructure:"fmtcompat"`
	Check      bool `mapstructure:"check"`
	Verbose    bool `mapstructure:"verbose"`
	Quiet      bool `mapstructure:"quiet"`
	Uncoloured bool `mapstructure:"uncoloured"`

	Fmt    FlagsFmt    `mapstructure:",squash"`
	Blocks FlagsBlocks `mapstructure:",squash"`
}

// FlagsFmt holds the flags shared by the fmt and diff commands.
type FlagsFmt struct {
	Pattern        string `mapstructure:"pattern"`
	FixFinishLines bool   `mapstructure:"fix-finish-lines"`
}

// FlagsBlocks holds the flags for the blocks command.
type FlagsBlocks struct {
	ZeroTerminated bool `mapstructure:"zero-terminated"`
	JSON           bool `mapstructure:"json"`
}

// flagEnvMap is the full set of viper-managed flags and the env var each one can be
// set with ("" = flag only: zero-terminated and json select per-invocation output
// framing for scripts, an env var would silently corrupt whatever is parsing the output).
var flagEnvMap = map[string]string{
	"fmtcompat":        "TERRAFMT_FMTCOMPAT",
	"check":            "TERRAFMT_CHECK",
	"verbose":          "TERRAFMT_VERBOSE",
	"quiet":            "TERRAFMT_QUIET",
	"uncoloured":       "TERRAFMT_UNCOLOURED",
	"pattern":          "TERRAFMT_PATTERN",
	"fix-finish-lines": "TERRAFMT_FIX_FINISH_LINES",
	"zero-terminated":  "",
	"json":             "",
}

func configureFlags(root *cobra.Command) error {
	pflags := root.PersistentFlags()
	pflags.BoolP("fmtcompat", "f", false, "enable format string (%s, %d etc) compatibility")
	pflags.BoolP("check", "c", false, "return an error during diff if formatting is required")
	pflags.BoolP("verbose", "v", false, "show files as they are processed & additional stats")
	pflags.BoolP("quiet", "q", false, "quiet mode, only shows block line numbers ")
	pflags.BoolP("uncoloured", "u", false, "disable coloured output")

	for name, env := range flagEnvMap {
		if env == "" {
			continue
		}
		if err := viper.BindEnv(name, env); err != nil {
			return fmt.Errorf("error binding '%s' to env '%s': %w", name, env, err)
		}
	}

	viper.SetConfigName(".terrafmt")
	viper.SetConfigType("env")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

// bindCommandFlags binds the executed command's flags (local + inherited persistent) to
// viper. This runs from the root PersistentPreRunE rather than configureFlags because
// fmt and diff each define their own local --pattern flag: a static bind would attach
// whichever instance was bound last regardless of which command is actually running.
func bindCommandFlags(cmd *cobra.Command) error {
	var errs error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if _, ok := flagEnvMap[f.Name]; !ok {
			return // not viper-managed (e.g. help)
		}
		if err := viper.BindPFlag(f.Name, f); err != nil && errs == nil {
			errs = fmt.Errorf("error binding '%s' flag: %w", f.Name, err)
		}
	})
	return errs
}

// GetFlags returns the fully populated FlagData. We unmarshal from Viper instead of
// reading pflags directly because pflags only parse command-line arguments; Viper merges
// environment variables and the .terrafmt config file on top of them.
func GetFlags() (*FlagData, error) {
	var f FlagData
	if err := viper.Unmarshal(&f); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &f, nil
}
