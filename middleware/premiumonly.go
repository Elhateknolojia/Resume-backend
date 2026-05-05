package middleware

import (
    "net/http"
)

func PremiumOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tier := r.Context().Value("tier").(string)
        if tier != "premium" {
            http.Error(w, "Upgrade required for premium features", http.StatusPaymentRequired)
            return
        }
        next.ServeHTTP(w, r)
    })
}
