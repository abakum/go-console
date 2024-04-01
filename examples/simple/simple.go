package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/abakum/go-console"
	"github.com/abakum/term"
)

// Using the `console` and raw mode of `os.Stdin` is the only way to call a shell in a loop for Windows.
func main() {
	var args []string

	fmt.Println("runtime.GOOS", runtime.GOOS)
	fmt.Println("runtime.GOARCH", runtime.GOARCH)

	var ioe *term.IOE

	for i := 0; i < 8; i++ {
		con, err := console.New(120, 60)
		if err != nil {
			panic(err)
		}
		defer con.Close()

		raw := i > 1
		shell := i > 3
		if runtime.GOOS == "windows" {
			args = []string{"cmd", "/c", "pause"}
			if shell {
				args = args[0:1]
				if i > 5 {
					args[0] = "powershell"
				}
			}
		} else {
			args = []string{"read", "-n1", "-rsp", "Press any key to continue . . ."}
			if shell {
				args = []string{"shell"}
			}
		}

		// case (cmd && raw && shell && shell=="cmd") is true then hang
		if i%2 == 1 && shell && args[0] == "cmd" {
			raw = false
		}
		m := "Press `Enter` then `Enter`"
		switch {
		case shell && !raw:
			m = "Type `exit` then press `Enter` then `Enter`"
		case shell:
			m = "Type `exit` then press `Enter`"
		case raw:
			m = "Press `Enter`"
		}
		fmt.Print(m, " raw ", raw, " shell ", shell)
		switch i % 2 {
		case 0:
			fmt.Println(" con", args[0])
			if raw {
				ioe = term.NewIOE()
			}
			if err := con.Start(args); err != nil {
				panic(err)
			}
			go func() {
				io.Copy(os.Stdout, con)
				if raw {
					ioe.Close()
				}
				fmt.Println("Stdout done")
			}()
			go func() {
				con.Wait()
				con.Close()
				if raw {
					fmt.Println()
				}
				fmt.Println("Wait done")
			}()

			if raw {
				io.Copy(con, ioe.ReadCloser())
			} else {
				io.Copy(con, os.Stdin)
			}

			fmt.Println("Stdin done")

		case 1:
			fmt.Println("", args[0])
			if raw {
				ioe = term.NewIOE()
			}
			cmd := exec.Command(args[0], args[1:]...)

			out, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}
			in, err := cmd.StdinPipe()
			if err != nil {
				panic(err)
			}

			err = cmd.Start()
			if err != nil {
				panic(err)
			}
			go func() {
				io.Copy(os.Stdout, out)

				if raw {
					ioe.Close()
				}
				in.Close()
				fmt.Println("Stdout done")
			}()
			go func() {
				cmd.Wait()
				fmt.Println("Wait done")
			}()

			if raw {
				io.Copy(in, ioe.ReadCloser())
			} else {
				io.Copy(in, os.Stdin)
			}

			fmt.Println("Stdin done")
		}
	}
}
