/*
Copyright © 2025 Matt Krueger <mkrueger@rstms.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Version: "0.1.1",
	Use:     "gcrond",
	Short:   "userland cron daemon",
	Long: `
Read a contab file and run until stopped, executing the task schedule
described in the file. 

The first 5 fields of the crontab file set the schedule, and the 
following fields are passed to a shell for execution.

On Unix-like systems the default shell is /bin/sh
On Windows the value of environment variable COMSPEC is used.

The --seconds flag is provided, an extra seconds field is expected
in each crontab line

The --exec flag may be used to provide a single crontab record on
the command line instead of reading a crontab file.

Example crontab line:
    * * * * * path_to_command -flag params >/var/log/output

`,
	Run: func(cmd *cobra.Command, args []string) {
		crontab := viper.GetString("crontab")
		shell := viper.GetString("shell")
		cmdFlag := viper.GetString("flag")
		exec := viper.GetString("exec")
		seconds := viper.GetBool("seconds")
		crond, err := NewCron(crontab, exec, shell, cmdFlag, seconds)
		cobra.CheckErr(err)
		err = crond.Run()
		cobra.CheckErr(err)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	cobra.OnInitialize(InitConfig)
	OptionString("logfile", "l", "", "log filename")
	OptionString("config", "c", "", "config file")
	OptionString("crontab", "f", "gcrontab", "crontab file")
	OptionString("exec", "e", "", "crontab entry")
	OptionString("shell", "S", "", "shell binary pathname")
	OptionString("flag", "F", "", "shell command flag")
	OptionSwitch("debug", "", "produce debug output")
	OptionSwitch("verbose", "v", "increase verbosity")
	OptionSwitch("seconds", "s", "use seconds field in crontab")
}
