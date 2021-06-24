package cmd

import (
	"net"

	ipfsLib "github.com/scala-network/scala-blockchain-downloader/src/ipfs"
)

func fetchFileHash() string {
	txtrecords, _ := net.LookupTXT("sbd-hash.scalaproject.io")
	var ret string;

	for _, txt := range txtrecords {
		ret = txt
	}

	return ret;
}

func DownloadAndImport(importToolPath string, dataDir string, downloadOnly bool, importVerify bool, force bool) {
	ipfsLib.DownloadChain(fetchFileHash(), importToolPath, dataDir, downloadOnly, importVerify, force)
}