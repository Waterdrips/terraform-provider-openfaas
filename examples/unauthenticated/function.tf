resource "openfaas_function" "function_test" {
  name            = "testaccopenfaasfunction-basic-y2uzjhk1q1"
  image           = "alpine:latest"
  f_process       = "sha512sum"
}