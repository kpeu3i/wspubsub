package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// HubSubscriptionChannelRequired returned when trying to subscribe with an empty channels list.
type HubSubscriptionChannelRequired struct {
	message string
}

// HubSubscriptionChannelRequired implements an error interface.
func (e *HubSubscriptionChannelRequired) Error() string {
	return fmt.Sprintf("wspubsub: %s", e.message)
}

// NewHubSubscriptionChannelRequiredError initializes a new HubSubscriptionChannelRequired.
func NewHubSubscriptionChannelRequiredError() *HubSubscriptionChannelRequired {
	return &HubSubscriptionChannelRequired{message: "at least one subscription channel is required"}
}

// IsHubSubscriptionChannelRequiredError checks if error type is HubSubscriptionChannelRequired.
func IsHubSubscriptionChannelRequiredError(err error) (*HubSubscriptionChannelRequired, bool) {
	v, ok := errors.Cause(err).(*HubSubscriptionChannelRequired)

	return v, ok
}
