package parser


type rule struct {
	Tag              string    `yaml:"tag"`
	FilterTags       []string  `yaml:"filter"`
	Condition        condition `yaml:"condition"`
	ExtractAttribute string    `yaml:"extract"`
}

type condition struct {
	Attribute string `yaml:"attribute"`
	Value     string `yaml:"value"`
}

func (r rule) getTag() string {
	return r.Tag
}

func (r rule) getFilterTags() []string {
	return r.FilterTags
}

func (r rule) getConditionAttribute() string {
	return r.Condition.getAttribute()
}

func (r rule) getConditionValue() string {
	return r.Condition.getValue()
}

func (r rule) getExtractAttribute() string {
	return r.ExtractAttribute
}

func (c condition) getAttribute() string {
	return c.Attribute
}

func (c condition) getValue() string {
	return c.Value
}
