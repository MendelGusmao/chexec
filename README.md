# chexec

Is a simple wrapper for `os/exec`.Command useful for running long-lived background processes.
It uses channels for reading from the child process' stdout and stderr.

## Example

```golang
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MendelGusmao/chexec"
)

func main() {
	cmd := chexec.Command("/bin/bash", "-c", "c=0; while true; do echo $c; date; ping thisdomaindoesnt.exist; c=$((c+1)); [ $c -eq 3 ] && exit; sleep 1; done")

	if err := cmd.Run(); err != nil {
		fmt.Println("run:", err)
		os.Exit(1)
	}

	done := cmd.Wait()
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			fmt.Println("process timed out")

			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("stopping:", err)
			}

			os.Exit(0)
		case v := <-cmd.Stdout:
			fmt.Printf("cmd.Stdout: %s\n", v)
		case v := <-cmd.Stderr:
			fmt.Printf("cmd.Stderr: %s\n", v)
		case err := <-done:
			fmt.Println("waiting:", err)
			os.Exit(1)
		}
	}
}
```
