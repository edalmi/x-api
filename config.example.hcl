"cache" "redis" {
  "address" = "localhost:12345"
  "password" = "abc"
  "db" = 0
}

"logger" "zap" {
  "format" = "json"
}

"db" "sqlite" {
  "path" = "/tmp/db.sqlite"
}

"serve" "admin" {
  "host" = "0.0.0.0"
  "port" = 12340
}

"serve" "metrics" {
  "host" = "0.0.0.0"
  "port" = 12341
}

"serve" "public" {
  "host" = "0.0.0.0"
  "port" = 12342
}

"serve" "healthz" {
  "host" = "0.0.0.0"
  "port" = 12343
}
