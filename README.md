# Conduit — E-Commerce Microservices Platform

> Status: ✅ **Production-ready (development complete)**

## Overview

**Conduit** is a fully implemented **e-commerce microservices platform** built with an **architecture-first, backend-driven approach**.

The system demonstrates modern e-commerce backend design: clear service boundaries, explicit communication contracts, and scalable infrastructure with minimal technical debt.

Each service is independently deployable and communicates over well-defined interfaces.

---

## Services

The platform consists of the following core services:

- **API Service**
  - Public-facing REST API
  - Entry point for all client requests
  - Delegates business logic to internal services

- **gRPC Service**
  - Central business logic layer
  - Handles domain operations (orders, users, sessions, etc.)
  - Exposes internal APIs via gRPC

- **Notification Service**
  - Asynchronous event-driven service
  - Handles notifications (e.g., order events)
  - Fully decoupled from core business logic

---

## Architecture

Conduit follows industry-proven backend design principles:

- **Microservices Architecture**
- **Clean Architecture**
- **Hexagonal (Ports & Adapters) Architecture**
- **Explicit service-to-service communication (REST + gRPC)**
- **Separation of infrastructure, domain, and application layers**

---

## Tech Stack

- **Go**
- **REST & gRPC**
- **Docker & Docker Compose**
- **MySQL**
- **Clean / Hexagonal Architecture**

---

## Project Goals

- Showcase real-world microservice design
- Avoid tight coupling and shared state
- Model production-grade backend systems
- Provide a foundation for future extensions (auth, payments, search, etc.)

---

## Status

All planned services are implemented and fully integrated.  
The project is **feature-complete** and ready for deployment or further extension.
