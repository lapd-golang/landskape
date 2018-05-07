package model

import (
	"cloud.google.com/go/datastore"
)

// System is the generic name for a IT landscape object.
// Examples are: Webservice, Database schema, Ftp server, Third party solution
type System struct {
	// populated from DBKey
	ID         string
	Attributes []Attribute `datastore:",flatten"`
	// internal
	DBKey *datastore.Key `datastore:"__key__" json:"-"`
	Journal
}

func NewSystem(id string) *System {
	key := datastore.NameKey("landskape.System", id, nil)
	key.Namespace = "landskape"
	return &System{DBKey: key}
}

func (s System) AttributeList() []Attribute { return s.Attributes }

func (s *System) DeleteAttribute(name string) {
	s.SetAttribute(name, "")
}

func (s *System) SetAttribute(name, value string) {
	if len(name) == 0 {
		return
	}
	if len(value) == 0 {
		// remove it
		without := []Attribute{}
		for _, each := range s.Attributes {
			if each.Name != name {
				without = append(without, each)
			}
		}
		s.Attributes = without
		return
	}
	// replace or add
	for i, each := range s.Attributes {
		if each.Name == name {
			s.Attributes[i] = Attribute{Name: name, Value: value}
			return
		}
	}
	// not found, add it
	s.Attributes = append(s.Attributes, Attribute{Name: name, Value: value})
}

// Attribute is a generic key-value pair of strings
// Each attribute has its own lifecyle to track value changes
type Attribute struct {
	Journal
	Name, Value string
}

// AttributeValue finds the value of an attribute for a given name, return empty string if not found
func AttributeValue(holder AttributesHolder, name string) string {
	for _, each := range holder.AttributeList() {
		if each.Name == name {
			return each.Value
		}
	}
	return ""
}

// Connection is the generic name for a logical connection between 2 IT landscape object.
// From and To refer to the Id of the System.
// Example of Type are:  http, https, aq, jdbc, ftp, smtp
type Connection struct {
	Journal
	From, To   string
	Type       string      `datastore:"Type,noindex"`
	Attributes []Attribute `datastore:",flatten"`
	// populated
	FromSystem, ToSystem System `datastore:"-" json:"-"`
	// internal
	DBKey *datastore.Key `datastore:"__key__" json:"-"`
}

func (c Connection) Validate() error {
	return nil // TODO
}

func (c Connection) AttributeList() []Attribute {
	return c.Attributes
}

func (c *Connection) SetAttribute(name, value string) {
	if len(name) == 0 {
		return
	}
	if len(value) == 0 {
		// remove it
		without := []Attribute{}
		for _, each := range c.Attributes {
			if each.Name != name {
				without = append(without, each)
			}
		}
		c.Attributes = without
		return
	}
	// replace or add
	for i, each := range c.Attributes {
		if each.Name == name {
			c.Attributes[i] = Attribute{Name: name, Value: value}
			return
		}
	}
	// not found, add it
	c.Attributes = append(c.Attributes, Attribute{Name: name, Value: value})
}

func (c *Connection) DeleteAttribute(name string) {
	c.SetAttribute(name, "")
}
