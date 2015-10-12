// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bind

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"math"
	"strings"
)

// TODO(hyangah): error code/domain propagation

type objcGen struct {
	*printer
	fset *token.FileSet
	pkg  *types.Package
	err  ErrorList

	prefix string // prefix arg passed by flag.

	// fields set by init.
	pkgName    string
	namePrefix string
	funcs      []*types.Func
	names      []*types.TypeName
	constants  []*types.Const
	vars       []*types.Var
}

func (g *objcGen) init() {
	g.pkgName = g.pkg.Name()
	g.namePrefix = g.prefix + strings.Title(g.pkgName)
	g.funcs = nil
	g.names = nil

	scope := g.pkg.Scope()
	hasExported := false
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}
		hasExported = true
		switch obj := obj.(type) {
		case *types.Func:
			if isCallable(obj) {
				g.funcs = append(g.funcs, obj)
			}
		case *types.TypeName:
			g.names = append(g.names, obj)
		case *types.Const:
			if _, ok := obj.Type().(*types.Basic); !ok {
				g.errorf("unsupported exported const for %s: %T", obj.Name(), obj)
				continue
			}
			g.constants = append(g.constants, obj)
		case *types.Var:
			g.vars = append(g.vars, obj)
		default:
			g.errorf("unsupported exported type for %s: %T", obj.Name(), obj)
		}
	}
	if !hasExported {
		g.errorf("no exported names in the package %q", g.pkg.Path())
	}
}

const objcPreamble = `// Objective-C API for talking to %[1]s Go package.
//   gobind %[2]s %[3]s
//
// File is generated by gobind. Do not edit.

`

func (g *objcGen) genH() error {
	g.init()

	g.Printf(objcPreamble, g.pkg.Path(), g.gobindOpts(), g.pkg.Path())
	g.Printf("#ifndef __Go%s_H__\n", strings.Title(g.pkgName))
	g.Printf("#define __Go%s_H__\n", strings.Title(g.pkgName))
	g.Printf("\n")
	g.Printf("#include <Foundation/Foundation.h>")
	g.Printf("\n\n")

	// @class names
	for _, obj := range g.names {
		named := obj.Type().(*types.Named)
		switch t := named.Underlying().(type) {
		case *types.Struct:
			g.Printf("@class %s%s;\n\n", g.namePrefix, obj.Name())
		case *types.Interface:
			if !makeIfaceSummary(t).implementable {
				g.Printf("@class %s%s;\n\n", g.namePrefix, obj.Name())
			}
		}
	}

	// @interfaces
	for _, obj := range g.names {
		named := obj.Type().(*types.Named)
		switch t := named.Underlying().(type) {
		case *types.Struct:
			g.genStructH(obj, t)
			g.Printf("\n")
		case *types.Interface:
			g.genInterfaceH(obj, t)
			g.Printf("\n")
		}
	}

	// const
	// TODO: prefix with k?, or use a class method?
	for _, obj := range g.constants {
		switch b := obj.Type().(*types.Basic); b.Kind() {
		case types.String, types.UntypedString:
			g.Printf("FOUNDATION_EXPORT NSString* const %s%s;\n", g.namePrefix, obj.Name())
		default:
			g.Printf("FOUNDATION_EXPORT const %s %s%s;\n", g.objcType(obj.Type()), g.namePrefix, obj.Name())
		}
	}
	if len(g.constants) > 0 {
		g.Printf("\n")
	}

	// var
	if len(g.vars) > 0 {
		g.Printf("@interface %s : NSObject \n", g.namePrefix)
		for _, obj := range g.vars {
			objcType := g.objcType(obj.Type())
			g.Printf("+ (%s) %s;\n", objcType, obj.Name())
			g.Printf("+ (void) set%s:(%s)v;\n", obj.Name(), objcType)
		}
		g.Printf("@end\n")
	}

	// static functions.
	for _, obj := range g.funcs {
		g.genFuncH(obj)
		g.Printf("\n")
	}

	// declare all named types first.
	g.Printf("#endif\n")

	if len(g.err) > 0 {
		return g.err
	}
	return nil
}

