package main

import (
	"context"
	"flag"
	"fmt"
	"po/internal/app"
	"po/internal/boot"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var port string
	flag.StringVar(&port, "p", "8000", "The application port")
	flag.Parse()

	if err := boot.Boot(ctx); err != nil {
		panic(err)
	}

	if err := app.Serve(ctx, port); err != nil {
		panic(err)
	}

	fmt.Println("Build something AMAZING!")
}

//config := vault.DefaultConfig()
//
//config.Address = "http://127.0.0.1:8200"
//
//client, err := vault.NewClient(config)
//
//if err != nil {
//	log.Fatalf("unable to initialize Vault client: %v", err)
//}
//
//client.SetToken("hvs.6ja7wvLDcakpK2EMh6WTUMeo")
//
//_, err = client.KVv2("secret").Put(context.Background(), "mysql_password", map[string]interface{}{
//	"password": "Hashi123",
//})
//
//if err != nil {
//	log.Fatalf("unable to write secret: %v", err)
//}
//
//fmt.Println("Secret written successfully.")
