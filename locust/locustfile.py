from locust import HttpUser, between, task
import os

API_USER = os.getenv("LOCUST_USER", "admin")
API_PASS = os.getenv("LOCUST_PASS", "admin123")


class ApiUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        self.token = ""
        self.login()

    def login(self):
        payload = {"username": API_USER, "password": API_PASS}
        with self.client.post("/api/v1/auth/login", json=payload, catch_response=True) as resp:
            if resp.status_code != 200:
                resp.failure("login failed")
                return
            data = resp.json().get("data", {})
            self.token = data.get("access_token", "")
            if not self.token:
                resp.failure("missing access_token")

    def auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"}

    @task(2)
    def list_users(self):
        if not self.token:
            self.login()
            return
        self.client.get("/api/v1/users", headers=self.auth_headers())
