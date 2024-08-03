package main

import (
	"context"
	"fmt"
	secure "github.com/Dentrax/obscure-go/types"
	immudb "github.com/codenotary/immudb/pkg/client"
	"log"
)

// use:
// https://github.com/codenotary/immudb-client-examples/tree/master/go/todos-sample-stdlib
// https://github.com/codenotary/immudb-client-examples/tree/master/go
// https://play.codenotary.com/?topic=cli%2Fselect&live=true

func main() {
	fmt.Println("Hello world!")
	//
	// key := "aloo"
	// DEFAULTS
	pass := secure.NewString("immudb")
	user := "immudb"
	opts := immudb.DefaultOptions().WithAddress("127.0.0.1").WithPort(3322)
	client := immudb.NewClient().WithOptions(opts)
	//
	err := client.OpenSession(context.Background(), []byte(user), []byte(pass.Get()), "defaultdb")
	if err != nil {
		log.Fatal(err)
	}
	//
	log.Println("1")
	//
	// // ensure connection is closed
	defer client.CloseSession(context.Background())
	//
	// // write an entry
	// // upon submission, the SDK validates proofs and updates the local state under the hood
	// hdr, err := client.VerifiedSet(context.Background(), []byte(key), []byte("immutable world"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Sucessfully set a verified entry: ('%s', '%s') @ tx %d\n", []byte(key), []byte("immutable world"), hdr.Id)
	//
	// // read an entry
	// // upon submission, the SDK validates proofs and updates the local state under the hood
	// entry, err := client.VerifiedGet(context.Background(), []byte(key))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Sucessfully got verified entry: ('%s', '%s') @ tx %d\n", entry.Key, entry.Value, entry.Tx)

}
