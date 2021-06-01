package gen

import (
	"go/ast"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/tools/go/loader"

	"github.com/artisanhe/tools/codegen/loaderx"
	"github.com/artisanhe/tools/courier"
)

var (
	courierPkgImportPath = "github.com/artisanhe/tools/courier"
	routerTypeString     = reflectTypeString(reflect.TypeOf(new(courier.Router)))
)

func isRouterType(tpe types.Type) bool {
	return strings.HasSuffix(tpe.String(), courierPkgImportPath+".Router")
}

func NewRouterScanner(program *loader.Program) *RouterScanner {
	routerScanner := &RouterScanner{
		program: program,
		routers: map[*types.Var]*Router{},
	}

	routerScanner.init()

	return routerScanner
}

type RouterScanner struct {
	program *loader.Program
	routers map[*types.Var]*Router
}

func (scanner *RouterScanner) Router(typeName *types.Var) *Router {
	return scanner.routers[typeName]
}

type OperatorTypeName struct {
	Path string
	*types.TypeName
}

type OperatorTypeNames []*OperatorTypeName

func OperatorTypeNameFromType(tpe types.Type) *OperatorTypeName {
	switch tpe.(type) {
	case *types.Named:
		return &OperatorTypeName{
			TypeName: tpe.(*types.Named).Obj(),
		}
	case *types.Pointer:
		return &OperatorTypeName{
			TypeName: tpe.(*types.Pointer).Elem().(*types.Named).Obj(),
		}
	}
	return nil
}

func FromArgs(pkgInfo *loader.PackageInfo, args ...ast.Expr) OperatorTypeNames {
	opTypeNames := OperatorTypeNames{}
	for _, arg := range args {
		opTypeName := OperatorTypeNameFromType(pkgInfo.TypeOf(arg))
		if opTypeName == nil {
			continue
		}
		if callExpr, ok := arg.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "Group" {
					if strings.Contains(pkgInfo.ObjectOf(selectorExpr.Sel).Pkg().String(), courierPkgImportPath) {
						switch v := callExpr.Args[0].(type) {
						case *ast.BasicLit:
							opTypeName.Path, _ = strconv.Unquote(v.Value)
						}
					}
				}
			}
		}
		opTypeNames = append(opTypeNames, opTypeName)

	}
	return opTypeNames
}

func (scanner *RouterScanner) init() {
	for _, pkgInfo := range scanner.program.AllPackages {
		for ident, obj := range pkgInfo.Defs {
			if typeVar, ok := obj.(*types.Var); ok {
				if typeVar != nil && !strings.HasSuffix(typeVar.Pkg().Path(), courierPkgImportPath) {
					if isRouterType(typeVar.Type()) {
						router := NewRouter()

						ast.Inspect(ident.Obj.Decl.(ast.Node), func(node ast.Node) bool {
							switch node.(type) {
							case *ast.CallExpr:
								callExpr := node.(*ast.CallExpr)
								router.AppendOperators(FromArgs(pkgInfo, callExpr.Args...)...)
								return false
							}
							return true
						})

						scanner.routers[typeVar] = router
					}
				}
			}
		}
	}

	for _, pkgInfo := range scanner.program.AllPackages {
		for selectExpr, selection := range pkgInfo.Selections {
			if selection.Obj() != nil {
				if typeFunc, ok := selection.Obj().(*types.Func); ok {
					recv := typeFunc.Type().(*types.Signature).Recv()
					if recv != nil && isRouterType(recv.Type()) {
						for typeVar, router := range scanner.routers {
							switch selectExpr.Sel.Name {
							case "Register":
								if typeVar == pkgInfo.ObjectOf(IdentOfCallExprSelectExpr(selectExpr)) {
									file := loaderx.FileOf(selectExpr, pkgInfo.Files...)
									ast.Inspect(file, func(node ast.Node) bool {
										switch node.(type) {
										case *ast.CallExpr:
											callExpr := node.(*ast.CallExpr)
											if callExpr.Fun == selectExpr {
												routerIdent := callExpr.Args[0]
												switch routerIdent.(type) {
												case *ast.Ident:
													argTypeVar := pkgInfo.ObjectOf(routerIdent.(*ast.Ident)).(*types.Var)
													if r, ok := scanner.routers[argTypeVar]; ok {
														router.Register(r)
													}
												case *ast.SelectorExpr:
													argTypeVar := pkgInfo.ObjectOf(routerIdent.(*ast.SelectorExpr).Sel).(*types.Var)
													if r, ok := scanner.routers[argTypeVar]; ok {
														router.Register(r)
													}
												case *ast.CallExpr:
													callExprForRegister := routerIdent.(*ast.CallExpr)
													router.With(FromArgs(pkgInfo, callExprForRegister.Args...)...)
												}
												return false
											}
										}
										return true
									})
								}
							}
						}
					}
				}
			}
		}
	}
}

func IdentOfCallExprSelectExpr(selectExpr *ast.SelectorExpr) *ast.Ident {
	switch selectExpr.X.(type) {
	case *ast.Ident:
		return selectExpr.X.(*ast.Ident)
	case *ast.SelectorExpr:
		return selectExpr.X.(*ast.SelectorExpr).Sel
	}
	return nil
}

func NewRouter(operators ...*OperatorTypeName) *Router {
	return &Router{
		operators: operators,
	}
}

type Router struct {
	parent    *Router
	operators []*OperatorTypeName
	children  map[*Router]bool
}

func (router *Router) AppendOperators(operators ...*OperatorTypeName) {
	router.operators = append(router.operators, operators...)
}

func (router *Router) With(operators ...*OperatorTypeName) {
	router.Register(NewRouter(operators...))
}

func (router *Router) Register(r *Router) {
	if router.children == nil {
		router.children = map[*Router]bool{}
	}
	r.parent = router
	router.children[r] = true
}

func (router *Router) Route(program *loader.Program) *Route {
	parent := router.parent
	operators := router.operators

	for parent != nil {
		operators = append(parent.operators, operators...)
		parent = parent.parent
	}

	route := Route{
		last:      router.children == nil,
		operators: operators,
	}

	route.SetMethod(program)
	route.SetPath(program)

	return &route
}

func (router *Router) Routes(program *loader.Program) (routes []*Route) {
	for child := range router.children {
		route := child.Route(program)
		if route.last {
			routes = append(routes, route)
		}
		if child.children != nil {
			routes = append(routes, child.Routes(program)...)
		}
	}
	return routes
}

type Route struct {
	Method    string
	Path      string
	last      bool
	operators []*OperatorTypeName
}

func (route *Route) SetPath(program *loader.Program) {
	p := "/"
	for _, operator := range route.operators {
		if operator.Path != "" {
			p += operator.Path
			continue
		}

		typeFunc := loaderx.MethodOf(operator.Type().(*types.Named), "Path")

		loaderx.ForEachFuncResult(program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
			if resultTypeAndValues[0].IsValue() {
				p += getConstVal(resultTypeAndValues[0].Value).(string)
			}
		})
	}
	route.Path = httprouter.CleanPath(p)
}

func (route *Route) SetMethod(program *loader.Program) {
	if len(route.operators) > 0 {
		operator := route.operators[len(route.operators)-1]
		typeFunc := loaderx.MethodOf(operator.Type().(*types.Named), "Method")

		loaderx.ForEachFuncResult(program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
			if resultTypeAndValues[0].IsValue() {
				route.Method = getConstVal(resultTypeAndValues[0].Value).(string)
			}
		})
	}
}
