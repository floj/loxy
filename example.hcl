frontend "http" {
  port = 8080
  # tls {
  #     cert = "./cert.pem"
  #     key = "./cert.key"
  # }

  middleware {
    proxy_headers {}
    logger {}
  }

  route "local" {
    backend = "local"

    match {
      path {
        has_any_prefix = ["/local"]
      }
    }

    modify {
      path {
        strip_prefix = ["/local"]
      }
    }
  }

  route "localfile" {
    backend = "localfile"
    match {
      path {
        has_any_prefix = ["/files"]
      }
    }
    modify {
      path {
        strip_prefix = ["/files"]
      }
    }
  }

  route "google" {
    backend = "google"
    match {
      path {
        has_any_prefix = ["/google"]
      }
    }

    modify {
      path {
        strip_prefix = ["/google"]
      }
    }
  }

}

frontend "http8888" {
  port = 8888
  # tls {
  #     cert = "./cert.pem"
  #     key = "./cert.key"
  # }

  middleware {
    proxy_headers {}
    logger {}
  }

  route "metrics" {
    backend = "metrics"

    match {
      path {
        is = "/metrics"
      }
    }
  }
}

backend "local" {
  reverse_proxy {
    targets = [
      "http://localhost:8088"
    ]
  }
}

backend "google" {
  reverse_proxy {
    targets = [
      "https://google.com"
    ]
  }
}

backend "localfile" {
  file_server {
    root = "./"
  }
}

backend "metrics" {
  prometheus {}
}
