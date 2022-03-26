package main

import (
	"CldResolver/engine/src/Cloudflare"
	"fmt"
	"log"
)

func main() {
	res, err := Cloudflare.Resolve("site.xyz", "subdomain")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)

}
