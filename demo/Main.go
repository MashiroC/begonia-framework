package main

import (
	"fmt"
)

//func main() {
//	app := begonia.Init()
//	app.Get("/getQuestion", func(ctx *begonia.Context) {
//		ctx.W.Header().Add("Access-Control-Allow-Origin", "*")
//		time.Sleep(time.Duration(3)*time.Second)
//		res:=""
//		for i:=0;i<500;i++ {
//			res+="css是魔鬼吧\n"
//		}
//		ctx.String(res)
//	})
//	app.Start(1234)
//}

type Tree struct {
	n *Node
}

type NodeType uint8

const (
	Normal NodeType = iota
	Root
	Param
	All
)

type Node struct {
	Path      string
	children  []*Node
	nodeType  NodeType
	quickFind string
	priority  int
	handle    interface{}
}

func (n *Node) AddChild(path string, h interface{}) {

	if len(n.Path) == 0 {
		n.Path = path
		n.priority = 1
		n.handle = h
		return
	}

	if n.Path == path {
		panic("两条相同路由")
	}

	//找到最长公共前缀
	minLength := Min(len(path), len(n.Path))
	i := 0
	for i < minLength && path[i] == n.Path[i] {
		i++
	}
	//如果该节点是最后一个节点
	if i != 0 && i == len(n.Path) {
		for j := 0; j < len(n.quickFind); j++ {
			if n.quickFind[j] == path[i] {
				//TODO 判断新加入节点对于参数节点的影响合不合法
				n.children[j].AddChild(path[i:], h)
				return
			}
		}
		if i < len(path) {
			n.insertNode(path[i:], h)
		} else {
			n.handle = h
		}
		return
		//panic("寻找公共前缀出了点问题 " + n.Path + " " + path)
	}
	//不是最后一个节点
	n.splitNode(i)
	//fmt.Println(n.Path, path, i, len(n.Path), len(path))
	if i == len(path) {
		fmt.Println(n.Path, path, i, len(n.Path), len(path))
		n.handle = h
	} else if i >= len(n.Path) {
		n.insertNode(path[i:], h)
	}
}

func (n *Node) splitNode(i int) (newNode *Node) {

	if n.nodeType != Normal {
		panic("只能分割normal结点,已有参数节点")
	}
	//原来的后面那部分
	newNode = CopyNode(n)
	newNode.Path = n.Path[i:]

	n.Path = n.Path[:i]
	n.priority = n.priority + 1
	n.quickFind = string([]byte{newNode.Path[0]})
	n.children = []*Node{newNode}
	n.handle = nil

	return
}

func (n *Node) insertNode(path string, handle interface{}) (newNode *Node) {
	//具有参数节点
	for i := 0; i < len(path); i++ {
		// /asd/:asd/zxcasd/zxcv
		// /asd/:asd/zxchdgfh
		if path[i] == ':' {

			//fmt.Println(path)
			if i != 0 {
				firstNode := NewNode(path[:i])
				n.children = append(n.children, firstNode)
				n.quickFind = n.quickFind + string([]byte{firstNode.Path[0]})
				n = firstNode
			}

			j := i
			for j < len(path) && path[j] != '/' {
				j++
			}

			if n.Path[0] == ':' {
				panic("已有参数节点")
			}

			secondNode := NewNode(path[i:j])
			secondNode.nodeType = Param
			secondNode.priority = 0

			n.children = append(n.children, secondNode)
			n.quickFind = n.quickFind + string([]byte{secondNode.Path[0]})

			if j < len(path) && j != i {
				secondNode.insertNode(path[j:], handle)
			} else {
				secondNode.handle = handle
				newNode = secondNode
			}
			return
		}
	}
	if n.nodeType == Param && path[0] != '/' {
		panic("一个节点下只应有一个参数节点")
	}
	//新添加的部分
	newNode = NewNode(path)
	newNode.handle = handle

	if n.children == nil {
		n.children = []*Node{newNode}
		n.priority = n.priority + 1
	} else {
		nodes := append(n.children, newNode)
		n.children = nodes
	}
	n.quickFind = n.quickFind + string([]byte{newNode.Path[0]})
	return
}

func NewNode(path string) (newNode *Node) {
	newNode = &Node{
		Path:     path,
		nodeType: Normal,
		priority: 1,
	}
	return
}

func CopyNode(n *Node) (newNode *Node) {
	newNode = &Node{
		Path:      n.Path,
		children:  n.children,
		nodeType:  n.nodeType,
		quickFind: n.quickFind,
		priority:  n.priority,
		handle:    n.handle,
	}
	return
}

func (n *Node) getValue(path string) {
	panic("implement me")
}

type TreeAction interface {
	addChild(path string, h interface{})
	getValue(path string)
}

func main() {
	n := Node{nodeType: Root, Path: "/"}
	n.AddChild("/test/asd/zxcv", 1)
	n.AddChild("/test/asd/zxc", 2)
	n.AddChild("/test/zxc/zxc", 3)
	n.AddChild("/test/:asd/zxc/:zxc", 5)
	//n.AddChild("/test/:zxc/",nil)
	n.printNode(0)
}

func (n *Node) printNode(i int) {
	for z := 0; z < i*3; z++ {
		fmt.Print(" ")
	}
	fmt.Print(n.Path+" ")
	if n.handle!=nil {
		fmt.Print(n.handle)
	}
	fmt.Println()
	for z := 0; z < len(n.children); z++ {
		n.children[z].printNode(i + 1)
	}
}

func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
