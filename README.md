# Concurrent Load-Balancing Reverse Proxy with Health Monitoring

---

## 🚀 Project Roadmap & Progress

### Phase 1: Core Architecture (The Foundation)
- [x] **1. API Design & Mocking**
    - [x] Create a single Backend Server instance.
    - [x] Implement standard CRUD operations (Create, Read, Update, Delete).
    - [x] Build a simple Client to send requests to this API.
- [x] **2. Basic Reverse Proxy Implementation**
    - [x] Establish a Proxy-to-Server connection.
    - [x] Route: `Client` → `Proxy` → `Server`.
    - [x] Route: `Server` → `Proxy` → `Client`.
    - [x] Ensure request/response headers are preserved during the hop.

### Phase 2: Multi-Server & Load Balancing
- [ ] **3. Scaling the Backend**
    - [x] Spin up multiple Backend Server instances.
    - [x] Implement **Round-Robin** selection logic.
    - [x] Implement **Least-Connections** selection logic.
    - [ ] *Optional:* Add performance benchmarking to compare strategy efficiency.
- [ ] **4. Active Health Monitoring**
    - [x] Build an Admin API to manually `Add`, `Remove`, or `Check` server status.
    - [x] Implement an automated "Pulse" check (Active Health Check) every 5 minutes.
    - [x] Implement **Failover Logic**: Automatically bypass down servers and redirect to functional ones.

### Phase 3: Advanced Traffic Management
- [ ] **5. Persistence & Hardware Optimization**
    - [ ] Implement **Sticky Sessions**: Use IP Hashing or Cookies to pin clients to specific backends.
    - [ ] Implement **Weighted Load Balancing**: Assign $Weight_i$ to backends based on simulated capacity.
- [ ] **6. Security & Hardening**
    - [ ] Enable **SSL Termination**: Proxy handles decryption (HTTPS → HTTP) and encryption.
    - [ ] Implement **Rate Limiting**: Prevent bot spam by limiting clients to $X$ requests per minute (e.g., 20 RPM).

### Phase 4: Observability & Real-World Application
- [ ] **7. Automation & Analytics**
    - [ ] Create a Traffic Simulator for automated client request testing.
    - [ ] Build an Analytics Engine to track:
        - [ ] Peak and Lowest traffic periods.
        - [ ] HTTP Status Code distribution (200s, 404s, 500s).
        - [ ] Hourly/Daily summary reports.
- [ ] **8. Real Case Scenario: The Bookstore API**
    - [ ] Migrate the simple CRUD API to a complex Bookstore system.
    - [ ] Stress test all features (SSL, Load Balancing, Health Checks) under realistic workloads.

---

## 📂 Project Structure

```text
/Concurrent-Load-Balancing-Reverse-Proxy
├── README.md               <-- Roadmap and Documentation
├── go.mod                  <-- Go dependencies and module info
│
├── /cmd                    <-- Main application entry points
│   ├── /proxy              <-- The Reverse Proxy application
│   │   └── main.go
│   └── /backend            <-- The Test Backend API (Phase 1)
│       └── main.go
│
├── /internal               <-- Private project logic
│   ├── /balancer           <-- Load balancing strategies (Step 3 & 5)
│   ├── /health             <-- Health check and failover logic (Step 4)
│   ├── /security           <-- SSL/TLS and Rate limiting (Step 6)
│   └── /analytics          <-- Metrics and logging (Step 7)
│
├── /certs                  <-- SSL/TLS Certificates
└── /scripts                <-- Automation and testing scripts