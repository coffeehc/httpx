package httpx

import (
  "bytes"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/gogo/protobuf/proto"
  "github.com/json-iterator/go"
)

var (
  BindingProtoBuf      = protobufBinding{}
  BindingJSON          = jsonBinding{}
  BindingXML           = binding.XML
  BindingForm          = binding.Form
  BindingQuery         = binding.Query
  BindingFormPost      = binding.FormPost
  BindingFormMultipart = binding.FormMultipart
  BindingMsgPack       = binding.MsgPack
  BindingYAML          = binding.YAML
  BindingUri           = binding.Uri
)

func Bind(c *gin.Context, obj interface{}) error {
  return c.ShouldBindWith(obj, GetBinding(c))
}

func GetBinding(c *gin.Context) binding.Binding {
  method, contentType := c.Request.Method, c.ContentType()
  if method == "GET" {
    return BindingForm
  }

  switch contentType {
  case binding.MIMEJSON:
    return BindingJSON
  case binding.MIMEXML, binding.MIMEXML2:
    return BindingXML
  case binding.MIMEPROTOBUF:
    return BindingProtoBuf
  case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
    return BindingMsgPack
  case binding.MIMEYAML:
    return BindingYAML
  default: // case MIMEPOSTForm, MIMEMultipartPOSTForm:
    return BindingForm
  }
}

type protobufBinding struct{}

func (protobufBinding) Name() string {
  return "protobuf"
}

func (b protobufBinding) Bind(req *http.Request, obj interface{}) error {
  buf, err := ioutil.ReadAll(req.Body)
  if err != nil {
    return err
  }
  return b.BindBody(buf, obj)
}

func (protobufBinding) BindBody(body []byte, obj interface{}) error {
  if err := proto.Unmarshal(body, obj.(proto.Message)); err != nil {
    return err
  }
  // Here it's same to return validate(obj), but util now we can't add
  // `binding:""` to the struct which automatically generate by gen-proto
  return nil
  // return validate(obj)
}

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonBinding struct{}

func (jsonBinding) Name() string {
  return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
  if req == nil || req.Body == nil {
    return fmt.Errorf("invalid request")
  }
  return decodeJSON(req.Body, obj)
}

func (jsonBinding) BindBody(body []byte, obj interface{}) error {
  return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
  decoder := json.NewDecoder(r)
  if EnableDecoderUseNumber {
    decoder.UseNumber()
  }
  if err := decoder.Decode(obj); err != nil {
    return err
  }
  return validate(obj)
}

func validate(obj interface{}) error {
  if binding.Validator == nil {
    return nil
  }
  return binding.Validator.ValidateStruct(obj)
}
