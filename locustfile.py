from locust import HttpUser, task

class ClientBFFUser(HttpUser):
    # TODO: understand wait_time function
    # wait_time = between(1, 3)
    common_resource = "create"

    # TODO: set request quantity threshold to adapt to experiment workload levels
    # TODO: set order to tasks to test each application individually

    @task
    def do_client_server_interaction_http(self):
        headers = {
            "Content-Type": "application/json"
        }
        payload = {
            "resource": self.common_resource,
            "request_quantity": 1
        }

        print(f"Payload: {payload}")

        res = self.client.post("/interact/http", json=payload, headers=headers)
        
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")

    @task
    def do_client_server_interaction_grpc(self):
        headers = {
            "Content-Type": "application/json"
        }
        payload = {
            "resource": self.common_resource,
            "request_quantity": 1
        }

        print(f"Payload: {payload}")

        res = self.client.post("/interact/grpc", json=payload, headers=headers)
        
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")

    @task
    def do_client_server_interaction_rabbit_mq(self):
        headers = {
            "Content-Type": "application/json"
        }
        payload = {
            "resource": self.common_resource,
            "request_quantity": 1
        }

        print(f"Payload: {payload}")

        res = self.client.post("/interact/rabbitmq", json=payload, headers=headers)
        
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")