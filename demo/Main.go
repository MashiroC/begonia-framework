package main

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.

type Params []Param

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
//获得请求中的参数
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
//对上述方法的包装，忽略掉是否查找到的返回值
func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

//方法对应的路由树
type methodTree struct {
	method string //方法
	root   *node  //路由树的根节点
}

type methodTrees []methodTree

//根据方法找路由树
func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

//计算路由内的路由参数的数量
func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= 255 {
		return 255
	}
	return uint8(n)
}

type nodeType uint8

const (
	static   nodeType = iota // default //普通节点
	root                     //根节点
	param                    //参数节点
	catchAll                 //TODO:还没看懂这个类型是啥
)

//路由树上的节点
type node struct {
	path      string   //节点的路由
	indices   string   //我称它为子节点快速索引，根据子节点的优先级，抽取节点首字母构成
	children  []*node  //子节点
	priority  uint32   //优先级 表示该节点下的(包括该节点本身和子节点)的路由数量
	nType     nodeType //节点类型
	maxParams uint8    //节点下路由的最大路由参数数量
	wildChild bool     //是否为一个路由参数的节点的父节点 如果是 那么该节点下一层不能有路由参数节点外其他节点存在
}

// increments priority of the given child and reorders if necessary.
//对子节点和节点的indices排序
//新的路由添加后 对子节点按照优先级 对children和indices进行排序
//优先级为该节点下的(包括该节点本身和子节点)的路由数量
//pos是新增的路由的节点下下标
//调用时如果是新添加的节点 那么新节点的path暂时为空 位于数组的最后一个元素
//如第一条路由是/helloworld 第二条是/hellogo
//n就是/hello节点 pos就是1
//如现有两条路由 /helloworld和/hellogo
//新添加一条路由 /hellowowowo
//这里会调用两次 第一次n是/hello pos是0 第二次调用n是wo pos是1
func (n *node) incrementChildPrio(pos int) int {

	//增加优先级
	n.children[pos].priority++
	prio := n.children[pos].priority

	// adjust position (move to front)
	//因为只变动了一个节点的优先级 其他节点是从大到小排序好的 使用冒泡的方法排序子节点数组
	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < prio {
		// swap node positions
		n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]
		newPos--
	}

	//排序好之后修改快速索引字符串
	// build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // unchanged prefix, might be empty
			n.indices[pos:pos+1] + // the index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // rest without char at 'pos'
	}

	//返回排序好之后新节点的下标
	return newPos
}

// addRoute adds a node with the given handle to the path.
// Not concurrency-safe!
//添加handle到路由上 非线程安全
func (n *node) addRoute(path string) {
	fullPath := path
	n.priority++
	numParams := countParams(path)

	// non-empty tree
	//如果树非空
	if len(n.path) > 0 || len(n.children) > 0 {
	walk:
		//循环取非公共的节点
		for {
			// Update maxParams of the current node
			//更新节点下面的路由的最大参数数量
			if numParams > n.maxParams {
				n.maxParams = numParams
			}

			// Find the longest common prefix.
			// This also implies that the common prefix contains no ':' or '*'
			// since the existing key can't contain those chars.
			//找到公共最长前缀的下标
			i := 0
			max := min(len(path), len(n.path))
			for i < max && path[i] == n.path[i] {
				i++
			}

			//分割路由，抽取公共前缀替代这个节点，剩余部分作为一个子节点
			// Split edge
			if i < len(n.path) {
				//这个是将原来的节点抽取出前缀之后的新节点
				//比如原本的路由/helloworld 新添加的路由是/hellogo
				//这个child是world
				child := node{
					path:      n.path[i:],
					wildChild: n.wildChild,
					indices:   n.indices,
					children:  n.children,
					priority:  n.priority - 1,
				}

				// Update maxParams (max of all children)
				//更新节点下的最大路由参数数量
				for i := range child.children {
					if child.children[i].maxParams > child.maxParams {
						child.maxParams = child.children[i].maxParams
					}
				}

				//抽离公共前缀后 创建子节点
				// 例中的world节点连接在/hello节点下
				n.children = []*node{&child}
				// []byte for proper unicode char conversion, see #65

				//将子节点的首字母放在父节点的indices下
				//上例中的'w'
				n.indices = string([]byte{n.path[i]})
				//将本节点的path改为公共前缀
				n.path = path[:i]
				n.wildChild = false
			}

			// Make new node a child of this node
			//如果公共前缀和新添加路由的path一样
			// 抽取非公共部分 创建新节点
			if i < len(path) {
				//path非公共部分 比如新路由/hellogo的go
				path = path[i:]

				//如果本节点是一个参数节点的父节点
				//就让n指向参数节点
				if n.wildChild {
					n = n.children[0]
					n.priority++

					//如果后面还有参数节点
					// Update maxParams of the child node
					if numParams > n.maxParams {
						n.maxParams = numParams
					}
					numParams--

					// Check if the wildcard matches
					//如果上面没有continue 那么路由有问题 报错
					//这里有个bug
					//如果添加的两条路由为/aaa/:bbb/ccc 和 /aaa/:bbb/ddd/:eee/fff 会panic出来 反之不会
					//如果上述的路由第二条变成/aaa/:bbb/ddd/:eee/fff/:ggg/hhh 则不会panic出来
					//TODO:好像没这个问题了 好像是我的问题
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] {
						// check for longer wildcard, e.g. :name and :names
						if len(n.path) >= len(path) || path[len(n.path)] == '/' {
							continue walk
						}
					}

					pathSeg := path
					if n.nType != catchAll {
						pathSeg = strings.SplitN(path, "/", 2)[0]
					}
					prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
					panic("'" + pathSeg +
						"' in new path '" + fullPath +
						"' conflicts with existing wildcard '" + n.path +
						"' in existing prefix '" + prefix +
						"'")
				}

				//后面没有参数节点了
				c := path[0]

				//如果本节点是参数节点
				// TODO:这个if没看懂 如果说一个节点是参数节点 那么c应该是':'但是这里是'/'才行 不太懂
				// slash after param
				if n.nType == param && c == '/' && len(n.children) == 1 {
					n = n.children[0]
					n.priority++
					continue walk
				}

				// Check if a child with the next path byte exists
				//如果和子节点首字母一样，那么具有公共部分，对increment和子节点重新排序 然后取子节点 继续循环
				for i := 0; i < len(n.indices); i++ {
					if c == n.indices[i] {
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue walk
					}
				}

				// Otherwise insert it
				//如果当前节点不是参数节点
				if c != ':' && c != '*' {
					// []byte for proper unicode char conversion, see #65
					//更新快速查找字符串
					n.indices += string([]byte{c})
					//新的节点
					//上例中/hellogo中的go
					child := &node{
						maxParams: numParams,
					}
					n.children = append(n.children, child)
					n.incrementChildPrio(len(n.indices) - 1)
					n = child
				}
				n.insertChild(numParams, path, fullPath)
				return

			} else if i == len(path) { // Make node a (in-path) leaf
				//向该节点添加handlesChain
			}
			return
		}
	} else { // Empty tree
		//如果是一棵空的树 那么就以刚添加的节点为root节点
		n.insertChild(numParams, path, fullPath)
		n.nType = root
	}
}

