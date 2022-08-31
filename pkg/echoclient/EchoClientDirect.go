package echoclient

import (
	"fmt"
	"log"
	"time"

	"github.com/hiveot/echorpc/pkg/echoservice"
)

// Invoke the echo service directly
func InvokeEchoDirect(text string, count int) {
	fmt.Println("Invoking echo directly")

	client := echoservice.NewEchoService()
	t1 := time.Now()
	for i := 0; i < count; i++ {
		response, err := client.Echo(text)
		if err != nil {
			log.Fatalf("error echo response: %s", err)
		}
		_ = response

		// fmt.Println("Response:", response)
	}
	d1 := time.Since(t1)
	fmt.Printf("Duration of %d calls using direct call: %d microsec\n", count, d1.Microseconds())

	// return response.Text, err
}
