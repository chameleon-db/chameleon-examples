package middleware

import "net/http"

// CORS middleware para permitir peticiones del frontend
// SOLO PARA DESARROLLO, NO USAR EN PRODUCCIÓN SIN CONFIGURAR ORÍGENES ESPECÍFICOS!
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir cualquier origen (para desarrollo)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Métodos permitidos
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Headers permitidos
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Permitir credenciales (cookies, auth headers)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Manejar preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
