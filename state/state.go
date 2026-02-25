package state

import "sync"

type Resource struct {
	Wood  int
	Brick int
	Sheep int
	Wheat int
	Ore   int
}

type State struct {
	PlayerScores    map[int]int
	PlayerResources map[int]Resource
	SharedResources Resource
	Turn            int
}

var (
	mu           sync.RWMutex
	currentState = State{
		PlayerScores:    map[int]int{},
		PlayerResources: map[int]Resource{},
	}
)

func GetState() State {
	mu.RLock()
	defer mu.RUnlock()
	return copyState(currentState)
}

func SetState(s State) {
	mu.Lock()
	defer mu.Unlock()
	currentState = copyState(s)
}

func UpdateState(updater func(*State)) {
	mu.Lock()
	defer mu.Unlock()
	updater(&currentState)
}

func copyState(s State) State {
	copyScores := make(map[int]int, len(s.PlayerScores))
	for k, v := range s.PlayerScores {
		copyScores[k] = v
	}

	copyResources := make(map[int]Resource, len(s.PlayerResources))
	for k, v := range s.PlayerResources {
		copyResources[k] = v
	}

	return State{
		PlayerScores:    copyScores,
		PlayerResources: copyResources,
		SharedResources: s.SharedResources,
		Turn:            s.Turn,
	}
}

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeFloat   FieldType = "float"
	FieldTypeObject  FieldType = "object"
	FieldTypeArray   FieldType = "array"
)

type SchemaField struct {
	Type   FieldType              `json:"type"`
	Fields map[string]SchemaField `json:"fields,omitempty"`
	Elem   *SchemaField           `json:"elem,omitempty"`
}

type StateSchema map[string]SchemaField

var (
	schemaMu     sync.RWMutex
	stateSchema  StateSchema
	schemaLocked bool
	dynamicState map[string]any
)

func ConfigureStateSchema(schema StateSchema, initial map[string]any) error {
	schemaMu.Lock()
	defer schemaMu.Unlock()
	if schemaLocked {
		return ErrSchemaLocked
	}
	if err := validateSchema(schema); err != nil {
		return err
	}
	stateSchema = schema
	schemaLocked = true
	if initial == nil {
		dynamicState = buildDefaultState(schema)
		return nil
	}
	if err := validateStateWithSchema(initial, schema); err != nil {
		return err
	}
	dynamicState = deepCopyMap(initial, schema)
	return nil
}

func GetStateSchema() (StateSchema, error) {
	schemaMu.RLock()
	defer schemaMu.RUnlock()
	if !schemaLocked {
		return nil, ErrSchemaNotSet
	}
	return copySchema(stateSchema), nil
}

func SetDynamicState(state map[string]any) error {
	schemaMu.Lock()
	defer schemaMu.Unlock()
	if !schemaLocked {
		return ErrSchemaNotSet
	}
	if err := validateStateWithSchema(state, stateSchema); err != nil {
		return err
	}
	dynamicState = deepCopyMap(state, stateSchema)
	return nil
}

func GetDynamicState() (map[string]any, error) {
	schemaMu.RLock()
	defer schemaMu.RUnlock()
	if !schemaLocked {
		return nil, ErrSchemaNotSet
	}
	return deepCopyMap(dynamicState, stateSchema), nil
}

func buildDefaultState(schema StateSchema) map[string]any {
	defaults := make(map[string]any, len(schema))
	for k, field := range schema {
		defaults[k] = buildDefaultValue(field)
	}
	return defaults
}

func buildDefaultValue(field SchemaField) any {
	switch field.Type {
	case FieldTypeString:
		return ""
	case FieldTypeInteger:
		return 0
	case FieldTypeFloat:
		return 0.0
	case FieldTypeObject:
		obj := make(map[string]any, len(field.Fields))
		for k, v := range field.Fields {
			obj[k] = buildDefaultValue(v)
		}
		return obj
	case FieldTypeArray:
		return []any{}
	default:
		return nil
	}
}

func validateSchema(schema StateSchema) error {
	if len(schema) == 0 {
		return ErrSchemaEmpty
	}
	for name, field := range schema {
		if name == "" {
			return ErrSchemaInvalid
		}
		if err := validateSchemaField(field); err != nil {
			return err
		}
	}
	return nil
}

func validateSchemaField(field SchemaField) error {
	switch field.Type {
	case FieldTypeString, FieldTypeInteger, FieldTypeFloat:
		return nil
	case FieldTypeObject:
		if len(field.Fields) == 0 {
			return ErrSchemaInvalid
		}
		for _, sub := range field.Fields {
			if err := validateSchemaField(sub); err != nil {
				return err
			}
		}
		return nil
	case FieldTypeArray:
		if field.Elem == nil {
			return ErrSchemaInvalid
		}
		return validateSchemaField(*field.Elem)
	default:
		return ErrSchemaInvalid
	}
}

func validateStateWithSchema(state map[string]any, schema StateSchema) error {
	if state == nil {
		return ErrStateInvalid
	}
	if len(state) != len(schema) {
		return ErrStateInvalid
	}
	for key, field := range schema {
		value, ok := state[key]
		if !ok {
			return ErrStateInvalid
		}
		if err := validateValue(value, field); err != nil {
			return err
		}
	}
	return nil
}

