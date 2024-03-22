package metadata

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"octopus/config"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type OptionBuilder[T any] func(*T)

func LoadOption[T any](t *T, builders ...OptionBuilder[T]) {
	for _, builder := range builders {
		builder(t)
	}
}

type MetaData struct {
	*HttpMeta
	*GrpcMeta
	Descriptor *Descriptor
	Logger     *zerolog.Logger
	Target     string
	Result     any
}

type HttpMeta struct {
	Request        *http.Request
	Payload        map[string]any
	Response       http.ResponseWriter
	Callbackheader *metadata.MD
}

type GrpcMeta struct {
	Header      *metadata.MD
	GrpcContext context.Context
}

type PathHandler func(string, *zerolog.Logger) (*URI, error)

type ErrorMeta struct {
	Error string `json:"error"`
}

type URI struct {
	HttpMethod  string
	ServiceName string
	Method      string
	Host        string
	Params      map[string]any
}

type Descriptor struct {
	*URI
	RequestMessage  string
	ResponseMessage string
}

func (d *Descriptor) convertToMessage(dic map[string]proto.Message) (proto.Message, proto.Message, error) {
	var (
		reqIn, resOut proto.Message
		ok            bool
	)
	if reqIn, ok = dic[d.RequestMessage]; !ok {
		return nil, nil, fmt.Errorf(config.NOPROTOMESSAGE, d.RequestMessage)
	} else if resOut, ok = dic[d.ResponseMessage]; !ok {
		return nil, nil, fmt.Errorf(config.NOPROTOMESSAGE, d.ResponseMessage)
	}
	req := reflect.New(reflect.TypeOf(reqIn).Elem()).Interface()
	res := reflect.New(reflect.TypeOf(resOut).Elem()).Interface()

	in := req.(proto.Message)
	out := res.(proto.Message)
	return in, out, nil
}

func (d *Descriptor) GetFullMethod() string {
	return fmt.Sprintf("/%v/%v", d.ServiceName, d.Method)
}

func (m *MetaData) FormatAll(handler PathHandler) error {
	path := m.Request.URL.Path
	path = path[strings.Index(path, "/")+1:]
	uri, err := handler(path, m.Logger)
	if err == nil {
		uri.HttpMethod = m.Request.Method
		m.Descriptor = &Descriptor{
			URI: uri,
		}
		m.FormatPayload()
	}
	return err
}

func DefaultPathHandler(path string, logger *zerolog.Logger) (*URI, error) {
	st := strings.Split(path, "/")
	if len(st) != 2 {
		return nil, fmt.Errorf(config.WRONGPATH, path)
	}
	return &URI{
		ServiceName: strings.Replace(st[0], "-", ".", 1),
		Method:      st[1],
	}, nil
}

func PathMatcher(key, path string, keytype config.ParamType) (*URI, error) {
	urlcontext := strings.Split(path, "/")
	if len(urlcontext) > 2 {
		return nil, fmt.Errorf(config.WRONGPATH, path)
	}
	st := strings.Split(urlcontext[0], "-")
	if len(st) != 3 {
		return nil, fmt.Errorf(config.WRONGPATH, path)
	}
	params := make(map[string]any)
	if len(urlcontext) == 2 && urlcontext[1] != "" {
		switch keytype {
		case config.Int:
			if paramvalue, err := strconv.Atoi(urlcontext[1]); err != nil {
				return nil, fmt.Errorf(config.WRONGPATH, path)
			} else {
				params[key] = paramvalue
			}
		case config.Float:
			if paramvalue, err := strconv.ParseFloat(urlcontext[1], 64); err != nil {
				return nil, fmt.Errorf(config.WRONGPATH, path)
			} else {
				params[key] = paramvalue
			}
		default:
			params[key] = urlcontext[1]
		}
	}

	return &URI{
		ServiceName: fmt.Sprintf("%v.%v", st[0], st[1]),
		Method:      st[2],
		Params:      params,
	}, nil
}

func (m *MetaData) FormatPayload() {
	m.Request.ParseForm()
	payload := make(map[string]any)
	fristload := make(map[string]any)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	for key, v := range m.Request.Form {
		var data map[string]any
		err := json.Unmarshal([]byte(key), &data)
		if err == nil {
			for kk, vv := range data {
				fristload[kk] = vv
			}
		} else {
			if len(v) > 0 {
				fristload[key] = v[0]
			} else {
				fristload[key] = ""
			}
		}
	}
	b, err := io.ReadAll(m.Request.Body)
	defer m.Request.Body.Close()
	if len(b) != 0 && err == nil {
		err = json.Unmarshal(b, &payload)
		if err != nil {
			m.Logger.Panic().Msg(config.WRONGPATHPATTERN)
		}
		for k, v := range fristload {
			payload[k] = v
		}
	} else {
		payload = fristload
	}

	for key, value := range m.Descriptor.Params {
		payload[key] = value
	}

	m.Payload = payload
}

func (m *MetaData) GetProtoMessage(dic map[string]proto.Message) (proto.Message, proto.Message, error) {
	reqIn, resOut, err := m.Descriptor.convertToMessage(dic)
	if err != nil {
		return nil, nil, err
	}
	if len(m.Payload) > 0 {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		b, err := json.Marshal(m.Payload)
		if err != nil {
			return nil, nil, err
		}
		err = json.Unmarshal(b, reqIn)

		if err != nil {
			return nil, nil, err
		}
	}

	return reqIn, resOut, nil
}

func (m *MetaData) Errorf(errormsg string, isDebug bool) {
	if !isDebug {
		errormsg = config.SYSTEMERROR
	} else {
		m.Logger.Error().Msg(errormsg)
	}
	m.Logger.Error().Msg(LoggerTrace())

	m.Result = ErrorMeta{
		Error: errormsg,
	}
}

func LoggerTrace() string {
	buf := make([]byte, 64<<10)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	return string(buf)
}

type ProtoTable map[string]proto.Message

func (table ProtoTable) AddProtoMessage(messageName string, logger *zerolog.Logger) error {
	if _, ok := table[messageName]; !ok {
		defer func() {
			if e := recover(); e != nil {
				logger.Error().Any("Panic", e).Msg(fmt.Sprintf("%v package is not exits", messageName))
				logger.Error().Msg(LoggerTrace())
			}
		}()
		messageType := reflect2.TypeByName(messageName).Type1()
		if message, ok := reflect.New(messageType).Interface().(proto.Message); ok {
			table[messageName] = message
			return nil
		} else {
			return fmt.Errorf("%v package is not exits", messageName)
		}
	}
	return nil
}
