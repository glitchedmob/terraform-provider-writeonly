// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestWoDataResourceImplementsResource(t *testing.T) {
	var _ resource.Resource = &WoDataResource{}
	var _ resource.ResourceWithImportState = &WoDataResource{}
}

func TestWoDataResourceModel(t *testing.T) {
	model := WoDataResourceModel{
		Id:             types.StringValue("test-id"),
		InputWo:        types.StringValue("test-input"),
		InputWoVersion: types.NumberValue(big.NewFloat(1)),
		Output:         types.StringValue("test-output"),
	}

	if model.Id.ValueString() != "test-id" {
		t.Errorf("expected id to be test-id, got %s", model.Id.ValueString())
	}
}

func TestWoDataResourceModelNull(t *testing.T) {
	model := WoDataResourceModel{
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

	if len(id1) != 24 {
		t.Errorf("expected id length 24, got %d", len(id1))
	}

	expectedPrefix := "wo_data_"
	if id1[:8] != expectedPrefix {
		t.Errorf("expected id to start with %s, got %s", expectedPrefix, id1[:8])
	}
}
