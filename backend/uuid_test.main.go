package main

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	var u pgtype.UUID
	u.Scan("c94fb118-eb87-43f1-b844-6a849767554d")
	fmt.Println(u.String())
}
