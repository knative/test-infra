/*
Copyright 2021 The Knative Authors

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

package main

import (
	"os"
	"reflect"

	"gopkg.in/yaml.v3"

	"knative.dev/hack/schema/docs"
	"knative.dev/hack/schema/example"
	"knative.dev/hack/schema/schema"
)

func Example_kindLoremIpsum() {
	docs.SetRoot("knative.dev/hack/schema")

	t := reflect.TypeOf(example.LoremIpsum{})
	s := schema.GenerateForType(t)

	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	_ = enc.Encode(s)
	// Output:
	// type: object
	// properties:
	//   spec:
	//     description: Spec defines the desired state.
	//     type: object
	//     properties:
	//       aaa:
	//         description: Aaa is the first way.
	//         type: object
	//         required:
	//           - ccc
	//         properties:
	//           ccc:
	//             description: Ccc shows loop protection.
	//             type: object
	//           lorem:
	//             description: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam pellentesque eget arcu eget porta. Morbi ex urna, tincidunt in odio eget, hendrerit mattis odio. Sed vel augue rhoncus, rhoncus mi eget, tempor nisi. Nullam eleifend scelerisque pellentesque. Fusce efficitur urna mauris, sed suscipit sapien rhoncus et. Nunc viverra porta libero, mattis venenatis orci. Pellentesque molestie egestas iaculis. Donec sodales tristique ex, eget consectetur elit rutrum sed. Proin mollis, tellus vitae lobortis pretium, lacus dolor rhoncus tellus, at ultrices elit mauris vel enim. Suspendisse tempor ligula a est posuere, in egestas eros vehicula. Nulla mi magna, cursus in ultrices eget, porttitor eu odio. Nunc augue nisi, molestie at laoreet ut, sagittis a libero. Ut ullamcorper leo lectus, vel placerat ipsum lacinia vitae. Morbi commodo nibh neque, in ornare diam sodales ac. Defaults to true.
	//             type: boolean
	//           praesent:
	//             description: Praesent pulvinar consectetur enim. Aenean lobortis, eros quis molestie euismod, nisl nunc mattis quam, et gravida risus diam at nulla. Donec interdum, tortor a semper tincidunt, nibh odio euismod orci, rhoncus rhoncus purus lacus pharetra mi. Suspendisse placerat dignissim magna convallis dictum. Nulla facilisi. Vivamus sed tristique turpis.
	//             type: string
	//       bbb:
	//         description: Bbb is the second way.
	//         type: object
	//         required:
	//           - ccc
	//         properties:
	//           ccc:
	//             description: Ccc shows loop protection.
	//             type: object
	//           lorem:
	//             description: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam pellentesque eget arcu eget porta. Morbi ex urna, tincidunt in odio eget, hendrerit mattis odio. Sed vel augue rhoncus, rhoncus mi eget, tempor nisi. Nullam eleifend scelerisque pellentesque. Fusce efficitur urna mauris, sed suscipit sapien rhoncus et. Nunc viverra porta libero, mattis venenatis orci. Pellentesque molestie egestas iaculis. Donec sodales tristique ex, eget consectetur elit rutrum sed. Proin mollis, tellus vitae lobortis pretium, lacus dolor rhoncus tellus, at ultrices elit mauris vel enim. Suspendisse tempor ligula a est posuere, in egestas eros vehicula. Nulla mi magna, cursus in ultrices eget, porttitor eu odio. Nunc augue nisi, molestie at laoreet ut, sagittis a libero. Ut ullamcorper leo lectus, vel placerat ipsum lacinia vitae. Morbi commodo nibh neque, in ornare diam sodales ac. Defaults to true.
	//             type: boolean
	//           praesent:
	//             description: Praesent pulvinar consectetur enim. Aenean lobortis, eros quis molestie euismod, nisl nunc mattis quam, et gravida risus diam at nulla. Donec interdum, tortor a semper tincidunt, nibh odio euismod orci, rhoncus rhoncus purus lacus pharetra mi. Suspendisse placerat dignissim magna convallis dictum. Nulla facilisi. Vivamus sed tristique turpis.
	//             type: string
	//       ccc:
	//         description: Ccc is the third way.
	//         type: string
	//       maecenas:
	//         description: Maecenas tristique lobortis turpis, nec varius mauris vestibulum nec. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Vivamus non dapibus magna.
	//         type: string
	//       sed:
	//         description: Sed euismod nunc ac sollicitudin ornare.
	//         type: string
	//       verboseTypes:
	//         description: VerboseTypes shows an example of a ton of types.
	//         type: object
	//         properties:
	//           float32:
	//             description: AFloat32 is a field with the type float32.
	//             type: number
	//             format: float
	//           float64:
	//             description: AFloat64 is a field with the type float64.
	//             type: number
	//             format: double
	//           int8:
	//             description: AInt8 is a field with the type int8.
	//             type: integer
	//             maximum: 127
	//           int16:
	//             description: AInt16 is a field with the type int16.
	//             type: integer
	//             maximum: 32767
	//           int32:
	//             description: AInt32 is a field with the type int32.
	//             type: integer
	//             format: int32
	//           int63:
	//             description: AInt64 is a field with the type int64.
	//             type: integer
	//             format: int64
	//           map:
	//             description: AMap is a field with the type map.
	//             type: object
	//             x-kubernetes-preserve-unknown-fields: true
	//           uint:
	//             description: AUint is a field with the type uint.
	//             type: integer
	//             minimum: 0
	//           uint8:
	//             description: Uint8 is a field with the type uint8.
	//             type: integer
	//             maximum: 255
	//             minimum: 0
	//           uint16:
	//             description: AUint16 is a field with the type uint16.
	//             type: integer
	//             maximum: 65535
	//             minimum: 0
	//           uint32:
	//             description: AUint32 is a field with the type uint32.
	//             type: integer
	//             format: int64
	//             maximum: 4.294967295e+09
	//             minimum: 0
	//           uint64:
	//             description: AUint64 is a field with the type uint64.
	//             type: integer
	//             format: int64
	//             minimum: 0
	//   status:
	//     description: Status represents the current state. This data may be out of date.
	//     type: object
	//     properties:
	//       aliquam:
	//         description: Aliquam consequat placerat ante, eu ullamcorper purus consectetur quis.
	//         type: array
	//         items:
	//           type: string
	//       donec:
	//         description: Donec mollis purus id ipsum varius, sit amet ultricies elit laoreet.
	//         type: boolean
	//       duis:
	//         description: Duis vulputate purus sed porta tristique.
	//         type: array
	//         items:
	//           type: object
	//           properties:
	//             aenean:
	//               description: Aenean a purus porttitor nulla rhoncus posuere.
	//               type: string
	//             maecenas:
	//               description: Maecenas sed velit ac velit fringilla dapibus.
	//               type: string
	//       luctus:
	//         description: Luctus leo vitae ipsum fermentum, vitae pellentesque sapien finibus.
	//         type: integer
	//         format: int32
	//       suspendisse:
	//         description: Suspendisse ipsum risus, porttitor a auctor vel, maximus eu mi.
	//         type: string
}
