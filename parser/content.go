package parser

var numberOfAdditionalFields = 2

type ParsedContent struct {
	id int
	website string
	content map[string]string
}

func newParsedContent() ParsedContent {
	return ParsedContent{
		content: make(map[string]string),
	}
}

func (p ParsedContent) setField(field string, value string) {
	p.content[field] = value
}

func (p *ParsedContent) setWebsite(website string) {
	p.website = website
}

func (p *ParsedContent) setId(id int) {
	p.id = id
}

func (p ParsedContent) size() int {
	return numberOfAdditionalFields + len(p.content)
}

func (p ParsedContent) GetFieldNamesAndValues() ([]string, []interface{}){
	size := p.size()
	values := make([]interface{}, 0, size)
	fieldNames := make([]string, 0, size)
	fieldNames = append(fieldNames, "id", "website")
	values = append(values, p.id, p.website)
	for fieldName, value := range p.content {
		fieldNames = append(fieldNames, fieldName)
		values = append(values, value)
	}
	return fieldNames, values
}
