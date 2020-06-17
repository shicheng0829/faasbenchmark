package aliyun

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
	"path/filepath"
)

type Stack struct {
	*sls.Stack
}

func (aliyun *Aliyun) NewStack(stackPath string) (stack.Stack, error) {
	slsYamlDirPath := filepath.Join(stackPath, aliyun.Name())
	stack, err := sls.New(aliyun.Name(), slsYamlDirPath)

	if err != nil {
		return nil, err
	}

	stack.Opts["creds"] = aliyun.credentials

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {

	var functions []stack.Function

	funcs := s.ListFunctionsFromYaml()

	for _, f := range funcs {
		functions = append(functions, &Function{name: f.Name,
			handler:     f.Handler,
			description: f.Description,
			memorySize:  f.MemorySize,
			runtime:     f.Runtime})
	}

	return functions
}
