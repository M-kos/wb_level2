package currenttime

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

func CurrentTime() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Println("error writing to stderr:", err)
		}
		os.Exit(1)
	}

	fmt.Println(time.String())
}
