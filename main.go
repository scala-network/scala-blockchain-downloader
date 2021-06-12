package main

import (
  "log"
  "os"
  "fmt"

  "github.com/urfave/cli/v2"
  "github.com/scala-network/scala-blockchain-downloader/src/cmd"
)

func main() {
  app := &cli.App{
    Flags: []cli.Flag {

      &cli.BoolFlag{
	Name: "download-only",
      },

    },
    Action: func(c *cli.Context) error {
      name := "for downloading only"
      cmd.Test(c.Bool("download-only"))
      fmt.Println(c.Bool("download-only"), name)
      return nil
    },
  }
  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
