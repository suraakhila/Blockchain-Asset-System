package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gorilla/mux"
    "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Asset represents the structure of an asset
type Asset struct {
    DealerID    string `json:"dealerId"`
    MSISDN      string `json:"msisdn"`
    MPIN        string `json:"mpin"`
    Balance     string `json:"balance"`
    Status      string `json:"status"`
    TransAmount string `json:"transAmount"`
    TransType   string `json:"transType"`
    Remarks     string `json:"remarks"`
}

// Function to connect to the Fabric network
func connectToNetwork() (*gateway.Gateway, error) {
    walletPath := filepath.Join("wallet")
    wallet, err := gateway.NewFileSystemWallet(walletPath)
    if err != nil {
        return nil, fmt.Errorf("failed to create wallet: %w", err)
    }

    ccpPath := filepath.Join("path", "to", "connection.json") // Update this path

    gw, err := gateway.Connect(
        gateway.WithConfig(gateway.WithConfigOption(gateway.ConfigPath(filepath.Clean(ccpPath)))),
        gateway.WithIdentity(wallet, "admin"), // Update with wallet identity
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to gateway: %w", err)
    }
    return gw, nil
}

// Handler to create a new asset
func createAssetHandler(w http.ResponseWriter, r *http.Request) {
    var asset Asset
    err := json.NewDecoder(r.Body).Decode(&asset)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    gw, err := connectToNetwork()
    if err != nil {
        log.Printf("Error connecting to network: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer gw.Close()

    network := gw.GetNetwork("mychannel") // Replace with your channel name
    contract := network.GetContract("assetContract") // Replace with your contract name

    _, err = contract.SubmitTransaction("CreateAsset", asset.DealerID, asset.MSISDN, asset.MPIN, asset.Balance, asset.Status, asset.TransAmount, asset.TransType, asset.Remarks)
    if err != nil {
        log.Printf("Failed to submit transaction: %v", err)
        http.Error(w, "Failed to create asset", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(asset)
}

// Handler to update an asset
func updateAssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dealerID := vars["dealerId"]

    var asset Asset
    err := json.NewDecoder(r.Body).Decode(&asset)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    gw, err := connectToNetwork()
    if err != nil {
        log.Printf("Error connecting to network: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer gw.Close()

    network := gw.GetNetwork("mychannel")
    contract := network.GetContract("assetContract")

    _, err = contract.SubmitTransaction("UpdateAsset", dealerID, asset.MSISDN, asset.MPIN, asset.Balance, asset.Status, asset.TransAmount, asset.TransType, asset.Remarks)
    if err != nil {
        log.Printf("Failed to submit transaction: %v", err)
        http.Error(w, "Failed to update asset", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(asset)
}

// Handler to read an asset
func getAssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dealerID := vars["dealerId"]

    gw, err := connectToNetwork()
    if err != nil {
        log.Printf("Error connecting to network: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer gw.Close()

    network := gw.GetNetwork("mychannel")
    contract := network.GetContract("assetContract")

    result, err := contract.EvaluateTransaction("ReadAsset", dealerID)
    if err != nil {
        log.Printf("Failed to evaluate transaction: %v", err)
        http.Error(w, "Asset not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(result)
}

// Handler to delete an asset
func deleteAssetHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    dealerID := vars["dealerId"]

    gw, err := connectToNetwork()
    if err != nil {
        log.Printf("Error connecting to network: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer gw.Close()

    network := gw.GetNetwork("mychannel")
    contract := network.GetContract("assetContract")

    _, err = contract.SubmitTransaction("DeleteAsset", dealerID)
    if err != nil {
        log.Printf("Failed to submit transaction: %v", err)
        http.Error(w, "Failed to delete asset", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Asset with DealerID %s deleted successfully", dealerID)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/assets", createAssetHandler).Methods("POST")
    r.HandleFunc("/assets/{dealerId}", updateAssetHandler).Methods("PUT")
    r.HandleFunc("/assets/{dealerId}", getAssetHandler).Methods("GET")
    r.HandleFunc("/assets/{dealerId}", deleteAssetHandler).Methods("DELETE")

    http.Handle("/", r)
    fmt.Println("API Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
