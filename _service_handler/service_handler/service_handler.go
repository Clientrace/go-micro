package _service_handler

/* Type Map Reference */
var ATTRIB_TYPE_MAP = map[string]interface{}{
	"string":  []string{"string"},
	"boolean": []string{"bool"},
	"number":  []string{"int", "float64"},
}

/* Spec Event Attribute */
type EventAttrib struct {
	name  string
	dType string
}

/* Required Event specification */
type ReqEventSpec struct {
	ReqEventAttributes []EventAttrib
}

/* Requester Identity */
type Identity struct {
	Email    string
	Username string
	Role     string
}

/* Request on Service Event */
type ServiceEvent struct {
	Identity
	RequestBody map[string]interface{}
	QueryParams map[string]interface{}
	PathParams  map[string]interface{}
}

type Service struct {
}

type ServiceHandler interface {
	ParseEvent(ServiceEvent)
}

/*
	@Exported Function
	Parse AWS Event to get the identity and requests objects
	event AWS  HTTP Event
	requestFmt Required event body format
*/
// func (sh ServiceEvent) ParseEvent(event events.APIGatewayProxyRequest, requestFmt map[string]interface{}) ServiceEvent {
// HTTP Request URL Endpiont
// requestEndpoint := event.RequestContext.ResourcePath
// }

/* @Internal Function */
func recursiveAttributeCheck(endpoint string, attributes map[string]interface{}, rEventSpec ReqEventSpec) bool {
	return true
}
