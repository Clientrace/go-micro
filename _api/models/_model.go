package _api

import (
	"encoding/json"
)


type ModelAttribInt interface {
	DType string, 			// Data type
	IsRequired bool,		// Is field required
}

type ModelAttribBool interface {
	DType string, 			// Data type
	IsRequired bool,
}

type Model interface {
	
}



