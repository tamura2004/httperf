package domain

type Config struct {
  Client   ClientConfig
  Scinario ScinarioConfig
  Target   TargetConfig
}

type ClientConfig struct {
  Bps      int
  Proxy    string
  Header   map[string]string
  UserName string
  Password string
}

type ScinarioConfig struct {
  Count     int
  Interval  string
  RampUp    string
  WorkerNum int
  Timeout   string
}

type TargetConfig struct {
  Url []string
}
