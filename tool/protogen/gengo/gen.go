package gengo

import (
	"fmt"
	"github.com/bobwong89757/protoplus/codegen"
	"github.com/bobwong89757/protoplus/gen"
)

func GenGo(ctx *gen.Context) error {

	gen := codegen.NewCodeGen("cmgo").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(FuncMap).
		ParseTemplate(goCodeTemplate, ctx).
		FormatGoCode()

	if gen.Error() != nil {
		fmt.Println(string(gen.Code()))
		return gen.Error()
	}

	return gen.WriteOutputFile(ctx.OutputFileName).Error()
}
