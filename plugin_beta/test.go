package plugin_beta

type A struct {
	user ITest
}

func(a *A) regUser(user ITest) {
	a.user = user
}

type ITest interface {
	Init(*A)
}

type B struct {
}

func (b *B) Init(a *A) {
	a.regUser(b)
}

