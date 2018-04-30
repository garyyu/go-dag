# phantom
BlockDAG algorithm's Go language simulation for paper "PHANTOM: A Scalable BlockDAG protocol".

Here is the [paper link on International Association for Cryptologic Research (IACR)](https://eprint.iacr.org/2018/104.pdf) by **Yonatan Sompolinsky** and **Aviv Zohar** in Feb. 2018.

They have setup a start-up company to develop BlockDAG since Q4 2017, their website: [https://www.daglabs.com]. And there's an official DAGlabs slack channel: [https://daglabs.slack.com].

---

# How to build

My Go Lang environment:
```bash
$ go version
go version go1.10.1 darwin/amd64
```
To run the simulation for the example on the paper page 7 Fig.3 or page 16 Fig.4, for algorithm 1 **Selection of a blue set**.

```bash
$ go test -run=Fig3
$ go test -run=Fig4
```

To add a new example DAG to see the DAG blue selection behaviour, it's quite easy. For example, to test a DAG in this figure 'Fig.4', just add a piece of codes like this:
![Fig.4](https://github.com/garyyu/go-phantom/blob/master/pics/Fig.4.jpg)

```golang
	ChainAddBlock("Genesis", []string{}, chain)

	ChainAddBlock("B", []string{"Genesis"}, chain)
	ChainAddBlock("C", []string{"Genesis"}, chain)
	ChainAddBlock("D", []string{"Genesis"}, chain)
	ChainAddBlock("E", []string{"Genesis"}, chain)

	ChainAddBlock("F", []string{"B","C"}, chain)
	ChainAddBlock("I", []string{"C","D"}, chain)
	ChainAddBlock("H", []string{"E"}, chain)

	ChainAddBlock("J", []string{"F","D"}, chain)
	ChainAddBlock("L", []string{"F"}, chain)
	ChainAddBlock("K", []string{"J","I","E"}, chain)
	ChainAddBlock("N", []string{"D","H"}, chain)

	ChainAddBlock("M", []string{"L","K"}, chain)
	ChainAddBlock("O", []string{"K"}, chain)
	ChainAddBlock("P", []string{"K"}, chain)
	ChainAddBlock("Q", []string{"N"}, chain)

	ChainAddBlock("R", []string{"O","P","N"}, chain)

	ChainAddBlock("S", []string{"Q"}, chain)
	ChainAddBlock("T", []string{"S"}, chain)
	ChainAddBlock("U", []string{"T"}, chain)
```
And run it in Go Playground:  .
The output will be like this:
```console
- Phantom Paper Simulation - Algorithm 1: Selection of a blue set. -
chainInitialize(): done. blocks= 20
blue set selection done. blue blocks = (G).B.C.D.F.I.J.K.O.P.M.R.(V). 	total blue: 13
```

To run the simulation for the example on page 3 Fig.2, page 8 "C. Step #2", for algorithm 2 **Ordering of the DAG**.

```bash
$ go test -run=Fig2
```

To run the benchmark test:

```bash
$ go test ./phantom -bench=Blocks -benchmem

$ go test ./ -test.bench BlueSelection -benchmem  -run=^$

$ go test ./ -test.bench BlockOrdering -benchmem  -run=^$
```

Please join us the BlockDAG discussion on [https://godag.github.io](https://godag.github.io).



