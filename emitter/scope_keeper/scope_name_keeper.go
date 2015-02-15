package scope_keeper

import (
	_ "fmt"
	"strconv"
	_ "strings"
)

type scopeType int

const (
	funcScope scopeType = iota
	sentenceScope
)

type scopeEntry struct {
	scopeType
	funcName    string
	sentenceNum int
}

type ScopeKeeper struct {
	scopes []scopeEntry
}

func NewScopeKeeper() *ScopeKeeper {
	return &ScopeKeeper{make([]scopeEntry, 0)}
}

func (scopeKeeper *ScopeKeeper) Copy() *ScopeKeeper {
	newScopeKeeper := ScopeKeeper{make([]scopeEntry, len(scopeKeeper.scopes))}

	copy(newScopeKeeper.scopes, scopeKeeper.scopes)

	return &newScopeKeeper
}

func (scopeKeeper *ScopeKeeper) AddFuncScope(funcName string) {
	scopeKeeper.scopes = append(scopeKeeper.scopes, scopeEntry{scopeType: funcScope, funcName: funcName})
}

func (scopeKeeper *ScopeKeeper) AddSentenceScope(sentenceNum int) {
	scopeKeeper.scopes = append(scopeKeeper.scopes, scopeEntry{scopeType: sentenceScope, sentenceNum: sentenceNum})
}

func (scopeKeeper *ScopeKeeper) PopLastSentenceScope() {

	if scopeKeeper.scopes[len(scopeKeeper.scopes)-1].scopeType == sentenceScope {
		scopeKeeper.scopes = scopeKeeper.scopes[0 : len(scopeKeeper.scopes)-1]
	}
}

func (scopeKeeper *ScopeKeeper) GetFuncName(funcName string) string {
	scopeName := ""

	for _, scopeEntry := range scopeKeeper.scopes {

		if scopeEntry.scopeType == funcScope {
			scopeName += scopeEntry.funcName
		} else {
			scopeName += strconv.Itoa(scopeEntry.sentenceNum)
		}
	}

	return scopeName + funcName
}

func (scopeKeeper *ScopeKeeper) GetAllScopes() []string {
	allScopes := make([]string, 0)

	if len(scopeKeeper.scopes) <= 1 {
		return allScopes
	}

	scope := scopeKeeper.scopes[0].funcName

	for _, scopeEntry := range scopeKeeper.scopes[1:] {

		//fmt.Printf("Scope: %s\n", scope)
		if scopeEntry.scopeType == funcScope {
			scope += scopeEntry.funcName
		} else {
			//fmt.Printf("Scope sentence: %s\n", strconv.Itoa(scopeEntry.sentenceNum))
			scope += strconv.Itoa(scopeEntry.sentenceNum)
			allScopes = append(allScopes, scope)
		}
	}

	return allScopes
}
