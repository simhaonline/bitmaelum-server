package core

import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "os"
)

// ParseOptions will parse the commandline options given by opts. It will exit when issues arise or help is wanted
func ParseOptions(opts interface{}) {
    parser := flags.NewParser(opts, flags.IgnoreUnknown)
    _, err := parser.Parse()
    if err != nil {
        flagsError, _ := err.(*flags.Error)
        if flagsError.Type == flags.ErrHelp {
            os.Exit(1)
        }

        fmt.Println()
        parser.WriteHelp(os.Stdout)
        fmt.Println()
        os.Exit(1)
    }
}
