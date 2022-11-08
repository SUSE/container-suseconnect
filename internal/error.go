// Copyright (c) 2022 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package containersuseconnect

const (
	// GetCredentialsError indicates a failure to retrieve or parse
	// credentials
	GetCredentialsError = iota
	// NetworkError is a placeholder for generic network communication
	// errors
	NetworkError
	// InstalledProductError signals issues with the installed products
	InstalledProductError
	// SubscriptionServerError means that the subscription server did
	// something unexpected
	SubscriptionServerError
	// SubscriptionError marks issues with the actual subscription
	SubscriptionError
	// RepositoryError indicates that there is something wrong with the
	// repository that the subscription server gave us
	RepositoryError
)

// SuseConnectError is a custom error type allowing us to distinguish between
// different error kinds via the `ErrorCode` field
type SuseConnectError struct {
	ErrorCode int
	message   string
}

func (s *SuseConnectError) Error() string {
	return s.message
}