//向节点插入子节点数据
//这里的n就是新节点了，而不是新节点的父节点
//就例如原本有一条/helloworld路由 新添加一条/hellogo路由
//这里的n不是/hello这个节点 而是在/hello下面新开的一个空节点
//(自我感觉这种做法有点怪)
//TODO:正在读
func (n *node) insertChild(numParams uint8, path string, fullPath string) {

	//root.addRoute("/hellogo")
	//root.addRoute("/helloworld")
	//root.addRoute("/hellowowowo")
	//root.addRoute("/hellowowoaa")
	//root.addRoute("/hellogin")
	//root.addRoute("/hellogs")

	var offset int // already handled bytes of the path

	fmt.Println("fullPath:"+fullPath+" path:"+path+" n:"+n.path)
	// find prefix until first wildcard (beginning with ':' or '*')
	//这个loop是解析参数路由
	for i, max := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// find wildcard end (either '/' or path end)
		end := i + 1
		for end < max && path[end] != '/' {
			switch path[end] {
			// the wildcard name must not contain ':' and '*'
			case ':', '*':
				panic("only one wildcard per path segment is allowed, has: '" +
					path[i:] + "' in path '" + fullPath + "'")
			default:
				end++
			}
		}

		// check if this Node existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > 0 {
			panic("wildcard route '" + path[i:end] +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if end-i < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if c == ':' { // param
			// split path at the beginning of the wildcard
			if i > 0 {
				n.path = path[offset:i]
				offset = i
			}

			child := &node{
				nType:     param,
				maxParams: numParams,
			}
			n.children = []*node{child}
			n.wildChild = true
			n = child
			n.priority++
			numParams--

			// if the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max {
				n.path = path[offset:end]
				offset = end

				child := &node{
					maxParams: numParams,
					priority:  1,
				}
				n.children = []*node{child}
				n = child
			}

		} else { // catchAll
			if end != max || numParams > 1 {
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			}

			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			}

			// currently fixed width 1 for '/'
			i--
			if path[i] != '/' {
				panic("no / before catch-all in path '" + fullPath + "'")
			}

			n.path = path[offset:i]

			// first node: catchAll node with empty path
			child := &node{
				wildChild: true,
				nType:     catchAll,
				maxParams: 1,
			}
			n.children = []*node{child}
			n.indices = string(path[i])
			n = child
			n.priority++

			// second node: node holding the variable
			child = &node{
				path:      path[i:],
				nType:     catchAll,
				maxParams: 1,
				priority:  1,
			}
			n.children = []*node{child}

			return
		}
	}

	// insert remaining path part and handle to the leaf
	n.path = path[offset:]
}