func (g *objcGen) gobindOpts() string {
	opts := []string{"-lang=objc"}
	if g.prefix != "Go" {
		opts = append(opts, "-prefix="+g.prefix)
	}
	return strings.Join(opts, " ")
}

func (g *objcGen) genM() error {
	g.init()

	g.Printf(objcPreamble, g.pkg.Path(), g.gobindOpts(), g.pkg.Path())
	g.Printf("#include %q\n", g.namePrefix+".h")
	g.Printf("#include <Foundation/Foundation.h>\n")
	g.Printf("#include \"seq.h\"\n")
	g.Printf("\n")
	g.Printf("static NSString* errDomain = @\"go.%s\";\n", g.pkg.Path())
	g.Printf("\n")

	g.Printf("@protocol goSeqRefInterface\n")
	g.Printf("-(GoSeqRef*) ref;\n")
	g.Printf("@end\n")
	g.Printf("\n")

	g.Printf("#define _DESCRIPTOR_ %q\n\n", g.pkgName)
	for i, obj := range g.funcs {
		g.Printf("#define _CALL_%s_ %d\n", obj.Name(), i+1)
	}
	g.Printf("\n")

	// struct, interface.
	var interfaces []*types.TypeName
	for _, obj := range g.names {
		named := obj.Type().(*types.Named)
		switch t := named.Underlying().(type) {
		case *types.Struct:
			g.genStructM(obj, t)
		case *types.Interface:
			if g.genInterfaceM(obj, t) {
				interfaces = append(interfaces, obj)
			}
		}
		g.Printf("\n")
	}

	// const
	for _, o := range g.constants {
		g.genConstM(o)
	}
	if len(g.constants) > 0 {
		g.Printf("\n")
	}

	// vars
	if len(g.vars) > 0 {
		g.Printf("@implementation %s\n", g.namePrefix)
		for _, o := range g.vars {
			g.genVarM(o)
		}
		g.Printf("@end\n\n")
	}

	// global functions.
	for _, obj := range g.funcs {
		g.genFuncM(obj)
		g.Printf("\n")
	}

	// register proxy functions.
	if len(interfaces) > 0 {
		g.Printf("__attribute__((constructor)) static void init() {\n")
		g.Indent()
		for _, obj := range interfaces {
			g.Printf("go_seq_register_proxy(\"go.%s.%s\", proxy%s%s);\n", g.pkgName, obj.Name(), g.namePrefix, obj.Name())
		}
		g.Outdent()
		g.Printf("}\n")
	}

	if len(g.err) > 0 {
		return g.err
	}

	return nil
}

func (g *objcGen) genVarM(o *types.Var) {
	varDesc := fmt.Sprintf("%q", g.pkg.Name()+"."+o.Name())
	objcType := g.objcType(o.Type())

	// setter
	s1 := &funcSummary{
		name:   "set" + o.Name(),
		ret:    "void",
		params: []paramInfo{{typ: o.Type(), name: "v"}},
	}
	g.Printf("+ (void) %s:(%s)v {\n", s1.name, objcType)
	g.Indent()
	g.genFunc(varDesc, "1", s1, false) // false: not instance method.
	g.Outdent()
	g.Printf("}\n\n")

	// getter
	s2 := &funcSummary{
		name:      o.Name(),
		ret:       objcType,
		retParams: []paramInfo{{typ: o.Type(), name: "ret"}},
	}
	g.Printf("+ (%s) %s {\n", s2.ret, s2.name)
	g.Indent()
	g.genFunc(varDesc, "2", s2, false)
	g.Outdent()
	g.Printf("}\n\n")
}

