package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"

    _ "github.com/lib/pq"
)

type CreateServiceRequest struct {
    ServiceName string  `json:"service_name"`
}

type CreateServiceResponse struct {
    ServiceID   int     `json:"service_id"`
    ServiceName string  `json:"service_name"`
}

type RegisterServiceNodeRequest struct {
    ServiceName string  `json:"service_name"`
    ServicePort int     `json:"service_port"`
}

type RegisterServiceNodeResponse struct {
    Message string  `json:"message"`
}

func createServiceHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var req CreateServiceRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("JSON decode error: %v", err)
            http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
            return
        }

        if req.ServiceName == "" {
            http.Error(w, "service_name is required", http.StatusBadRequest)
            return
        }

        var newId int
        insertSQL := `INSERT INTO services (service_name) VALUES ($1) RETURNING service_id;`
        err := db.QueryRow(insertSQL, req.ServiceName).Scan(&newId)
        if err != nil {
            log.Printf("DB insert error: %v", err)
            http.Error(w, "Failed to insert record", http.StatusInternalServerError)
            return
        }

        resp := CreateServiceResponse {
            ServiceID:      newId,
            ServiceName:    req.ServiceName,
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("JSON encode error: %v", err)
        }

    }
}

func registerServiceNodeHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        var req RegisterServiceNodeRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("JSON decode error: %v", err)
            http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
            return
        }

        if req.ServiceName == "" {
            http.Error(w, "service_name is required", http.StatusBadRequest)
            return
        }

        var serviceId int
        queryStr := `SELECT service_id FROM services WHERE service_name LIKE $1`
        err := db.QueryRow(queryStr, req.ServiceName).Scan(&serviceId)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, fmt.Sprintf("No service found with name: %s", req.ServiceName), http.StatusBadRequest)
                return
            } else {
                log.Printf("SQL Error: %v\n", err)
                http.Error(w, "Failed to query service", http.StatusInternalServerError)
            }
        }

        nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)

        endpoint := ServiceEndpoint {
            IPAddress:  nodeIP,
            Port:       req.ServicePort,
        }

        serviceMapLock.RLock()
        _, ok := serviceMap[serviceId]
        serviceMapLock.RUnlock()
        if !ok {
            serviceMapLock.Lock()
            serviceMap[serviceId] = &ServiceEndpointList{
                endpoints:  make([]ServiceEndpoint, 0),
            }
            serviceMapLock.Unlock()
        }

        serviceMap[serviceId].lock.Lock()
        serviceMap[serviceId].endpoints = append(serviceMap[serviceId].endpoints, endpoint)
        serviceMap[serviceId].lock.Unlock()

        resp := RegisterServiceNodeResponse {
            Message: "ok",
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("JSON encode error: %v", err)
        }

    }
}
