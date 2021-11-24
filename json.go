package snippets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func unknownType() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	fmt.Println(m["Parents"])  // 读取 json 内容
	fmt.Println(m["a"] == nil) // 判断 key 是否存在
}

func main() {
	person1 := Person{"张三", 24}
	bytes1, err := json.Marshal(&person1)
	if err == nil {
		// 返回字节数组 []byte
		fmt.Println("json.Marshal 编码结果:", string(bytes1)) // {"name":"张三","age":24}
	}

	str := `{"name": "李四", "age": 25}`
	// Unmarshal 接收字节数组参数
	bytes2 := []byte(str)
	var person2 Person
	if json.Unmarshal(bytes2, &person2) == nil {
		fmt.Println("json.Unmarshal 解码结果:", person2.Name, person2.Age) // 李四 25
	}

	// 使用 json.NewEncoder 编码
	person3 := Person{"王五", 30}
	// 编码结果暂存到 buffer
	bytes3 := new(bytes.Buffer)
	_ = json.NewEncoder(bytes3).Encode(person3)
	if err == nil {
		fmt.Print("json.NewEncoder 编码结果:", string(bytes3.Bytes())) // {"name":"王五","age":30}
	}

	// 使用 json.NewEncoder 解码
	str4 := `{"name":"赵六","age":28}`
	var person4 Person
	// 创建 string reader 作为参数
	err = json.NewDecoder(strings.NewReader(str4)).Decode(&person4)
	if err == nil {
		fmt.Println("json.NewDecoder 解码结果:", person4.Name, person4.Age) // 赵六 28
	}
}