func (g *objcGen) genConstM(o *types.Const) {
	cName := fmt.Sprintf("%s%s", g.namePrefix, o.Name())
	cType := g.objcType(o.Type())

	switch b := o.Type().(*types.Basic); b.Kind() {
	case types.Bool, types.UntypedBool:
		v := "NO"
		if constant.BoolVal(o.Val()) {
			v = "YES"
		}
		g.Printf("const BOOL %s = %s;\n", cName, v)

	case types.String, types.UntypedString:
		g.Printf("NSString* const %s = @%s;\n", cName, o.Val())

	case types.Int, types.Int8, types.Int16, types.Int32:
		g.Printf("const %s %s = %s;\n", cType, cName, o.Val())

	case types.Int64, types.UntypedInt:
		i, exact := constant.Int64Val(o.Val())
		if !exact {
			g.errorf("const value %s for %s cannot be represented as %s", o.Val(), o.Name(), cType)
			return
		}
		if i == math.MinInt64 {
			// -9223372036854775808LL does not work because 922337203685477508 is
			// larger than max int64.
			g.Printf("const int64_t %s = %dLL-1;\n", cName, i+1)
		} else {
			g.Printf("const int64_t %s = %dLL;\n", cName, i)
		}

	case types.Float32, types.Float64, types.UntypedFloat:
		f, _ := constant.Float64Val(o.Val())
		if math.IsInf(f, 0) || math.Abs(f) > math.MaxFloat64 {
			g.errorf("const value %s for %s cannot be represented as double", o.Val(), o.Name())
			return
		}
		g.Printf("const %s %s = %g;\n", cType, cName, f)

	default:
		g.errorf("unsupported const type %s for %s", b, o.Name())
	}
}

type funcSummary struct {
	name              string
	ret               string
	params, retParams []paramInfo
}

type paramInfo struct {
	typ  types.Type
	name string
}

func (g *objcGen) funcSummary(obj *types.Func) *funcSummary {
	s := &funcSummary{name: obj.Name()}

	sig := obj.Type().(*types.Signature)
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		v := paramInfo{
			typ:  p.Type(),
			name: paramName(params, i),
		}
		s.params = append(s.params, v)
	}

	res := sig.Results()
	switch res.Len() {
	case 0:
		s.ret = "void"
	case 1:
		p := res.At(0)
		if isErrorType(p.Type()) {
			s.retParams = append(s.retParams, paramInfo{
				typ:  p.Type(),
				name: "error",
			})
			s.ret = "BOOL"
		} else {
			name := p.Name()
			if name == "" || paramRE.MatchString(name) {
				name = "ret0_"
			}
			typ := p.Type()
			s.retParams = append(s.retParams, paramInfo{typ: typ, name: name})
			s.ret = g.objcType(typ)
		}
	case 2:
		name := res.At(0).Name()
		if name == "" || paramRE.MatchString(name) {
			name = "ret0_"
		}
		s.retParams = append(s.retParams, paramInfo{
			typ:  res.At(0).Type(),
			name: name,
		})

		if !isErrorType(res.At(1).Type()) {
			g.errorf("second result value must be of type error: %s", obj)
			return nil
		}
		s.retParams = append(s.retParams, paramInfo{
			typ:  res.At(1).Type(),
			name: "error", // TODO(hyangah): name collision check.
		})
		s.ret = "BOOL"
	default:
		// TODO(hyangah): relax the constraint on multiple return params.
		g.errorf("too many result values: %s", obj)
		return nil
	}

	return s
}

func (s *funcSummary) asFunc(g *objcGen) string {
	var params []string
	for _, p := range s.params {
		params = append(params, g.objcType(p.typ)+" "+p.name)
	}
	if !s.returnsVal() {
		for _, p := range s.retParams {
			params = append(params, g.objcType(p.typ)+"* "+p.name)
		}
	}
	return fmt.Sprintf("%s %s%s(%s)", s.ret, g.namePrefix, s.name, strings.Join(params, ", "))
}

func (s *funcSummary) asMethod(g *objcGen) string {
	var params []string
	for i, p := range s.params {
		var key string
		if i != 0 {
			key = p.name
		}
		params = append(params, fmt.Sprintf("%s:(%s)%s", key, g.objcType(p.typ), p.name))
	}
	if !s.returnsVal() {
		for _, p := range s.retParams {
			var key string
			if len(params) > 0 {
				key = p.name
			}
			params = append(params, fmt.Sprintf("%s:(%s)%s", key, g.objcType(p.typ)+"*", p.name))
		}
	}
	return fmt.Sprintf("(%s)%s%s", s.ret, s.name, strings.Join(params, " "))
}

