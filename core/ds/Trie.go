package ds

// trieNode 表示前缀树中的一个节点
type trieNode struct {
	children map[rune]*trieNode  // 子节点（按字符分支）
	ids      map[uint64]struct{} // 当前路径代表的前缀下所有 ID（例如 srcIP、协议等记录的 ID）
}

// Trie 是前缀树的根结构
type Trie struct {
	root *trieNode
}

// NewTrie 初始化一棵新的 Trie
func NewTrie() *Trie {
	return &Trie{root: &trieNode{
		children: make(map[rune]*trieNode),
		ids:      make(map[uint64]struct{}),
	}}
}

// Insert 将某个 key（例如字符串字段）插入到 Trie 中，并绑定对应的 ID
func (t *Trie) Insert(key string, id uint64) {
	node := t.root
	for _, r := range key {
		if _, ok := node.children[r]; !ok {
			node.children[r] = &trieNode{
				children: make(map[rune]*trieNode),
				ids:      make(map[uint64]struct{}),
			}
		}
		node = node.children[r]
		node.ids[id] = struct{}{} // 在每一层都记录这个 ID（支持前缀查询）
	}
}

// Delete 删除某个 key 对应的 id（仅从前缀路径上清除这个 id）
func (t *Trie) Delete(key string, id uint64) {
	node := t.root
	for _, r := range key {
		if n, ok := node.children[r]; ok {
			node = n
			delete(node.ids, id) // 从该层节点中移除 id
		} else {
			return // 不存在路径，忽略
		}
	}
}

// QueryPrefix 返回以 prefix 开头的所有 ID（即位于该前缀路径末尾的节点）
func (t *Trie) QueryPrefix(prefix string) map[uint64]struct{} {
	node := t.root
	for _, r := range prefix {
		if n, ok := node.children[r]; ok {
			node = n
		} else {
			return nil // 前缀不存在
		}
	}
	return node.ids // 返回该前缀路径末尾节点收集的所有 id
}