func validateValue(value any, field SchemaField) error {
	switch field.Type {
	case FieldTypeString:
		if _, ok := value.(string); !ok {
			return ErrStateInvalid
		}
		return nil
	case FieldTypeInteger:
		if isInteger(value) {
			return nil
		}
		return ErrStateInvalid
	case FieldTypeFloat:
		if isFloat(value) {
			return nil
		}
		return ErrStateInvalid
	case FieldTypeObject:
		obj, ok := value.(map[string]any)
		if !ok {
			return ErrStateInvalid
		}
		if len(obj) != len(field.Fields) {
			return ErrStateInvalid
		}
		for key, subField := range field.Fields {
			v, ok := obj[key]
			if !ok {
				return ErrStateInvalid
			}
			if err := validateValue(v, subField); err != nil {
				return err
			}
		}
		return nil
	case FieldTypeArray:
		arr, ok := value.([]any)
		if !ok {
			return ErrStateInvalid
		}
		if field.Elem == nil {
			return ErrSchemaInvalid
		}
		for _, v := range arr {
			if err := validateValue(v, *field.Elem); err != nil {
				return err
			}
		}
		return nil
	default:
		return ErrStateInvalid
	}
}

func isInteger(value any) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case float32:
		return float32(int64(v)) == v
	case float64:
		return float64(int64(v)) == v
	default:
		return false
	}
}

func isFloat(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

func deepCopyMap(state map[string]any, schema StateSchema) map[string]any {
	copyState := make(map[string]any, len(schema))
	for key, field := range schema {
		copyState[key] = deepCopyValue(state[key], field)
	}
	return copyState
}

func deepCopyValue(value any, field SchemaField) any {
	switch field.Type {
	case FieldTypeString:
		if s, ok := value.(string); ok {
			return s
		}
		return ""
	case FieldTypeInteger:
		return value
	case FieldTypeFloat:
		return value
	case FieldTypeObject:
		obj, _ := value.(map[string]any)
		copyObj := make(map[string]any, len(field.Fields))
		for k, subField := range field.Fields {
			copyObj[k] = deepCopyValue(obj[k], subField)
		}
		return copyObj
	case FieldTypeArray:
		arr, _ := value.([]any)
		copyArr := make([]any, 0, len(arr))
		for _, v := range arr {
			copyArr = append(copyArr, deepCopyValue(v, *field.Elem))
		}
		return copyArr
	default:
		return nil
	}
}

func copySchema(schema StateSchema) StateSchema {
	copySchema := make(StateSchema, len(schema))
	for k, v := range schema {
		copySchema[k] = copySchemaField(v)
	}
	return copySchema
}

func copySchemaField(field SchemaField) SchemaField {
	copyField := SchemaField{Type: field.Type}
	if len(field.Fields) > 0 {
		copyField.Fields = make(map[string]SchemaField, len(field.Fields))
		for k, v := range field.Fields {
			copyField.Fields[k] = copySchemaField(v)
		}
	}
	if field.Elem != nil {
		elemCopy := copySchemaField(*field.Elem)
		copyField.Elem = &elemCopy
	}
	return copyField
}

func ParseStateSchema(input map[string]any) (StateSchema, error) {
	if input == nil || len(input) == 0 {
		return nil, ErrSchemaEmpty
	}
	schema := make(StateSchema, len(input))
	for key, raw := range input {
		fieldMap, ok := raw.(map[string]any)
		if !ok {
			return nil, ErrSchemaInvalid
		}
		field, err := parseSchemaField(fieldMap)
		if err != nil {
			return nil, err
		}
		schema[key] = field
	}
	if err := validateSchema(schema); err != nil {
		return nil, err
	}
	return schema, nil
}

func parseSchemaField(input map[string]any) (SchemaField, error) {
	rawType, ok := input["type"].(string)
	if !ok || rawType == "" {
		return SchemaField{}, ErrSchemaInvalid
	}
	field := SchemaField{Type: FieldType(rawType)}
	switch field.Type {
	case FieldTypeString, FieldTypeInteger, FieldTypeFloat:
		return field, nil
	case FieldTypeObject:
		rawFields, ok := input["fields"].(map[string]any)
		if !ok || len(rawFields) == 0 {
			return SchemaField{}, ErrSchemaInvalid
		}
		field.Fields = make(map[string]SchemaField, len(rawFields))
		for key, raw := range rawFields {
			subMap, ok := raw.(map[string]any)
			if !ok {
				return SchemaField{}, ErrSchemaInvalid
			}
			subField, err := parseSchemaField(subMap)
			if err != nil {
				return SchemaField{}, err
			}
			field.Fields[key] = subField
		}
		return field, nil
	case FieldTypeArray:
		rawElem, ok := input["elem"].(map[string]any)
		if !ok {
			return SchemaField{}, ErrSchemaInvalid
		}
		elem, err := parseSchemaField(rawElem)
		if err != nil {
			return SchemaField{}, err
		}
		field.Elem = &elem
		return field, nil
	default:
		return SchemaField{}, ErrSchemaInvalid
	}
}
