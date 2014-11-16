package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

import (
	"tokens"
)

const (
	tab = " "
)

func PrintUnit(unit *Unit) string {

	resStr := printBuiltinsFuncs(0, unit.Builtins)
	resStr += printGlobalFuncs(0, unit.GlobMap)

	return resStr
}

func genTabs(depth int) string {
	return strings.Repeat(tab, depth)
}

func printBuiltinsFuncs(depth int, funcs map[string]bool) string {
	usedFuncs := make([]string, 0)

	for funcName, used := range funcs {
		if used {
			usedFuncs = append(usedFuncs, funcName)
		}
	}

	return fmt.Sprintf("%s Builtins: %s\n", genTabs(depth), strings.Join(usedFuncs, ", "))
}

func printGlobalFuncs(depth int, funcs map[string]*Function) string {
	globalFuncs := make([]string, 0)

	for _, funcInfo := range funcs {
		globalFuncs = append(globalFuncs, printFuncInfo(depth, funcInfo))
	}

	return fmt.Sprintf("%s", strings.Join(globalFuncs, "\n"))
}

func printFuncInfo(depth int, funcInfo *Function) string {

	funcInfoStr := fmt.Sprintf("%s\n", printFuncName(depth, funcInfo))
	funcInfoStr += fmt.Sprintf("%s\n", printFuncHeader(depth+1, funcInfo))
	funcInfoStr += fmt.Sprintf("%s\n", printScope(depth+1, funcInfo.Params))
	funcInfoStr += fmt.Sprintf("%s\n", printFuncSentences(depth+1, funcInfo.Sentences))

	return funcInfoStr
}

func printFuncName(depth int, funcInfo *Function) string {
	tabs := genTabs(depth)

	if funcInfo.HasName {
		return fmt.Sprintf("%s%s", tabs, funcInfo.FuncName)
	}

	return fmt.Sprintf("%s%s", tabs, "Anonymous")
}

func printFuncHeader(depth int, funcInfo *Function) string {
	tabs := genTabs(depth)

	funcHeaderStr := fmt.Sprintf("%sIsEntry: %s\n", tabs, strconv.FormatBool(funcInfo.IsEntry))
	funcHeaderStr += fmt.Sprintf("%sIsIdent: %s\n", tabs, strconv.FormatBool(funcInfo.IsIdent))
	funcHeaderStr += fmt.Sprintf("%sIsSe: %s\n", tabs, strconv.FormatBool(funcInfo.IsSe))
	funcHeaderStr += fmt.Sprintf("%sRollback: %s\n", tabs, strconv.FormatBool(funcInfo.Rollback))

	return funcHeaderStr
}

func printScope(depth int, funcParams Scope) string {
	tabs := genTabs(depth)
	scopeStr := fmt.Sprintf("%sNested funcs num: %d\n", tabs, len(funcParams.FuncMap))

	exists := "exists"

	if funcParams.Parent == nil {
		exists = "not exists"
	}

	scopeStr += fmt.Sprintf("%sParent scope is %s\n", tabs, exists)
	scopeStr += fmt.Sprintf("%sVars: \n", tabs)

	for varType, vars := range funcParams.VarMap {
		scopeStr += fmt.Sprintf("%s%s%s: ", tabs, tab, tokens.VarType(varType).String())

		varPairs := make([]string, 0)

		for varName, varNum := range vars {
			varPairs = append(varPairs, fmt.Sprintf("{%s.%s, %d}", tokens.VarType(varType).String(), varName, varNum))
		}

		scopeStr += fmt.Sprintf("%s\n", strings.Join(varPairs, ", "))
	}

	return scopeStr
}

//type Sentence struct {
//	coords.Fragment
//	Scope
//	Pattern Expr
//	Actions []*Action
//}

