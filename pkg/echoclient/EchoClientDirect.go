package echoclient

import (
	"fmt"
	"log"
	"time"

	"github.com/hiveot/echorpc/pkg/echoservice"
)

// Invoke the upper service directly
func InvokeUpperDirect(text string, count int) {
	fmt.Println("Invoking upper directly")

	client := echoservice.NewEchoService()
	t1 := time.Now()
	for i := 0; i < count; i++ {
		response, err := client.Upper(text)
		if err != nil {
			log.Fatalf("error upper response: %s", err)
		}
		_ = response

		// fmt.Println("Response:", response)
	}
	d1 := time.Since(t1)
	fmt.Printf("Duration of %d calls using direct call: %d microsec\n", count, d1.Microseconds())

	// return response.Text, err
}
