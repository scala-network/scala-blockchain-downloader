package main

import (
  "fmt"
  "log"
  "os"
  "path/filepath"
  "bufio"
  "runtime"
  "os/exec"

  "github.com/urfave/cli/v2"
  cmd "github.com/scala-network/scala-blockchain-downloader/src/cmd"
)

func main() {

fmt.Printf(`
  ███████╗ ██████╗ █████╗ ██╗      █████╗ 
  ██╔════╝██╔════╝██╔══██╗██║     ██╔══██╗
  ███████╗██║     ███████║██║     ███████║
  ╚════██║██║     ██╔══██║██║     ██╔══██║
  ███████║╚██████╗██║  ██║███████╗██║  ██║
  ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝  
                      BLOCKCHAIN DOWNLOADER
`)

fmt.Println("\n")

  var workingDir string
  var errW error

  workingDir, errW = os.Executable()
  if errW != nil {
    fmt.Println(errW)
    fmt.Print("Press enter to continue...")
    _, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
    os.Exit(0)
  }
  workingDir = filepath.Dir(workingDir)


  cli.VersionFlag = &cli.BoolFlag{
    Name: "version",
    Aliases: []string{"V"},
    Usage: "print only the version",
  }

  importToolPath := filepath.Join(workingDir, "scala-blockchain-import")

  if runtime.GOOS == "linux" {
    fmt.Printf("Fixing UDP issue")
    exec.Command("sysctl -w net.core.rmem_max=2500000").Output()
  }

  if runtime.GOOS == "windows" {
    importToolPath = filepath.Join(workingDir, "scala-blockchain-import.exe")
  }

  app := &cli.App{
    Version: "v1.0.0",
    Name: "Scala Blockchain Downloader",
    Usage: "A simple tool to download and import the latest Scala blockchain database",
    Flags: []cli.Flag {

      &cli.StringFlag{
        Name: "data-dir",
        Value: "",
        Usage: "Set a custom blockchain path",
        Required: false,
      },

      &cli.BoolFlag{
        Name: "without-import-verification",
        Value: true,
        Usage: "If --dangerous-unverified-import 1 should be used on import",
      },

      &cli.BoolFlag{
	      Name: "download-only",
        Value: false,
        Usage: "Only download the blockchain data from IPFS, don't import it.",
      },

      &cli.StringFlag{
        Name: "import-tool-path",
        Value: importToolPath,
        Usage: "Set a custom path for scala-blockchain-import binary",
        Required: false,
      },

    },
    Action: func(c *cli.Context) error {
      cmd.DownloadAndImport(c.String("import-tool-path"), c.String("data-dir"), c.Bool("download-only"), c.Bool("without-import-verification"));
      return nil
    },
  }
  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}