func printFuncSentences(depth int, sentences []*Sentence) string {
	tabs := genTabs(depth)
	funcSentencesStr := ""

	for ind, sentece := range sentences {

		funcSentencesStr += fmt.Sprintf("%s<---Sentence #%d--->\n", tabs, ind)
		funcSentencesStr += fmt.Sprintf("%s\n", printScope(depth+1, sentece.Scope))
		funcSentencesStr += fmt.Sprintf("%s\n", printPattern(depth+1, sentece.Pattern))
		funcSentencesStr += fmt.Sprintf("%s\n", printActions(depth+1, sentece.Actions))
		funcSentencesStr += "\n"
	}

	return funcSentencesStr
}

//type Expr struct {
//	coords.Fragment
//	Terms []*Term
//}

func printPattern(depth int, pattern Expr) string {
	tabs := genTabs(depth)

	return fmt.Sprintf("%sPattern:\n%s", tabs, printExpr(depth, pattern))
}

func printExprs(depth int, exprs []*Expr) string {
	tabs := genTabs(depth)
	exprsStr := make([]string, 0)

	for _, expr := range exprs {
		exprsStr = append(exprsStr, printExpr(0, *expr))
	}

	return fmt.Sprintf("%s%s", tabs, strings.Join(exprsStr, "; "))
}

func printExpr(depth int, pattern Expr) string {
	termStrs := make([]string, 0)

	for _, term := range pattern.Terms {
		termStrs = append(termStrs, printTerm(depth+1, term))
	}

	return strings.Join(termStrs, "\n")
}

//type Term struct {
//	Comment string
//	coords.Fragment
//	TermTag
//	tokens.Value
//	Exprs []*Expr
//	*Function
//}

func printTerm(depth int, term *Term) string {
	tabs := genTabs(depth)
	var termStr string

	switch term.TermTag {
	case L:
	case R:
		termStr = TermTag(term.TermTag).String()
		break

	case STR:
		termStr = string(term.Value.Str)
		break

	case COMP:
		termStr = fmt.Sprintf("{Symbol:%s}", term.Value.Name)
		break

	case INT:
		termStr = term.Value.Int.String()
		break

	case FLOAT:
		termStr = "Float number"
		//termStr = strconv.FormatFloat(term.tokens.Value.Float, 10)
		break

	case VAR:
		termStr = fmt.Sprintf("%s.%s", tokens.VarType(term.Value.VarType).String(), term.Value.Name)
		break

	case EXPR:
		//TO FIX:
		//termStr = fmt.Sprintf("(%s)", printPattern(0, ))
		break

	case BRACED_EXPR:
		termStr = TermTag(term.TermTag).String()
		break

	case BRACKETED_EXPR:
		termStr = TermTag(term.TermTag).String()
		break

	case ANGLED_EXPR:
		termStr = TermTag(term.TermTag).String()
		break

	case EVAL:
		//termStr = printEvalExpr(0, term)
		termStr = TermTag(term.TermTag).String()
		break

	case FUNC:
		termStr = fmt.Sprintf("Func term:\n%s", printFuncInfo(depth+1, term.Function))
		break
	}

	return fmt.Sprintf("%s%s", tabs, termStr)
}

func printEvalExpr(depth int, term *Term) string {
	tabs := genTabs(depth)

	return fmt.Sprintf("%s<%s %s>", tabs, term.Function.FuncName, printExprs(depth, term.Exprs))
}

func printActions(depth int, actions []*Action) string {
	tabs := genTabs(depth)
	actionsStr := fmt.Sprintf("%sActions:\n", tabs)

	for _, action := range actions {
		actionsStr += fmt.Sprintf("%s\n", printAction(depth+1, action))
	}

	return actionsStr
}

//type Action struct {
//	Comment string
//	coords.Fragment
//	ActionOp
//	Expr
//}

func printAction(depth int, action *Action) string {
	tabs := genTabs(depth)

	return fmt.Sprintf("%s%s:\n%s", tabs, action.ActionOp.String(), printExpr(depth+1, action.Expr))
}
