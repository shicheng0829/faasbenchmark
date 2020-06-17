package tencent

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
	"path/filepath"
)

type Stack struct {
	*sls.Stack
}

func (tencent *Tencent) NewStack(stackPath string) (stack.Stack, error) {
	slsYamlDirPath := filepath.Join(stackPath, tencent.Name())
	stack, err := sls.New(tencent.Name(), slsYamlDirPath)

	if err != nil {
		return nil, err
	}

	stack.Opts["creds"] = tencent.credentials

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
