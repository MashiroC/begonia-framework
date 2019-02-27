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
	for i < minLength && path[i] == n.Path[i] && path[i] != ':' {
		i++
	}
	//如果该节点是最后一个节点
	if i != 0 && i == len(n.Path) {

		for j := 0; j < len(n.children); j++ {
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
	n.insertNode(path[i:], h)
}

func (n *Node) splitNode(i int) (newNode *Node) {

	if n.nodeType != Normal {
		panic("只能分割normal结点")
	}
	//原来的后面那部分
	newNode = CopyNode(n)
	newNode.Path = n.Path[i:]

	n.Path = n.Path[:i]
	n.priority = n.priority + 1
	n.quickFind = string([]byte{newNode.Path[0]})
	n.children = []*Node{newNode}

	return
}

func (n *Node) insertNode(path string, handle interface{}) (newNode *Node) {
	//具有参数节点
	for i := 0; i < len(path); i++ {
		// /asd/:asd/zxcasd/zxcv
		// /asd/:asd/zxchdgfh
		if path[i] == ':' {
			//fmt.Println(n.children[0])
			if i != 0 {
				firstNode := NewNode(path[:i])
				n.children = append(n.children, firstNode)
				n = firstNode
			}
			j := i
			for j < len(path) && path[j] != '/' {
				j++
			}
			secondNode := NewNode(path[i:j])
			secondNode.nodeType = Param
			secondNode.priority=0
			n.children = append(n.children,secondNode)

			if j < len(path) && j != i {
				n = secondNode
				path = path[j:]
				break
			} else {
				newNode = secondNode
				return
			}
		}
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
	n.AddChild("/test/asd/zxcv", nil)
	n.AddChild("/test/:asd/", nil)
	n.AddChild("/test/:zxc/asd",nil)
	n.printNode(0)
}

func (n *Node) printNode(i int) {
	for z := 0; z < i*3; z++ {
		fmt.Print(" ")
	}
	fmt.Println(n.Path)
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
