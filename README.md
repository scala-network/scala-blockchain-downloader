# Scala Blockchain Downloader

A simple tool to download and import the latest Scala blockchain file. It uses IPFS to fetch the files, by running an ephemeral node underneath.

![img](https://i.gyazo.com/1c3b6a58323e4efaba6428b68755dfd5.png)

## Common uses

To download and import the blockchain on first start

`./scala-blockchain-downloader`

## Compiling

The tool is written in Go and can be cross-compiled to Linux, Windows and MacOS.

### Linux

- Install Go

https://golang.org/dl/

- Clone the repository

```
git clone https://github.com/scala-network/scala-blockchain-downloader
cd scala-blockchain-downloader
```

* Make for all platforms

```
make
```

All the binaries will be found in the `bin/` directory.

## License

Please see [LICENSE](https://github.com/scala-network/scala-blockchain-downloader/blob/master/LICENSE) file.