func (s *funcSummary) callMethod(g *objcGen) string {
	var params []string
	for i, p := range s.params {
		var key string
		if i != 0 {
			key = p.name
		}
		params = append(params, fmt.Sprintf("%s:%s", key, p.name))
	}
	if !s.returnsVal() {
		for _, p := range s.retParams {
			var key string
			if len(params) > 0 {
				key = p.name
			}
			params = append(params, fmt.Sprintf("%s:&%s", key, p.name))
		}
	}
	return fmt.Sprintf("%s%s", s.name, strings.Join(params, " "))
}

func (s *funcSummary) returnsVal() bool {
	return len(s.retParams) == 1 && !isErrorType(s.retParams[0].typ)
}

func (g *objcGen) genFuncH(obj *types.Func) {
	if s := g.funcSummary(obj); s != nil {
		g.Printf("FOUNDATION_EXPORT %s;\n", s.asFunc(g))
	}
}

func (g *objcGen) seqType(typ types.Type) string {
	s := seqType(typ)
	if s == "String" {
		// TODO(hyangah): non utf-8 strings.
		s = "UTF8"
	}
	return s
}

func (g *objcGen) genFuncM(obj *types.Func) {
	s := g.funcSummary(obj)
	if s == nil {
		return
	}
	g.Printf("%s {\n", s.asFunc(g))
	g.Indent()
	g.genFunc("_DESCRIPTOR_", fmt.Sprintf("_CALL_%s_", s.name), s, false)
	g.Outdent()
	g.Printf("}\n")
}

func (g *objcGen) genGetter(desc string, f *types.Var) {
	t := f.Type()
	if isErrorType(t) {
		t = types.Typ[types.String]
	}
	s := &funcSummary{
		name:      f.Name(),
		ret:       g.objcType(t),
		retParams: []paramInfo{{typ: t, name: "ret_"}},
	}

	g.Printf("- %s {\n", s.asMethod(g))
	g.Indent()
	g.genFunc(desc+"_DESCRIPTOR_", desc+"_FIELD_"+f.Name()+"_GET_", s, true)
	g.Outdent()
	g.Printf("}\n\n")
}

func (g *objcGen) genSetter(desc string, f *types.Var) {
	t := f.Type()
	if isErrorType(t) {
		t = types.Typ[types.String]
	}
	s := &funcSummary{
		name:   "set" + f.Name(),
		ret:    "void",
		params: []paramInfo{{typ: t, name: "v"}},
	}

	g.Printf("- %s {\n", s.asMethod(g))
	g.Indent()
	g.genFunc(desc+"_DESCRIPTOR_", desc+"_FIELD_"+f.Name()+"_SET_", s, true)
	g.Outdent()
	g.Printf("}\n\n")
}

