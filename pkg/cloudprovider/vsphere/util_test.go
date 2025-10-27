/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vsphere

import (
	"strings"
	"testing"
)

func TestShouldProcessNode(t *testing.T) {
	tests := []struct {
		name       string
		providerID string
		expected   bool
	}{
		{
			name:       "empty ProviderID should be processed",
			providerID: "",
			expected:   true,
		},
		{
			name:       "vsphere ProviderID should be processed",
			providerID: "vsphere://422e4956-ad22-1139-6d72-59cc8f26bc90",
			expected:   true,
		},
		{
			name:       "aws ProviderID should not be processed",
			providerID: "aws:///us-west-2a/i-1234567890abcdef0",
			expected:   false,
		},
		{
			name:       "azure ProviderID should not be processed",
			providerID: "azure:///subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/my-vm",
			expected:   false,
		},
		{
			name:       "gcp ProviderID should not be processed",
			providerID: "gce://my-project/us-central1-a/my-instance",
			expected:   false,
		},
		{
			name:       "custom ProviderID should not be processed",
			providerID: "custom://some-node-id",
			expected:   false,
		},
		{
			name:       "malformed ProviderID should not be processed",
			providerID: "not-a-valid-provider-id",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldProcessNode(tt.providerID)
			if result != tt.expected {
				t.Errorf("ShouldProcessNode(%q) = %v, expected %v", tt.providerID, result, tt.expected)
			}
		})
	}
}

func TestInvalidProviderID(t *testing.T) {
	providerID := ""

	UUID := GetUUIDFromProviderID(providerID)

	if UUID != "" {
		t.Errorf("Should return an empty string")
	}
}

func TestUpperUUIDFromProviderID(t *testing.T) {
	tmpUUID := strings.ToUpper("423740e7-c66e-05e3-9d0b-9e1205b24d43")
	providerID := ProviderPrefix + tmpUUID

	UUID := GetUUIDFromProviderID(providerID)

	if UUID != "423740e7-c66e-05e3-9d0b-9e1205b24d43" {
		t.Errorf("Failed to extract UUID")
	}
}

func TestUUIDFromProviderID(t *testing.T) {
	providerID := "vsphere://423740e7-c66e-05e3-9d0b-9e1205b24d43"

	UUID := GetUUIDFromProviderID(providerID)

	if UUID != "423740e7-c66e-05e3-9d0b-9e1205b24d43" {
		t.Errorf("Failed to extract UUID")
	}
}

func TestUUIDFromUUID(t *testing.T) {
	UUIDOrg := "423740e7-c66e-05e3-9d0b-9e1205b24d43"

	UUIDNew := GetUUIDFromProviderID(UUIDOrg)

	if UUIDOrg != UUIDNew {
		t.Errorf("Failed to just return the UUID")
	}
}

func TestUUIDConvertInvalid(t *testing.T) {
	k8sUUID := ""

	biosUUID := ConvertK8sUUIDtoNormal(k8sUUID)

	if biosUUID != "" {
		t.Errorf("Should return empty string")
	}
}

func TestUUIDConvert(t *testing.T) {
	k8sUUID := "56492e42-22ad-3911-6d72-59cc8f26bc90"

	biosUUID := ConvertK8sUUIDtoNormal(k8sUUID)

	if biosUUID != "422e4956-ad22-1139-6d72-59cc8f26bc90" {
		t.Errorf("Failed to translate UUID")
	}
}

func TestUpperUUIDConvert(t *testing.T) {
	k8sUUID := strings.ToUpper("422e4956-ad22-1139-6d72-59cc8f26bc90")

	biosUUID := ConvertK8sUUIDtoNormal(k8sUUID)

	if biosUUID != "56492e42-22ad-3911-6d72-59cc8f26bc90" {
		t.Errorf("Failed to translate UUID")
	}
}

func TestUUIDConvertAndRevert(t *testing.T) {
	k8sUUID := "42278c9d-79fb-f2af-b060-d7f167fa261c"

	//converts
	tmpUUID := ConvertK8sUUIDtoNormal(k8sUUID)

	//reverts to original
	orgUUID := ConvertK8sUUIDtoNormal(tmpUUID)

	if orgUUID != "42278c9d-79fb-f2af-b060-d7f167fa261c" {
		t.Errorf("Failed to revert UUID")
	}
}

func TestArrayContainsCaseInsensitive(t *testing.T) {
	arr := []string{"First", "second", "THIRD"}

	if !ArrayContainsCaseInsensitive(arr, "First") {
		t.Errorf("Failed to find First")
	}

	if !ArrayContainsCaseInsensitive(arr, "firsT") {
		t.Errorf("Failed to find firsT")
	}

	if ArrayContainsCaseInsensitive(arr, "firs") {
		t.Errorf("Found firs")
	}

	if !ArrayContainsCaseInsensitive(arr, "second") {
		t.Errorf("Failed to find second")
	}

	if !ArrayContainsCaseInsensitive(arr, "Second") {
		t.Errorf("Failed to find Second")
	}

	if ArrayContainsCaseInsensitive(arr, "SecondInLine") {
		t.Errorf("Found SecondInLine")
	}

	if !ArrayContainsCaseInsensitive(arr, "THIRD") {
		t.Errorf("Failed to find THIRD")
	}

	if !ArrayContainsCaseInsensitive(arr, "third") {
		t.Errorf("Failed to find third")
	}

	if ArrayContainsCaseInsensitive(arr, "ThirdMakesACrowd") {
		t.Errorf("Found ThirdMakesACrowd")
	}
}
