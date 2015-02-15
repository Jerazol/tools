package main

import (
	"fmt"
        "os"
	"os/exec"
        "log"

)

func main() {
        // find Downloads/ -type f -exec stat --printf "%s\n" "{}" \;
        out, err := exec.Command("/home/tommygi/repos/tools/findfilesizes.sh", os.Args[1]).Output()
        if err != nil {
                log.Fatal(err)
        }

        fmt.Printf("%s", out[1])
}