func (g *objcGen) genFunc(pkgDesc, callDesc string, s *funcSummary, isMethod bool) {
	g.Printf("GoSeq in_ = {};\n")
	g.Printf("GoSeq out_ = {};\n")
	if isMethod {
		g.Printf("go_seq_writeRef(&in_, self.ref);\n")
	}
	for _, p := range s.params {
		st := g.seqType(p.typ)
		if st == "Ref" {
			g.Printf("if ([(id<NSObject>)(%s) isKindOfClass:[%s class]]) {\n", p.name, g.refTypeBase(p.typ))
			g.Indent()
			g.Printf("id<goSeqRefInterface> %[1]s_proxy = (id<goSeqRefInterface>)(%[1]s);\n", p.name)
			g.Printf("go_seq_writeRef(&in_, %s_proxy.ref);\n", p.name)
			g.Outdent()
			g.Printf("} else {\n")
			g.Indent()
			g.Printf("go_seq_writeObjcRef(&in_, %s);\n", p.name)
			g.Outdent()
			g.Printf("}\n")
		} else {
			g.Printf("go_seq_write%s(&in_, %s);\n", st, p.name)
		}
	}
	g.Printf("go_seq_send(%s, %s, &in_, &out_);\n", pkgDesc, callDesc)

	if s.returnsVal() {
		p := s.retParams[0]
		if seqTyp := g.seqType(p.typ); seqTyp != "Ref" {
			g.Printf("%s %s = go_seq_read%s(&out_);\n", g.objcType(p.typ), p.name, g.seqType(p.typ))
		} else {
			ptype := g.objcType(p.typ)
			g.Printf("GoSeqRef* %s_ref = go_seq_readRef(&out_);\n", p.name)
			g.Printf("%s %s = %s_ref.obj;\n", ptype, p.name, p.name)
			g.Printf("if (%s == NULL) {\n", p.name)
			g.Indent()
			g.Printf("%s = [[%s alloc] initWithRef:%s_ref];\n", p.name, g.refTypeBase(p.typ), p.name)
			g.Outdent()
			g.Printf("}\n")
		}
	} else {
		for _, p := range s.retParams {
			if isErrorType(p.typ) {
				g.Printf("NSString* _%s = go_seq_readUTF8(&out_);\n", p.name)
				g.Printf("if ([_%s length] != 0 && %s != nil) {\n", p.name, p.name)
				g.Indent()
				g.Printf("NSMutableDictionary* details = [NSMutableDictionary dictionary];\n")
				g.Printf("[details setValue:_%s forKey:NSLocalizedDescriptionKey];\n", p.name)
				g.Printf("*%s = [NSError errorWithDomain:errDomain code:1 userInfo:details];\n", p.name)
				g.Outdent()
				g.Printf("}\n")
			} else if seqTyp := g.seqType(p.typ); seqTyp != "Ref" {
				g.Printf("%s %s_val = go_seq_read%s(&out_);\n", g.objcType(p.typ), p.name, g.seqType(p.typ))
				g.Printf("if (%s != NULL) {\n", p.name)
				g.Indent()
				g.Printf("*%s = %s_val;\n", p.name, p.name)
				g.Outdent()
				g.Printf("}\n")
			} else {
				g.Printf("GoSeqRef* %s_ref = go_seq_readRef(&out_);\n", p.name)
				g.Printf("if (%s != NULL) {\n", p.name)
				g.Indent()
				g.Printf("*%s = %s_ref.obj;\n", p.name, p.name)
				g.Printf("if (*%s == NULL) {\n", p.name)
				g.Indent()
				g.Printf("*%s = [[%s alloc] initWithRef:%s_ref];\n", p.name, g.refTypeBase(p.typ), p.name)
				g.Outdent()
				g.Printf("}\n")
				g.Outdent()
				g.Printf("}\n")
			}
		}
	}

	g.Printf("go_seq_free(&in_);\n")
	g.Printf("go_seq_free(&out_);\n")
	if n := len(s.retParams); n > 0 {
		p := s.retParams[n-1]
		if isErrorType(p.typ) {
			g.Printf("return ([_%s length] == 0);\n", p.name)
		} else {
			g.Printf("return %s;\n", p.name)
		}
	}
}

func (g *objcGen) genInterfaceInterface(obj *types.TypeName, summary ifaceSummary, isProtocol bool) {
	g.Printf("@interface %[1]s%[2]s : NSObject", g.namePrefix, obj.Name())
	if isProtocol {
		g.Printf(" <%[1]s%[2]s>", g.namePrefix, obj.Name())
	}
	g.Printf(" {\n}\n")
	g.Printf("@property(strong, readonly) id ref;\n")
	g.Printf("\n")
	g.Printf("- (id)initWithRef:(id)ref;\n")
	for _, m := range summary.callable {
		s := g.funcSummary(m)
		g.Printf("- %s;\n", s.asMethod(g))
	}
	g.Printf("@end\n")
	g.Printf("\n")
}

func (g *objcGen) genInterfaceH(obj *types.TypeName, t *types.Interface) {
	summary := makeIfaceSummary(t)
	if !summary.implementable {
		g.genInterfaceInterface(obj, summary, false)
		return
	}
	g.Printf("@protocol %s%s\n", g.namePrefix, obj.Name())
	for _, m := range makeIfaceSummary(t).callable {
		s := g.funcSummary(m)
		g.Printf("- %s;\n", s.asMethod(g))
	}
	g.Printf("@end\n")
}

