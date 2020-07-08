package aliyun

import (
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/nuweba/faasbenchmark/stack"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Stack struct {
	stackPath string
	serviceName string
	client *fc.Client
}

func (aliyun *Aliyun) NewStack(stackPath string) (stack.Stack, error) {
	client, _ := fc.NewClient(os.Getenv("ENDPOINT"), "2016-08-15", os.Getenv("ACCESS_KEY_ID"), os.Getenv("ACCESS_KEY_SECRET"))
	s := Stack{stackPath: stackPath, client: client, serviceName: "faasbenchmark"}
	return s, nil
}

func (s Stack) DeployStack() error {
	fmt.Println("Creating service")
	createServiceOutput, err := s.client.CreateService(fc.NewCreateServiceInput().
		WithServiceName(s.serviceName).
		WithDescription("this is a smoke test for go sdk"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if createServiceOutput != nil {
		fmt.Printf("CreateService response: %s \n", createServiceOutput)
	}
	fmt.Println("Creating function1")
	createFunctionInput1 := fc.NewCreateFunctionInput(s.serviceName).WithFunctionName("testf1").
		WithDescription("go sdk test function").
		WithHandler("index.handler").WithRuntime("nodejs12").
		WithCode(fc.NewCode().WithFiles(path.Join(s.stackPath,"aliyun" ,"index.js"))).
		WithTimeout(5).WithMemorySize(512)
	createFunctionOutput, err := s.client.CreateFunction(createFunctionInput1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("CreateFunction response: %s \n", createFunctionOutput)
	}
	createHTTPTriggerConfigOutput := fc.NewHTTPTriggerConfig().WithMethods("GET").WithAuthType(fc.AuthAnonymous)
	createTriggerInput := fc.NewCreateTriggerInput(s.serviceName, "testf1").WithTriggerType(fc.TRIGGER_TYPE_HTTP).
		WithTriggerName("httptrigger").
		WithTriggerConfig(createHTTPTriggerConfigOutput)
	createTriggerOutput, err := s.client.CreateTrigger(createTriggerInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("CreateTrigger response: %s \n", createTriggerOutput)
	}
	return err
}

func (s Stack) RemoveStack() error {
	fmt.Println("Deleting functions")
	listFunctionsOutput, err := s.client.ListFunctions(fc.NewListFunctionsInput(s.serviceName).WithLimit(10))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("ListFunctions response: %s \n", listFunctionsOutput)
		for _, fuc := range listFunctionsOutput.Functions {
			fmt.Printf("Deleting function %s \n", *fuc.FunctionName)
			if output, err := s.client.DeleteTrigger(fc.NewDeleteTriggerInput(s.serviceName, *fuc.FunctionName, "httptrigger")); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Printf("DeleteFunction response: %s \n", output)
			}
			if output, err := s.client.DeleteFunction(fc.NewDeleteFunctionInput(s.serviceName, *fuc.FunctionName)); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Printf("DeleteFunction response: %s \n", output)
			}

		}
	}
	fmt.Println("Deleting service")
	deleteServiceOutput, err := s.client.DeleteService(fc.NewDeleteServiceInput(s.serviceName))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("DeleteService response: %s \n", deleteServiceOutput)
	}
	return err
}

func (s Stack) StackId() string {
	_, fileName := filepath.Split(s.stackPath)
	if fileName == "sleep" {
		fileName += "func"
	}
	return fileName
}

func (s Stack) Project() string {

	return ""
}

func (s Stack) Stage() string {
	return "test"
}

func (s Stack) ListFunctions() []stack.Function {
	var functions []stack.Function
	fmt.Println("Listing functions")
	listFunctionsOutput, err := s.client.ListFunctions(fc.NewListFunctionsInput(s.serviceName).WithPrefix("test"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("ListFunctions response: %s \n", listFunctionsOutput)
	}
	for _, v := range listFunctionsOutput.Functions {
		functions = append(functions, &Function{name: *v.FunctionName,
			handler:     *v.Handler,
			description: *v.Description,
			memorySize:  strconv.Itoa(int(*v.MemorySize)),
			runtime:     *v.Runtime} )

	}
	return functions
}
