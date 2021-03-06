variable "key" {}
variable "secret" {}
variable "key_pair" {}


// hostnames are used as a source
variable "hostnames" {
  type = "list"
  default = [
    "huey",
    "dewey",
    "louie"
  ]
}

variable "zone" {
  default = "ch-dk-2"
}

variable "template" {
  default = "Linux Ubuntu 17.10 64-bit"
}

