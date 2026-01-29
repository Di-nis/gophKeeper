package data

// RawGetter - интерфейс получения raw значения.
type RawGetter interface {
	GetRaw() string
}

// AliasSetter - интерфейс установки алиаса.
type AliasSetter interface {
	SetAlias(string)
}

// // getRawValue - приведение к единому типу данных.
func getRawValue(data any) (string, error) {
	v, ok := data.(RawGetter)
	if !ok {
		return "", ErrGetRawValue
	}
	raw := v.GetRaw()
	return raw, nil
}

// updateData - обновление данных.
func updateData(data any, alias string) error {
	v, ok := data.(AliasSetter)
	if !ok {
		return ErrDataUnsupportedType
	}
	v.SetAlias(alias)
	return nil
}
