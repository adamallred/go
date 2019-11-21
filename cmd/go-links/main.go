package main

import (
	"log"

	"github.com/kellegous/go/context"
	"github.com/kellegous/go/web"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("addr", ":8080", "default bind port")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	ctx, err := context.Open()
	if err != nil {
		log.Panic(err)
	}
	defer ctx.Close()

	log.Panic(web.ListenAndServe(ctx))
}
