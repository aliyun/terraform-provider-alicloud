package fc

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"context"
	"time"

	"github.com/gorilla/websocket"
)

const (
	listInstancesPathName          = "/instances"
	listInstancesSourcesPath       = singleFunctionPath + listInstancesPathName
	listInstancesWithQualifierPath = singleFunctionWithQualifierPath + listInstancesPathName

	instanceExecPathName          = "/instances/%s/exec"
	instanceExecSourcesPath       = singleFunctionPath + instanceExecPathName
	instanceExecWithQualifierPath = singleFunctionWithQualifierPath + instanceExecPathName

	messageStdin  byte = 0
	messageStdout byte = 1
	messageStderr byte = 2
	serverErr     byte = 3
	reSize        byte = 4
	serverClose   byte = 5
)

type MessageCallbackFunction func(data []byte)

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}


type InstanceExecOutput struct {
	Header              http.Header
	WebsocketConnection *websocket.Conn
	ctx                 context.Context
	cancel              context.CancelFunc
	errChan             chan error
	lock                sync.Mutex
}

func (o *InstanceExecOutput) start(input *InstanceExecInput) error {
	o.ctx, o.cancel = context.WithCancel(context.Background())
	if o.WebsocketConnection == nil {
		return fmt.Errorf("WebSocket is not initialized")
	}
	o.errChan = make(chan error)
	go o.reader(input.onStdout, input.onStderr)
	o.WebsocketConnection.SetPingHandler(func(data string) error {
		o.lock.Lock()
		defer o.lock.Unlock()
		return o.WebsocketConnection.WriteControl(websocket.PingMessage, []byte(data), time.Now().Add(time.Second*5))
	})
	return nil
}

func (o *InstanceExecOutput) writeMessage(messageType int, data []byte) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	return o.WebsocketConnection.WriteMessage(messageType, data)
}

func (o *InstanceExecOutput) reader(stdout, stderr MessageCallbackFunction) {
	for {
		select {
		case <-o.ctx.Done():
			return
		default:
			_, msg, err := o.WebsocketConnection.ReadMessage()
			if err != nil {
				o.reportError(err)
			}
			if len(msg) > 0 {
				var out MessageCallbackFunction
				switch msg[0] {
				case messageStdout:
					out = stdout
				case messageStderr:
					out = stderr
				case serverErr:
					o.reportError(fmt.Errorf(string(msg)))
				default:
					o.reportError(fmt.Errorf("unknown message type %d, message: %s", msg[0], msg))
				}
				data := msg[1:]
				if out != nil && len(data) > 0 {
					out(data)
				}
			}
		}
	}
}

func (o *InstanceExecOutput) reportError(err error) {
	if err != nil {
		select {
		case o.errChan <- err:
		default:
		}
		o.cancel()
	}
}

func (o *InstanceExecOutput) Closed() bool {
	select {
	case <-o.ctx.Done():
		return true
	default:
		return false
	}
}

func (o *InstanceExecOutput) ErrorChannel() chan error {
	return o.errChan
}

func (o *InstanceExecOutput) WriteStdin(msg []byte) error {
	if o.Closed() {
		return fmt.Errorf("WebSocket has closed")
	}
	return o.writeMessage(websocket.TextMessage, append([]byte{messageStdin}, msg...))
}

// InstanceExecInput define publish layer version response
type InstanceExecInput struct {
	Header       http.Header
	ServiceName  *string
	FunctionName *string
	Qualifier    *string
	InstanceID   *string

	Command     []string
	Stdin       bool
	Stdout      bool
	Stderr      bool
	TTY         bool
	IdleTimeout *int

	onStdout MessageCallbackFunction
	onStderr MessageCallbackFunction
}

func NewInstanceExecInput(
	serviceName, functionName, instanceID string, command []string,
) *InstanceExecInput {
	return &InstanceExecInput{
		ServiceName:  &serviceName,
		FunctionName: &functionName,
		InstanceID:   &instanceID,
		Command:      command,
	}
}

func (i *InstanceExecInput) WithServiceName(serviceName string) *InstanceExecInput {
	i.ServiceName = &serviceName
	return i
}

func (i *InstanceExecInput) WithFunctionName(functionName string) *InstanceExecInput {
	i.FunctionName = &functionName
	return i
}

func (i *InstanceExecInput) WithQualifier(qualifier string) *InstanceExecInput {
	i.Qualifier = &qualifier
	return i
}

func (i *InstanceExecInput) WithInstance(instanceID string) *InstanceExecInput {
	i.InstanceID = &instanceID
	return i
}

func (i *InstanceExecInput) WithCommand(command []string) *InstanceExecInput {
	i.Command = command
	return i
}

func (i *InstanceExecInput) WithStdin(stdin bool) *InstanceExecInput {
	i.Stdin = stdin
	return i
}

func (i *InstanceExecInput) WithStdout(stdout bool) *InstanceExecInput {
	i.Stdout = stdout
	return i
}

func (i *InstanceExecInput) WithStderr(stderr bool) *InstanceExecInput {
	i.Stderr = stderr
	return i
}

func (i *InstanceExecInput) WithTTY(tty bool) *InstanceExecInput {
	i.TTY = tty
	return i
}

func (i *InstanceExecInput) WithIdleTimeout(idleTimeout int) *InstanceExecInput {
	i.IdleTimeout = &idleTimeout
	return i
}

func (i *InstanceExecInput) OnStdout(callback MessageCallbackFunction) *InstanceExecInput {
	i.onStdout = callback
	return i
}

func (i *InstanceExecInput) OnStderr(callback MessageCallbackFunction) *InstanceExecInput {
	i.onStderr = callback
	return i
}

func (i *InstanceExecInput) GetQueryParams() url.Values {
	queries := make(url.Values)
	for _, cmd := range i.Command {
		queries.Add("command", cmd)
	}
	queries.Add("stdin", boolToString(i.Stdin))
	queries.Add("stdout", boolToString(i.Stdout))
	queries.Add("stderr", boolToString(i.Stderr))
	queries.Add("tty", boolToString(i.TTY))
	if i.IdleTimeout != nil {
		queries.Add("idleTimeout", fmt.Sprint(*i.IdleTimeout))
	}
	return queries
}

func (i *InstanceExecInput) GetPath() string {
	if i.Qualifier != nil {
		return fmt.Sprintf(instanceExecWithQualifierPath, pathEscape(*i.ServiceName), pathEscape(*i.Qualifier), pathEscape(*i.FunctionName), pathEscape(*i.InstanceID))
	} else {
		return fmt.Sprintf(instanceExecSourcesPath, pathEscape(*i.ServiceName), pathEscape(*i.FunctionName), pathEscape(*i.InstanceID))
	}
}

func (i *InstanceExecInput) GetHeaders() Header {
	return make(Header)
}

func (i *InstanceExecInput) GetPayload() interface{} {
	return nil
}

func (i *InstanceExecInput) Validate() error {
	if IsBlank(i.ServiceName) {
		return fmt.Errorf("Service name is required but not provided")
	}
	if IsBlank(i.FunctionName) {
		return fmt.Errorf("Function name is required but not provided")
	}
	if IsBlank(i.InstanceID) {
		return fmt.Errorf("InstanceID is required but not provided")
	}
	if i.Command == nil || len(i.Command) == 0 {
		return fmt.Errorf("Command is required but not provided")
	}
	return nil
}