func (g *objcGen) genInterfaceM(obj *types.TypeName, t *types.Interface) bool {
	summary := makeIfaceSummary(t)

	desc := fmt.Sprintf("_GO_%s_%s", g.pkgName, obj.Name())
	g.Printf("#define %s_DESCRIPTOR_ \"go.%s.%s\"\n", desc, g.pkgName, obj.Name())
	for i, m := range summary.callable {
		g.Printf("#define %s_%s_ (0x%x0a)\n", desc, m.Name(), i+1)
	}
	g.Printf("\n")

	if summary.implementable {
		// @interface Interface -- similar to what genStructH does.
		g.genInterfaceInterface(obj, summary, true)
	}

	// @implementation Interface -- similar to what genStructM does.
	g.Printf("@implementation %s%s {\n", g.namePrefix, obj.Name())
	g.Printf("}\n")
	g.Printf("\n")
	g.Printf("- (id)initWithRef:(id)ref {\n")
	g.Indent()
	g.Printf("self = [super init];\n")
	g.Printf("if (self) { _ref = ref; }\n")
	g.Printf("return self;\n")
	g.Outdent()
	g.Printf("}\n")
	g.Printf("\n")

	for _, m := range summary.callable {
		s := g.funcSummary(m)
		g.Printf("- %s {\n", s.asMethod(g))
		g.Indent()
		g.genFunc(desc+"_DESCRIPTOR_", desc+"_"+m.Name()+"_", s, true)
		g.Outdent()
		g.Printf("}\n\n")
	}
	g.Printf("@end\n")
	g.Printf("\n")

	// proxy function.
	if summary.implementable {
		g.Printf("static void proxy%s%s(id obj, int code, GoSeq* in, GoSeq* out) {\n", g.namePrefix, obj.Name())
		g.Indent()
		g.Printf("switch (code) {\n")
		for _, m := range summary.callable {
			g.Printf("case %s_%s_: {\n", desc, m.Name())
			g.Indent()
			g.genInterfaceMethodProxy(obj, g.funcSummary(m))
			g.Outdent()
			g.Printf("} break;\n")
		}
		g.Printf("default:\n")
		g.Indent()
		g.Printf("NSLog(@\"unknown code %%x for %s_DESCRIPTOR_\", code);\n", desc)
		g.Outdent()
		g.Printf("}\n")
		g.Outdent()
		g.Printf("}\n")
	}

	return summary.implementable
}

