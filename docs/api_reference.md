# Model Context Protocol API Reference

This document provides detailed API reference for the Model Context Protocol (MCP) Go SDK.

## Core Package

### ModelRequest

```go
type ModelRequest struct {
    ID         string                 `json:"id"`
    ModelData  map[string]interface{} `json:"modelData"`
    Parameters []Parameter            `json:"parameters"`
}
```

The `ModelRequest` represents a request to process a model. It contains:

- `ID`: A unique identifier for the request
- `ModelData`: A map containing model-specific data
- `Parameters`: A slice of parameters for the request

### ModelResponse

```go
type ModelResponse struct {
    ID           string                 `json:"id"`
    Success      bool                   `json:"success"`
    ErrorMessage string                 `json:"errorMessage,omitempty"`
    Results      map[string]interface{} `json:"results"`
    Timestamp    time.Time              `json:"timestamp"`
}
```

The `ModelResponse` represents the response from processing a model. It contains:

- `ID`: The identifier of the request this response relates to
- `Success`: Whether the request was processed successfully
- `ErrorMessage`: An optional error message when Success is false
- `Results`: A map containing the results of model processing
- `Timestamp`: When the response was generated

### Parameter

```go
type Parameter struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"`
}
```

The `Parameter` represents a parameter for model processing. It contains:

- `Name`: The name of the parameter
- `Value`: The value of the parameter
- `Type`: The data type of the parameter

### Status

```go
type Status string

const (
    StatusIdle      Status = "idle"
    StatusStarting  Status = "starting"
    StatusRunning   Status = "running"
    StatusStopping  Status = "stopping"
    StatusError     Status = "error"
)
```

The `Status` represents the state of an MCP component.

### StatusChangeEvent

```go
type StatusChangeEvent struct {
    OldStatus Status
    NewStatus Status
    Error     error
}
```

The `StatusChangeEvent` represents a change in component status. It contains:

- `OldStatus`: The previous status
- `NewStatus`: The new status
- `Error`: An optional error that caused the status change

### Component

```go
type Component interface {
    Start() error
    Stop() error
    Status() Status
    OnStatusChange(func(StatusChangeEvent))
}
```

The `Component` interface defines the basic lifecycle methods for MCP components.

## Client Package

### Client

```go
type Client struct {
    // Fields omitted for brevity
}

func New(options ...Option) *Client
func (c *Client) Start() error
func (c *Client) Stop() error
func (c *Client) Status() core.Status
func (c *Client) OnStatusChange(func(core.StatusChangeEvent))
func (c *Client) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error)
```

The `Client` is the main entry point for MCP clients. It implements the `core.Component` interface.

### Options

```go
type Options struct {
    ServerHost           string
    ServerPort           int
    ConnectionTimeout    time.Duration
    AutoReconnect        bool
    MaxReconnectAttempts int
    ReconnectDelay       time.Duration
    EnableTLS            bool
}

func DefaultOptions() Options
func WithServerHost(host string) Option
func WithServerPort(port int) Option
func WithConnectionTimeout(timeout time.Duration) Option
func WithAutoReconnect(enabled bool) Option
func WithMaxReconnectAttempts(attempts int) Option
func WithReconnectDelay(delay time.Duration) Option
func WithTLS(enabled bool) Option
```

The `Options` provide configuration for an MCP client.

## Server Package

### Server

```go
type Server struct {
    // Fields omitted for brevity
}

func New(options ...Option) *Server
func (s *Server) Start() error
func (s *Server) Stop() error
func (s *Server) Status() core.Status
func (s *Server) OnStatusChange(func(core.StatusChangeEvent))
func (s *Server) RegisterHandler(handler Handler) error
```

The `Server` is the main entry point for MCP servers. It implements the `core.Component` interface.

### Handler

```go
type Handler interface {
    Methods() []string
}
```

The `Handler` interface defines the basic methods that all handlers must implement.

### ModelHandler

```go
type ModelHandler interface {
    Handler
    ProcessModel(context.Context, *core.ModelRequest) (*core.ModelResponse, error)
}
```

The `ModelHandler` interface defines a handler for model processing requests.

### DefaultModelHandler

```go
type DefaultModelHandler struct{}

func NewDefaultModelHandler() *DefaultModelHandler
func (h *DefaultModelHandler) Methods() []string
func (h *DefaultModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error)
```

The `DefaultModelHandler` provides a simple implementation of the `ModelHandler` interface.

### Options

```go
type Options struct {
    Host                 string
    Port                 int
    MaxConcurrentClients int
    ConnectionTimeout    time.Duration
    EnableTLS            bool
    CertificatePath      string
    CertificateKeyPath   string
}

func DefaultOptions() Options
func WithHost(host string) Option
func WithPort(port int) Option
func WithMaxConcurrentClients(max int) Option
func WithConnectionTimeout(timeout time.Duration) Option
func WithTLS(enabled bool) Option
func WithCertificatePath(path string) Option
func WithCertificateKeyPath(path string) Option
```

The `Options` provide configuration for an MCP server.