// getValue returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getValue(path string, po Params, unescape bool) (p Params, tsr bool) {
	p = po
walk: // Outer loop for walking the tree
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					c := path[0]
					for i := 0; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = path == "/"
					return
				}

				// handle wildcard child
				n = n.children[0]
				switch n.nType {
				case param:
					// find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// save param value
					if cap(p) < int(n.maxParams) {
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[1:]
					val := path[:end]
					if unescape {
						var err error
						if p[i].Value, err = url.QueryUnescape(val); err != nil {
							p[i].Value = val // fallback, in case of error
						}
					} else {
						p[i].Value = val
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						tsr = len(path) == end+1
						return
					}

					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						tsr = n.path == "/"
					}

					return

				case catchAll:
					// save param value
					if cap(p) < int(n.maxParams) {
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[2:]
					if unescape {
						var err error
						if p[i].Value, err = url.QueryUnescape(path); err != nil {
							p[i].Value = path // fallback, in case of error
						}
					} else {
						p[i].Value = path
					}

					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == n.path {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.

			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == '/' {
					n = n.children[i]
					tsr = len(n.path) == 1 || n.nType == catchAll
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
				path == n.path[:len(n.path)-1])
		return
	}
}

// findCaseInsensitivePath makes a case-insensitive lookup of the given path and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) {
	ciPath = make([]byte, 0, len(path)+1) // preallocate enough memory

	// Outer loop for walking the tree
	for len(path) >= len(n.path) && strings.ToLower(path[:len(n.path)]) == strings.ToLower(n.path) {
		path = path[len(n.path):]
		ciPath = append(ciPath, n.path...)

		if len(path) > 0 {
			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.wildChild {
				r := unicode.ToLower(rune(path[0]))
				for i, index := range n.indices {
					// must use recursive approach since both index and
					// ToLower(index) could exist. We must check both.
					if r == unicode.ToLower(index) {
						out, found := n.children[i].findCaseInsensitivePath(path, fixTrailingSlash)
						if found {
							return append(ciPath, out...), true
						}
					}
				}

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				found = fixTrailingSlash && path == "/"
				return
			}

			n = n.children[0]
			switch n.nType {
			case param:
				// find param end (either '/' or path end)
				k := 0
				for k < len(path) && path[k] != '/' {
					k++
				}

				// add param value to case insensitive path
				ciPath = append(ciPath, path[:k]...)

				// we need to go deeper!
				if k < len(path) {
					if len(n.children) > 0 {
						path = path[k:]
						n = n.children[0]
						continue
					}

					// ... but we can't
					if fixTrailingSlash && len(path) == k+1 {
						return ciPath, true
					}
					return
				}

				//handles不为nil
				if fixTrailingSlash && len(n.children) == 1 {
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.children[0]
					if n.path == "/" {
						return append(ciPath, '/'), true
					}
				}
				return

			case catchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			}
		} else {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash {
				for i := 0; i < len(n.indices); i++ {
					if n.indices[i] == '/' {
						n = n.children[i]
						if len(n.path) == 1 ||
							(n.nType == catchAll) {
							return append(ciPath, '/'), true
						}
						return
					}
				}
			}
			return
		}
	}

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash {
		if path == "/" {
			return ciPath, true
		}
		if len(path)+1 == len(n.path) && n.path[len(path)] == '/' &&
			strings.ToLower(path) == strings.ToLower(n.path[:len(path)]) {
			return append(ciPath, n.path...), true
		}
	}
	return
}

func main() {
	root := &node{nType: root}
	//root.addRoute("/sdfbvcx/zcv",nil)
	//root.addRoute("/zxcvzxcv/asd",nil)
	//root.addRoute("/zxcv/:zxc/a",nil)
	//root.addRoute("/cvb/:zxc",nil)
	//如果添加的两条路由为/aaa/:bbb/ccc 和 /aaa/:bbb/ddd/:eee/fff 会panic出来
	//root.addRoute("/aaa/:bbb/ccc")

	//root.addRoute("/zxcv/a", nil)

	//root.addRoute("/aaa/:bbb/ccc/:eee/fff")
	root.addRoute("/hellogo")
	root.addRoute("/helloworld")
	root.addRoute("/hellowowowo")
	root.addRoute("/hellowowoaa")
	root.addRoute("/hellogin")
	root.addRoute("/hellogs")
	//fmt.Println(root.children[0].children[0].path)
	//fmt.Println(root.wildChild)
	//root.addRoute("/zxcv/baefds",nil)
	//fmt.Println(root.children[0].path)
	//fmt.Println()
	//fmt.Println(root.wildChild)
	//root.addRoute("/t",nil)
	//root.addRoute("/asdasd/zxc",nil)
	//fmt.Println(root.children[0].path)
	//fmt.Println(root.children[1].path)
	//fmt.Println(root.children[2].path)
}
