terraform {
  required_providers {
    movies = {
      source = "hashicorp.com/edu/movies"
    }
  }
}

provider "movies" {
  host = "localhost"
  port = "8080"
}