func (g *objcGen) genInterfaceMethodProxy(obj *types.TypeName, s *funcSummary) {
	g.Printf("id<%[1]s%[2]s> o = (id<%[1]s%[2]s>)(obj);\n", g.namePrefix, obj.Name())
	// read params from GoSeq* inseq
	for _, p := range s.params {
		stype := g.seqType(p.typ)
		ptype := g.objcType(p.typ)
		if stype == "Ref" {
			g.Printf("GoSeqRef* %s_ref = go_seq_readRef(in);\n", p.name)
			g.Printf("%s %s = %s_ref.obj;\n", ptype, p.name, p.name)
			g.Printf("if (%s == NULL) {\n", p.name)
			g.Indent()
			g.Printf("%s = [[%s alloc] initWithRef:%s_ref];\n", p.name, g.refTypeBase(p.typ), p.name)
			g.Outdent()
			g.Printf("}\n")
		} else {
			g.Printf("%s %s = go_seq_read%s(in);\n", ptype, p.name, stype)
		}
	}

	// call method
	if !s.returnsVal() {
		for _, p := range s.retParams {
			if isErrorType(p.typ) {
				g.Printf("NSError* %s = NULL;\n", p.name)
			} else {
				g.Printf("%s %s;\n", g.objcType(p.typ), p.name)
			}
		}
	}

	if s.ret == "void" {
		g.Printf("[o %s];\n", s.callMethod(g))
	} else {
		g.Printf("%s returnVal = [o %s];\n", s.ret, s.callMethod(g))
	}

	// write result to GoSeq* outseq
	if len(s.retParams) == 0 {
		return
	}
	if s.returnsVal() { // len(s.retParams) == 1 && s.retParams[0] != error
		p := s.retParams[0]
		if stype := g.seqType(p.typ); stype == "Ref" {
			g.Printf("if ([(id<NSObject>)(returnVal) isKindOfClass:[%s class]]) {\n", g.refTypeBase(p.typ))
			g.Indent()
			g.Printf("id<goSeqRefInterface>retVal_proxy = (id<goSeqRefInterface>)(returnVal);\n")
			g.Printf("go_seq_writeRef(out, retVal_proxy.ref);\n")
			g.Outdent()
			g.Printf("} else {\n")
			g.Indent()
			g.Printf("go_seq_writeRef(out, returnVal);\n")
			g.Outdent()
			g.Printf("}\n")
		} else {
			g.Printf("go_seq_write%s(out, returnVal);\n", stype)
		}
		return
	}
	for i, p := range s.retParams {
		if isErrorType(p.typ) {
			if i == len(s.retParams)-1 { // last param.
				g.Printf("if (returnVal) {\n")
			} else {
				g.Printf("if (%s == NULL) {\n", p.name)
			}
			g.Indent()
			g.Printf("go_seq_writeUTF8(out, NULL);\n")
			g.Outdent()
			g.Printf("} else {\n")
			g.Indent()
			g.Printf("NSString* %[1]sDesc = [%[1]s localizedDescription];\n", p.name)
			g.Printf("if (%[1]sDesc == NULL || %[1]sDesc.length == 0) {\n", p.name)
			g.Indent()
			g.Printf("%[1]sDesc = @\"gobind: unknown error\";\n", p.name)
			g.Outdent()
			g.Printf("}\n")
			g.Printf("go_seq_writeUTF8(out, %sDesc);\n", p.name)
			g.Outdent()
			g.Printf("}\n")
		} else if seqTyp := g.seqType(p.typ); seqTyp == "Ref" {
			// TODO(hyangah): NULL.
			g.Printf("if ([(id<NSObject>)(%s) isKindOfClass:[%s class]]) {\n", p.name, g.refTypeBase(p.typ))
			g.Indent()
			g.Printf("id<goSeqRefInterface>%[1]s_proxy = (id<goSeqRefInterface>)(%[1]s);\n", p.name)
			g.Printf("go_seq_writeRef(out, %s_proxy.ref);\n", p.name)
			g.Outdent()
			g.Printf("} else {\n")
			g.Indent()
			g.Printf("go_seq_writeObjcRef(out, %s);\n", p.name)
			g.Outdent()
			g.Printf("}\n")
		} else {
			g.Printf("go_seq_write%s(out, %s);\n", seqTyp, p.name)
		}
	}
}

func (g *objcGen) genStructH(obj *types.TypeName, t *types.Struct) {
	g.Printf("@interface %s%s : NSObject {\n", g.namePrefix, obj.Name())
	g.Printf("}\n")
	g.Printf("@property(strong, readonly) id ref;\n")
	g.Printf("\n")
	g.Printf("- (id)initWithRef:(id)ref;\n")

	// accessors to exported fields.
	for _, f := range exportedFields(t) {
		name, typ := f.Name(), g.objcFieldType(f.Type())
		g.Printf("- (%s)%s;\n", typ, name)
		g.Printf("- (void)set%s:(%s)v;\n", name, typ)
	}

	// exported methods
	for _, m := range exportedMethodSet(types.NewPointer(obj.Type())) {
		s := g.funcSummary(m)
		g.Printf("- %s;\n", s.asMethod(g))
	}
	g.Printf("@end\n")
}

func (g *objcGen) genStructM(obj *types.TypeName, t *types.Struct) {
	fields := exportedFields(t)
	methods := exportedMethodSet(types.NewPointer(obj.Type()))

	desc := fmt.Sprintf("_GO_%s_%s", g.pkgName, obj.Name())
	g.Printf("#define %s_DESCRIPTOR_ \"go.%s.%s\"\n", desc, g.pkgName, obj.Name())
	for i, f := range fields {
		g.Printf("#define %s_FIELD_%s_GET_ (0x%x0f)\n", desc, f.Name(), i)
		g.Printf("#define %s_FIELD_%s_SET_ (0x%x1f)\n", desc, f.Name(), i)
	}
	for i, m := range methods {
		g.Printf("#define %s_%s_ (0x%x0c)\n", desc, m.Name(), i)
	}

	g.Printf("\n")
	g.Printf("@implementation %s%s {\n", g.namePrefix, obj.Name())
	g.Printf("}\n\n")
	g.Printf("- (id)initWithRef:(id)ref {\n")
	g.Indent()
	g.Printf("self = [super init];\n")
	g.Printf("if (self) { _ref = ref; }\n")
	g.Printf("return self;\n")
	g.Outdent()
	g.Printf("}\n\n")

	for _, f := range fields {
		g.genGetter(desc, f)
		g.genSetter(desc, f)
	}

	for _, m := range methods {
		s := g.funcSummary(m)
		g.Printf("- %s {\n", s.asMethod(g))
		g.Indent()
		g.genFunc(desc+"_DESCRIPTOR_", desc+"_"+m.Name()+"_", s, true)
		g.Outdent()
		g.Printf("}\n\n")
	}
	g.Printf("@end\n")
}

