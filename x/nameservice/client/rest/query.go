package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/yxuco/cosmos-nameservice/x/nameservice/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, querierRoute string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", querierRoute), namesHandler(cliCtx, restName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", querierRoute, restName), resolveNameHandler(cliCtx, querierRoute)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", querierRoute, restName), whoIsHandler(cliCtx, querierRoute)).Methods("GET")
}

func resolveNameHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryResolve, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func whoIsHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", querierRoute, types.QueryWhois, name), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func namesHandler(cliCtx context.CLIContext, querierRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", querierRoute, types.QueryNames), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
