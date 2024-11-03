# Blockchain-Asset-System 🏦

## Overview

This project implements a **blockchain-based asset management system** using **Hyperledger Fabric**. The system allows for **creating**, **updating**, **querying**, and **tracking assets** securely and transparently.

## Table of Contents 📚
- [Prerequisites](#prerequisites)
- [Setup Hyperledger Fabric Test Network](#setup-hyperledger-fabric-test-network)
- [Deploy the Smart Contract](#deploy-the-smart-contract)
- [Develop REST API](#develop-rest-api)
- [Dockerize REST API](#dockerize-rest-api)
- [License](#license)

## Prerequisites ✅

Before you start, make sure you have the following installed on your computer:

- [**Docker**](https://docs.docker.com/get-docker/) 
  
- [**Docker Compose**](https://docs.docker.com/compose/install/)
  
- [**Go Programming Language**](https://golang.org/doc/install) 
  
- [**Node.js** (optional for REST API)](https://nodejs.org/en/download/) 

## Setup Hyperledger Fabric Test Network 🛠️

Run the following commands in your terminal to set up the test network:

```bash
# Clone the Fabric Samples repository
git clone https://github.com/hyperledger/fabric-samples.git

# Navigate into the directory
cd fabric-samples

# Start the test network
cd test_network
./network.sh up

# Create a channel
./network.sh createChannel

# Deploy the chaincode
./network.sh deployCC 
