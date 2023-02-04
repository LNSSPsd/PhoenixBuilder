package collaborate

const (
	INTERFACE_GEN_STRING_LIST_HINT_RESOLVER            = "INTERFACE_GEN_STRING_LIST_HINT_RESOLVER"
	INTERFACE_GEN_STRING_LIST_HINT_RESOLVER_WITH_INDEX = "INTERFACE_GEN_STRING_LIST_HINT_RESOLVER_WITH_INDEX"
	INTERFACE_GEN_INT_RANGE_RESOLVER                   = "INTERFACE_GEN_INT_RANGE_RESOLVER"
	INTERFACE_GEN_YES_NO_RESOLVER                      = "INTERFACE_GEN_YES_NO_RESOLVER"
	INTERFACE_QUERY_FOR_PLAYER_NAME                    = "INTERFACE_QUERY_FOR_PLAYER_NAME"
)

type GEN_STRING_LIST_HINT_RESOLVER func(available []string) (string, func(params []string) (selection int, cancel bool, err error))

type GEN_STRING_LIST_HINT_RESOLVER_WITH_INDEX func(_available []string) (string, func(params []string) (selection int, cancel bool, err error))

type GEN_INT_RANGE_RESOLVER func(min int, max int) (string, func(params []string) (selection int, cancel bool, err error))

type GEN_YES_NO_RESOLVER func() (string, func(params []string) (bool, error))

type QUERY_FOR_PLAYER_NAME func(src string, dst string, searchFn FUNCTYPE_GET_POSSIBLE_NAME) (name string, cancel bool)
