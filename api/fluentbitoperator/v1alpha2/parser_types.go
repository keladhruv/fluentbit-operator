/*

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

package v1alpha2

import (
	"bytes"
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/parser"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ParserSpec defines the desired state of Parser
type ParserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// JSON defines json parser configuration.
	JSON *parser.JSON `json:"json,omitempty"`
	// Regex defines regex parser configuration.
	Regex *parser.Regex `json:"regex,omitempty"`
	// LTSV defines ltsv parser configuration.
	LTSV *parser.LSTV `json:"ltsv,omitempty"`
	// Logfmt defines logfmt parser configuration.
	Logfmt *parser.Logfmt `json:"logfmt,omitempty"`
	// Decoders are a built-in feature available through the Parsers file, each Parser definition can optionally set one or multiple decoders.
	// There are two type of decoders type: Decode_Field and Decode_Field_As.
	Decoders []Decorder `json:"decoders,omitempty"`
}

type Decorder struct {
	// If the content can be decoded in a structured message,
	// append that structure message (keys and values) to the original log message.
	DecodeField string `json:"decodeField,omitempty"`
	// Any content decoded (unstructured or structured) will be replaced in the same key/value,
	// no extra keys are added.
	DecodeFieldAs string `json:"decodeFieldAs,omitempty"`
}

// +kubebuilder:object:root=true
// +genclient

// Parser is the Schema for the parsers API
type Parser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ParserSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ParserList contains a list of Parser
type ParserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Parser `json:"items"`
}

func (list ParserList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[PARSER]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", item.Name))
			buf.WriteString(fmt.Sprintf("    Format    %s\n", p.Name()))

			kvs, err := p.Params(sl)
			if err != nil {
				return err
			}
			buf.WriteString(kvs.String())

			for _, decorder := range item.Spec.Decoders {
				if decorder.DecodeField != "" {
					buf.WriteString(fmt.Sprintf("    Decode_Field    %s\n", decorder.DecodeField))
				}
				if decorder.DecodeFieldAs != "" {
					buf.WriteString(fmt.Sprintf("    Decode_Field_As    %s\n", decorder.DecodeFieldAs))
				}
			}
			return nil
		}

		for i := 0; i < reflect.ValueOf(item.Spec).NumField()-1; i++ {
			p, _ := reflect.ValueOf(item.Spec).Field(i).Interface().(plugins.Plugin)
			if err := merge(p); err != nil {
				return "", err
			}
		}
	}

	return buf.String(), nil
}

func init() {
	SchemeBuilder.Register(&Parser{}, &ParserList{})
}
