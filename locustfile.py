from locust import HttpUser, task

class BFFUser(HttpUser):
    @task
    def do_client_server_interaction(self):
        headers = {
            "Content-Type": "application/json"
        }
        payload = {
            "resource": "create",
            "request_quantity": 1
        }

        print(f"Payload: {payload}")

        res = self.client.post("/interact/http", json=payload, headers=headers)
        
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")

    # @task()
    # def do_client_server_interaction_grpc(self):
    #     headers = {
    #         "Content-Type": "application/json"
    #     }
    #     payload = {
    #         "resource": self.common_resource,
    #         "request_quantity": 1
    #     }

    #     print(f"Payload: {payload}")

    #     res = self.client.post("/interact/grpc", json=payload, headers=headers)
        
    #     print(f"Status Code: {res.status_code}")
    #     print(f"Response Content: {res.content}")

    #     self.request_count += 1
    #     if self.request_count >= self.request_threshold:
    #         self.environment.runner.quit()

    # @task()
    # def do_client_server_interaction_rabbit_mq(self):
    #     headers = {
    #         "Content-Type": "application/json"
    #     }
    #     payload = {
    #         "resource": self.common_resource,
    #         "request_quantity": 1
    #     }

    #     print(f"Payload: {payload}")

    #     res = self.client.post("/interact/rabbitmq", json=payload, headers=headers)
        
    #     print(f"Status Code: {res.status_code}")
    #     print(f"Response Content: {res.content}")

    #     self.request_count += 1
    #     if self.request_count >= self.request_threshold:
    #         self.environment.runner.quit()