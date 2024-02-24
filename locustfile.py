from locust import HttpUser, task

class BFFUser(HttpUser):
    @task
    def do_client_server_interaction(self):
        headers = {
            "Content-Type": "application/json"
        }
        """
        The "resource" key specify the feature on the server application.
        The "request_quantity" key specify the number of requests to be made to the server.
        """
        payload = {
            "resource": "getByDiscount",
            "request_quantity": 1
        }

        """
        The /interact/<method> endpoint is used trigger BFF and specified server.
        The options to /interact/<method> are http|grpc|rabbitmq
        """
        res = self.client.post("/interact/grpc", json=payload, headers=headers)
        
        """
        The response content will be printed to the console.
        Contains useful info for debugging.
        """
        print(f"Status Code: {res.status_code}")
        print(f"Response Content: {res.content}")