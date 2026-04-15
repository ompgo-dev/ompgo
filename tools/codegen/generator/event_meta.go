package generator

import (
	"sort"
	"strings"

	"github.com/ompgo-dev/ompgo/tools/codegen/mapping"
	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

type EventArgMeta struct {
	Name       string
	Type       string
	FieldName  string
	ParamName  string
	ParamLower string
}

type EventMeta struct {
	Name         string
	StructName   string
	MethodName   string
	HandlerName  string
	RegisterName string
	ExportName   string
	HandlerVar   string
	BadRet       string
	Args         []EventArgMeta
}

func buildEventMeta(event model.Event) EventMeta {
	structName := mapping.EventStructName(event.Name)
	methodName := mapping.EventMethodName(structName)
	handlerName := mapping.HandlerName(event.Name)
	register := mapping.RegisterName(event.Name)
	exportName := mapping.ExportName(event.Name)
	handlerVar := mapping.HandlerVarName(event.Name)

	args := make([]EventArgMeta, 0, len(event.Args))
	for _, arg := range event.Args {
		args = append(args, EventArgMeta{
			Name:       arg.Name,
			Type:       arg.Type,
			FieldName:  mapping.Capitalize(mapping.SanitizeName(arg.Name)),
			ParamName:  mapping.SafeGoName(arg.Name),
			ParamLower: strings.ToLower(arg.Name),
		})
	}

	return EventMeta{
		Name:         event.Name,
		StructName:   structName,
		MethodName:   methodName,
		HandlerName:  handlerName,
		RegisterName: register,
		ExportName:   exportName,
		HandlerVar:   handlerVar,
		BadRet:       event.BadRet,
		Args:         args,
	}
}

func orderedEvents(eventGroups map[string][]model.Event) []model.Event {
	groupNames := make([]string, 0, len(eventGroups))
	count := 0
	for group, events := range eventGroups {
		groupNames = append(groupNames, group)
		count += len(events)
	}
	sort.Strings(groupNames)

	ordered := make([]model.Event, 0, count)
	for _, group := range groupNames {
		ordered = append(ordered, eventGroups[group]...)
	}
	return ordered
}

func orderedEventMetas(eventGroups map[string][]model.Event) []EventMeta {
	ordered := orderedEvents(eventGroups)
	items := make([]EventMeta, 0, len(ordered))
	for _, event := range ordered {
		items = append(items, buildEventMeta(event))
	}
	return items
}
