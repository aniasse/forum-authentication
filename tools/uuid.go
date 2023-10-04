package tools

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

func Generate_Id(name string) uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("❌ error while generating id ", err)
	}

	if name == "" {
		return id
	}

	new_id := uuid.NewV5(id, name)
	fmt.Println("✅ Id has been generated successfully")
	return new_id

}

func Id_from_str(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		log.Fatalln("❌ error while generating id from string", err)
	}
	fmt.Println("✅ Id has been generated successfully from string")
	return id

}
