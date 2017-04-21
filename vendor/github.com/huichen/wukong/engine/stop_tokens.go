package engine

import (
	"bufio"
	"log"
	"os"
)

type StopTokens struct {
	stopTokens map[string]bool
}

// 从stopTokenFile中读入停用词，一个词一行
// 文档索引建立时会跳过这些停用词
func (st *StopTokens) Init(stopTokenFile string) {
	st.stopTokens = make(map[string]bool)
	if stopTokenFile == "" {
		return
	}

	file, err := os.Open(stopTokenFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			st.stopTokens[text] = true
		}
	}

}

func (st *StopTokens) IsStopToken(token string) bool {
	_, found := st.stopTokens[token]
	return found
}
