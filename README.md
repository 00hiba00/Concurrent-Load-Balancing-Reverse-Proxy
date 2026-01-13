# Concurrent Load-Balancing Reverse Proxy with Health Monitoring

---

## 🚀 Project Roadmap & Progress

### Phase 1: Core Architecture (The Foundation)
- [ ] **1. API Design & Mocking**
    - [ ] Create a single Backend Server instance.
    - [ ] Implement standard CRUD operations (Create, Read, Update, Delete).
    - [ ] Build a simple Client to send requests to this API.
- [ ] **2. Basic Reverse Proxy Implementation**
    - [ ] Establish a Proxy-to-Server connection.
    - [ ] Route: `Client` → `Proxy` → `Server`.
    - [ ] Route: `Server` → `Proxy` → `Client`.
    - [ ] Ensure request/response headers are preserved during the hop.

### Phase 2: Multi-Server & Load Balancing
- [ ] **3. Scaling the Backend**
    - [ ] Spin up multiple Backend Server instances.
    - [ ] Implement **Round-Robin** selection logic.
    - [ ] Implement **Least-Connections** selection logic.
    - [ ] *Optional:* Add performance benchmarking to compare strategy efficiency.
- [ ] **4. Active Health Monitoring**
    - [ ] Build an Admin API to manually `Add`, `Remove`, or `Check` server status.
    - [ ] Implement an automated "Pulse" check (Active Health Check) every 5 minutes.
    - [ ] Implement **Failover Logic**: Automatically bypass down servers and redirect to functional ones.

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
/concurrent-proxy
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