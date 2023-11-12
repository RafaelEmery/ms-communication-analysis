from locust import HttpUser, task

class ClientBFFUser(HttpUser):
    # TODO: understand wait_time function
    # wait_time = between(1, 3)

    @task
    def do_client_server_interaction(self):
        headers = {
            "Content-Type": "application/json"
        }
        payload = {
            "resource": "http",
            "communication_method": "create",
            "request_quantity": 10
        }

        print(f"Payload: {payload}")

        res = self.client.post("/interact", json=payload, headers=headers)
        
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")