package router

import (
	"net/http"
	execsRoutes "restapi/internal/api/router/routes/execs"
	rootRoutes "restapi/internal/api/router/routes/root"
	studentsRoutes "restapi/internal/api/router/routes/students"
	teachersRoutes "restapi/internal/api/router/routes/teachers"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()

	rootRoutes.Register(mux)
	teachersRoutes.Register(mux)
	studentsRoutes.Register(mux)
	execsRoutes.Register(mux)

	return mux
}
