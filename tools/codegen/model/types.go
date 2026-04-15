package model

// APIFunction represents a function in the open.mp C API
type APIFunction struct {
	Ret    string     `json:"ret"`
	Name   string     `json:"name"`
	Params []APIParam `json:"params"`
}

// APIParam represents a parameter of a C API function
type APIParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Event represents an event in the open.mp event system
type Event struct {
	Name   string     `json:"name"`
	BadRet string     `json:"badret"`
	Args   []EventArg `json:"args"`
}

// EventArg represents an argument of an event
type EventArg struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
