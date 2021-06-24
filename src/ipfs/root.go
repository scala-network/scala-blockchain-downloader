package ipfs

import (
	"context"
	"fmt"
	"time"
	"io"
	"io/ioutil"
	//"log"
	"os/exec"
	"os"
	"path/filepath"
	"bufio"
	"bytes"
	//"strings"
	//"sync"
	//"net"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	//peerstore "github.com/libp2p/go-libp2p-peerstore"
	//ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/TwinProduction/go-color"
	"github.com/cheggaaa/pb/v3"
	//"github.com/libp2p/go-libp2p-core/peer"
)

type path struct {
	path string
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}


func createTempRepo(ctx context.Context) (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func spawnEphemeral(ctx context.Context) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	return createNode(ctx, repoPath)
}

func printDownloadedSize(totalSize int) (err error){
	bar := pb.New(totalSize)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	bar.Start()

	var oldSize int = 0
	var count int = 0

	for {
		fi, err := os.Stat("blockchain.raw")
		if err != nil {
			return err
		}
		size := fi.Size()

		if count == 0 {
			bar.Add((int(size)))
			oldSize = int(size)
			count++
		}else{
			bar.Add((int(size) - oldSize))
			oldSize = int(size)
			count++
		}

		time.Sleep(time.Millisecond)

		if ((totalSize - int(size)) / 100000) <= 50  {
			bar.Finish()
			break
		}
	}
	return nil
}

func DownloadChain(hash string, importToolPath string, dataDir string, downloadOnly bool, importVerify bool) {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start ephemeral IPFS node
	ipfs, err := spawnEphemeral(ctx)

	if err != nil {
        panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
    }

	fmt.Println(color.Ize(color.Green, "Started an ephemeral IPFS node"))

	outputBasePath := "./"
	outputPathFile := outputBasePath + "blockchain.raw"

	ipfsHash := icorepath.New(hash)

	fmt.Println(color.Ize(color.Green, "Start downloading blockchain data"))
	fmt.Printf("\n")
	fileStat, err := ipfs.Object().Stat(ctx, ipfsHash)

	if err != nil {
		panic(fmt.Errorf("Could not object stat with error: %s", err))
	}

	go printDownloadedSize(fileStat.CumulativeSize)

	rootNodeFile, err := ipfs.Unixfs().Get(ctx, ipfsHash)
	if err != nil {
		panic(fmt.Errorf("Could not get file with CID: %s", err))
	}

	err = files.WriteTo(rootNodeFile, outputPathFile)
	if err != nil {
		panic(fmt.Errorf("Could not write out the fetched CID: %s", err))
	}

	fmt.Println("\n")
	fmt.Println(color.Ize(color.Green, "Downloaded blockchain data\n"))

	if downloadOnly == false {
		_, err := os.Stat(importToolPath)
		if os.IsNotExist(err) {
			fmt.Printf(color.Ize(color.Red,`
The blockchain import tool 'scala-blockchain-import' does 
not exist in the current directory. 

Please execute this tool from the same path as the 'scala-blockchain-import' 
tool or set the flag --import-tool-path to the correct location
`))
			fmt.Print("\n")
			fmt.Print("Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		}else{
			importArgs := []string{
				"--input-file",
				outputPathFile,
			}

			if importVerify == true {
				importArgs = append(importArgs, "--dangerous-unverified-import=1")
			}

			if dataDir != "" {
				dataDirArg := fmt.Sprintf(`--data-dir %v`, dataDir)
				importArgs = append(importArgs, dataDirArg)
			}

			
			importCommand := exec.Command(
				importToolPath,
				importArgs...)


			stdoutIn, _ := importCommand.StdoutPipe()
			stderrIn, _ := importCommand.StderrPipe()	

			var errStdout, errStderr error
			var stdoutBuf, stderrBuf bytes.Buffer
			stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
			stderr := io.MultiWriter(os.Stderr, &stderrBuf)
			err = importCommand.Start()
			if err != nil {
				fmt.Printf("Unable to start the import tool: %s\n", err)
				fmt.Print("Press enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				os.Exit(0)
			}
	
			go func() {
				_, errStdout = io.Copy(stdout, stdoutIn)
			}()
	
			go func() {
				_, errStderr = io.Copy(stderr, stderrIn)
			}()
	
			err = importCommand.Wait()
			if err != nil {
				fmt.Printf("Failed to import the downloaded blockchain: %s\n", err)
				fmt.Print("Press enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				os.Exit(0)
			}
			if errStdout != nil || errStderr != nil {
				fmt.Printf("Unable to capture the import tool's output: %s, %s\n",
					errStdout,
					errStderr)
				fmt.Print("Press enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				os.Exit(0)
			}
	
			err = os.Remove(outputPathFile)
			if err != nil {
				fmt.Printf("The downloaded file '%s' could not be removed: %s\n",
				outputPathFile,
					err,
				)
			}	

			fmt.Printf(`
Imported downloaded blockchain file successfully.
You may now start 'scalad' or your wallet.
Thank you for using the Scala Blockchain Downloader.
`)
			

		}
	}else{
		s1 := fmt.Sprintf(`
You selected the download only option. You may use the
downloaded blockchain file now.

The location of the downloaded file is: %v

`, outputPathFile)

		fmt.Printf(color.Ize(color.Green, s1))
		os.Exit(0)	
	}
}