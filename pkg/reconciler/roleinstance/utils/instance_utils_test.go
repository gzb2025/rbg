/*
Copyright 2026 The RBG Authors.

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

package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/rbgs/api/workloads/constants"
)

func TestGetPodComponentID(t *testing.T) {
	tests := []struct {
		name     string
		pod      *corev1.Pod
		expected int32
	}{
		{
			name: "valid id from label",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.ComponentIDLabelKey: "123",
					},
				},
			},
			expected: 123,
		},
		{
			name: "invalid non-numeric id in label",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.ComponentIDLabelKey: "abc",
					},
				},
			},
			expected: -1,
		},
		{
			name: "id in label exceeds int32 max",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.ComponentIDLabelKey: "2147483648",
					},
				},
			},
			expected: -1,
		},
		{
			name: "valid id from pod name",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "instance-web-789",
				},
			},
			expected: 789,
		},
		{
			name: "double hyphen pod name uses the final segment",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "instance-db--987",
				},
			},
			expected: 987,
		},
		{
			name: "invalid non-numeric id in pod name",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "instance-cache-xyz",
				},
			},
			expected: -1,
		},
		{
			name: "id in pod name exceeds int32 max",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "instance-store-2147483648",
				},
			},
			expected: -1,
		},
		{
			name: "int32 max value in label",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.ComponentIDLabelKey: "2147483647",
					},
				},
			},
			expected: math.MaxInt32,
		},
		{
			name: "int32 max value in pod name",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "instance-max-2147483647",
				},
			},
			expected: math.MaxInt32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetPodComponentID(tt.pod))
		})
	}
}
