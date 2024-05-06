package codegen

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type RouterGroup struct {
	Name                 string
	Description          string
	OperationDefinitions []OperationDefinition
}

// RouterGroupDefinitions provides router groups and sets back reference in the grouped operation definitions.
func RouterGroups(swagger *openapi3.T, ops []OperationDefinition) ([]RouterGroup, error) {
	routerGroups := make([]RouterGroup, 0)

	for _, tag := range swagger.Tags {
		if extension, ok := tag.Extensions[extGoRouterGroup]; ok {
			goRouterGroup, err := extParseGoRouterGroup(extension)

			if err != nil {
				return nil, fmt.Errorf("invalid value for %q: %w", extGoRouterGroup, err)
			}

			if goRouterGroup {
				routerGroups = append(routerGroups, RouterGroup{Name: tag.Name, Description: tag.Description})
			}
		}
	}

	if len(routerGroups) == 0 {
		return routerGroups, nil
	}

	for i := 0; i < len(ops); i++ {
		found := false

		for _, opTag := range ops[i].Spec.Tags {
			for j := 0; j < len(routerGroups); j++ {
				if opTag == routerGroups[j].Name {
					if found {
						return nil, fmt.Errorf("%s OperationId has two tags with %s extension set", ops[i].OperationId, extGoRouterGroup)
					}

					routerGroups[j].OperationDefinitions = append(routerGroups[j].OperationDefinitions, ops[i])
					ops[i].RouterGroup = &routerGroups[j]

					found = true
				}
			}
		}
	}

	return routerGroups, nil
}
