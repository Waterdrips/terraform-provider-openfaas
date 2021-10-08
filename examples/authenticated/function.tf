resource "openfaas_function" "function_test" {
  name      = "test-function"
  image     = "alpine:latest"
  f_process = "env"
  labels = {
    Group       = "London"
    Environment = "Test"
  }

  limits {
    memory = "20m"
    cpu    = "100m"
  }

  secrets = {
    superSecret = "secret value!!"
  }

  annotations = {
    CreatedDate = "Mon Sep  3 07:15:55 BST 2018"
  }
}

