package ticket

import (
	"fmt"
	"math/rand"
)

func generateTicket(pfx string, l int) string {
	var TicketRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	t := make([]rune, l)
	for i := range t {
		t[i] = TicketRunes[rand.Intn(len(TicketRunes))]
	}

	return fmt.Sprintf("%s-%s", pfx, string(t))
}
