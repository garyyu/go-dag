# GoDAG
BlockDAG algorithms Go Lang simulation.

BlockChain (for example Bitcoin, Etherum, etc.) is just a 'k=0' special subtype of BlockDAG, that's why they suffer from the highly restrictive throughput. DAG is the future!

Thanks **Yonatan Sompolinsky** and **Aviv Zohar** for the most important contributions on blockDAG research, and their great paper "PHANTOM: A Scalable BlockDAG protocol" on [International Association for Cryptologic Research (IACR)](https://eprint.iacr.org/2018/104.pdf) in Feb. 2018.

They setup a start-up company to develop BlockDAG since Q4 2017, their website: [https://www.daglabs.com]. And there's an official DAGlabs slack channel: [https://daglabs.slack.com].

---

# How to build this simulation

Here is my Go Lang environment, the other Go version such as 1.9 should also be OK.
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
![Fig.4](https://github.com/garyyu/go-dag/blob/master/pics/Fig.4.jpg)

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

# Click to Run the demo

And *Run it in Go Playground* :   ![#c5f015](https://placehold.it/15/c5f015/000000?text=+)[Click to Run](https://play.golang.org/p/n8ckn-8X0CA)![#c5f015](https://placehold.it/15/c5f015/000000?text=+). You can **edit** your blockDAG example and **run** the 'blue selection algorithm' online!
The output will be like this:
```console
- BlockDAG Algorithms Simulation - Algorithm 1: Selection of a blue set. -
chainInitialize(): done. blocks= 20
blue set selection done. blue blocks = (G).B.C.D.F.I.J.K.O.P.M.R.(V). 	total blue: 13
```
# Other Tests

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
# Engagement

Please join us the BlockDAG discussion on [https://godag.github.io](https://godag.github.io).



