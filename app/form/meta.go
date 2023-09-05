package form

import "sismo-datagroup-service/app/model"

type MetaForm struct {
	Name              string              `json:"group_name"`
	Description       string              `json:"description,omitempty"`
	Spec              string              `json:"spec,omitempty"`
	GenerateFrequency string              `json:"generate_frequency,omitempty"`
	Handler           string              `json:"handler"`
	TwitterConfig     model.TwitterConfig `json:"config"`
}
