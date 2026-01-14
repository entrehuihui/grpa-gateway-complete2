package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

var grpcPath = "./service/myrpc/proto"

func main() {
	// 读取grpc文件夹文件
	grpcFileInfo, err := os.ReadDir(grpcPath)
	if err != nil {
		log.Fatalln(err)
	}

	grpcFileList := make([]fs.DirEntry, 0)
	for _, v := range grpcFileInfo {
		if strings.Contains(v.Name(), "_grpc.pb.go") {
			grpcFileList = append(grpcFileList, v)
		}
	}

	// 读取mod包名
	modFile, err := os.Open("./go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer modFile.Close()
	scanner := bufio.NewScanner(modFile)
	var modName string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "module") {
			modName = strings.Split(line, " ")[1]
			break
		}
	}

	// 读取函数方法
	for _, v := range grpcFileList {
		filename := v.Name()
		filename = filename[:len(filename)-11]

		readFuncMethod(filename, modName)
	}

	readServiceFile(grpcFileList, modName)
	readServerFile(grpcFileList, modName)
}

// 处理sServerF文件
func readServerFile(grpcFileList []fs.DirEntry, modName string) {
	path := "./service/myrpc/server.go"
	body, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(body)
	bodyList := strings.Split(bodyString, "// 注入GRPC服务")
	if len(bodyList) != 3 {
		log.Fatal("server文件有误,识别不出两个<// 注入GRPC服务>标识")
	}
	bodyGwList := strings.Split(bodyList[2], "// 注入GW服务")
	if len(bodyGwList) != 3 {
		log.Fatal("server文件有误,识别不出两个<// 注入GW服务>标识")
	}
	bodyList = append(bodyList[:2], bodyGwList...)

	// 注入GW服务
	// 注入GRPC服务
	gwData := "// 注入GW服务\n"
	grpcData := "// 注入GRPC服务\n"
	for _, v := range grpcFileList {
		name := v.Name()
		name = name[:len(name)-11]
		name = strings.ToUpper(name[:1]) + name[1:]
		regisertGw := fmt.Sprintf(`	err = proto.Register%sHandlerFromEndpoint(ctx, mux, cfg.grpcPort, opts)
	if err != nil {
		log.Fatal("启动GW错误:", err)
	}`, name)
		gwData += regisertGw + "\n"
		regisertGrpc := fmt.Sprintf(`	proto.Register%sServer(gs, server)`, name)
		grpcData += regisertGrpc + "\n"
	}

	bodyList[1] = grpcData + "	// 注入GRPC服务"
	bodyList[3] = gwData + "	// 注入GW服务"
	bodyString = strings.Join(bodyList, "")

	// 判断是否已经引用包名
	if !strings.Contains(bodyString, fmt.Sprintf(`"%s/service/myrpc/middleware"`, modName)) {
		importStr := fmt.Sprintf(`"%s/service/myrpc/middleware"`, modName)
		bodyString = strings.ReplaceAll(bodyString, "import (", "import (\n\t"+importStr)
	}
	if !strings.Contains(bodyString, fmt.Sprintf(`"%s/service/myrpc/proto"`, modName)) {
		importStr := fmt.Sprintf(`"%s/service/myrpc/proto"`, modName)
		bodyString = strings.ReplaceAll(bodyString, "import (", "import (\n\t"+importStr)
	}
	if !strings.Contains(bodyString, fmt.Sprintf(`"%s/service/myrpc/service"`, modName)) {
		importStr := fmt.Sprintf(`"%s/service/myrpc/service"`, modName)
		bodyString = strings.ReplaceAll(bodyString, "import (", "import (\n\t"+importStr)
	}

	log.Println("写入service服务")
	os.WriteFile(path, []byte(bodyString), 0666)
}

// 处理service文件
func readServiceFile(grpcFileList []fs.DirEntry, modName string) {
	if len(grpcFileList) == 0 {
		return
	}

	path := "./service/myrpc/service/service.go"
	body := `package service

import "%s/service/myrpc/proto"

// Service .
type Service struct {
	// ##继承
%s	// ##继承
}

// NewService .
func NewService() *Service {
	s := new(Service)
	return s
}
	`
	data := ""
	for _, v := range grpcFileList {
		name := v.Name()
		name = name[:len(name)-11]
		name = strings.ToUpper(name[:1]) + name[1:]
		name = "\tproto.Unimplemented" + name + "Server\n"
		data += name
	}
	body = fmt.Sprintf(body, modName, data)
	log.Println("写入service服务")
	os.WriteFile(path, []byte(body), 0666)
}

