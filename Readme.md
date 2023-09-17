<div align="center">
<h1 align="center">
<br>Tezos-Delegation-Service
</h1>
<h3>◦ Track and display delegations</h3>
<h3>◦ Developed with the software and tools listed below.</h3>

<p align="center">
<img src="https://img.shields.io/badge/containerd-575757.svg?style&logo=containerd&logoColor=white" alt="containerd" />
<img src="https://img.shields.io/badge/Docker-2496ED.svg?style&logo=Docker&logoColor=white" alt="Docker" />
<img src="https://img.shields.io/badge/Go-00ADD8.svg?style&logo=Go&logoColor=white" alt="Go" />
</p>
</div>

---

## 📒 Table of Contents
- [📒 Table of Contents](#-table-of-contents)
- [📍 Overview](#-overview)
- [📂 Project Structure](#-project-structure)
- [🚀 Getting Started](#-getting-started)
  - [✔️ Prerequisites](#️-prerequisites)
  - [📦 Installation](#-installation)
  - [🎮 Using Tezos-Delegation-Service](#-using-tezos-delegation-service)
  - [🧪 Running Tests](#-running-tests)
  - [🧪 Stop the service](#-stop-the-service)
  - [🧪 Cleaning/Uninstalling](#-cleaninguninstalling)

---


## 📍 Overview

The Tezos Delegation Service project aims to provide an API for retrieving and managing Tezos protocol delegations. It's polling the Tezos blockchain, retrieving and storing delegation data in a PostgreSQL database, and exposing this data through API endpoints.

You can have a look at the [swagger](https://tezos-delegation-api.ew.r.appspot.com/swagger/index.html#/default/get-delegations)

---


## 📂 Project Structure

```bash
Tezos-api
├── cmd
│   ├── api
│   │   └── main.go
│   └── cron
├── config
├── domain
│   ├── adapter
│   ├── entity
│   ├── repository
│   └── usecase
├── infrastructure
│   ├── adapter
│   │   └── tezos
│   └── repository
└── migration
```

---

## 🚀 Getting Started

### ✔️ Prerequisites

Before you begin, ensure that you have the following prerequisites installed:
> - `ℹ️ docker`
> - `ℹ️ docker-compose`
> - `ℹ️ go1.20` in order to launch tests
> - `ℹ️ httpie/ curl`
> - `ℹ️ make`

### 📦 Installation

1. Clone the Tezos-Delegation-Service repository:
```sh
git clone https://github.com/frisk038/Tezos-Delegation-Service
```

2. Change to the project directory:
```sh
cd Tezos-Delegation-Service
```

3. Install the dependencies:
```sh
docker-compose up
```
It will pop a postgres db container, will build the service inside a go container. And start listening on port **:8080**

### 🎮 Using Tezos-Delegation-Service

```sh
http localhost:8080/xtz/delegations
```
This command will return the last delegations on tezos blockchain.

### 🧪 Running Tests
```sh
make test
```
This one launch every test of the project

### 🧪 Stop the service
```sh
make down
```

### 🧪 Cleaning/Uninstalling
```sh
make clean
```
This command will stop/remove the containers and images