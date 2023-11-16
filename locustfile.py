from locust import HttpUser, task, between, TaskSet

class BFFUserTaskSet(TaskSet):
    def on_start(self):
        self.common_resource = "create"
        self.request_count = 0
        self.request_threshold = 1000

    @task(1)
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
        
        # TODO: fix the quit for threshold
        self.request_count += 1
        if self.request_count >= self.request_threshold:
            self.user.environment.runner.quit()
    
    # @task(2)
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
    #         self.user.environment.runner.quit()

    # @task(3)
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
    #         self.user.environment.runner.quit()


class BFFUser(HttpUser):
    tasks = [BFFUserTaskSet]
    wait_time = between(1, 3) # Simulate time between tasks

    # TODO: set request quantity threshold to adapt to experiment workload levels
    # TODO: set order to tasks to test each application individually