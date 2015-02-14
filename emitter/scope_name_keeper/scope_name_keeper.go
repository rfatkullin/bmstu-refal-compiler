package scope_name_keeper

type FullPathKeeper struct {
	nameChains []string
}

func NewFullPathKeeper() *FullPathKeeper {
	return &FullPathKeeper{make([]string, 8)}
}

func (pathKeeper *FullPathKeeper) PushNameElement(el string) {
	pathKeeper.nameChains = append(pathKeeper.nameChains, el)
}

func (pathKeeper *FullPathKeeper) PushIntElement(el int) {
	pathKeeper.nameChains = append(pathKeeper.nameChains, string(el))
}

func (pathKeeper *FullPathKeeper) PopElement() {
	if pathKeeper.nameChains != nil {
		pathKeeper.nameChains = pathKeeper.nameChains[0 : len(pathKeeper.nameChains)-1]
	}
}

func (pathKeeper *FullPathKeeper) WithLastName(name string) string {
	return pathKeeper.String() + name
}

func (pathKeeper *FullPathKeeper) String() string {
	fullName := ""

	if pathKeeper.nameChains != nil {
		for _, el := range pathKeeper.nameChains {
			fullName += el
		}
	}

	return fullName
}
