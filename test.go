// package main

// var domainDict map[string]domainNode

// type Node struct {
// 	a int
// }

// type domainNode struct {
// 	time int
// 	urls map[string]Node
// }

// func main() {
// 	dict := make(map[string]domainNode)
// 	dict["google.com"] = domainNode{1, make(map[string]Node)}
// 	dict["google.com"].urls["1"] = Node{1}

// 	dict["baidu.com"] = domainNode{1, make(map[string]Node)}
// 	dict["baidu.com"].urls["2"] = Node{2}

// 	println("~~~~~~~~~~~~~~~~~~")
// 	for host, v := range dict {
// 		println("HOST", host)
// 		for k, s := range v.urls {
// 			println("url", k, s.a)
// 		}
// 	}
// 	println("~~~~~~~~~~~~~~~~~~")
// }
