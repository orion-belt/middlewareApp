// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NetworkProbeDestination Network Probe Destination
//
// swagger:model network_probe_destination
type NetworkProbeDestination struct {

	// destination details
	// Required: true
	DestinationDetails *NetworkProbeDestinationDetails `json:"destination_details"`

	// destination id
	// Required: true
	DestinationID NetworkProbeDestinationID `json:"destination_id"`
}

// Validate validates this network probe destination
func (m *NetworkProbeDestination) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDestinationDetails(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDestinationID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkProbeDestination) validateDestinationDetails(formats strfmt.Registry) error {

	if err := validate.Required("destination_details", "body", m.DestinationDetails); err != nil {
		return err
	}

	if m.DestinationDetails != nil {
		if err := m.DestinationDetails.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destination_details")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destination_details")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkProbeDestination) validateDestinationID(formats strfmt.Registry) error {

	if err := validate.Required("destination_id", "body", NetworkProbeDestinationID(m.DestinationID)); err != nil {
		return err
	}

	if err := m.DestinationID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("destination_id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("destination_id")
		}
		return err
	}

	return nil
}

// ContextValidate validate this network probe destination based on the context it is used
func (m *NetworkProbeDestination) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDestinationDetails(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDestinationID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkProbeDestination) contextValidateDestinationDetails(ctx context.Context, formats strfmt.Registry) error {

	if m.DestinationDetails != nil {
		if err := m.DestinationDetails.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destination_details")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destination_details")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkProbeDestination) contextValidateDestinationID(ctx context.Context, formats strfmt.Registry) error {

	if err := m.DestinationID.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("destination_id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("destination_id")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *NetworkProbeDestination) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetworkProbeDestination) UnmarshalBinary(b []byte) error {
	var res NetworkProbeDestination
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