type FuncInfo struct {
	FuncName   string
	DetailInfo string
	FuncNames  string
}

func readFuncMethod(fileName string, modName string) {
	file, err := os.Open(grpcPath + "/" + fileName + "_grpc.pb.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	clientName := fmt.Sprintf("type %s%sClient interface", strings.ToUpper(fileName[:1]), fileName[1:])
	scanner := bufio.NewScanner(file)
	searchStart := false
	detailInfo := ""
	funcInfo := make([]FuncInfo, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, clientName) {
			searchStart = true
			continue
		}
		if line == "}" {
			break
		}
		if searchStart {
			if strings.Contains(line, "//") {
				detailInfo = strings.ReplaceAll(line, "\t", "")
				continue
			}
			if strings.Contains(line, "ctx context.Context") {
				funcName := strings.ReplaceAll(line, ", opts ...grpc.CallOption", "")
				funcName = strings.ReplaceAll(funcName, "\t", "")
				funcName = strings.ReplaceAll(funcName, "*", "*proto.")
				funcNames := strings.ReplaceAll(strings.Split(funcName, "ctx context.Context,")[0], "\t", "") + "ctx context.Context"
				funcInfo = append(funcInfo, FuncInfo{
					FuncName:   funcName,
					DetailInfo: detailInfo,
					FuncNames:  funcNames,
				})
				detailInfo = ""
			}
		}

	}

	// 处理operate文件夹
	operateFile := "./service/operate" + "/" + fileName + ".go"
	operateBody := readService(operateFile)
	checkOperateFunc(operateFile, string(operateBody), funcInfo, modName)

	// 处理service文件夹
	serceFile := "./service/myrpc/service" + "/" + fileName + ".go"
	serviceBody := readService(serceFile)
	checkServiceFunc(serceFile, string(serviceBody), funcInfo, modName)
}

func checkServiceFunc(fileName string, body string, funcList []FuncInfo, modName string) {
	if body == "" {
		// 如果是空的  直接写入文件
		emptyBody(fileName, funcList)
		return
	}
	// 如果不是空
	funcInfo := make([]FuncInfo, 0)
	for _, v := range funcList {
		if !strings.Contains(body, v.FuncNames) {
			funcInfo = append(funcInfo, v)
		}
	}
	if len(funcInfo) == 0 {
		return
	}
	for _, v := range funcInfo {
		body += createFunc(v)
	}
	// 写入文件
	os.WriteFile(fileName, []byte(body), 0666)
}

func emptyBody(fileName string, funcInfo []FuncInfo) {
	body := `package service

import (
	"context"

	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
	"github.com/entrehuihui/grpa-gateway-complete2/service/operate"
)
`
	for _, v := range funcInfo {
		body += createFunc(v)
	}

	os.WriteFile(fileName, []byte(body), 0777)
}

func createFunc(v FuncInfo) string {
	reFuncs := strings.Split(v.FuncName, "(ctx context.Context,")[0]
	body := fmt.Sprintf(`
%s
func (s Service) %s {
	return operate.%s(ctx, in)
}
`, v.DetailInfo, v.FuncName, reFuncs)
	return body
}

func checkOperateFunc(fileName string, body string, funcList []FuncInfo, modName string) {
	if body == "" {
		// 如果是空的  直接写入文件
		emptyBodyOperate(fileName, funcList)
		return
	}
	// 如果不是空
	funcInfo := make([]FuncInfo, 0)
	for _, v := range funcList {
		if !strings.Contains(body, v.FuncNames) {
			funcInfo = append(funcInfo, v)
		}
	}
	log.Println("funcInfo========>>", len(funcInfo))
	if len(funcInfo) == 0 {
		return
	}
	for _, v := range funcInfo {
		body += "\n" + createOpetateFunc(v)
	}
	// 写入文件
	os.WriteFile(fileName, []byte(body), 0666)
}

func emptyBodyOperate(fileName string, funcInfo []FuncInfo) {
	body := `package operate

import (
	"context"

	"github.com/entrehuihui/grpa-gateway-complete2/service/myrpc/proto"
)
`
	for _, v := range funcInfo {
		body += createOpetateFunc(v)
	}

	os.WriteFile(fileName, []byte(body), 0777)
}

func createOpetateFunc(v FuncInfo) string {
	body := fmt.Sprintf(`
%s
func %s {
	return nil, nil
}
`, v.DetailInfo, v.FuncName)

	return body
}

func readService(fileName string) []byte {
	if !osStat(fileName) {
		// 如果不存在 创建文件
		return make([]byte, 0)
	}

	body, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
func osStat(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Fatal(err)
	return false
}