func (g *objcGen) errorf(format string, args ...interface{}) {
	g.err = append(g.err, fmt.Errorf(format, args...))
}

func (g *objcGen) refTypeBase(typ types.Type) string {
	switch typ := typ.(type) {
	case *types.Pointer:
		if _, ok := typ.Elem().(*types.Named); ok {
			return g.objcType(typ.Elem())
		}
	case *types.Named:
		n := typ.Obj()
		if n.Pkg() == g.pkg {
			switch typ.Underlying().(type) {
			case *types.Interface, *types.Struct:
				return g.namePrefix + n.Name()
			}
		}
	}

	// fallback to whatever objcType returns. This must not happen.
	panic(fmt.Sprintf("wtf: %+T", typ))
	return g.objcType(typ)
}

func (g *objcGen) objcFieldType(t types.Type) string {
	if isErrorType(t) {
		return "NSString*"
	}
	return g.objcType(t)
}

func (g *objcGen) objcType(typ types.Type) string {
	if isErrorType(typ) {
		return "NSError*"
	}

	switch typ := typ.(type) {
	case *types.Basic:
		switch typ.Kind() {
		case types.Bool, types.UntypedBool:
			return "BOOL"
		case types.Int:
			return "int"
		case types.Int8:
			return "int8_t"
		case types.Int16:
			return "int16_t"
		case types.Int32, types.UntypedRune: // types.Rune
			return "int32_t"
		case types.Int64, types.UntypedInt:
			return "int64_t"
		case types.Uint8:
			// byte is an alias of uint8, and the alias is lost.
			return "byte"
		case types.Uint16:
			return "uint16_t"
		case types.Uint32:
			return "uint32_t"
		case types.Uint64:
			return "uint64_t"
		case types.Float32:
			return "float"
		case types.Float64, types.UntypedFloat:
			return "double"
		case types.String, types.UntypedString:
			return "NSString*"
		default:
			g.errorf("unsupported type: %s", typ)
			return "TODO"
		}
	case *types.Slice:
		elem := g.objcType(typ.Elem())
		// Special case: NSData seems to be a better option for byte slice.
		if elem == "byte" {
			return "NSData*"
		}
		// TODO(hyangah): support other slice types: NSArray or CFArrayRef.
		// Investigate the performance implication.
		g.errorf("unsupported type: %s", typ)
		return "TODO"
	case *types.Pointer:
		if _, ok := typ.Elem().(*types.Named); ok {
			return g.objcType(typ.Elem()) + "*"
		}
		g.errorf("unsupported pointer to type: %s", typ)
		return "TODO"
	case *types.Named:
		n := typ.Obj()
		if n.Pkg() != g.pkg {
			g.errorf("type %s is in package %s; only types defined in package %s is supported", n.Name(), n.Pkg().Name(), g.pkg.Name())
			return "TODO"
		}
		switch t := typ.Underlying().(type) {
		case *types.Interface:
			if makeIfaceSummary(t).implementable {
				return "id<" + g.namePrefix + n.Name() + ">"
			} else {
				return g.namePrefix + n.Name() + "*"
			}
		case *types.Struct:
			return g.namePrefix + n.Name()
		}
		g.errorf("unsupported, named type %s", typ)
		return "TODO"
	default:
		g.errorf("unsupported type: %#+v, %s", typ, typ)
		return "TODO"
	}
}
