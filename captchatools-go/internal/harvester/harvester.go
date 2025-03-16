package harvester

import "context"

// Harvester defines the interface for interacting with captcha solving services.
// This interface enables a consistent API across different captcha solving
// providers, allowing for easy switching between services while maintaining
// the same code structure.
//
// The interface follows a design pattern similar to the Python implementation
// of the Captcha-Tools library, where each solving service implements this
// common interface, enabling polymorphic behavior.
type Harvester interface {
	// GetToken requests a captcha solution from the solving service using
	// the configured parameters and any additional options provided.
	//
	// Parameters:
	//   options - Variable number of TokenOption functions that customize the request
	//
	// Returns:
	//   - A CaptchaAnswer containing the solution, or
	//   - An error if the solution could not be obtained
	GetToken(...TokenOption) (CaptchaAnswer, error)

	// GetTokenWithContext is similar to GetToken but accepts a context for
	// cancellation, timeout control, and other context-aware operations.
	//
	// Parameters:
	//   ctx     - Context for controlling cancellation and timeouts
	//   options - Variable number of TokenOption functions that customize the request
	//
	// Returns:
	//   - A pointer to a CaptchaAnswer containing the solution, or
	//   - An error if the solution could not be obtained or the context was canceled
	GetTokenWithContext(context.Context, ...TokenOption) (CaptchaAnswer, error)

	// GetBalance retrieves the current account balance from the captcha solving service.
	//
	// Returns:
	//   - The current balance as a float32, or
	//   - An error if the balance could not be retrieved
	GetBalance() (float32, error)
}
