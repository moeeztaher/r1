package ErrorHandler

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"

	"encoding/json"
    "net/http"
    "fmt"
    "r1/r1/Apis"
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"

    "github.com/gorilla/mux"

)

var ctx = context.Background()

// Initialize Redis client
var rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379", // Use default Redis port
})

var supportedContentTypes = []string{
    "application/json",
}

const (
    limit     = 5           // Maximum number of requests
    window    = 1 * time.Minute // Time window
)

const MaxPayloadSize = 1 * 1024 * 1024 // 1 MB

// For unauthorized access (401)
func checkAuthorization(w http.ResponseWriter, r *http.Request) bool {
    // Implement authorization logic here
    // Handle MongoDB authorization error (HTTP 500)
    errorResponse := Apis.ErrorResponse{
        Type:             "https://example.com/errors/mongodb-authorization",
        Title:            "MongoDB Authorization Error",
        Status:           http.StatusUnauthorized,
        Detail:           "Unauthorized Access",
        Cause:            "err.Error()",
        SupportedFeatures: "string",
    }

    writeErrorResponse(w, errorResponse, http.StatusUnauthorized)
    return true
}

// For forbidden access (403)
func checkPermissions(w http.ResponseWriter, r *http.Request) bool {
    // Implement permission check logic here
    errorResponse := Apis.ErrorResponse{
        Type:             "https://example.com/errors/mongodb-permission",
        Title:            "MongoDB Permission Error",
        Status:           http.StatusForbidden,
        Detail:           "Wrong permissions, forbidden access",
        Cause:            "err.Error()",
        SupportedFeatures: "string",
    }
    writeErrorResponse(w, errorResponse, http.StatusForbidden)
    return true
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    errorResponse := Apis.ErrorResponse{
        Type:             "Not Found",
        Title:            "route not found",
        Status:           http.StatusNotFound,
        Detail:           "route is not found",
        Cause:            "err.Error()",
        SupportedFeatures: "string",
    }
    writeErrorResponse(w, errorResponse, http.StatusNotFound)
}

func CheckContentLength(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Length") == "" {
            errorResponse := Apis.ErrorResponse{
                Type:             "Content Length",
                Title:            "MissingHeader",
                Status:           http.StatusLengthRequired,
                Detail:           "Content-Length header is required",
                Cause:            "err.Error()",
                SupportedFeatures: "string",
            }
            writeErrorResponse(w, errorResponse, http.StatusLengthRequired)

            return
        }
        next.ServeHTTP(w, r)
    })
}

func CheckPayloadSize(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, MaxPayloadSize)
        if err := r.ParseForm(); err != nil {
            errorResponse := Apis.ErrorResponse{
                Type:             "Content Length Too Large",
                Title:            "Large Payload Length",
                Status:           http.StatusRequestEntityTooLarge,
                Detail:           "Content- Payload Length header is too large",
                Cause:            "err.Error()",
                SupportedFeatures: "string",
            }
            writeErrorResponse(w, errorResponse, http.StatusRequestEntityTooLarge)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func CheckContentType(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        contentType := r.Header.Get("Content-Type")
        for _, ct := range supportedContentTypes {
            if contentType == ct {
                next.ServeHTTP(w, r)
                return
            }
        }
        errorResponse := Apis.ErrorResponse{
            Type:             "UnsupportedMediaType",
            Title:            "Unsupported Media Type",
            Status:           http.StatusUnsupportedMediaType,
            Detail:           "Content- Media Type is not supported",
            Cause:            "err.Error()",
            SupportedFeatures: "string",
        }
        writeErrorResponse(w, errorResponse, http.StatusUnsupportedMediaType)
    })
}

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr // Use remote address as the key

        // Increment request count for this IP
        count, err := rdb.Incr(ctx, ip).Result()
        if err != nil {
			errorResponse := Apis.ErrorResponse{
				Type:             "RateLimit Error",
				Title:            "RedisError",
				Status:           http.StatusInternalServerError,
				Detail:           "Error incrementing request count",
				Cause:            "err.Error()",
				SupportedFeatures: "string",
			}
			writeErrorResponse(w, errorResponse, http.StatusInternalServerError)

            return
        }

        // Set expiration if it's the first request
        if count == 1 {
            rdb.Expire(ctx, ip, window)
        }

        // Check if the request count exceeds the limit
        if count > limit {
            w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))

			errorResponse := Apis.ErrorResponse{
				Type:             "RateLimitExceeded",
				Title:            "RedisError",
				Status:           http.StatusTooManyRequests,
				Detail:           "Too many requests",
				Cause:            "err.Error()",
				SupportedFeatures: "string",
			}
			writeErrorResponse(w, errorResponse, http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Helper function to write error response
func writeErrorResponse(w http.ResponseWriter, errResp Apis.ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResp)
}
