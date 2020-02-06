package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type HubSubscriptionChannelRequired struct {
	message string
}

func (e *HubSubscriptionChannelRequired) Error() string {
	return fmt.Sprintf("wspubsub: %s", e.message)
}

func NewHubSubscriptionChannelRequiredError() *HubSubscriptionChannelRequired {
	return &HubSubscriptionChannelRequired{message: "at least one subscription channel is required"}
}

func IsHubSubscriptionChannelRequiredError(err error) (*HubSubscriptionChannelRequired, bool) {
	v, ok := errors.Cause(err).(*HubSubscriptionChannelRequired)

	return v, ok
}
