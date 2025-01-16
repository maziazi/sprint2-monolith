# FitByte

Track your progress, show your growth!

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)


## Introduction
This project is a Tracking Activities with API. Built using Go and initialized using GoLand.

## Features
- **Authentication & Authorization**
  - **POST** /v1/login
  - **POST** /v1/register
- **Profile Management**
   - **GET** /v1/user
   - **PATCH** /v1/user
- **File Upload**
   - **POST** /v1/file
- **Activity**
   - **POST** /v1/activity
   - **GET** /v1/activity
   - **PATCH** /v1/activity/:activityId
   - **DELETE** /v1/activity/:activityId
## Installation

### Prerequisites
Make sure you have the following installed:
- [Go](https://golang.org/dl/) (version 1.x or higher)
- GoLand IDE (optional but recommended)

### Steps
1. Clone this repository:
   ```bash
   git clone https://github.com/maziazi/sprint2-monolith.git
   cd FitByte
   ```
   
2. Kalau Perlu Init Dulu
    ```bash 
    go mod init fitbyte
    ```

3. Initialize Go modules (if not already initialized):
   ```bash
   go mod tidy
   ```
4. Install dependencies (if any):
   ```bash
   go get
   ```

## Usage

Run the project using the command:
```bash
cd cmd
go run main.go
```

