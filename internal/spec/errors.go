package spec

import (
	"fmt"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
)

const (
	ErrKindInvalidRequest           ftag.Kind = "invalid_request"
	ErrKindInvalidClient            ftag.Kind = "invalid_client"
	ErrKindInvalidGrant             ftag.Kind = "invalid_grant"
	ErrKindUnauthorizedClient       ftag.Kind = "unauthorized_client"
	ErrKindUnsupportedResponseType  ftag.Kind = "unsupported_response_type"
	ErrKindUnsupportedGrantType     ftag.Kind = "unsupported_grant_type"
	ErrKindInvalidScope             ftag.Kind = "invalid_scope"
	ErrKindInsufficientScope        ftag.Kind = "insufficient_scope"
	ErrKindAccessDenied             ftag.Kind = "access_denied"
	ErrKindInvalidRequestURI        ftag.Kind = "invalid_request_uri"
	ErrKindInvalidRequestObject     ftag.Kind = "invalid_request_object"
	ErrKindRequestNotSupported      ftag.Kind = "request_not_supported"
	ErrKindRequestURINotSupported   ftag.Kind = "request_uri_not_supported"
	ErrKindRegistrationNotSupported ftag.Kind = "registration_not_supported"
	ErrKindResourceNotFound         ftag.Kind = "resource_not_found"
	ErrKindLoginRequired            ftag.Kind = "login_required"
	ErrKindSelectAccountRequired    ftag.Kind = "account_selection_required"
	ErrKindConsentRequired          ftag.Kind = "consent_required"
	ErrKindInteractionRequired      ftag.Kind = "interaction_required"
	ErrKindServerError              ftag.Kind = "server_error"
)

// GetErrorKind extracts the closest ftag.Kind tagged on the error. If not tagged, defaults to ErrKindServerError.
func GetErrorKind(err error) ftag.Kind {
	kind := ftag.Get(err)
	if len(kind) == 0 {
		return ErrKindServerError
	}
	return kind
}

// GetErrorStatus returns the corresponding HTTP status code to the tagged ftag.Kind.
func GetErrorStatus(kind ftag.Kind) int {
	switch kind {
	case ErrKindInvalidRequest,
		ErrKindInvalidGrant,
		ErrKindUnauthorizedClient,
		ErrKindInvalidScope,
		ErrKindInvalidRequestURI,
		ErrKindInvalidRequestObject,
		ErrKindRequestNotSupported,
		ErrKindRequestURINotSupported,
		ErrKindRegistrationNotSupported,
		ErrKindUnsupportedResponseType,
		ErrKindUnsupportedGrantType,
		ErrKindLoginRequired,
		ErrKindSelectAccountRequired,
		ErrKindConsentRequired,
		ErrKindInteractionRequired:
		return 400
	case ErrKindInvalidClient:
		return 401
	case ErrKindAccessDenied, ErrKindInsufficientScope:
		return 403
	case ErrKindResourceNotFound:
		return 404
	case ErrKindServerError:
		return 500
	default:
		panic(fmt.Sprintf("undefined status for error kind [%s]", kind))
	}
}

// GetErrorMessage returns the tagged external error message on the error, or a default message associated with the
// ftag.Kind tagged on the error. The returned value might be empty if the error has neither external error message
// nor supported ftag.Kind associated with it.
func GetErrorMessage(err error) string {
	issue := fmsg.GetIssue(err)
	if len(issue) > 0 {
		return issue
	}

	switch GetErrorKind(err) {
	case ErrKindInvalidRequest:
		return "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed."
	case ErrKindInvalidClient:
		return "Client authentication failed."
	case ErrKindInvalidGrant:
		return "The provided authorization grant or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."
	case ErrKindUnauthorizedClient:
		return "The client is not authorized to request an authorization code using this method."
	case ErrKindUnsupportedResponseType:
		return "The authorization server does not support obtaining an authorization code using this method."
	case ErrKindUnsupportedGrantType:
		return "The authorization grant type is not supported by the authorization server."
	case ErrKindInvalidScope:
		return "The requested scope is invalid, unknown, malformed, or exceeds the scope granted by the resource owner."
	case ErrKindInsufficientScope:
		return "The protected resource requires one or more scopes that exceeds the extent of grant."
	case ErrKindAccessDenied:
		return "he resource owner or authorization server denied the request."
	case ErrKindInvalidRequestURI:
		return "The request_uri in the Authorization Request returns an error or contains invalid data."
	case ErrKindInvalidRequestObject:
		return "The request parameter contains an invalid Request Object."
	case ErrKindRequestNotSupported:
		return "The OP does not support use of the request parameter."
	case ErrKindRequestURINotSupported:
		return "The OP does not support use of the request_uri parameter."
	case ErrKindResourceNotFound:
		return "The requested resource is not found on the server."
	case ErrKindRegistrationNotSupported:
		return "The OP does not support use of the registration parameter."
	case ErrKindLoginRequired:
		return "The Authorization Server requires End-User authentication."
	case ErrKindSelectAccountRequired:
		return "The End-User is REQUIRED to select a session at the Authorization Server."
	case ErrKindConsentRequired:
		return "The Authorization Server requires End-User consent."
	case ErrKindInteractionRequired:
		return "The Authorization Server requires End-User interaction of some form to proceed."
	case ErrKindServerError:
		return "The authorization server encountered an unexpected condition that prevented it from fulfilling the request."
	default:
		return ""
	}
}
