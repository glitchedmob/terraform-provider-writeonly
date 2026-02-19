// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestWriteonlyDataResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = &WriteonlyDataResource{}
	var _ resource.ResourceWithImportState = &WriteonlyDataResource{}
}

func TestWriteonlyDataResourceModel(t *testing.T) {
	model := WriteonlyDataResourceModel{
		Id:             types.StringValue("test-id"),
		InputWo:        types.StringValue("test-input"),
		InputWoVersion: types.NumberValue(big.NewFloat(1)),
		Output:         types.StringValue("test-output"),
	}

	if model.Id.ValueString() != "test-id" {
		t.Errorf("expected id to be test-id, got %s", model.Id.ValueString())
	}
}

func TestWriteonlyDataResourceModelNull(t *testing.T) {
	model := WriteonlyDataResourceModel{
		Id:             types.StringValue("test-id"),
		InputWo:        types.StringNull(),
		InputWoVersion: types.NumberNull(),
		Output:         types.StringValue(""),
	}

	if !model.InputWo.IsNull() {
		t.Error("input_wo should be null")
	}
}

func TestGenerateRandomID(t *testing.T) {
	id1 := generateRandomID()
	id2 := generateRandomID()

	if id1 == id2 {
		t.Error("generated IDs should be unique")
	}

	if len(id1) != 31 {
		t.Errorf("expected id length 31, got %d", len(id1))
	}

	expectedPrefix := "writeonly_data_"
	if id1[:15] != expectedPrefix {
		t.Errorf("expected id to start with %s, got %s", expectedPrefix, id1[:15])
	}
}
