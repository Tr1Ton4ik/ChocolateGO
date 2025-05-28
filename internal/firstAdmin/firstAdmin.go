package firstAdmin

import (
	"chocolate/internal/passwords"
	"os"
)

var Name = os.Getenv("NAME")
var Password = passwords.Hash(os.Getenv("PASSWORD"))
