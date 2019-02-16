package begonia

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

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
	children  []Node
	nodeType  NodeType
	quickFind string
	priority  int
	handle    interface{}
}

func (n *Node) addChild(path string, h interface{}) {

	if n.Path == path {
		panic("两条相同路由")
	}

	//找到最长公共前缀
	minLength := min(len(path), len(n.Path))
	i := 0
	for i < minLength && path[i] == n.Path[i] {
		i++
	}

	//如果该节点不是最后一个节点
	if i == len(n.Path) {
		for j := 0; j < len(n.children); j++ {
			if n.quickFind[j] == path[0] {
				//TODO 判断新加入节点对于参数节点的影响合不合法
				n.children[j].addChild(path[i:], h)
				return
			}
		}
		panic("寻找公共前缀出了点问题 " + n.Path + " " + path)
	}

	n.insertNode(i, path, h)

}

func (n *Node) insertNode(i int, path string, handle interface{}) {
	//原来的后面那部分
	secondHalf := Node{
		Path:      n.Path[:i],
		children:  n.children,
		nodeType:  n.nodeType,
		quickFind: n.quickFind,
		priority:  n.priority,
		handle:    n.handle,
	}
	//新添加的部分
	newNode := Node{
		Path:     path[i:],
		nodeType: Normal,
		priority: 1,
	}
	if n.children == nil {
		n.children = []Node{secondHalf, newNode}
		n.priority = 2
	} else {
		append(n.children, secondHalf)
		append(n.children, newNode)
	}
}

func (n *Node) getValue(path string) {
	panic("implement me")
}

type TreeAction interface {
	addChild(path string, h interface{})
	getValue(path string)
}
